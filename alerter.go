package main

import (
	"encoding/json"
    "fmt"
    "net/url"
    "time"
    "bytes"

	"github.com/boltdb/bolt"
)

func AlertWorker(alerts <-chan string, db *bolt.DB) {
    for u := range alerts {
        required, err := IsAlertRequired(u, db)
        if err != nil {
            //TODO: add logging here
            continue
        }

        if required {
            fmt.Printf("Imitating alert for %s", u)
        }
    }
}

func IsAlertRequired(u string, db *bolt.DB) (bool, error) {
    r, err := url.Parse(u)
    if err != nil {
        return false, err
    }

    domain := r.Host
    count := 0
    fail_count := 0

    db.View(func(tx *bolt.Tx) error {
        c := tx.Bucket([]byte(domain)).Cursor()

        t := time.Now()
        ts := t.Add(-2 * time.Minute)

        min := []byte(ts.Format(time.RFC3339))
        max := []byte(t.Format(time.RFC3339))

        for k,v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
            count += 1
            var sr StatRecord
            if err := json.Unmarshal(v, &sr); err != nil {
                //TODO: add logging
                c.Next()
            }

            if sr.Rstatus > 400 {
                fail_count += 1
            }
        }

        return nil
    })

    if count - fail_count > 1 {
        return true, nil
    }

    return false, nil
}

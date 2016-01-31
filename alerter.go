package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/boltdb/bolt"
)

func AlertWorker(urls []string, db *bolt.DB) {
    alerts := make(map[string]int)

	for {
        for _, u := range urls {
		    status, err := IsOkUrl(u, db)
		    if err != nil {
			    //TODO: add logging here
			    continue
		    }

            alert_for_url := alerts[u]

		    if !status && alert_for_url == 0{
		        fmt.Printf("Imitating alert for %s\n", u)
                alerts[u] = 1
		    } else if status && alert_for_url > 0 {
                fmt.Printf("Imitating recovery for %s\n", u)
                alerts[u] = 0
            }

        }

        time.Sleep(30 * time.Second)
	}
}

// Check last m minutes for url in database
// return true if StatusCode < 400 for m time
func IsOkUrl(u string, db *bolt.DB) (bool, error) {
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

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
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

    // at least two checks
	if count > 2 && count - fail_count <= 1 {
		return false, nil
	}

	return true, nil
}

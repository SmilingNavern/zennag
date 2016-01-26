package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/boltdb/bolt"
)

type StatRecord struct {
	Uri     string
	Rstatus int
	Rtime   time.Duration
}

func (sr StatRecord) String() string {
	return fmt.Sprintf("%s %d %f",
		sr.Uri,
		sr.Rstatus,
		sr.Rtime.Seconds())
}

func OpenDb() *bolt.DB {
	db, err := bolt.Open("stat.db", 0600, nil)
	if err != nil {
		panic("Can't open db for statistic")
	}

	return db
}

func SaveStatus(db *bolt.DB, u *url.URL, resp *http.Response, d time.Duration) error {
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(u.Host))
		if err != nil {
			return err
		}

		// dirty hack for /
		var uPath string
		if u.Path == "" {
			uPath = "/"
		} else {
			uPath = u.Path
		}
		sr := &StatRecord{uPath, resp.StatusCode, d}
		data, err := json.Marshal(sr)
		if err != nil {
			return err
		}

		timeFormat := "2006-01-02 15:04:05"
		cur := time.Now()
		t := cur.Format(timeFormat)

		if err := b.Put([]byte(t), data); err != nil {
			return err
		}

		return nil

	}); err != nil {
		return err
	}

	return nil

}

func ShowStatus(db *bolt.DB, u string) error {
	ur, err := url.Parse(u)
	if err != nil {
		return err
	}

	fmt.Printf("Host %s status:\n", ur.Host)
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ur.Host))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var sr StatRecord
			err := json.Unmarshal(v, &sr)
			if err != nil {
				fmt.Println(err)
				c.Next()
			}
			fmt.Printf("%s => %s\n", k, sr)
		}
		fmt.Printf("\n\n")
		return nil
	})
	return nil
}

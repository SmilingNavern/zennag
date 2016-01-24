package main

import (
	"fmt"

	"github.com/boltdb/bolt"
)

func OpenDb() *bolt.DB {
	db, err := bolt.Open("stat.db", 0600, nil)
	if err != nil {
		panic("Can't open db for statistic")
	}

	return db
}

func SaveStatus(db *bolt.DB, url string, data string) error {
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Stats"))
		if err != nil {
			return err
		}

		if err := b.Put([]byte(url), []byte(data)); err != nil {
			return err
		}

		return nil

	}); err != nil {
		return err
	}

	return nil

}

func ShowStatus(db *bolt.DB) error {
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Stats"))
		b.ForEach(func(k, v []byte) error {
			fmt.Printf("The %s status: %s\n", k, v)
			return nil
		})
		return nil
	})
	return nil
}

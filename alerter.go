package main

import (
    "fmt"
    "net/url"

	"github.com/boltdb/bolt"
)

func AlertWorker(urls []string, db *bolt.DB) {
    for {
        for _, u := range urls {
            result, err := CheckDb(u)
        }
    }
}

func CheckDb(u string) bool, error{
    r, err := url.Parse(u)
    if err != nil {
        return err
    }

    domain := r.Host
}

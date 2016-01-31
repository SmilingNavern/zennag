package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

func Worker(jobs <-chan string, alerts chan<- string, db *bolt.DB) {
	for j := range jobs {
		ts := time.Now()
		resp, err := http.Get(j)
		if err != nil {
			fmt.Println(err)
			continue
		}

		//get HTTP request time
		te := time.Now()
		dur := te.Sub(ts)

		defer resp.Body.Close()

		u, err := url.Parse(j)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if err := SaveStatus(db, u, resp, dur); err != nil {
			fmt.Println(err)
		}
	}
}

func main() {

	workerPoolSize := 3

	config := ParseConfig()
	urls := config.Urls
	timeout := config.Timeout
	db := OpenDb()
    err := PrepareDb(db, urls)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

	defer db.Close()

	jobs := make(chan string, 100)
	alerts := make(chan string, 100)

	// only show stored info in db
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		for _, u := range urls {
			if err := ShowStatus(db, u); err != nil {
				fmt.Println(err)
			}
		}
		os.Exit(0)
	}

	for w := 1; w <= workerPoolSize; w++ {
		go Worker(jobs, alerts, db)
	}

	go AlertWorker(urls, db)

	for {
		for i := 0; i < len(urls); i++ {
			jobs <- urls[i]
		}

		time.Sleep(timeout * time.Second)
	}

	close(jobs)

	for _, u := range urls {
		ShowStatus(db, u)
	}
}

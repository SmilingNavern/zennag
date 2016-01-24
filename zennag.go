package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

func Worker(jobs <-chan string, answers chan<- string, db *bolt.DB) {
	for j := range jobs {
		resp, err := http.Get(j)
		if err != nil {
			fmt.Println(err)
			continue
		}

		defer resp.Body.Close()

		if err := SaveStatus(db, j, resp.Status); err != nil {
			fmt.Println(err)
		}
		answers <- fmt.Sprintf("%s: %s\n", j, resp.Status)
	}
}

func main() {

	workerPoolSize := 3

	config := ParseConfig()
	urls := config.Urls
	db := OpenDb()

	defer db.Close()

	jobs := make(chan string, 100)
	answers := make(chan string, 100)

	// only show stored info in db
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		if err := ShowStatus(db); err != nil {
			fmt.Println(err)
		}
		os.Exit(0)
	}

	for w := 1; w <= workerPoolSize; w++ {
		go Worker(jobs, answers, db)
	}

	for {
		for i := 0; i < len(urls); i++ {
			jobs <- urls[i]
		}

		for a := 1; a <= len(urls); a++ {
			fmt.Println(<-answers)
		}

		time.Sleep(30 * time.Second)
	}

	close(jobs)

	ShowStatus(db)
}

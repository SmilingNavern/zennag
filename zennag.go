package main

import (
    "fmt"
    "net/http"
    "time"
)

func worker(w int, jobs<-chan string, answers chan<- string) {
    for j := range jobs {
        resp, err := http.Get(j)
        if err != nil {
            fmt.Println(err)
            continue
        }

        defer resp.Body.Close()
        answers <- fmt.Sprintf("%s: %s\n", j, resp.Status)
    }
}

func main() {
    workerPoolSize := 3;

    config := ParseConfig()
    urls := config.Urls

    jobs := make(chan string, 100)
    answers := make(chan string, 100)

    for w := 1; w <= workerPoolSize; w++ {
        go worker(w, jobs, answers)
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

}

package main

import (
    "fmt"
    "net/http"
    //"io/ioutil"
)

func worker(w int, jobs chan string, answers chan string) {
    for j := range jobs {
        resp, err := http.Get(j)
        if err != nil {
            fmt.Println(err)
        }

        defer resp.Body.Close()
        //body, _ := ioutil.ReadAll(resp.Body)

        answers <- fmt.Sprintf("%s: %s\n", j, resp.Status)
    }
}

func main() {
    workerPoolSize := 3;

    var urls [3]string
    urls[0] = "https://beget.ru"
    urls[1] = "http://habrahabr.ru"
    urls[2] = "https://google.ru"

    jobs := make(chan string)
    answers := make(chan string)

    for w := 1; w <= workerPoolSize; w++ {
        go worker(w, jobs, answers)
    }

    for i := 1; i <= len(urls); i++ {
        jobs <- urls[i]
    }

    close(jobs)

    for a := 1; a <= len(urls); a++ {
        fmt.Println(<-answers)
    }
}

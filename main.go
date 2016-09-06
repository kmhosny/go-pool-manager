package main

import (
  "fmt"
  "strconv"
  "os/exec"
  )

func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Println("worker", id, "processing job", j)
        out, err := exec.Command("./foo.sh", strconv.Itoa(j)).Output()
        if err != nil {
          // TODO: capture this in a log file
          // log.Fatal(err)
        }
        fmt.Printf("= %s\n", out)
        results <- j * 2
    }
}

func main() {
    concurrency := 8
    total := 12

    jobs := make(chan int, 100)
    results := make(chan int, 100)

    for w := 1; w <= concurrency; w++ {
        go worker(w, jobs, results)
    }

    for j := 1; j <= total; j++ {
        jobs <- j
    }
    close(jobs)

    // Collect results of the work done
    for r := 1; r <= total; r++ {
        <- results
    }
}

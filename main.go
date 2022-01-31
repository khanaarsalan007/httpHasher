package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type result struct {
	url string
	md5 string
}
type task struct {
	url string
}

func getArgs() (*int, []string) {
	var (
		parallel = flag.Int("parallel", 10, "number of requests to be made in parallel")
	)

	flag.Parse()
	args := flag.Args()

	return parallel, args

}

func main() {
	parallel, urls := getArgs()
	process := make(chan bool)
	go getMD5(parallel, urls, process)
	<-process
}

func getMD5(parallel *int, urls []string, isDone chan bool) {
	tasks := make(chan task)
	go func() {
		for _, url := range urls {
			tasks <- task{url: url}
		}
		close(tasks)
	}()

	results := make(chan result)
	var wg sync.WaitGroup
	wg.Add(*parallel)
	go func() {
		wg.Wait()
		close(results)
	}()

	for i := 0; i < *parallel; i++ {
		go func() {
			defer wg.Done()
			for t := range tasks {
				r, err := fetch(t.url)
				if err != nil {
					log.Printf("could not fetch %v: %v", t.url, err)
					continue
				}
				res := result{url: t.url, md5: r}
				results <- res
			}
		}()
	}

	for result := range results {
		fmt.Println(result)
	}

	isDone <- true
}

func fetch(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("could not get %s: %v", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusTooManyRequests {
			return "", fmt.Errorf("you are being rate limited")
		}

		return "", fmt.Errorf("bad response from server: %v", http.StatusNotFound)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error in reading the response")

	}
	checkSum := md5.Sum(bytes)
	return fmt.Sprintf("%x", checkSum), nil
}

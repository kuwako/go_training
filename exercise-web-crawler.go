package main

import (
	"fmt"
	"sync"
	"time"
)

var fetched = make(map[string]bool)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type crawlResult struct {
	url  string
	body string
	urls []string
	err  error
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	ch := make(chan crawlResult, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Wait()
		close(ch)
	}()
	go crawlWorker(url, depth, fetcher, ch, &wg)

	for r := range ch {
		if r.err == nil {
			fmt.Printf("found: %s %q\n", r.url, r.body)
		} else {
			fmt.Printf("not found: %s\n", r.url)
		}
	}
}

func crawlWorker(url string, depth int, fetcher Fetcher, ch chan crawlResult, wg *sync.WaitGroup) {
	defer wg.Done()

	// This implementation doesn't do either:
	if depth <= 0 {
		return
	}
	if fetched[url] {
		return
	}
	fetched[url] = true
	body, urls, err := fetcher.Fetch(url)
	r := crawlResult{url, body, urls, err}
	ch <- r
	if err != nil {
		return
	}
	for _, u := range urls {
		wg.Add(1)
		go crawlWorker(u, depth-1, fetcher, ch, wg)
	}
	return
}

func main() {
	Crawl("http://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	time.Sleep(300 * time.Millisecond)
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}

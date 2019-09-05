package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type linkEntry struct {
	count int
	url   string
}

func getRequest(url string) (*http.Response, error) {
	// Disable the SSL verification
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := http.Client{Transport: transport}

	// Get response
	resp, err := client.Get(url)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func computeLinksEntries(linkEntrySlice []linkEntry) []linkEntry {
	count := make(map[string]int)
	list := []linkEntry{}

	for _, entry := range linkEntrySlice {
		count[entry.url]++
	}

	for u, c := range count {
		list = append(list, linkEntry{c, u})
	}

	return list
}

func getLinks(httpBody io.Reader, baseURL string) []linkEntry {
	links := []linkEntry{}

	pageTokenizer := html.NewTokenizer(httpBody)

	for {
		tokenType := pageTokenizer.Next()

		if tokenType == html.ErrorToken {
			err := pageTokenizer.Err()
			if err == io.EOF {
				return computeLinksEntries(links) //end of the file, return the links array
			}
			log.Fatalf("error tokenizing HTML: %v", pageTokenizer.Err())
		}

		token := pageTokenizer.Token()

		if tokenType == html.StartTagToken && token.Data == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" {

					if strings.HasPrefix(attr.Val, "/") {
						resolvedURL := fmt.Sprintf("%s%s", baseURL, attr.Val)
						links = append(links, linkEntry{1, resolvedURL})
					}

					if strings.HasPrefix(attr.Val, baseURL) {
						links = append(links, linkEntry{1, attr.Val})
					}
				}
			}
		}
	}
}

func parseStartURL(u string) string {
	parsed, _ := url.Parse(u)
	return fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
}

func crawl(url string, baseURL string, queue chan string) {
	resp, err := getRequest(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	fmt.Printf("\n%d | %s\n", resp.StatusCode, url)

	links := getLinks(resp.Body, baseURL)

	fmt.Println("Links |")
	for _, link := range links {
		fmt.Printf("      |- %d x %s\n", link.count, link.url)
		l := link.url
		go func() { queue <- l }()
	}
}

func main() {
	crawledURL := flag.String("url", "https://julien-gonzalez.fr", "URL that will be crawled")
	flag.Parse()

	queue := make(chan string, 10)
	seen := make(map[string]bool)
	baseURL := parseStartURL(*crawledURL)

	go func() {
		queue <- *crawledURL
	}()

	for url := range queue {
		if !seen[url] {
			seen[url] = true
			crawl(url, baseURL, queue)
		}
	}
}

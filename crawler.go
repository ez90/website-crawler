package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ez90/crawler/htmlParser"
)

type linkEntry struct {
	count int
	url   string
	text  string
}

func getRequest(url string) (*http.Response, error) {
	// Disable the SSL verification
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := http.Client{Transport: transport}

	// Get response
	resp, err := client.Get(url)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func computeLinksEntries(linkEntrySlice []linkEntry) []linkEntry {
	keys := make(map[string]bool)
	count := make(map[string]int)
	list := []linkEntry{}

	for _, entry := range linkEntrySlice {
		count[entry.url]++

		if _, value := keys[entry.url]; !value {
			keys[entry.url] = true
			list = append(list, entry)
		}
	}

	for i := range list {
		list[i].count = count[list[i].url]
	}

	return list
}

func isInternalURL(url string, baseURL string) bool {
	if strings.HasPrefix(url, "/") {
		return true
	}

	if strings.HasPrefix(url, baseURL) {
		return true
	}

	return false
}

func normalizeURL(url string, baseURL string) string {
	if strings.HasPrefix(url, "/") {
		resolvedURL := fmt.Sprintf("%s%s", baseURL, url)
		return resolvedURL
	}

	return url
}

func normalizeBaseURL(u string) string {
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

	doc, _ := htmlParser.NewDocumentfromReader(resp.Body)
	nodes := htmlParser.FindTag("a", doc)
	links := []linkEntry{}

	for _, node := range nodes {
		url := htmlParser.FindAttr("href", node)
		if isInternalURL(url, baseURL) {
			url = normalizeURL(url, baseURL)
			links = append(links, linkEntry{0, url, htmlParser.GetContent(node)})
		}
	}

	links = computeLinksEntries(links)

	fmt.Println("Links |")
	for _, link := range links {
		fmt.Printf("      |- %d x %s : %s\n", link.count, link.url, link.text)
		l := link.url
		go func() { queue <- l }()
	}
}

func main() {
	inputURL := flag.String("url", "http://sparkk.fr", "URL that will be crawled")
	flag.Parse()

	baseURL := normalizeBaseURL(*inputURL)
	queue := make(chan string, 10)
	seen := make(map[string]bool)

	go func() {
		queue <- baseURL
	}()

	for url := range queue {
		if !seen[url] {
			seen[url] = true
			crawl(url, baseURL, queue)
		}
	}
}

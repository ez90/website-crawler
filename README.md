
#Simple crawler

Crawl all internal links for a specified URL

Command will display all crawled links with their response status and a list of discovered internal links.
It will never crawl the same URL twice.

## Build
```
go build crawler.go
```

## Usage
```
crawler --url <my_url>
```

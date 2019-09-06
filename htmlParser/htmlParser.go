package htmlParser

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// NewDocumentfromReader return the root node
func NewDocumentfromReader(r io.Reader) (*html.Node, error) {
	doc, err := html.Parse(r)

	if err != nil {
		return nil, err
	}

	return doc, nil
}

// FindTag return array of nodes with specific tag name
func FindTag(tag string, n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == tag {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, FindTag(tag, c)...)
	}
	return ret
}

// FindAttr return value of a node attribute with specific name
// the first match is returned by default
func FindAttr(nodeAttr string, n *html.Node) string {
	var ret string
	for _, attr := range n.Attr {
		if attr.Key == nodeAttr {
			ret = attr.Val
			break
		}
	}
	return ret
}

// GetContent return the text content of a nodes
// if node content is not text (like Elements), function will retur empty string
func GetContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += GetContent(c)
	}
	return strings.Join(strings.Fields(ret), " ")
}

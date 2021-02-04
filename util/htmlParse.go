package util

import (
	"strings"

	"golang.org/x/net/html"
)

var (
	REGIONS = []string{"Americas", "EMEA", "APAC"}
)

type Matcher func(x string) bool

//TODO: pass in a function that matches a node with a string
//func StringMatcher

func GetBody(doc *html.Node) *html.Node {
	var body *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "body" {
			body = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if body != nil {
		return body
	}
	return nil
}

func GetDivByClass(doc *html.Node, class string) *html.Node {
	var target *html.Node
	var find func(*html.Node)

	find = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "div" {
			// test match
			for _, attr := range node.Attr {
				if attr.Key == "class" {
					if attr.Val == class {
						target = node
						return
					}
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			find(c)
		}
	}
	find(doc)
	if target != nil {
		return target
	}
	return nil
}

func GetNodeFromText(node *html.Node, text string) *html.Node {
	var target *html.Node
	var find func(*html.Node)

	find = func(node *html.Node) {
		if node.Type == html.TextNode {
			// test match
			t := RemoveSpecialChars(node.Data)
			t = strings.Trim(t, " ")
			if text == t {
				target = node
				return
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			find(c)
		}
	}
	find(node)
	if target != nil {
		return target
	}
	return nil
}

func GetChildrenOfType(node *html.Node, nodeType string) []html.Node {
	return nil
}
func GetSiblingsOfType(start *html.Node, nodeType string) []*html.Node {
	var targets []*html.Node
	var find func(node *html.Node)

	find = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == nodeType {
			targets = append(targets, node)
		}
		if node.NextSibling == nil {
			return
		}
		find(node.NextSibling)
	}
	find(start.NextSibling)
	if len(targets) == 0 {
		return nil
	}
	return targets
}
func GetParentOfType(start *html.Node, nodeType string) *html.Node {
	var target *html.Node
	var find func(node *html.Node)

	find = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == nodeType {
			target = node
			return
		}
		if node.Parent == nil {
			return
		}
		find(node.Parent)
	}
	find(start.Parent)

	if target != nil {
		return target
	}
	return nil
}

// todo: parse non qwerty chars from a string
func RemoveSpecialChars(candidate string) string {
	var r []byte
	for _, c := range []byte(candidate) {
		// 126 is ~, pretty high up in UTF-8
		if c < 126 {
			r = append(r, c)
		}
	}
	return string(r[:])
}

func ParseFundingFromNode(node *html.Node) {

}

// TODO: determine why this SOs
func GetBodyNR(doc *html.Node) *html.Node {
	if doc.Type == html.ElementNode && doc.Data == "body" {
		return doc
	}
	for child := doc.FirstChild; child != nil; child = child.NextSibling {
		return GetBodyNR(doc)
	}
	return nil
}

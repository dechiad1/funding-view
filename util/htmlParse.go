package util

import (
	"strings"

	"golang.org/x/net/html"
)

var dfs func(*html.Node)

// pass in pointer to slice, to fill with nodes during this DFS operation
func GetDFSNodes(doc *html.Node, r *[]*html.Node) {
	dfs = func(node *html.Node) {
		*r = append(*r, node)
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			dfs(child)
		}
	}
	dfs(doc)
}

func HtmlToString(node *html.Node) string {
	text := ""
	dfs = func(n *html.Node) {
		if n.Type == html.TextNode {
			text = text + n.Data
		} else {
			text = text + "<" + n.Data + ">"
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			dfs(child)
		}
		if n.Type == html.ElementNode {
			text = text + "</" + n.Data + ">"
		}
	}
	dfs(node)
	return text
}

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

func contains(t string, s ...string) bool {
	exists := make(map[string]bool)
	for _, ss := range s {
		exists[ss] = true
	}

	if exists[t] {
		return true
	}
	return false

}

func GetNodesOfType(node *html.Node, nodeTypes ...string) []*html.Node {
	targets := make([]*html.Node, 0)
	var find func(*html.Node)

	find = func(t *html.Node) {
		if t.Type == html.ElementNode && contains(t.Data, nodeTypes...) {
			targets = append(targets, t)
		}
		for c := t.FirstChild; c != nil; c = c.NextSibling {
			find(c)
		}
	}
	find(node)
	return targets
}

func NodeToStringSlice(node *html.Node) []string {
	var result []string
	var find func(*html.Node)

	find = func(t *html.Node) {
		if t.Type == html.TextNode {
			data := RemoveSpecialChars(t.Data)
			result = append(result, data)
		}

		for c := t.FirstChild; c != nil; c = c.NextSibling {
			find(c)
		}
	}
	find(node)
	return result
}

func GetTextFromNode(node *html.Node) string {
	var result string
	var find func(*html.Node)

	find = func(t *html.Node) {
		if t.Type == html.TextNode {
			result = result + t.Data
		}

		for c := t.FirstChild; c != nil; c = c.NextSibling {
			find(c)
		}
	}
	find(node)
	return result
}

type matchMany func(*html.Node, map[string]bool)

func matchNodeTypes(node *html.Node, types map[string]bool) bool {
	if node.Type == html.ElementNode && types[node.Data] {
		return true
	}
	return false
}

func GetSiblingsOfType(start *html.Node, types ...string) []*html.Node {
	var targets []*html.Node
	var find func(node *html.Node)

	typeMap := make(map[string]bool)
	for _, t := range types {
		typeMap[t] = true
	}

	find = func(node *html.Node) {
		if matchNodeTypes(node, typeMap) {
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

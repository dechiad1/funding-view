package datasource

import (
	"fmt"
	"time"

	"golang.org/x/net/html"
)

type FlCompanyDTO struct {
	Name          string
	description   string
	funding       string
	investors     string
	location      string
	industry      string
	employeeCount string
	date          time.Time
}

// any full company (starting with position 0)
// 0 - company name: p, FirstChild: strong, FirstChild: a, FirstChild: text
// 0 - description: p, SecondChild: text
// 1 - new funding: p, FirstChild: text
// 2 - round investors: p, FirstChild: text
// 3 - press: p, list of a's for each press link :/
// 4 - HQ location: p, FirstChild: text
// 5 - industry: p, FirstChild: text
// 6 - employees: p, FirstChild: text

func ParseCompany(company []*html.Node) *FlCompanyDTO {
	if len(company) != 7 {
		fmt.Printf("invalid slice received, expecting 7 got %d\n", len(company))
		return nil
	}

	name := parseName(company[0])
	description := parseDescription(company[0])
	funding := getTextFromFirstChild(company[1], "funding")
	investors := getTextFromFirstChild(company[2], "investors")
	location := getTextFromFirstChild(company[4], "location")
	industry := getTextFromFirstChild(company[5], "industry")
	employees := getTextFromFirstChild(company[6], "employees")

	return &FlCompanyDTO{
		Name:          name,
		description:   description,
		funding:       funding,
		investors:     investors,
		location:      location,
		industry:      industry,
		employeeCount: employees,
	}
}

func parseName(node *html.Node) string {
	var text string
	var find func(n *html.Node)

	find = func(n *html.Node) {
		if n.Type == html.TextNode {
			text = n.Data
			return
		}
		if n.FirstChild != nil {
			find(n.FirstChild)
		} else {
			fmt.Println("Warn: can not find text within name p tag")
		}
		return
	}
	find(node)
	return text
}

func parseDescription(node *html.Node) string {
	if node.FirstChild != nil {
		fc := node.FirstChild
		if fc.NextSibling != nil {
			sib := fc.NextSibling
			if sib.Type == html.TextNode {
				return sib.Data
			}
			fmt.Println("Warn: second child is not a text node")
			return ""
		}
		fmt.Println("Warn: description node has no second child")
		return ""
	}
	fmt.Println("Warn: description node has no first child")
	return ""
}

func getTextFromFirstChild(node *html.Node, field string) string {
	if node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
		return node.FirstChild.Data
	}
	fmt.Printf("Warn: node %v for field %s has no first child\n", node, field)
	return ""
}

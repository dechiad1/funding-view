package parser

import (
	"company-funding/dto"
	"company-funding/util"
	"strings"

	"golang.org/x/net/html"
)

type FlParser struct {
	doc string
}

func (p *FlParser) Parse(node *html.Node) {
	var find func(*html.Node)
	find = func(node *html.Node) {
		for fc := node.FirstChild; fc != nil; fc = fc.NextSibling {

		}
	}
	find(node)
}

func (p *FlParser) IsSeparator(node *html.Node) bool {
	if node.Type == html.TextNode && node.Data == "—" {
		return true
	} else if node.Type == html.ElementNode && node.Data == "hr" {
		return true
	}
	return false
}

func IsAttribute(node *html.Node, company *dto.FlCompanyDTO) {
	// All except company name & description
	if node.Type == html.TextNode {
		data := node.Data
		dataType := strings.Split(data, ":")

		if len(dataType) == 1 {
			return
		}

		switch t := dataType[0]; t {

		case "New Funding Raised":
			company.Funding = strings.TrimSpace(dataType[1])
		case "Round Investors":
			company.Investors = strings.TrimSpace(dataType[1])
		case "Press":
		case "HQ":
			company.Location = strings.TrimSpace(dataType[1])
		case "Industry":
			company.Industry = strings.TrimSpace(dataType[1])
		case "Employee Count":
			company.EmployeeCount = strings.TrimSpace(dataType[1])
		default:
			return
		}
	} else {
		// get the nodes children as a string & test if its likely to be name & description
		text := util.HtmlToString(node)
		if strings.Count(text, "<p>") != 1 {
			return
		}
		if strings.Count(text, "<a>") > 4 {
			return
		}
		if !strings.Contains(text, "LinkedIn") {
			return
		}
		// at this point, its probably the right element
		content := util.NodeToStringSlice(node)
		nameSet := false
		description := ""
		for _, c := range content {
			if c == "" {
				continue
			} else {
				if !nameSet {
					company.Name = c
					nameSet = true
				} else {
					description = description + c
				}
			}
		}
		company.Description = description
	}

}

package parser

import (
	"company-funding/dto"
	"company-funding/repository"
	"company-funding/util"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// any full company (starting with position 0)
// 0 - company name: p, FirstChild: strong, FirstChild: a, FirstChild: text
// 0 - description: p, SecondChild: text
// 1 - new funding: p, FirstChild: text
// 2 - round investors: p, FirstChild: text
// 3 - press: p, list of a's for each press link :/
// 4 - HQ location: p, FirstChild: text
// 5 - industry: p, FirstChild: text
// 6 - employees: p, FirstChild: text

func Parse(start *html.Node) {
	db := repository.Connect()
	node := util.GetNodeFromText(start, "Americas")
	if node == nil {
		fmt.Println("no body!")
		os.Exit(1)
	}

	// get notable html.Node near data
	blockquote := util.GetParentOfType(node, "blockquote")

	// get nodes that contain relevant information
	sibs := util.GetSiblingsOfType(blockquote, "p", "div")

	// loop through sibs & create a company in between each div
	var companies []*dto.FlCompanyDTO
	var companyAttributes []*html.Node
	for _, sib := range sibs {
		if sib.Data == "p" {
			companyAttributes = append(companyAttributes, sib)
		} else if sib.Data == "div" {
			if len(companyAttributes) > 1 {
				c := ParseCompany(companyAttributes)
				companies = append(companies, c)
			}
			companyAttributes = nil
		}
	}

	for _, c := range companies {
		company := repository.ConvertCompanyDto(c)
		db.Create(company)
		fmt.Printf("Saving funding for %s\n", c.Name)
	}
}

func ParseCompany(company []*html.Node) *dto.FlCompanyDTO {
	if len(company) < 2 {
		fmt.Printf("invalid slice received, got %d attributes\n", len(company))
		return nil
	}

	name := parseName(company[0])
	description := parseDescription(company[0])

	flCompanyDTO := &dto.FlCompanyDTO{
		Name:        name,
		Description: description,
	}

	for i := 1; i < len(company); i++ {
		getTextFromFirstChild(company[i], flCompanyDTO)
	}
	return flCompanyDTO
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

func getTextFromFirstChild(node *html.Node, company *dto.FlCompanyDTO) {
	if node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
		data := node.FirstChild.Data
		dataType := strings.Split(data, ":")

		switch t := dataType[0]; t {

		case "New Funding Raised":
			company.Funding = dataType[1]
		case "Round Investors":
			company.Investors = dataType[1]
		case "Press":
		case "HQ":
			company.Location = dataType[1]
		case "Industry":
			company.Industry = dataType[1]
		case "Employee Count":
			company.EmployeeCount = dataType[1]
		default:
			fmt.Printf("unrecognized company field: %s\n", t)
		}
	} else {
		fmt.Printf("Warn: node %v has no first text child\n", node)
	}
}

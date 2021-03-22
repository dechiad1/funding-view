package parser

import (
	"company-funding/dto"
	"company-funding/repository"
	"company-funding/util"
	"fmt"
	"strings"

	"golang.org/x/net/html"
	"gorm.io/gorm"
)

type FlParser struct {
	Dev        bool
	Db         *gorm.DB
	CurrentDoc string
}

func (p *FlParser) Parse(node *html.Node) {
	companies := make([]*dto.FlCompanyDTO, 0)
	p.parse(node, &companies)

	if p.Dev {
		fmt.Println(len(companies))
		for _, c := range companies {
			c.PrintDto()
			fmt.Println()
		}
	} else {
		for i, c := range companies {
			dao := repository.CompanyFunding{}
			c.ConvertFlCompanyDto(&dao)
			if c.Name == "" {
				fmt.Printf("Company number %d from %s is null", i, p.CurrentDoc)
			}
			c.PrintDto()
			p.Db.Create(&dao)
		}
	}
}

func (p *FlParser) parse(node *html.Node, companies *[]*dto.FlCompanyDTO) {
	dfsnodes := make([]*html.Node, 0)
	util.GetDFSNodes(node, &dfsnodes)
	company := &dto.FlCompanyDTO{}
	for _, node := range dfsnodes {
		if p.IsSeparator(node) {
			if company.IsSaveAble() {
				*companies = append(*companies, company)
			}
			company = &dto.FlCompanyDTO{}
		} else {
			IsAttribute(node, company)
		}
	}

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
		if node.Type == html.ElementNode && node.Data == "a" {
			// First link is always the company name
			if company.IsEmpty() {
				if node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
					company.Name = node.FirstChild.Data
				} else {
					fmt.Printf("First child is not the company")
				}
			}
		}
	}
}

func contains(text string, words ...string) bool {
	contains := false
	for _, w := range words {
		if strings.Contains(text, w) {
			contains = true
		}
	}
	return contains
}

// Depreciate approach collecting name & description?
func getNameAndDescription(node *html.Node, company *dto.FlCompanyDTO) {
	// get the nodes children as a string & test if its likely to be name & description
	text := util.HtmlToString(node)
	if strings.Contains(text, "<br>") {
		// TODO: add printline here for the date with this blocking item
		return
	}
	if strings.Count(text, "<p>") != 1 {
		return
	}
	if strings.Count(text, "<a>") > 4 {
		return
	}
	if !contains(text, "LinkedIn", "Linkedin", "linkedin") {
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

package parser

import (
	"company-funding/dto"
	"company-funding/repository"
	"company-funding/util"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func StarterNode(doc *html.Node, searchStrings ...string) *html.Node {
	var node *html.Node
	count := 0
	for node == nil && count < len(searchStrings) {
		node = util.GetNodeFromText(doc, searchStrings[count])
		count = count + 1
	}

	if node == nil {
		fmt.Println("Could not find node with contents listed in starterNodes")
		return nil
	}

	blockquote := util.GetParentOfType(node, "blockquote")
	return blockquote
}

func Parse(start *html.Node) int {
	db := repository.Connect()
	blockquote := StarterNode(start, "Americas", "Forwarded this Letter?")
	if blockquote == nil {
		return -1
	}
	// get nodes that contain relevant information
	sibs := util.GetSiblingsOfType(blockquote, "p", "div")

	// loop through sibs & create a company in between each div
	var companies []*dto.FlCompanyDTO
	var companyAttributes []*html.Node
	for _, sib := range sibs {
		valid, _ := validateAttribute(sib)
		if valid {
			companyAttributes = append(companyAttributes, sib)
		} else if isSeparator(sib) {
			c := ParseCompany(companyAttributes)
			if c != nil {
				companies = append(companies, c)
			}
			companyAttributes = make([]*html.Node, 0)
		}
	}

	for _, c := range companies {
		c.PrintDto()
		if c.IsSaveAble() {
			company := repository.ConvertCompanyDto(c)
			db.Create(company)
			fmt.Printf("Saving funding for %s\n", c.Name)
		}
	}
	return 0
}

func isSeparator(node *html.Node) bool {
	if node.Type == html.ElementNode && node.Data == "p" {
		if node.FirstChild == node.LastChild && node.FirstChild.Type == html.TextNode && node.FirstChild.Data == "—" {
			return true
		}
	} else if node.Type == html.ElementNode && node.Data == "div" {
		return true
	}
	return false
}

func ParseCompany(company []*html.Node) *dto.FlCompanyDTO {
	if len(company) < 2 {
		fmt.Printf("invalid slice received, got %d attributes\n", len(company))
		return nil
	}

	flCompanyDTO := &dto.FlCompanyDTO{}

	for i := 0; i < len(company); i++ {
		_, t := validateAttribute(company[i])
		if t == 1 {
			getTextFromFirstChild(company[i], flCompanyDTO)
		} else if t == 2 {
			parseName(company[i], flCompanyDTO)
		}
	}
	return flCompanyDTO
}

func parseName(node *html.Node, company *dto.FlCompanyDTO) {
	name := ""
	var targetParent *html.Node
	var find func(n *html.Node)
	var findLink func(n *html.Node)
	// find the first a || strong node. starting at elements first child & moving right
	findLink = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "strong" || n.Data == "a") {
			targetParent = n
			return
		}
		if n.NextSibling != nil {
			findLink(n.NextSibling)
		}
	}
	findLink(node.FirstChild)
	// get the text from it
	find = func(n *html.Node) {
		if n == nil {
			return
		}
		if n.Type == html.TextNode {
			name = n.Data
			return
		}
		if n.FirstChild != nil {
			find(n.FirstChild)
		} else {
			fmt.Println("Warn: can not find text within name p tag")
		}
		return
	}
	find(targetParent)

	if name != "" {
		company.Name = strings.TrimSpace(name)
		// description
		if targetParent.NextSibling != nil {
			sib := targetParent.NextSibling
			if sib.Type == html.TextNode {
				company.Description = strings.TrimSpace(sib.Data)
			}
		}
	}
}

func getTextFromFirstChild(node *html.Node, company *dto.FlCompanyDTO) {
	if node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
		data := node.FirstChild.Data
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
			fmt.Printf("unrecognized company field: %s\n", t)
		}
	} else {
		fmt.Printf("Warn: node %v has no first text child\n", node)
	}
}

var attributeIndicators = [6]string{
	"New Funding Raised",
	"Round Investors",
	"Press",
	"HQ",
	"Industry",
	"Employee Count",
}

// TODO: return attribute type: 1 for name description, 2 otherwise. can manually normalize companies without 1s
func validateAttribute(node *html.Node) (bool, int) {
	fc := node.FirstChild
	// skip emojis as first text in a 'p'
	if fc != nil {

		if fc.Type == html.TextNode {
			for _, s := range attributeIndicators {
				if strings.Contains(fc.Data, s) {
					return true, 1
				}
			}
		}

		// For the name, description element
		linkNodes := util.GetNodesOfType(node, "a", "strong")
		if len(linkNodes) > 0 {
			return true, 2
		}
	}
	return false, 0
}

func DebugParse(doc *html.Node) {
	blockquote := StarterNode(doc, "Americas", "Forwarded this Letter?")
	sibs := util.GetSiblingsOfType(blockquote, "p", "div")

	for _, s := range sibs {
		text := util.GetTextFromNode(s)
		fmt.Println(text)
	}

}

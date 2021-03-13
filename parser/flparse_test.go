package parser

import (
	"company-funding/dto"
	"company-funding/util"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

var testValidationNodes = `<!doctype html>
<html>
	<head/>
	<body>
		<p>-</p>
		<p>here is a really long string that doesnt actually contain any meaninful scontent therefore we dwa to to make it return false</p>
		<p><strong>bold</strong>does some bold stuff<a href="blah">blah</a></p>
		<p><strong><a href="http://www.enevate.com">Enevate</a></strong> develops Li-ion battery charging solutions using a pure silicon-dominant battery technology, enabling electric vehicle owners to charge their cars as fast as refueling a gas car. <strong>View on: <a href="https://www.linkedin.com/company/enevate-corporation">LinkedIn.com</a> | <a href="https://www.linkedin.com/sales/company/2024574">Sales Navigator</a>.</strong></p>
		<p>HQ: nj</p>
		<p><strong><a href="https://ro.co/">Ro</a></strong> is a mission-driven healthcare technology company where doctors, pharmacists and engineers are working together to reinvent the way the healthcare system works. <a href="https://www.linkedin.com/company/getroman/">View on LinkedIn</a>.</p>
		<p>🦄<a href="https://www.outreach.io/">Outreach</a> is the leading Sales Engagement Platform, accelerates revenue growth by optimizing every interaction throughout the customer lifecycle. <a href="https://www.linkedin.com/company/outreach-saas/">View on LinkedIn</a>. 🦄</p>
	</body>
</html>
`

func Test_validateAttribute(t *testing.T) {
	doc, _ := html.Parse(strings.NewReader(testValidationNodes))
	ps := util.GetNodesOfType(doc, "p")

	cases := []struct {
		node *html.Node
		want bool
	}{
		{ps[0], false},
		{ps[1], false},
		{ps[2], true},
		{ps[3], true},
		{ps[4], true},
		{ps[5], true},
		{ps[6], true},
	}

	for i, c := range cases {
		n := c.node

		valid, _ := validateAttribute(n)
		if valid != c.want {
			t.Errorf("attribute not validated, got %v but wanted %v, for input node at index %d and data of %s\n", valid, c.want, i, n.Data)
		}
	}
}

var testCompanyNodes = `<!doctype html>
<html>
	<head/>
	<body>
		<p><strong><a href="https://ro.co/">Ro</a></strong> is a mission-driven healthcare technology company where doctors, pharmacists and engineers are working together to reinvent the way the healthcare system works. <a href="https://www.linkedin.com/company/getroman/">View on LinkedIn</a>.</p>
		<p>New Funding Raised: $85M, Series B</p>
		<p>HQ: New York, New York</p>
		<p>Industry: Hospital &amp; Health Care</p>
		<p>Employee Count: 632</p>
		
		<p><strong><a href="http://www.enevate.com">Enevate</a></strong> develops Li-ion battery charging solutions using a pure silicon-dominant battery technology, enabling electric vehicle owners to charge their cars as fast as refueling a gas car. <strong>View on: <a href="https://www.linkedin.com/company/enevate-corporation">LinkedIn.com</a> | <a href="https://www.linkedin.com/sales/company/2024574">Sales Navigator</a>.</strong></p>
		<p>New Funding Raised: $81M, Series E</p>
		<p>Round Investors: Fidelity Management and Research Company (lead), Infinite Potential Technologies, Mission Ventures</p>
		<p>Press: <a href="https://www.businesswire.com/news/home/20210210005278/en/Fidelity-Leads-81M-Investment-in-Enevate-to-Accelerate-Commercialization-of-Fast-Charging-Electric-Vehicle-Battery-Technology">Business Wire</a>, <a href="https://www.aftermarketnews.com/fidelity-leads-81m-investment-in-enevate/">Aftermarket News</a>, <a href="https://vcnewsdaily.com/enevate/venture-capital-funding/jygxgstbhj">VC News Daily</a></p>
		<p>HQ: Irvine, CA</p>
		<p>Industry: Renewables &amp; Environment</p>
		<p>Employee Count: 54</p>

		<p>🦄<a href="https://www.outreach.io/">Outreach</a> is the leading Sales Engagement Platform, accelerates revenue growth by optimizing every interaction throughout the customer lifecycle. <a href="https://www.linkedin.com/company/outreach-saas/">View on LinkedIn</a>. 🦄</p>
		<p>New Funding Raised: $114M, Series E</p>
		<p>Round Investors: Lone Pine Capital (lead), DFJ Growth, Four Rivers Group, Lemonade Capital, M12, Mayfield Fund, Meritech Capital Partners, Sapphire Ventures, Spark Capital, Trinity Ventures</p>
		<p>HQ: Seattle, WA</p>
		<p>Industry:&nbsp; Internet</p>
		<p>Employee Count: 578</p>
	</body>
</html>
`

func Test_ParseCompany(t *testing.T) {
	doc, _ := html.Parse(strings.NewReader(testCompanyNodes))
	ps := util.GetNodesOfType(doc, "p")

	expect := []dto.FlCompanyDTO{
		{
			Name:          "Ro",
			Description:   "is a mission-driven healthcare technology company where doctors, pharmacists and engineers are working together to reinvent the way the healthcare system works.",
			Location:      "New York, New York",
			Funding:       "$85M, Series B",
			Industry:      "Hospital & Health Care",
			EmployeeCount: "632",
		},
		{
			Name:          "Enevate",
			Description:   "develops Li-ion battery charging solutions using a pure silicon-dominant battery technology, enabling electric vehicle owners to charge their cars as fast as refueling a gas car.",
			Funding:       "$81M, Series E",
			Investors:     "Fidelity Management and Research Company (lead), Infinite Potential Technologies, Mission Ventures",
			Industry:      "Renewables & Environment",
			Location:      "Irvine, CA",
			EmployeeCount: "54",
		},
		{
			Name:          "Outreach",
			Description:   "is the leading Sales Engagement Platform, accelerates revenue growth by optimizing every interaction throughout the customer lifecycle.",
			Location:      "Seattle, WA",
			Funding:       "$114M, Series E",
			Investors:     "Lone Pine Capital (lead), DFJ Growth, Four Rivers Group, Lemonade Capital, M12, Mayfield Fund, Meritech Capital Partners, Sapphire Ventures, Spark Capital, Trinity Ventures",
			Industry:      "Internet",
			EmployeeCount: "578",
		},
	}

	cases := []struct {
		company []*html.Node
		want    dto.FlCompanyDTO
	}{
		{ps[0:5], expect[0]},
		{ps[5:12], expect[1]},
		{ps[12:], expect[2]},
	}

	for _, c := range cases {
		wantedCompany := c.want
		valid := ParseCompany(c.company)
		if valid.Name != wantedCompany.Name {
			t.Errorf("name not validated, got %s but wanted %s\n", valid.Name, wantedCompany.Name)
		}
		if valid.Description != wantedCompany.Description {
			t.Errorf("description not validated, got %s but wanted %s\n", valid.Description, wantedCompany.Description)
		}
		if valid.Funding != wantedCompany.Funding {
			t.Errorf("funding not validated, got %s but wanted %s\n", valid.Funding, wantedCompany.Funding)
		}
		if valid.Investors != wantedCompany.Investors {
			t.Errorf("Investors not validated, got %s but wanted %s\n", valid.Investors, wantedCompany.Investors)
		}
		if valid.Industry != wantedCompany.Industry {
			t.Errorf("Industry not validated, got %s but wanted %s\n", valid.Industry, wantedCompany.Industry)
		}
		if valid.Location != wantedCompany.Location {
			t.Errorf("Location not validated, got %s but wanted %s\n", valid.Location, wantedCompany.Location)
		}
		if valid.EmployeeCount != wantedCompany.EmployeeCount {
			t.Errorf("Employee Count not validated, got %s but wanted %s\n", valid.EmployeeCount, wantedCompany.EmployeeCount)
		}
	}

}

var testValidSeparators = `<!doctype html>
<html>
	<head/>
	<body>
		<p>-</p>
		<p>—</p>
		<p>here is a really long string that doesnt actually contain any meaninful scontent therefore we dwa to to make it return false</p>
		<div>hey there</div>
		<p><strong><a href="http://www.enevate.com">Enevate</a></strong> develops Li-ion battery charging solutions using a pure silicon-dominant battery technology, enabling electric vehicle owners to charge their cars as fast as refueling a gas car. <strong>View on: <a href="https://www.linkedin.com/company/enevate-corporation">LinkedIn.com</a> | <a href="https://www.linkedin.com/sales/company/2024574">Sales Navigator</a>.</strong></p>
		<div/>
	</body>
</html>
`

func Test_isSeparator(t *testing.T) {
	doc, _ := html.Parse(strings.NewReader(testValidSeparators))
	elems := util.GetNodesOfType(doc, "p", "div")

	cases := []struct {
		node *html.Node
		want bool
	}{
		{elems[0], false},
		{elems[1], true},
		{elems[2], false},
		{elems[3], true},
		{elems[4], false},
		{elems[5], true},
	}

	for i, c := range cases {
		want := c.want
		node := c.node
		got := isSeparator(node)
		if want != got {
			t.Errorf("wanted %v but got %v for element %d of test case\n", want, got, i)
		}
	}
}

var testStarterNode = `<!doctype html>
<html>
	<head/>
	<body>
		<div>
			<blockquote>Weird String</blockquote>
			<div>hey there</div>
			<p>P!</p>
			<li>
				<blockquote>anothaone</blockquote>
				<div>HOWABOUTTHIS</div>
			</li>
		</div>
	</body>
</html>
`

func Test_StarterNode(t *testing.T) {
	doc, _ := html.Parse(strings.NewReader(testStarterNode))

	cases := []struct {
		searchString   string
		wantedNodeType string
	}{
		{"Weird String", "blockquote"},
		{"anothaone", "blockquote"},
	}

	for _, c := range cases {
		search := c.searchString
		wanted := c.wantedNodeType
		got := StarterNode(doc, search)
		if wanted != got.Data {
			t.Errorf("wanted %s but got %s for search string %s\n", wanted, got.Data, c.searchString)
		}
	}
}

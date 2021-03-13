package parser

import (
	"company-funding/dto"
	"company-funding/util"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

var doc = `<!doctype html>
<html>
	<head/>
	<body>
		<p>-</p>
		<p>—</p>
		<p>here is a really long string that doesnt actually contain any meaninful scontent therefore we dwa to to make it return false</p>
		<div><hr></div>
		<p><strong><a href="http://www.enevate.com">Enevate</a></strong> develops Li-ion battery charging solutions using a pure silicon-dominant battery technology, enabling electric vehicle owners to charge their cars as fast as refueling a gas car. <strong>View on: <a href="https://www.linkedin.com/company/enevate-corporation">LinkedIn.com</a> | <a href="https://www.linkedin.com/sales/company/2024574">Sales Navigator</a>.</strong></p>
		<p>New Funding Raised: $81M, Series E</p>
		<p>Round Investors: Fidelity Management and Research Company (lead), Infinite Potential Technologies, Mission Ventures</p>
		<p>Press: <a href="https://www.businesswire.com/news/home/20210210005278/en/Fidelity-Leads-81M-Investment-in-Enevate-to-Accelerate-Commercialization-of-Fast-Charging-Electric-Vehicle-Battery-Technology">Business Wire</a>, <a href="https://www.aftermarketnews.com/fidelity-leads-81m-investment-in-enevate/">Aftermarket News</a>, <a href="https://vcnewsdaily.com/enevate/venture-capital-funding/jygxgstbhj">VC News Daily</a></p>
		<p>HQ: Irvine, CA</p>
		<p>Industry: Renewables &amp; Environment</p>
		<p>Employee Count: 54</p>
		<div><hr></div>
		<p>🦄<a href="https://www.outreach.io/">Outreach</a> is the leading Sales Engagement Platform, accelerates revenue growth by optimizing every interaction throughout the customer lifecycle. <a href="https://www.linkedin.com/company/outreach-saas/">View on LinkedIn</a>. 🦄</p>
		<p>New Funding Raised: $114M, Series E</p>
		<p>Round Investors: Lone Pine Capital (lead), DFJ Growth, Four Rivers Group, Lemonade Capital, M12, Mayfield Fund, Meritech Capital Partners, Sapphire Ventures, Spark Capital, Trinity Ventures</p>
		<p>HQ: Seattle, WA</p>
		<p>Industry:&nbsp; Internet</p>
		<p>Employee Count: 578</p>
		<p>—</p>
	</body>
</html>
`

func Test_isElement_name_description(t *testing.T) {
	doc, _ := html.Parse(strings.NewReader(doc))
	pNodes := util.GetNodesOfType(doc, "p")

	if len(pNodes) != 17 {
		t.Errorf("len pNodes is not 17. Got %d", len(pNodes))
	}

	cases := []struct {
		Node      *html.Node
		Valid     bool
		Want      string
		Attribute int
	}{
		{pNodes[0], false, "", -1},
		{pNodes[3], true, "Enevate", 0},
		{pNodes[10], true, "Outreach", 0},
	}

	for i, c := range cases {
		nc := dto.FlCompanyDTO{}
		IsAttribute(c.Node, &nc)

		if c.Valid {
			attrValue := reflect.ValueOf(&nc).Elem().Field(c.Attribute).Interface()
			v := attrValue.(string)
			if v != c.Want {
				t.Errorf("Got %s, but wanted %s from the %d test case\n", v, c.Want, i)
			}
		} else {
			if nc.IsSaveAble() {
				t.Errorf("Got saveable from the %d test case", i)
			}
		}
	}
}

func Test_IsSeparator(t *testing.T) {
	doc, _ := html.Parse(strings.NewReader(doc))
	p := &FlParser{}
	elems := make([]*html.Node, 0)
	util.GetDFSNodes(doc, &elems)

	// debug: 89 elems from the parser doc above
	// for i, e := range elems {
	// 	fmt.Printf("%d: data of elem: %s\n", i, e.Data)
	// }
	// fmt.Println(len(elems))

	cases := []struct {
		node *html.Node
		want bool
	}{
		{elems[0], false},
		{elems[1], false},
		{elems[2], false},
		{elems[3], false},
		{elems[4], false},
		{elems[5], false},
		{elems[6], false},
		{elems[7], false},
		{elems[8], false},
		{elems[9], false},
		{elems[10], false},
		{elems[11], true}, //—
		{elems[12], false},
		{elems[13], false},
		{elems[14], false},
		{elems[15], false},
		{elems[16], false},
		{elems[17], true}, //hr
		{elems[18], false},
		{elems[19], false},
		{elems[30], false},
		{elems[31], false},
		{elems[32], false},

		{elems[59], false},
		{elems[60], true}, //hr

		{elems[86], false},
		{elems[87], true}, //-
		{elems[88], false},
	}

	for i, c := range cases {
		//fmt.Printf("data of elem: %s\n", c.node.Data)
		want := c.want
		node := c.node
		got := p.IsSeparator(node)
		if want != got {
			t.Errorf("wanted %v but got %v for element %d of test case\n", want, got, i)
		}
	}
}

package util

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

var doc = `<p><strong><a href="http://www.enevate.com">Enevate</a></strong> develops Li-ion battery charging solutions using a pure silicon-dominant battery technology, enabling electric vehicle owners to charge their cars as fast as refueling a gas car. <strong>View on: <a href="https://www.linkedin.com/company/enevate-corporation">LinkedIn.com</a> | <a href="https://www.linkedin.com/sales/company/2024574">Sales Navigator</a>.</strong></p>`

func Test_HtmlToString(t *testing.T) {
	node, _ := html.Parse(strings.NewReader(doc))

	//get the p node
	p := GetNodesOfType(node, "p")

	cases := []struct {
		node *html.Node
		want string
	}{
		{p[0], "<p><strong><a>Enevate</a></strong> develops Li-ion battery charging solutions using a pure silicon-dominant battery technology, enabling electric vehicle owners to charge their cars as fast as refueling a gas car. <strong>View on: <a>LinkedIn.com</a> | <a>Sales Navigator</a>.</strong></p>"},
	}

	for _, c := range cases {
		want := c.want
		node := c.node
		got := HtmlToString(node)
		if want != got {
			t.Errorf("wanted: %s \n but got: %s \n", want, got)
		}
	}
}

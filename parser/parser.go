package parser

import "golang.org/x/net/html"

type TableParser interface {
	IsSeparator() bool //TODO: pass in separator(s) struct
	Parse(*html.Node)
}

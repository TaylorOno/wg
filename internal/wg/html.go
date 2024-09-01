package wg

import (
	"golang.org/x/net/html"
)

type Action func(node *html.Node)

func ApplyToAll(n *html.Node, args []string, action Action) {
	var f func(*html.Node, []string, bool)
	f = func(n *html.Node, args []string, uni bool) {
		if uni == true {
			if n.Type == html.ElementNode && matchElementName(n, args[0]) {
				if len(args) > 1 && len(args) < 4 {
					for i := 0; i < len(n.Attr); i++ {
						attr := n.Attr[i]
						searchAttrName := args[1]
						searchAttrVal := args[2]
						if attributeAndValueEquals(attr, searchAttrName, searchAttrVal) {
							action(n)
						}
					}
				} else if len(args) == 1 {
					action(n)
				}
			}
		}
		uni = true
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, args, true)
		}
	}
	f(n, args, false)
}

// Using depth first search to find the first occurrence and return
func findOnce(n *html.Node, args []string, uni bool) (*html.Node, bool) {
	if uni == true {
		if n.Type == html.ElementNode && matchElementName(n, args[0]) {
			if len(args) > 1 && len(args) < 4 {
				for i := 0; i < len(n.Attr); i++ {
					attr := n.Attr[i]
					searchAttrName := args[1]
					searchAttrVal := args[2]
					if attributeAndValueEquals(attr, searchAttrName, searchAttrVal) {
						return n, true
					}
				}
			} else if len(args) == 1 {
				return n, true
			}
		}
	}
	uni = true
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p, q := findOnce(c, args, true)
		if q != false {
			return p, q
		}
	}
	return nil, false
}

func matchElementName(n *html.Node, name string) bool {
	return name == "" || name == n.Data
}

// attributeAndValueEquals reports when the html.Attribute attr has the same attribute name and value as from
// provided arguments
func attributeAndValueEquals(attr html.Attribute, attribute, value string) bool {
	return attr.Key == attribute && attr.Val == value
}

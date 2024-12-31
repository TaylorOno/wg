package drinks

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

// AddDrinkExtensions injects the drink rating javascript into the Williams and Graham page
func AddDrinkExtensions() (*html.Node, error) {
	menu, err := http.DefaultClient.Get("https://williamsandgraham.com/drinks/")
	if err != nil {
		return nil, fmt.Errorf("failed to read menu: %v", err)
	}
	defer menu.Body.Close()

	doc, err := html.Parse(menu.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse a menu: %v", err)
	}

	head, ok := findOnce(doc, []string{"head"}, false)
	if ok {
		head.AppendChild(&html.Node{
			Type: html.ElementNode,
			Data: "link",
			Attr: []html.Attribute{
				{Key: "rel", Val: "stylesheet"},
				{Key: "href", Val: "https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css"},
			},
		})
		// inline the css, this could optimize this with a CDN
		head.AppendChild(&html.Node{
			Type: html.ElementNode,
			Data: "style",
			Attr: []html.Attribute{
				{Key: "type", Val: "text/css"},
			},
			FirstChild: &html.Node{
				Type: html.TextNode,
				Data: styles,
			},
		})
		// inline the javascript, this could optimize this with a CDN
		head.AppendChild(&html.Node{
			Type: html.ElementNode,
			Data: "script",
			Attr: []html.Attribute{
				{Key: "type", Val: "text/javascript"},
			},
			FirstChild: &html.Node{
				Type: html.TextNode,
				Data: javascript,
			},
		})

	}

	return doc, nil
}

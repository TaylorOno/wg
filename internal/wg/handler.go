package wg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"strings"
)

func Handle(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("In Handle")
	resp := InvokeResponse{
		Outputs:     map[string]resData{"res": {}},
		Logs:        []string{},
		ReturnValue: "",
	}

	fmt.Println("calling menu")
	menu, err := http.DefaultClient.Get("https://williamsandgraham.com/drinks/")
	if err != nil {
		resp.Logs = append(resp.Logs, fmt.Sprintf("failed to read menu: %v", err))
		writeJSON(w, resp, http.StatusInternalServerError)
		return
	}
	defer menu.Body.Close()

	doc, err := html.Parse(menu.Body)
	if err != nil {
		resp.Logs = append(resp.Logs, fmt.Sprintf("failed to parse a menu: %v", err))
		writeJSON(w, resp, http.StatusInternalServerError)
		return
	}

	head, ok := findOnce(doc, []string{"head"}, false)
	if ok {
		//head.AppendChild(&html.Node{
		//	Type: html.ElementNode,
		//	Data: "link",
		//	Attr: []html.Attribute{
		//		{Key: "rel", Val: "stylesheet"},
		//		{Key: "href", Val: "./wg.css"},
		//	},
		//})
		//head.AppendChild(&html.Node{
		//	Type: html.ElementNode,
		//	Data: "script",
		//	Attr: []html.Attribute{
		//		{Key: "id", Val: "drink-tracker"},
		//		{Key: "type", Val: "text/javascript"},
		//		{Key: "src", Val: "./script.js"},
		//	},
		//})
		head.AppendChild(&html.Node{
			Type: html.ElementNode,
			Data: "style",
			Attr: []html.Attribute{
				{Key: "type", Val: "text/css"},
			},
			FirstChild: &html.Node{
				Type: html.TextNode,
				Data: CheckBoxCSS,
			},
		})
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

	ApplyToAll(doc, []string{"h4", "class", "white"}, func(node *html.Node) {
		drinkName := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(node.FirstChild.Data), " ", ""))
		node.InsertBefore(&html.Node{
			Type: html.ElementNode,
			Data: "span",
			FirstChild: &html.Node{
				Type: html.ElementNode,
				Data: "input",
				Attr: []html.Attribute{
					{Key: "type", Val: "checkbox"},
					{Key: "class", Val: "check-box"},
					{Key: "id", Val: drinkName},
				},
			},
		}, node.FirstChild)
	})

	out := bytes.Buffer{}
	if err = html.Render(&out, doc); err != nil {
		resp.Logs = append(resp.Logs, fmt.Sprintf("failed to parse a menu: %v", err))
		writeJSON(w, resp, http.StatusInternalServerError)
		return
	}

	writeHTML(w, resp, http.StatusOK, out)
	return
}

type resData = map[string]interface{}

type InvokeResponse struct {
	Outputs     map[string]resData
	Logs        []string
	ReturnValue interface{}
}

func writeJSON(w http.ResponseWriter, result InvokeResponse, code int) {
	result.Outputs["res"]["statuscode"] = code
	result.Outputs["res"]["headers"] = map[string]string{"Content-Type": "application/json"}
	result.Outputs["res"]["body"] = fmt.Sprintf("{\"status\":\"%v\"}", code)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("failed to write response: %v", err)
		return
	}
}

func writeHTML(w http.ResponseWriter, result InvokeResponse, code int, out bytes.Buffer) {
	result.Outputs["res"]["statuscode"] = code
	result.Outputs["res"]["headers"] = map[string]string{"Content-Type": "text/html"}
	result.Outputs["res"]["body"] = string(out.Bytes())

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("failed to write response: %v", err)
		return
	}
}

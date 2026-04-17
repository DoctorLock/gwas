package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

type getVar struct {
	key      string
	required bool
}
type PageElement struct {
	ElementType string
	attributes  map[string]string
}
type Page struct {
	Title       string
	Html        string // name of html file
	RequireAuth bool
	RequestVars map[string]*RequestVar
	PageData    map[string]string
}

func LoadGetRequest(w http.ResponseWriter, p *Page) {
	data := make(map[string]interface{})

	// Add Page fields
	data["Title"] = p.Title

	// Add PageData fields
	for k, v := range p.PageData {
		data[k] = v
	}
	var content, header, footer *template.Template
	content, _ = template.ParseFiles(fmt.Sprintf("./html/public/%s.html", p.Html))
	header, _ = template.ParseFiles("./html/elements/header.html")
	footer, _ = template.ParseFiles("./html/elements/footer.html")
	header.Execute(w, data)
	content.Execute(w, data)
	footer.Execute(w, data)
}

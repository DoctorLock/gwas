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
	attributes map[string]string
	htmlFile   string
}

func (pe *PageElement) Render(w http.ResponseWriter, data map[string]string) bool {
	var content *template.Template
	content, _ = template.ParseFiles(fmt.Sprintf("./html/public/%s.html", pe.htmlFile))
	content.Execute(w, data)
	return true
}

type Page struct {
	Title       string
	Html        string // name of html file
	RequireAuth bool
	RequestVars map[string]*RequestVar
	// ProcessRequest is an optional per-instance hook to populate PageData from the incoming request.
	ProcessRequest func(p *Page, r *http.Request) error
	PageData       map[string]string
}

func (p *Page) processRequest(r *http.Request) error {
	if p.ProcessRequest != nil {
		return p.ProcessRequest(p, r)
	}
	return nil
}

func (p *Page) renderHtml(w http.ResponseWriter, data map[string]interface{}) {

	var content, header, footer *template.Template
	content, _ = template.ParseFiles(fmt.Sprintf("./html/public/%s.html", p.Html))
	header, _ = template.ParseFiles("./html/elements/header.html")
	footer, _ = template.ParseFiles("./html/elements/footer.html")
	header.Execute(w, data)
	content.Execute(w, p.PageData)
	footer.Execute(w, data)

}

func (p *Page) LoadGetRequest(w http.ResponseWriter, r *http.Request) {
	// Ensure PageData map exists
	if p.PageData == nil {
		p.PageData = make(map[string]string)
	}

	// Let the instance populate PageData from the request
	_ = p.processRequest(r)

	// If the instance set a redirect in PageData, perform it and stop processing
	if p.PageData != nil {
		if redirect, ok := p.PageData["redirect"]; ok && redirect != "" {
			http.Redirect(w, r, redirect, http.StatusFound)
			return
		}
	}

	// Prepare template data by merging PageData
	data := make(map[string]interface{})

	// Add Page fields
	data["Title"] = p.Title

	p.renderHtml(w, data)
}

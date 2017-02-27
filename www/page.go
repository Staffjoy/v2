package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
)

type page struct {
	Title        string // Used in <title>
	Description  string // SEO matters
	TemplateName string // e.g. home.html
	CSSId        string // e.g. 'careers'
	Version      string // e.g. master-1, for cachebusting
	CsrfField    template.HTML
}

func (p *page) Handler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "text/html; charset=UTF-8")

	if p.Description == "" {
		p.Description = defaultDescription
	}

	p.CsrfField = csrf.TemplateField(req)

	err := tmpl.ExecuteTemplate(res, p.TemplateName, p)

	if err != nil {
		logger.Panicf("Unable to render page %s - %s", p.Title, err)
	}
}

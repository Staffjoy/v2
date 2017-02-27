// Package apidocs serves a Swagger-page on an api at the path /page/
package apidocs

import (
	"fmt"
	"html/template"
	"mime"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/elazarl/go-bindata-assetfs"
)

const (
	// Prefix is the http path that renders the page
	Prefix = "/ui/"
	// All templates in this foilder will be loaded
	homeTemplate = "assets/templates/index.tmpl"
)

var (
	tmpl *template.Template
	// asset folders that are served directly
	assetPaths = []string{"css", "js", "images", "lang", "lib"}
)

type page struct {
	logger *logrus.Entry
}

func init() {
	tmplData, err := Asset(homeTemplate)
	if err != nil {
		panic("Unable to locate index template for swaggerpage")
	}
	tmpl, err = template.New(homeTemplate).Parse(string(tmplData))
	if err != nil {
		panic("Unable to parse swaggerpage template")
	}

}

// Serve runs the page using the path `/page/`
func Serve(mux *http.ServeMux, logger *logrus.Entry) {
	mime.AddExtensionType(".svg", "image/svg+xml")
	p := page{logger: logger} // todo - pass in option
	docMux, err := p.newMux()
	if err != nil {
		logger.Fatalf("Unable to process swagger page - %v", err)
	}
	mux.Handle(Prefix, docMux)
}

func (p *page) newMux() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	// Register asset folders we want served externally
	for _, path := range assetPaths {
		urlPath := fmt.Sprintf("%s%s/", Prefix, path) // Wrap in slashes
		mux.Handle(urlPath, http.StripPrefix(urlPath, http.FileServer(
			&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "assets/" + path})))
	}

	mux.HandleFunc(Prefix, func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Header().Set("Content-Type", "text/html; charset=UTF-8")
		err := tmpl.ExecuteTemplate(res, homeTemplate, p)
		if err != nil {
			p.logger.Panicf("Unable to render swaggerpage index %v", err)
		}
	})

	return mux, nil
}

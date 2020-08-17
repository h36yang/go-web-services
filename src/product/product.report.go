/*
Package product - This file contains code for the Templating module
*/
package product

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
)

// ReportFilter class
type ReportFilter struct {
	NameFilter         string `json:"productName"`
	ManufacturerFilter string `json:"manufacturer"`
	SKUFilter          string `json:"sku"`
}

func handleProductReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// generate report by filters
		var productFilter ReportFilter
		err := json.NewDecoder(r.Body).Decode(&productFilter)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		products, err := searchForProductData(productFilter)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		templateName := "report.gotmpl"
		funcMap := template.FuncMap{"mod": func(x, y int) bool {
			return x%y == 0
		}}
		t := template.New(templateName).Funcs(funcMap)
		t, err = t.ParseFiles(path.Join("templates", templateName))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		var tmpl bytes.Buffer
		if len(products) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		err = t.Execute(&tmpl, products)
		reader := bytes.NewReader(tmpl.Bytes())
		w.Header().Set("Content-Disposition", "attachment")
		w.Header().Set("Content-Type", "text/html")
		http.ServeContent(w, r, "report.html", time.Now(), reader)

	case http.MethodOptions:
		// Check CORS headers - do nothing here; let Middleware handler handle this

	default:
		// Unhandled methods
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

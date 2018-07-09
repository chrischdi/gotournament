package frontend

import (
	"fmt"
	"net/http"
)

func head(w http.ResponseWriter, title string) error {
	vars := map[string]interface{}{
		"Title": title,
		"Nav":   MenuLinks,
	}
	return writeTemplate(w, "head.tmpl", vars)
}

func foot(w http.ResponseWriter) error {
	return writeTemplate(w, "foot.tmpl", nil)
}

func errs(w http.ResponseWriter, e error) {
	fmt.Fprintf(w, "<div style=\"font-size: 1.5em; color: red;\">Error: %v</div>\n", e)
}

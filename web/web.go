package web

import (
	"fmt"
	"net/http"
)

type Assignments struct{}

func (h *Assignments) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

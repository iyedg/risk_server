package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	template *template.Template
	filename string
	once     sync.Once
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.template = template.Must(
		template.New(t.filename).Delims("[[", "]]").ParseFiles(filepath.Join(".", t.filename)),
	)
	err := t.template.Execute(w, &ip{IP: localIP()})
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

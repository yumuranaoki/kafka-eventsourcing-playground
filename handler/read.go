package main

import (
	"fmt"
	"net/http"

	"github.com/lovoo/goka"
)

func read(view *goka.View) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		val, _ := view.Get("id")
		if val == nil {
			fmt.Fprintf(w, "%s not found!", "id")
			return
		}

		id, ok := val.(string)
		if ok {
			fmt.Fprintf(w, "hello %s", id)
		} else {
			fmt.Fprintln(w, "not ok")
		}
	}
}

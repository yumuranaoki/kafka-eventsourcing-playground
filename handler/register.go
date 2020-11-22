package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/lovoo/goka"
)

type message struct {
	Id string `json:"id"`
}

func register(emitter *goka.Emitter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var m message
		b, err := ioutil.ReadAll(r.Body)

		if err != nil {
			fmt.Fprintf(w, "error: %v", err)
			return
		}

		err = json.Unmarshal(b, &m)
		if err != nil {
			fmt.Fprintf(w, "error: %v", err)
			return
		}

		id := m.Id

		if err = emitter.EmitSync("id", id); err != nil {
			fmt.Fprintf(w, "error: %s", err)
		}

		fmt.Fprintf(w, "hello %s", id)
	}
}

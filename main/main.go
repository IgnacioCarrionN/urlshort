package main

import (
	"fmt"
	"net/http"
	"urlshort"
)

const PORT  = ":8082"

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathToUrls := map[string]string {
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc": "https://godoc.org/gopkg.in/yaml.v2",
	}

	mapHandler := urlshort.MapHandler(pathToUrls, mux)

	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`

	jsonString := "[{\"path\": \"/tame\", \"url\": \"https://institutotame.com\"}]"

	_, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	jsonHandler, err := urlshort.JsonHandler(jsonString, mapHandler)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Server listening on port %s", PORT)
	http.ListenAndServe(PORT, jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}
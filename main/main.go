package main

import (
	"fmt"
	"net/http"
	"urlshort"
)

const PORT = ":8082"

func main() {
	mux := defaultMux()
	mux.HandleFunc("/new", urlshort.NewHandler)

	// Build the MapHandler using the mux as the fallback
	pathToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	mapHandler := urlshort.MapHandler(pathToUrls, mux)

	/*
		jsonString := "[{\"path\": \"/tame\", \"url\": \"https://institutotame.com\"}]"

		jsonHandler, err := urlshort.JsonHandler(jsonString, mapHandler)
		if err != nil {
			panic(err)
		}
	*/
	sqlHandler, err := urlshort.SqlHandler(mapHandler)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Server listening on port %s", PORT)

	err = http.ListenAndServe(PORT, sqlHandler)
	if err != nil {
		panic(err)
	}

	/*
		certmagic.Agreed = true
		certmagic.Email = "ignacio@institutotame.com"
		certmagic.CA = certmagic.LetsEncryptStagingCA

		err = certmagic.HTTPS([]string{"example.com"}, jsonHandler)
		if err != nil {
			panic(err)
		}

	*/
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Address: %s\n", r.RemoteAddr)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintln(w, "Hello, World!")
	if err != nil {
		panic(err)
	}
}

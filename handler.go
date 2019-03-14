package urlshort

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func MapHandler(pathToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if dest, ok := pathToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func NewHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	var url string
	var path string
	if len(query["url"]) > 0 {
		url = query["url"][0]
	} else {
		fmt.Fprint(w, "Error parsing Url")
	}
	if len(query["path"]) > 0 {
		path = "/" + query["path"][0]
	} else {
		fmt.Fprint(w, "Error parsing path")
	}

	err := newUrl(path, url)
	if err != nil {
		fmt.Fprintf(w, "Error saving data: %s", err)
	}
}

func JsonHandler(json string, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseJson(json)
	if err != nil {
		return nil, err
	}
	pathToUrls := buildMap(pathUrls)
	return MapHandler(pathToUrls, fallback), nil
}

func SqlHandler(fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseSql()
	if err != nil {
		return nil, err
	}
	pathToUrls := buildMap(pathUrls)
	return MapHandler(pathToUrls, fallback), nil
}

func buildMap(pathUrls []pathUrl) map[string]string {
	pathToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathToUrls[pu.Path] = pu.URL
	}
	return pathToUrls
}

func parseJson(jsonString string) ([]pathUrl, error) {
	var pathUrls []pathUrl
	err := json.Unmarshal([]byte(jsonString), &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func parseSql() ([]pathUrl, error) {
	database := openDatabase()
	rows, _ := database.Query("SELECT Path, URL FROM urls")
	pathUrls, err := sqlToStruct(rows)
	if err != nil {
		return nil, err
	}
	err = database.Close()
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func sqlToStruct(rows *sql.Rows) ([]pathUrl, error) {
	pathUrls := make([]pathUrl, 0)
	for rows.Next() {
		pathUrl := pathUrl{}
		err := rows.Scan(&pathUrl.Path, &pathUrl.URL)
		if err != nil {
			return nil, err
		}
		pathUrls = append(pathUrls, pathUrl)
		fmt.Println(pathUrl)
	}
	return pathUrls, nil
}

func newUrl(path string, url string) error {
	pathUrl := pathUrl{path, url}
	database := openDatabase()
	statement, err := database.Prepare("INSERT INTO urls (path, url) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(pathUrl.Path, pathUrl.URL)
	if err != nil {
		return err
	}
	err = database.Close()
	if err != nil {
		return err
	}
	return nil
}

func openDatabase() *sql.DB {
	database, _ := sql.Open("sqlite3", "urls.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS urls (id INTEGER PRIMARY KEY, Path TEXT, URL TEXT)")
	_, e := statement.Exec()
	if e != nil {
		fmt.Println(e)
	}
	return database
}

type pathUrl struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

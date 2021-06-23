package main

import (
	"database/sql"
	"flag"
	"github.com/dop251/goja"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"jhp"
	"log"
	"net/http"
	"os"
	"path"
)
var contentRoot = "content"


func connectDb() *sql.DB {
	p := os.Getenv("DB")
	if p == "" {
		p = path.Join(contentRoot, "babynames-gendered-2015.sqlite")
	}
	db , err := sql.Open("sqlite3", p)
	if err != nil {
		panic(err)
	}
	return db
}

func lookupPage(p string) (string, error) {

	c, err := ioutil.ReadFile(path.Join(contentRoot, p))
	if err != nil {
		return "Uh Oh, path not found.", err
	}

	return string(c), nil
}

func main() {
	flag.StringVar(&contentRoot, "content", contentRoot, "The place where your jhp and js lives")
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		req := jhp.RequestReader{
			DB:     connectDb(),
			Params: jhp.ToParamMap(r),
		}
		vm := goja.New()

		log.Println("has page..", r.URL.Path)
		var content string
		var p string

		if r.URL.Path == "/" {
			p = "index.jhp"
		} else {
			p = r.URL.Path
		}

		content, err := lookupPage(p)
		if err != nil {
			log.Println("ERROR:", err)
			http.Error(w, err.Error(), 500)
			return
		}

		err = jhp.Render(vm, req, w, content)

		if err != nil {
			log.Println("ERROR:", err)
			http.Error(w, err.Error(), 500)
		}
	})
	log.Println("server running on :9494")
	err := http.ListenAndServe(":9494", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

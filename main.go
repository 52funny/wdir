package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"path/filepath"

	"github.com/52funny/wdir/controller"
	"github.com/valyala/fasthttp"
)

func main() {
	port := flag.String("port", "8080", "the server port")
	path := flag.String("p", "/Users/52funny", "the server path")
	tPath := flag.String("t", "compress", "the template path")
	flag.Parse()
	t, err := template.ParseFiles(
		filepath.Join(*tPath, "index.html"),
		filepath.Join(*tPath, "bulma.min.css.html"),
		filepath.Join(*tPath, "main.css.html"),
		filepath.Join(*tPath, "header.html"),
		filepath.Join(*tPath, "main.js.html"),
	)
	if err != nil {
		log.Println(err)
	}

	handler := controller.HandleFastHTTP(*path, t, *tPath)
	fmt.Printf("Listen on http://localhost:%v\n", *port)
	if err := fasthttp.ListenAndServe(":"+*port, handler); err != nil {
		panic(err)
	}
}

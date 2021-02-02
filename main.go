package main

import (
	"flag"
	"html/template"
	"path/filepath"

	"github.com/52funny/wdir/config"
	"github.com/52funny/wdir/controller"
	"github.com/52funny/wdir/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	configName := flag.String("c", "config", "the config name")
	flag.Parse()
	err := config.ReadConfig(*configName)
	if err != nil {
		panic(err)
	}
	utils.InitLogger(config.LogPath)
}

func main() {

	t, err := template.ParseFiles(
		filepath.Join(config.Template, "index.html"),
		filepath.Join(config.Template, "bulma.min.css.html"),
		filepath.Join(config.Template, "main.css.html"),
		filepath.Join(config.Template, "header.html"),
		filepath.Join(config.Template, "main.js.html"),
	)
	if err != nil {
		utils.Log.Println(err)
	}

	handler := controller.HandleFastHTTP(config.Path, t, config.Template)
	utils.Log.Printf("Listen on http://localhost:%v\n", config.Port)
	if err := fasthttp.ListenAndServe(":"+config.Port, handler); err != nil {
		utils.Log.Println(err)
		panic(err)
	}
}

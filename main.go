package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"

	"github.com/52funny/wdir/config"
	"github.com/52funny/wdir/controller"
	"github.com/52funny/wdir/utils"
	"github.com/valyala/fasthttp"
)

//go:embed static/*
var embedF embed.FS

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	configName := flag.String("c", "config.yaml", "the config name")
	flag.Parse()
	err := config.ReadConfig(*configName)
	if err != nil {
		panic(err)
	}
	utils.InitLogger(config.Config.LogPath)
}

func main() {
	t, err := template.ParseFS(embedF,
		"static/index.html",
		"static/bulma.min.css.html",
		"static/main.css.html",
		"static/header.html",
		"static/main.js.html",
	)
	if err != nil {
		utils.Log.Fatal(err)
	}

	fsH := fasthttp.FSHandler(config.Config.Path, 0)
	handler := controller.HandleFastHTTP(fsH, t, &embedF, config.Config.Path)
	address := utils.GetNetAddress()

	fmt.Println("Version:", version, "Commit:", commit)
	fmt.Println("You can now view list in the browser.")
	fmt.Printf("  Local:%10c  http://localhost:%v\n", ' ', config.Config.Port)
	for _, addr := range address {
		fmt.Printf("  On Your NetWork:  http://%v:%v\n", addr, config.Config.Port)
	}

	if err := fasthttp.ListenAndServe("0.0.0.0:"+config.Config.Port, handler); err != nil {
		utils.Log.Println(err)
		panic(err)
	}
}

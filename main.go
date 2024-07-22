package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/52funny/wdir/config"
	"github.com/52funny/wdir/controller"
	"github.com/52funny/wdir/utils"
	"github.com/valyala/fasthttp"
)

//go:embed assets/*
var embedF embed.FS

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	configName := flag.String("c", "config.yaml", "the config name")
	port := flag.String("p", "80", "server port")
	flag.Parse()

	path := "."
	if len(flag.Args()) > 0 {
		path = flag.Arg(0)
	}
	err := config.ReadConfig(*configName, *port, path)
	if err != nil {
		panic(err)
	}
	utils.InitLogger(config.Config.LogPath)
}

func main() {
	t, err := template.ParseFS(embedF, "assets/index.html", "assets/header.html")
	if err != nil {
		utils.Log.Fatal(err)
	}

	rootPath, err := filepath.Abs(config.Config.Path)
	if err != nil {
		panic(err)
	}
	fsH := fasthttp.FSHandler(config.Config.Path, 0)
	handler := fasthttp.CompressHandler(controller.HandleFastHTTP(fsH, t, &embedF, rootPath, commit))
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

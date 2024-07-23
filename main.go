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

var certFile string
var keyFile string

func init() {
	configName := flag.String("c", "config.yaml", "the config name")
	port := flag.String("p", "80", "server port")
	flag.StringVar(&certFile, "tls-cert", "", "Path to SSL/TLS certificate")
	flag.StringVar(&keyFile, "tls-key", "", "Path to SSL/TLS certificate's private key")
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
	addr := fmt.Sprintf("0.0.0.0:%s", config.Config.Port)
	if len(certFile) > 0 && len(keyFile) > 0 {
		displayAddress(address, true)
		err = fasthttp.ListenAndServeTLS(addr, certFile, keyFile, handler)
	} else {
		displayAddress(address, false)
		err = fasthttp.ListenAndServe(addr, handler)
	}
	if err != nil {
		utils.Log.Fatal(err)
	}
}

func displayAddress(address []string, https bool) {
	scheme := "http"
	if https {
		scheme = "https"
	}
	fmt.Println("You can now view list in the browser.")
	fmt.Printf("  Local:%10c  %s://localhost:%v\n", ' ', scheme, config.Config.Port)
	for _, addr := range address {
		fmt.Printf("  On Your NetWork:  %s://%v:%v\n", scheme, addr, config.Config.Port)
	}
}

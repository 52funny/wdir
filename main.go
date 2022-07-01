package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"net"

	"github.com/52funny/wdir/config"
	"github.com/52funny/wdir/controller"
	"github.com/52funny/wdir/utils"
	"github.com/valyala/fasthttp"
)

//go:embed static/*
var embedF embed.FS

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

	handler := controller.HandleFastHTTP(config.Config.Path, t, &embedF)
	address := getNetAddress()
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

func getNetAddress() []string {
	netS := make([]string, 0)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				netS = append(netS, ipNet.IP.String())
			}
		}
	}
	return netS
}

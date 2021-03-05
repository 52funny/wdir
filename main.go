package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"net"
	"path/filepath"

	"github.com/52funny/wdir/config"
	"github.com/52funny/wdir/controller"
	"github.com/52funny/wdir/utils"
	"github.com/valyala/fasthttp"
)

//go:embed compress/icon.woff
var embedF embed.FS

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

	handler := controller.HandleFastHTTP(config.Path, t, config.Template, &embedF)
	address := getNetAddress()
	fmt.Println("You can now view list in the browser.")
	fmt.Printf("  Local:%10c  http://localhost:%v\n", ' ', config.Port)
	for _, addr := range address {
		fmt.Printf("  On Your NetWork:  http://%v:%v\n", addr, config.Port)
	}
	if err := fasthttp.ListenAndServe("0.0.0.0:"+config.Port, handler); err != nil {
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

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
	"golang.org/x/sync/errgroup"
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

func RedirectHTTPS(fsH fasthttp.RequestHandler, tls bool) fasthttp.RequestHandler {
	if !tls {
		return fsH
	}
	return func(ctx *fasthttp.RequestCtx) {
		if !ctx.IsTLS() {
			ctx.Redirect("https://"+string(ctx.Host())+string(ctx.RequestURI()), fasthttp.StatusPermanentRedirect)
			return
		}
		fsH(ctx)
	}
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
	tls := len(certFile) > 0 && len(keyFile) > 0
	fsH := fasthttp.FSHandler(config.Config.Path, 0)
	handler := RedirectHTTPS(fasthttp.CompressHandler(controller.HandleFastHTTP(fsH, t, &embedF, rootPath, commit)), tls)
	ipv4Addrs, ipv6Addrs := utils.GetNetIPv4Address(), utils.GetNetIPv6Address()
	addrs := append(ipv4Addrs, ipv6Addrs...)

	fmt.Println("Version:", version, "Commit:", commit)

	displayAddress(addrs, tls)

	ipv4Addr := fmt.Sprintf(":%v", config.Config.Port)
	ipv6Addr := fmt.Sprintf("[::]:%v", config.Config.Port)

	srv := &fasthttp.Server{
		Handler: handler,
	}
	v4Ln, err := net.Listen("tcp4", ipv4Addr)
	if err != nil {
		utils.Log.Fatal(err)
	}
	v6Ln, err := net.Listen("tcp6", ipv6Addr)
	if err != nil {
		utils.Log.Fatal(err)
	}
	defer v4Ln.Close()
	defer v6Ln.Close()

	errGroup := new(errgroup.Group)
	errGroup.Go(func() error {
		return srv.Serve(v4Ln)
	})
	errGroup.Go(func() error {
		return srv.Serve(v6Ln)
	})
	if tls {
		errGroup.Go(func() error {
			return srv.ServeTLS(v4Ln, certFile, keyFile)
		})
		errGroup.Go(func() error {
			return srv.ServeTLS(v6Ln, certFile, keyFile)
		})
	}
	if err := errGroup.Wait(); err != nil {
		utils.Log.Fatalln(err)
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

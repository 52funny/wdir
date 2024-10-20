package main

import (
	"embed"
	"fmt"
	"html/template"
	"net"
	"os"
	"path/filepath"

	"github.com/52funny/wdir/controller"
	"github.com/52funny/wdir/model"
	"github.com/52funny/wdir/utils"
	"github.com/jessevdk/go-flags"
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

var username string
var password string

type Options struct {
	Config         string   `short:"c" long:"config" description:"the config name" default:"config.yaml" required:"false"`
	Port           string   `short:"p" long:"port" description:"server port" default:"80" env:"PORT"`
	TLSCert        string   `long:"tls-cert" description:"Path to SSL/TLS certificate" required:"false"`
	TLSKey         string   `long:"tls-key" description:"Path to SSL/TLS certificate's private key" required:"false"`
	Username       []string `long:"user" description:"username for basic auth" required:"false"`
	Password       []string `long:"pass" description:"password for basic auth" required:"false"`
	LogPath        string   `long:"log" description:"log path" default:"log" required:"false" env:"LOGPATH"`
	HiddenDotFiles bool     `long:"show" description:"show hidden files" required:"false" env:"SHOWHIDDENFILES"`
}

var opts Options
var parser = flags.NewParser(&opts, flags.Default)

func initFlags() string {
	args, err := parser.Parse()
	if err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			fmt.Println(flagsErr)
			os.Exit(1)
		}
	}

	path := "."
	if len(args) > 0 {
		path = args[0]
	}
	utils.InitLogger(opts.LogPath)
	return path
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
	path := initFlags()
	t, err := template.ParseFS(embedF, "assets/index.html", "assets/header.html")
	if err != nil {
		utils.Log.Fatal(err)
	}

	rootPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	tls := len(certFile) > 0 && len(keyFile) > 0
	fsH := fasthttp.FSHandler(path, 0)
	basic := model.NewBasicAuth(opts.Username, opts.Password)
	handler := RedirectHTTPS(fasthttp.CompressHandler(controller.BasicAuth(basic, controller.HandleFastHTTP(fsH, t, &embedF, rootPath, commit))), tls)
	ipv4Addrs, ipv6Addrs := utils.GetNetIPv4Address(), utils.GetNetIPv6Address()
	addrs := append(ipv4Addrs, ipv6Addrs...)

	fmt.Println("Version:", version, "Commit:", commit)

	displayAddress(addrs, tls)

	ipv4Addr := fmt.Sprintf(":%v", opts.Port)

	ipv6Addr := fmt.Sprintf("[::]:%v", opts.Port)

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
	fmt.Printf("  Local:%10c  %s://localhost:%v\n", ' ', scheme, opts.Port)
	for _, addr := range address {
		fmt.Printf("  On Your NetWork:  %s://%v:%v\n", scheme, addr, opts.Port)
	}
}

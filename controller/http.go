package controller

import (
	"embed"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/52funny/wdir/config"
	"github.com/52funny/wdir/model"
	"github.com/52funny/wdir/utils"
	"github.com/h2non/filetype"
	"github.com/valyala/fasthttp"
)

// HandleFastHTTP is handle path matching
func HandleFastHTTP(path string, t *template.Template, embedF *embed.FS) fasthttp.RequestHandler {
	iconWoff, err := embedF.ReadFile("static/icon.woff")
	if err != nil {
		panic(err)
	}
	// for file
	fileHandler := fasthttp.FSHandler(path, 0)
	return func(ctx *fasthttp.RequestCtx) {
		method := string(ctx.Method())

		urlPath := string(ctx.Path())
		utils.Log.Println(method, urlPath)

		// if not show hidden files
		if !config.Config.ShowHiddenFiles {
			items := strings.Split(urlPath, "/")
			if len(items) > 0 {
				for _, it := range items {
					if len(it) > 0 && it[0:1] == "." {
						ctx.Error("no files", http.StatusNotFound)
						return
					}
				}
			}
		}
		if urlPath == "/icon.woff" {
			ctx.Write(iconWoff)
			return
		}

		state, err := os.Stat(filepath.Join(path, urlPath))
		if err != nil {
			if os.IsNotExist(err) {
				ctx.Error(urlPath+" Not Found", http.StatusNotFound)
				return
			}
			ctx.Error(err.Error(), http.StatusInternalServerError)
			return
		}

		// is Dir
		if state.IsDir() {
			files, err := ioutil.ReadDir(filepath.Join(path, urlPath))
			if err != nil {
				ctx.Error(err.Error(), http.StatusInternalServerError)
				return
			}

			directoryArray := make([]model.File, 0, len(files)/2)
			fileArray := make([]model.File, 0, len(files)/2)

			for _, item := range files {
				if !config.Config.ShowHiddenFiles && item.Name()[0:1] == "." {
					continue
				}
				types := ""
				fileInfo, _ := os.Stat(filepath.Join(path, urlPath, item.Name()))
				if fileInfo.IsDir() {
					types = "folder"
					directoryArray = append(directoryArray, model.File{
						Path: filepath.Join(urlPath, item.Name()),
						Fileinfo: model.FileInfo{
							Name: item.Name(),
							Type: types,
							Date: item.ModTime().Format("2006-01-02 15:04:05"),
							Size: utils.ConvertSize(item.Size()),
						},
					})
				} else {
					kind, _ := filetype.MatchFile(filepath.Join(path, urlPath, item.Name()))
					types = kind.Extension
					fileArray = append(fileArray, model.File{
						Path: filepath.Join(urlPath, item.Name()),
						Fileinfo: model.FileInfo{
							Name: item.Name(),
							Type: types,
							Date: item.ModTime().Format("2006-01-02 15:04:05"),
							Size: utils.ConvertSize(item.Size()),
						},
					})
				}

			}
			t.Execute(ctx, append(directoryArray, fileArray...))
			ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
		} else {
			fileHandler(ctx)
		}
	}
}

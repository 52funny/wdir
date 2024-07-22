package controller

import (
	"embed"
	"fmt"
	"html/template"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/52funny/wdir/config"
	"github.com/52funny/wdir/model"
	"github.com/52funny/wdir/utils"
	"github.com/h2non/filetype"
	"github.com/valyala/fasthttp"
)

// HandleFastHTTP is handle path matching
func HandleFastHTTP(fsH fasthttp.RequestHandler, t *template.Template, embedF *embed.FS, rootPath string, commit string) fasthttp.RequestHandler {
	assetsPath := fmt.Sprintf("/%s", commit)

	return func(ctx *fasthttp.RequestCtx) {
		method := string(ctx.Method())
		urlPath := string(ctx.Path())

		// retrieve some files in the asset directory
		if strings.HasPrefix(urlPath, assetsPath) {
			relPath, _ := filepath.Rel(assetsPath, urlPath)
			actualPath := fmt.Sprintf("assets/%s", relPath)
			bs, err := embedF.ReadFile(actualPath)
			if err != nil {
				ctx.Error(fmt.Sprintf("Error get %s", urlPath), http.StatusNotFound)
				return
			}
			mime := mime.TypeByExtension(actualPath[strings.LastIndex(actualPath, "."):])
			ctx.Response.Header.Set("Content-Type", mime)
			ctx.Write(bs)
			return
		}

		// log url path
		utils.Log.Println(method, urlPath)

		// check whether it is a hidden directory
		dstPath := filepath.Join(rootPath, urlPath)
		if !config.Config.ShowHiddenFiles && utils.PathHidden(dstPath) {
			ctx.Error(urlPath+"Not Found", http.StatusNotFound)
			return
		}

		stat, err := os.Stat(dstPath)
		if err != nil {
			if os.IsNotExist(err) {
				ctx.Error(urlPath+" Not Found", http.StatusNotFound)
				return
			}
			ctx.Error(err.Error(), http.StatusInternalServerError)
			return
		}

		switch stat.IsDir() {
		case true:
			RenderDir(ctx, t, dstPath, urlPath, commit)
		case false:
			fsH(ctx)
		}
	}
}

func RenderDir(ctx *fasthttp.RequestCtx, t *template.Template, dirPath, urlPath, commit string) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	directoryList := make([]model.File, 0, len(files))
	fileList := make([]model.File, 0, len(files))

	for _, item := range files {
		if !config.Config.ShowHiddenFiles && utils.FileHidden(item.Name()) {
			continue
		}
		fStat, _ := os.Stat(filepath.Join(dirPath, item.Name()))
		f := model.File{
			Path: filepath.Join(urlPath, item.Name()),
			FileInfo: model.FileInfo{
				Name: item.Name(),
				Date: fStat.ModTime().Format(time.DateTime),
				Size: utils.ConvertSize(fStat.Size()),
			},
		}
		switch fStat.IsDir() {
		case true:
			f.FileInfo.Type = "folder"
			directoryList = append(directoryList, f)
		case false:
			kind, _ := filetype.MatchFile(filepath.Join(dirPath, item.Name()))
			f.FileInfo.Type = kind.Extension
			fileList = append(fileList, f)
		}
	}
	m := make(map[string]interface{})
	m["Data"] = append(directoryList, fileList...)
	m["Assets"] = commit
	t.Execute(ctx, m)
	ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
}

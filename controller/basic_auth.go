package controller

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/52funny/wdir/model"
	"github.com/valyala/fasthttp"
)

func BasicAuth(b *model.BasicAuth, req fasthttp.RequestHandler) fasthttp.RequestHandler {
	if b.Len() == 0 {
		return req
	}
	return func(ctx *fasthttp.RequestCtx) {
		auth := ctx.Request.Header.Peek("Authorization")
		token := ctx.Request.URI().QueryArgs().Peek("token")

		if len(auth) == 0 && len(token) == 0 {
			sendAuthReq(ctx)
			return
		}

		if !b.Auth(auth) && !b.AuthToken(token) {
			ctx.WriteString("Unauthorized")
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			return
		}
		if len(auth) > 6 {
			sha := sha256.Sum256(auth[6:])
			ctx.SetUserValue("basicToken", hex.EncodeToString(sha[:]))
		}
		req(ctx)
	}
}

func sendAuthReq(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.DisableNormalizing()
	ctx.Response.Header.SetStatusCode(fasthttp.StatusUnauthorized)
	ctx.Response.Header.Add("WWW-Authenticate", `Basic realm="Come in to make fasthttp great again"`)
}

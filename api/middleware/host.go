package middleware

import (
	"context"
	"mianshiba/pkg/ctxcache"
	"mianshiba/types/consts"

	"github.com/cloudwego/hertz/pkg/app"
)

func SetHostMW() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		ctxcache.Store(c, consts.HostKeyInCtx, string(ctx.Host()))
		ctxcache.Store(c, consts.RequestSchemeKeyInCtx, string(ctx.GetRequest().Scheme()))
		ctx.Next(c)
	}
}

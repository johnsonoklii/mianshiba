package middleware

import (
	"context"
	"mianshiba/api/internal/httputil"
	"mianshiba/application/user"
	"mianshiba/pkg/ctxcache"
	"mianshiba/pkg/errorx"
	"mianshiba/pkg/jwt"
	"mianshiba/pkg/logs"
	"mianshiba/types/consts"
	"mianshiba/types/errno"

	"github.com/cloudwego/hertz/pkg/app"
)

var noNeedSessionCheckPath = map[string]bool{
	"/api/user/register": true,
	"/api/user/login":    true,
}

func SessionAuthMW() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		requestAuthType := ctx.GetInt32(RequestAuthTypeStr)
		if requestAuthType != int32(RequestAuthTypeWebAPI) {
			ctx.Next(c)
			return
		}

		if noNeedSessionCheckPath[string(ctx.GetRequest().URI().Path())] {
			ctx.Next(c)
			return
		}

		jwtToken := string(ctx.Cookie(jwt.TokenKey))
		if len(jwtToken) == 0 {
			logs.Errorf("[SessionAuthMW] jwt token is nil")
			httputil.InternalError(c, ctx,
				errorx.New(errno.ErrUserAuthenticationFailed, errorx.KV("reason", "missing jwt_key in cookie")))
			return
		}

		claims, err := jwt.ParseToken(jwtToken)
		if err != nil {
			logs.Errorf("[SessionAuthMW] parse jwt token failed, err: %v", err)
			httputil.InternalError(c, ctx, err)
			return
		}

		// redis
		redisToken, err := user.UserApplicationSVC.UserDomainSVC.GetJwtToken(c, claims.UserID)
		if err != nil {
			logs.Errorf("[SessionAuthMW] get session failed, err: %v", err)
			httputil.InternalError(c, ctx, err)
			return
		}

		if redisToken != jwtToken {
			logs.Errorf("[SessionAuthMW] jwt token not match, err: %v", err)
			httputil.InternalError(c, ctx,
				errorx.New(errno.ErrUserAuthenticationFailed, errorx.KV("reason", "jwt token not match")))
			return
		}

		if claims != nil {
			ctxcache.Store(c, consts.SessionDataKeyInCtx, claims)
		}

		ctx.Next(c)
	}
}

package mianshiba

import (
	"context"
	"mianshiba/api/internal/httputil"

	"github.com/cloudwego/hertz/pkg/app"
)

func invalidParamRequestResponse(c *app.RequestContext, errMsg string) {
	httputil.BadRequest(c, errMsg)
}

func internalServerErrorResponse(ctx context.Context, c *app.RequestContext, err error) {
	httputil.InternalError(ctx, c, err)
}

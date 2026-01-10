package middleware

import (
	"regexp"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

const HeaderAuthorizationKey = "Authorization"

// TODO
var needAuthPath = map[string]bool{}

var needAuthFunc = map[string]bool{}

func parseBearerAuthToken(authHeader string) string {
	if len(authHeader) == 0 {
		return ""
	}
	parts := strings.Split(authHeader, "Bearer")
	if len(parts) != 2 {
		return ""
	}

	token := strings.TrimSpace(parts[1])
	if len(token) == 0 {
		return ""
	}

	return token
}

func isNeedOpenapiAuth(c *app.RequestContext) bool {
	isNeedAuth := false

	uriPath := c.URI().Path()

	for rule, res := range needAuthFunc {
		if regexp.MustCompile(rule).MatchString(string(uriPath)) {
			isNeedAuth = res
			break
		}
	}

	if needAuthPath[string(c.GetRequest().URI().Path())] {
		isNeedAuth = true
	}

	return isNeedAuth
}

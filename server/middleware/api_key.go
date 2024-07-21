package middleware

import (
	"github.com/go-chi/render"
	"jsin/config"
	error2 "jsin/pkg/common/error"
	"jsin/pkg/constants"
	"net/http"
)

func ApiKeyValidateMiddleware(cfg config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get(constants.HeaderAPIKey)
			if len(apiKey) == 0 {
				_ = render.Render(w, r, error2.ErrRenderNoPermissionRequest("Missing X-Api-Key"))
				return
			}

			if apiKey != cfg.Server.ApiKey {
				_ = render.Render(w, r, error2.ErrRenderNoPermissionRequest("Invalid X-Api-Key"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"mariadiezmaback/pkg/response"
)

func Recovery(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					stack := string(debug.Stack())
					logger.Error("Panic recovered",
						slog.Any("error", rvr),
						slog.String("stack", stack),
						slog.String("path", r.URL.Path),
					)
					response.InternalServerError(w, fmt.Sprintf("Panic occurred: %v", rvr))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

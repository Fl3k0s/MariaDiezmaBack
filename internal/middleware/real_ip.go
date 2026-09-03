package middleware

import (
	"net"
	"net/http"
	"strings"
)

// RealIP middleware resolves the real client IP address and normalizes RemoteAddr
func RealIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := GetClientIP(r)
		if clientIP != "" {
			r.RemoteAddr = clientIP
		}
		next.ServeHTTP(w, r)
	})
}

// GetClientIP inspects X-Forwarded-For, X-Real-IP and falls back to RemoteAddr
func GetClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		parts := strings.Split(xForwardedFor, ",")
		for _, part := range parts {
			ip := strings.TrimSpace(part)
			if parsedIP := net.ParseIP(ip); parsedIP != nil {
				return ip
			}
		}
	}

	xRealIP := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if xRealIP != "" {
		if parsedIP := net.ParseIP(xRealIP); parsedIP != nil {
			return xRealIP
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && ip != "" {
		return ip
	}

	return r.RemoteAddr
}

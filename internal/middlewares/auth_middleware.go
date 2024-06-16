package middlewares

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type PublicUrl struct {
	Url string
}

var publicUrls = []PublicUrl{
	{Url: "/programs/all"},
	{Url: "/programs/get-by-id"},
	{Url: "/programs/get-by-name"},
	{Url: "/programs/get-by-category"},
	{Url: "/schedules/all"},
	{Url: "/schedules/get-by-program-id"},
	{Url: "/schedules/get-by-id"},
	{Url: "/schedules/get-by-day"},
	{Url: "/schedules/get-by-date"},
}

const (
	successMessage      = "Operation completed successfully"
	noAuthHeaderMessage = "Authorization header missing"
	noBearerMessage     = "Invalid Authorization header format"
	invalidTokenMessage = "Invalid token"
)

func isPathPublicUrl(path string) bool {
	for _, publicUrl := range publicUrls {
		if publicUrl.Url == path {
			return true
		}
	}
	return false
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		method := r.Method
		path := r.URL.Path
		userAgent := r.UserAgent()
		remoteAddr := r.RemoteAddr

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("Method: %s, Path: %s, User-Agent: %s, RemoteAddr: %s, Timestamp: %s, Status: %s",
				method, path, userAgent, remoteAddr, start.Format(time.RFC3339), noAuthHeaderMessage)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		privateKey := os.Getenv("PRIVATE_KEY")
		publicApiKey := os.Getenv("PUBLIC_API_KEY")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Printf("Method: %s, Path: %s, User-Agent: %s, RemoteAddr: %s, Timestamp: %s, Status: %s",
				method, path, userAgent, remoteAddr, start.Format(time.RFC3339), noBearerMessage)
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		clientKey := strings.TrimPrefix(authHeader, "Bearer ")

		if isPathPublicUrl(r.URL.Path) {
			if clientKey != publicApiKey {
				log.Printf("Method: %s, Path: %s, User-Agent: %s, RemoteAddr: %s, Timestamp: %s, Status: %s",
					method, path, userAgent, remoteAddr, start.Format(time.RFC3339), invalidTokenMessage)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			log.Printf("Method: %s, Path: %s, User-Agent: %s, RemoteAddr: %s, Timestamp: %s, Status: %s",
				method, path, userAgent, remoteAddr, start.Format(time.RFC3339), successMessage)
			next.ServeHTTP(w, r)
			return
		}

		if clientKey != privateKey {
			log.Printf("Method: %s, Path: %s, User-Agent: %s, RemoteAddr: %s, Timestamp: %s, Status: %s",
				method, path, userAgent, remoteAddr, start.Format(time.RFC3339), invalidTokenMessage)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		log.Printf("Method: %s, Path: %s, User-Agent: %s, RemoteAddr: %s, Timestamp: %s, Status: %s",
			method, path, userAgent, remoteAddr, start.Format(time.RFC3339), successMessage)

		next.ServeHTTP(w, r)
	})
}

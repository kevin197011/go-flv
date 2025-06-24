package middleware

import (
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Username string
	Password string
}

// GetAuthConfig returns authentication configuration from environment variables
func GetAuthConfig() AuthConfig {
	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")

	// Default credentials if not set in environment
	if username == "" {
		username = "admin"
	}
	if password == "" {
		password = "admin123"
	}

	return AuthConfig{
		Username: username,
		Password: password,
	}
}

// SessionMiddleware configures session management
func SessionMiddleware(secret string) gin.HandlerFunc {
	if secret == "" {
		secret = "go-flv-secret-key-change-in-production"
	}
	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24, // 24 hours
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	})
	return sessions.Sessions("go-flv-session", store)
}

// AuthMiddleware protects routes requiring authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		authenticated := session.Get("authenticated")

		if authenticated != true {
			// If it's an AJAX request, return JSON
			if c.GetHeader("Content-Type") == "application/json" ||
				c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":     "Authentication required",
					"login_url": "/admin/login",
				})
			} else {
				// Redirect to login page
				c.Redirect(http.StatusSeeOther, "/admin/login")
			}
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireLogin checks if user is authenticated
func RequireLogin(c *gin.Context) bool {
	session := sessions.Default(c)
	authenticated := session.Get("authenticated")
	return authenticated == true
}

// Login authenticates user and creates session
func Login(c *gin.Context, username, password string) bool {
	config := GetAuthConfig()

	if username == config.Username && password == config.Password {
		session := sessions.Default(c)
		session.Set("authenticated", true)
		session.Set("username", username)
		session.Save()
		return true
	}

	return false
}

// Logout destroys user session
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}

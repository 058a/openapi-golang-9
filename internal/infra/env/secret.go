package env

import "os"

func GetJwtSecret() string {
	JwtSecret := os.Getenv("JWT_SECRET")
	if JwtSecret == "" {
		JwtSecret = "secret"
	}
	return JwtSecret
}

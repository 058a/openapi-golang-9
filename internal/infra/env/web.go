package env

import "os"

func GetServiceUrl() string {
	serviceUrl := os.Getenv("SERVICE_URL")
	if serviceUrl == "" {
		serviceUrl = "http://localhost:1323"
	}
	return serviceUrl
}

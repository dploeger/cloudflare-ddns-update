package main

import (
	"github.com/dploeger/cloudflare-ddns-update/pkg"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	_, isDebug := os.LookupEnv("DEBUG")

	if _, ok := os.LookupEnv("CLOUDFLARE_API_TOKEN"); !ok {
		log.Fatalf("Missing CLOUDFLARE_API_TOKEN env variable")
	}
	if _, ok := os.LookupEnv("CLOUDFLARE_ZONE"); !ok {
		log.Fatalf("Missing CLOUDFLARE_ZONE env variable")
	}
	if _, ok := os.LookupEnv("AUTH_USERNAME"); !ok {
		log.Fatalf("Missing AUTH_USERNAME env variable")
	}
	if _, ok := os.LookupEnv("AUTH_PASSWORD"); !ok {
		log.Fatalf("Missing AUTH_USERNAME env variable")
	}

	var api pkg.API
	if a, err := pkg.NewAPI(os.Getenv("CLOUDFLARE_API_TOKEN"), os.Getenv("CLOUDFLARE_ZONE"), isDebug); err != nil {
		log.Fatalf(err.Error())
	} else {
		api = *a
	}

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "running",
		})
	})
	v3 := r.Group("v3", gin.BasicAuth(gin.Accounts{os.Getenv("AUTH_USERNAME"): os.Getenv("AUTH_PASSWORD")}))
	v3.GET("/update", api.DDNSUpdate)
	if err := r.Run(); err != nil {
		log.Fatalf("Can not run service: %s", err)
	}
}

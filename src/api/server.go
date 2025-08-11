package api

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-auth-otp-service/src/config"
	"golang.org/x/sync/errgroup"
	"log"
	"time"
)

var (
	configs = config.GetInstance()
	g       errgroup.Group
)

func Init() (err error) {
	g.Go(func() error {
		return initUserServer()
	})

	if err = g.Wait(); err != nil {
		log.Fatalln(err)
		return err
	}

	return err
}

func getNewRouter() *gin.Engine {
	// set gin to release mode.
	gin.SetMode(gin.ReleaseMode)

	// Initialize new app.
	router := gin.New()

	// Attach CORS middleware.

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "Accept-Language", "User-Agent", "Cache-Control", "Set-Cookie"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
	}))

	// Attach logger middleware.
	router.Use(gin.Logger())

	// Attach recovery middleware.
	router.Use(gin.Recovery())

	return router
}

func initUserServer() error {
	router := getNewRouter()

	//v1 := router.Group("api/v1")
	{
		//users.AuthenticationRouter(v1)
	}

	// Run App.
	if err := router.Run(
		fmt.Sprintf(":%s", configs.Get("APP_PORT")),
	); err != nil {
		return err
	}

	return nil
}

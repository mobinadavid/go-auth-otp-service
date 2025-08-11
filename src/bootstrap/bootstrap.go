package bootstrap

import (
	"context"
	"go-auth-otp-service/src/api"
	"go-auth-otp-service/src/cache"
	"go-auth-otp-service/src/database"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Init() (err error) {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// Initialize Database
	err = database.Init()
	if err != nil {
		log.Fatalf("Database: Service: Failed to Initialize: %v.", err)
	}
	log.Printf("Database: Service: Database  initial successfully. \n")

	// Initialize Cache
	err = cache.Init()
	if err != nil {
		log.Fatal("Failed to Initialize", zap.String("Service", "Cache"), zap.Error(err), zap.Time("timestamp", time.Now()))
	}

	//Initialize api
	go func() {
		err = api.Init()
		if err != nil {
			log.Fatal("Failed to Initialize.", zap.String("Service", "API"), zap.Error(err), zap.Time("timestamp", time.Now()))
		}
		log.Println("Initialized Successfully.", zap.String("Service", "API"), zap.Time("timestamp", time.Now()))
	}()
	log.Println("Application is now running.\nPress CTRL-C to exit")

	// app started ...
	time.Sleep(50 * time.Millisecond)
	log.Printf("Application is now running.Press CTRL-C to exit.\n")
	<-sc

	// Shutting down application
	log.Printf("Application shutting down....   \n")

	// Close Database
	err = database.GetInstance().Close()
	if err != nil {
		log.Fatalf("Databasde Service: Failed to close database. %v. \n", err)
	}
	log.Printf("Databasde Service: database close sucessfully. \n")

	return
}

package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"loaders-online/config"
	"loaders-online/pkg/database"
	"net/http"
	"os"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {

	dbConfig, err := config.New()
	if err != nil {
		logrus.Fatal(err)
	}

	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     dbConfig.DB.Host,
		Port:     dbConfig.DB.Port,
		Username: dbConfig.DB.Username,
		DBName:   dbConfig.DB.Name,
		SSLMode:  dbConfig.DB.SSLMode,
		Password: dbConfig.DB.Password,
	})

	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	srv := &http.Server{
		Addr: ":8080",
	}

	fmt.Println(dbConfig.DB)

	logrus.Info("SERVER STARTED")

	if err := srv.ListenAndServe(); err != nil {
		logrus.Fatal(err)
	}
}

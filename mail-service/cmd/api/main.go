package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Mailer Mail
}

const webport = "8080"

func main() {
	app := Config{
		Mailer: createMail(),
	}

	log.Println("Starting nail service on", webport)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webport),
		Handler: app.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func createMail() Mail {
	Port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        Port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("MAIL_NAME"),
		FromAddress: os.Getenv("MAIL_ADDRESS"),
	}
	return m
}

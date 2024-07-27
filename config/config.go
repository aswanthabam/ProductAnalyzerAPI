package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	MONGODB_URI    string
	MONGODB_DB     string
	SECRET_KEY     string
	RESEND_API_KEY string
	MAIL_ID        string
	HOST           string
}

var Config Configuration

func (cn *Configuration) Load() error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	notFound := ""
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		notFound += "MONGODB_URI, "
	}
	dbname := os.Getenv("MONGODB_DB")
	if dbname == "" {
		notFound += "MONGODB_DB, "
	}
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		notFound += "SECRET_KEY, "
	}
	resend_api_key := os.Getenv("RESEND_API_KEY")
	if resend_api_key == "" {
		notFound += "RESEND_API_KEY, "
	}
	mail_id := os.Getenv("MAIL_ID")
	if mail_id == "" {
		notFound += "MAIL_ID, "
	}
	host := os.Getenv("HOST")
	if host == "" {
		notFound += "HOST, "
	}
	if notFound != "" {
		return fmt.Errorf("ENVIRONMENT VARIABLES NOT FOUND: %s", notFound)
	}
	cn.MONGODB_URI = uri
	cn.MONGODB_DB = dbname
	cn.SECRET_KEY = secret
	cn.RESEND_API_KEY = resend_api_key
	cn.MAIL_ID = mail_id
	cn.HOST = host
	return nil
}

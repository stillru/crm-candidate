package config

import (
	"log"
	"os"
)

func InitLogger() *log.Logger {
	return log.New(os.Stdout, "CRM: ", log.LstdFlags|log.Lshortfile)
}

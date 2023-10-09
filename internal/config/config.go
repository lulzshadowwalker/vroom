package config

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("ERROR: cannot load app config %q\n", err)
	}
}

package main

import (
	"log"

	"github.com/lulzshadowwalker/vroom/internal/database/migrations"
)

func main() {
	err := migrations.Migrate()
	if err != nil {
		log.Fatalf("cannot complete migrations %q", err)
	}
}

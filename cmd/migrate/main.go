package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/simwawatim/book-api/internal/config"
	"github.com/simwawatim/book-api/internal/database"
	"github.com/simwawatim/book-api/migrations"
)

func main() {
	if len(os.Args) < 2 || (os.Args[1] != "up" && os.Args[1] != "down") {
		fmt.Fprintln(os.Stderr, "usage: migrate up|down")
		os.Exit(2)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if os.Args[1] == "up" {
		err = database.MigrateUp(ctx, db, migrations.FS)
	} else {
		err = database.MigrateDown(ctx, db, migrations.FS)
	}
	if err != nil {
		log.Fatal(err)
	}
}

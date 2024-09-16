package main

import (
	"github.com/iamyusuf/go-simple-service/internal/env"
)

const version = "1.0.0"

func main() {
	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		environment: env.GetString("ENV", "development"),
	}

	app := &application{
		config: cfg,
	}

	mux := app.mount()
	_ = app.run(mux)
}

package main

import (
	"fmt"

	"github.com/chatex-com/dependency-injection"
)

type Config struct {
	Listen      string
	DatabaseDSN string
}

func main() {
	container := di.NewContainer()

	container.Set(&Config{
		Listen:      ":8080",
		DatabaseDSN: "host=localhost port=5432 user=postgres dbname=postgres password=postgres",
	})

	var cfg *Config
	err := container.Load(&cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg.Listen, cfg.DatabaseDSN)
}

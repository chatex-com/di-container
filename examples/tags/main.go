package main

import (
	"fmt"

	"github.com/chatex-com/di-container"
)

type Config struct {
	Listen      string
	DatabaseDSN string
}

type Service struct {
	Foo        string  `di:"-"`
	ProdConfig *Config `di:"prod"`
	TestConfig *Config `di:"test,optional"`
}

func main() {
	container := di.NewContainer()

	container.Set(&Config{
		Listen:      ":8080",
		DatabaseDSN: "host=localhost port=5432 user=postgres dbname=postgres password=postgres",
	}, "prod")

	var service Service
	err := container.Resolve(&service)
	if err != nil {
		panic(err)
	}

	fmt.Println(service.Foo, service.ProdConfig, service.TestConfig)
}

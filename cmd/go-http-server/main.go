package main

import (
	"fmt"

	"github.com/aman-void/go-http-server/internal/config"
)

func main() {

	config.MustLoad()
	fmt.Println("Hello from GO HTTP server")
}

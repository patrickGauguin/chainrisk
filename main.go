package main

import (
	"github.com/joho/godotenv"
	"github.com/patrickGauguin/chainrisk/cmd"
)

func main() {
	godotenv.Load()
	cmd.Execute()
}

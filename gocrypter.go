package main

import (
	"gocrypter/cmd"
	"gocrypter/log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Panic(err)
	}
}

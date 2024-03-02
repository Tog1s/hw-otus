package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 3 {
		log.Fatalln("Arguments required. Usage: go-envdir <path> <command> <arguments>")
	}

	env, err := ReadDir(args[0])
	if err != nil {
		log.Fatalln(err)
	}

	statusCode := RunCmd(args[1:], env)
	os.Exit(statusCode)
}

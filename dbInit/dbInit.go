package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	binary, lookErr := exec.LookPath("psql")
	if lookErr != nil {
		panic(lookErr)
	}

	args := []string{
		"psql",
		os.Getenv("DATABASE_URL"),
		"-a",
		"-f",
		"./dbInit/dbInit.sql",
	}

	env := os.Environ()

	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
}

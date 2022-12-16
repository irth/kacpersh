package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	shell := os.Getenv("SHELL")
	if len(shell) == 0 {
		log.Fatalf("Please set the SHELL variable to a supported shell.")
	}

	os.Setenv("IN_KACPERSH", "1")
	term := Term{
		Command: exec.Command(shell, "-l"),
		BufSize: 16,
	}

	outCh := make(chan []byte)

	go func() {
		for buf := range outCh {
			os.Stdout.Write(buf)
		}
	}()

	if err := term.Spawn(outCh); err != nil {
		log.Fatalf("%s", err)
	}
}

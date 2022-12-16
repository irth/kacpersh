package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	os.Setenv("IN_KACPERSH", "1")
	term := Term{
		Command: exec.Command("zsh", "-l"),
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

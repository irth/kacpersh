package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
)

func main() {
	shell := os.Getenv("SHELL")
	if len(shell) == 0 {
		log.Fatalf("Please set the SHELL variable to a supported shell.")
	}

	tempDir, err := createTempDir()
	if err != nil {
		log.Fatalf("creating a temp dir: %s", err)
	}
	defer os.RemoveAll(tempDir)

	socketPath := path.Join(tempDir, "control")

	control := ControlServer{SocketPath: socketPath}
	go func() {
		if err := control.ListenAndServe(); err != nil {
			log.Fatalf("control server: %s", err)
		}
	}()

	os.Setenv("IN_KACPERSH", "1")
	os.Setenv("KACPERSH_SOCK", socketPath)
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

func createTempDir() (string, error) {
	username := "unknown"
	currentUser, err := user.Current()
	if err == nil {
		username = currentUser.Username
	}

	tempDirPattern := fmt.Sprintf("kacpersh-%s-*", username)
	tempDir, err := ioutil.TempDir("", tempDirPattern)
	if err != nil {
		return "", err
	}

	return tempDir, nil
}

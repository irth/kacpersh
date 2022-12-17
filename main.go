package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
)

//go:embed init.zsh
var initZSH []byte

func usage() {
	fmt.Printf("usage %s: kacpersh [init zsh]", os.Args[0])
	os.Exit(1)
}

func main() {
	// TODO: replace with actual argument parsing
	if len(os.Args) > 1 {
		if os.Args[1] == "init" {
			if len(os.Args) != 3 {
				usage()
			}
			switch os.Args[2] {
			case "zsh":
				os.Stdout.Write(initZSH)
			default:
				fmt.Printf("shell \"%s\" is not supported\n", os.Args[2])
				usage()
			}
		} else {
			usage()
		}
		os.Exit(0)
	}

	if path, ok := os.LookupEnv("KACPERSH_DEBUG"); ok {
		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		log.SetOutput(f)
	} else {
		log.SetOutput(io.Discard)
	}

	outCh := make(chan []byte)

	outByteCh := make(chan byte)
	go func() {
		for buf := range outCh {
			os.Stdout.Write(buf)
			for _, ch := range buf {
				outByteCh <- ch
			}
		}
	}()

	recorder := NewRecorder(outByteCh, 32*1024*1024)
	go recorder.Run()

	tempDir, err := createTempDir()
	if err != nil {
		log.Fatalf("creating a temp dir: %s", err)
	}
	defer os.RemoveAll(tempDir)

	socketPath := path.Join(tempDir, "control")

	control := ControlServer{SocketPath: socketPath, Recorder: recorder}
	go func() {
		if err := control.ListenAndServe(); err != nil {
			log.Fatalf("control server: %s", err)
		}
	}()

	os.Setenv("KACPERSH_SOCK", socketPath)

	shell := os.Getenv("SHELL")
	if len(shell) == 0 {
		log.Fatalf("Please set the SHELL variable to a supported shell.")
	}
	term := Term{
		Command: exec.Command(shell, "-l"),
		BufSize: 1, // has to be 1 until we implement in-band signaling
	}

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
	tempDir, err := os.MkdirTemp("", tempDirPattern)
	if err != nil {
		return "", err
	}

	return tempDir, nil
}

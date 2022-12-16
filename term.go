package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/term"
)

type Term struct {
	Command *exec.Cmd
	BufSize int
}

func (t *Term) Spawn(output chan<- []byte) error {
	ptmx, err := pty.Start(t.Command)
	if err != nil {
		return fmt.Errorf("starting pty: %w", err)
	}
	defer ptmx.Close()

	winch := make(chan os.Signal, 1)
	signal.Notify(winch, syscall.SIGWINCH)
	defer func() { signal.Stop(winch); close(winch) }()
	winch <- syscall.SIGWINCH // initial resize
	// TODO: handle sigwinch in select

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("making the terminal raw: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// we don't need stdin where we're going, so we just copy it straight to the pty
	go io.Copy(ptmx, os.Stdin)

	pumpErrors := pump(output, ptmx, t.BufSize)

	for {
		select {
		case err := <-pumpErrors:
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("stdout pump error: %w", err)
		case <-winch:
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				return fmt.Errorf("error resizing pty: %w", err)
			}
		}
	}
}

func pump(ch chan<- []byte, r io.Reader, bufSize int) <-chan error {
	errCh := make(chan error, 1)
	if bufSize == 0 {
		errCh <- fmt.Errorf("buffer size too small (0)")
		return errCh
	}

	go func() {
		for {
			// creating new buf every time so we don't have to wait for old buffer to be processed
			buf := make([]byte, bufSize)
			_, err := r.Read(buf)
			if err != nil {
				close(ch)
				errCh <- err
				return
			}
			ch <- buf
		}
	}()

	return errCh
}

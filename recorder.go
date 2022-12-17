package main

import "log"

type ControlMessage int

const (
	Unknown        ControlMessage = 0
	StartRecording ControlMessage = 1
	StopRecording  ControlMessage = 2
)

type Recorder struct {
	control chan ControlMessage
	sync    chan []byte

	stdout <-chan byte

	recording bool
	current   []byte
}

func NewRecorder(stdout <-chan byte) *Recorder {
	return &Recorder{
		control: make(chan ControlMessage),
		sync:    make(chan []byte),

		stdout: stdout,

		recording: false,
		current:   make([]byte, 0),
	}
}

func (r *Recorder) Start() {
	r.control <- StartRecording
	<-r.sync
}

func (r *Recorder) Stop() []byte {
	r.control <- StopRecording
	return <-r.sync
}

func (r *Recorder) Run() {
	for {
		select {
		case ch := <-r.stdout:
			if r.recording {
				r.current = append(r.current, ch)
			}
		case cmd := <-r.control:
			switch cmd {
			case StartRecording:
				log.Println("starting recording")
				r.recording = true
				r.sync <- nil // let the shell know it can continue
			case StopRecording:
				r.recording = false

				// make sure r.last has the correct size for copying
				// this is fine, always, as they have the same underlying capacity
				last := make([]byte, len(r.current))

				// save last command's output before cleaning the current buffer
				copy(last, r.current)

				// clean the buffer for the new recording
				r.current = r.current[:0]
				log.Println("stopping recording, got", len(last), "bytes")
				r.sync <- last
			}
		}
	}
}

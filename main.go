package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

const (
	DAEMON_STARTED = iota
	DAEMON_OK
	DAEMON_NG
)

var mount *string

func main() {
	child := flag.Bool("child", false, "")
	mount = flag.String("mount", "", "")
	flag.Parse()

	if *child {
		childMain()
	} else {
		if err := parentMain(); err != nil {
			log.Fatal(err)
		}
	}
}

func childCommand(w *os.File) *exec.Cmd {
	args := []string{"--child"}
	cmd := exec.Command(os.Args[0], args...)
	cmd.ExtraFiles = []*os.File{w}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func parentMain() (err error) {
	// create pipe
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}

	// create command
	cmd := childCommand(w)

	// start command
	if err := cmd.Start(); err != nil {
		return err
	}

	status := waitChild(r)

	// return result
	if status == DAEMON_OK {
		return nil
	} else {
		return fmt.Errorf("Failed to start child")
	}
}

func waitChild(r *os.File) int {
	// async: read child status
	var status int = DAEMON_STARTED
	go func() {
		buf := make([]byte, 1)
		r.Read(buf)
		status = int(buf[0])
	}()

	// wait child
	i := 0
	for i < 10 {
		if status != DAEMON_STARTED {
			break
		}
		time.Sleep(500 * time.Millisecond)
		i += 1
	}

	return status
}

func childMain() {
	// notify its status
	pipe := os.NewFile(uintptr(3), "pipe")
	if pipe != nil {
		defer pipe.Close()
		pipe.Write([]byte{DAEMON_OK})
	}

	// new session
	signal.Ignore(syscall.SIGCHLD)
	syscall.Close(0)
	syscall.Close(1)
	syscall.Close(2)
	syscall.Setsid()
	syscall.Umask(022)
	syscall.Chdir("/")

	for {
		time.Sleep(1000 * time.Millisecond)
	}
}

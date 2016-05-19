package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
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
	DAEMON_BAD
)

var mount *string

func main() {
	child := flag.Bool("child", false, "")
	mount = flag.String("mount", "", "")
	flag.Parse()
	defer glog.Flush()

	if len(flag.Args()) != 1 {
		log.Fatalf("%v [opts...] <path/to/config.yml>", os.Args[0])
	}

	if *child {
		glog.Info("childMain()")
		if err := childMain(); err != nil {
			log.Fatal(err)
		}
	} else {
		glog.Info("parentMain()")
		if err := parentMain(); err != nil {
			log.Fatal(err)
		}
	}
}

func childCommand(mnt string, w *os.File) *exec.Cmd {
	args := []string{"--child", "--mount", mnt}
	args = append(args, os.Args[1:]...)
	cmd := exec.Command(os.Args[0], args...)
	cmd.ExtraFiles = []*os.File{w}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func parentMain() error {
	// create tmpdir
	mnt, err := ioutil.TempDir("", "docker-memoize")
	if err != nil {
		return err
	}

	// create pipe
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}

	// create command
	cmd := childCommand(mnt, w)

	// start command
	if err := cmd.Start(); err != nil {
		return err
	}

	status := waitChild(r)

	// return result
	if status == DAEMON_OK {
		fmt.Println(mnt)
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

func childMain() error {
	// notify its status
	pipe := os.NewFile(uintptr(3), "pipe")

	// read config
	yaml, err := ioutil.ReadFile(flag.Args()[0])
	if err != nil {
		// notify BAD
		if pipe != nil {
			defer pipe.Close()
			pipe.Write([]byte{DAEMON_BAD})
		}
		return err
	}
	conf := Parse(string(yaml))

	// new session
	syscall.Close(0) // stdout
	syscall.Close(1) // stdin
	// syscall.Close(2) // stderr
	syscall.Setsid()
	syscall.Umask(022)
	syscall.Chdir("/")

	// mount fs
	server, err := MountFileSystem(conf, *mount)
	if err != nil {
		return err
	}

	// unmount filesystem
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	// async: wait signal
	go func() {
		<-sigchan
		glog.Info("server.Unmount()")
		server.Unmount()
	}()

	// notify OK
	if pipe != nil {
		defer pipe.Close()
		pipe.Write([]byte{DAEMON_OK})
	}
	server.Serve()

	// terminate
	signal.Stop(sigchan)
	glog.Flush()

	return nil
}

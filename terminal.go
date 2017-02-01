package termtest

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// Terminal holds a connection to a running tmux instance.
type Terminal struct {
	dir        string
	tmuxSocket string
}

// New starts a new tmux instance running cmd with the given dimensions.
func New() (*Terminal, error) {
	tmpdir, err := ioutil.TempDir("", "term-ui-test-")
	if err != nil {
		return nil, err
	}

	term := &Terminal{
		dir:        tmpdir,
		tmuxSocket: filepath.Join(tmpdir, "tmux-socket"),
	}

	// start tmux on the given socket, then create a new detached session with
	// the status bar disabled
	_, err = term.tmux("new-session", "-d", ";", "set", "-g", "status", "off")
	if err != nil {
		return nil, err
	}

	return term, nil
}

// tmux runs tmux with arguments and returns stdout.
func (term *Terminal) tmux(args ...string) ([]byte, error) {
	defaultArgs := []string{"-f", "/dev/null", "-S", term.tmuxSocket}
	args = append(defaultArgs, args...)

	cmd := exec.Command("tmux", args...)
	cmd.Stderr = os.Stderr
	return cmd.Output()
}

// Run runs the given command and returns a screen capture at the end of the
// program, just after exit. The terminal is started with the given dimensions.
func (term *Terminal) Run(x, y int, command string) ([]byte, error) {
	name := fmt.Sprintf("test-%d", rand.Int63())

	// start new session with random name with command and given dimensions
	sx := strconv.FormatInt(int64(x), 10)
	sy := strconv.FormatInt(int64(y), 10)

	// add tmux commands to run after the command we've been given
	command = fmt.Sprintf(`%s ; tmux capture-pane -b %s ; tmux wait-for -S %s`, command, name, name)

	_, err := term.tmux("new-session", "-d", "-x", sx, "-y", sy, "-s", name, command)
	if err != nil {
		return nil, err
	}

	// wait for buffer
	_, err = term.tmux("wait-for", name)
	if err != nil {
		return nil, err
	}

	// read buffer
	buf, err := term.tmux("show-buffer", "-b", name)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// Exit kills the tmux instance and all remaining processes running in it.
func (term *Terminal) Exit() error {
	_, err := term.tmux("kill-server")
	if err != nil {
		return err
	}

	return os.RemoveAll(term.dir)
}

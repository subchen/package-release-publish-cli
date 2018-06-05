package sh

import (
	"os"
	"os/exec"

	"github.com/subchen/storm/pkg/config"
)

type Command struct {
	Command string
	Args    []string
	Shell   bool
	PipeOut bool
	Dir     string
	Env     []string
}

func (c *Command) cmd() *exec.Cmd {
	var cmd *exec.Cmd
	if c.Shell {
		cmd = exec.Command(getShell(), c.Command)
	} else {
		cmd = exec.Command(c.Command, c.Args...)
	}

	if c.PipeOut {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	cmd.Env = os.Environ()

	if c.Dir != "" {
		cmd.Dir = c.Dir
		cmd.Env = append(cmd.Env, []string{
			"PWD=" + c.Dir,
		}...)
	}

	if len(c.Env) > 0 {
		cmd.Env = append(cmd.Env, c.Env...)
	}

	return cmd
}

func (c *Command) Run() error {
	return c.cmd().Run()
}

func (c *Command) RunOut() (string, error) {
	out, err := c.cmd().CombinedOutput()
	return string(out), err
}

func RunCmd(command string, args ...string) {
	RunCmdAt(config.ProjectWorkdir, command, args...)
}

func RunCmdAt(dir string, command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Dir = dir

	cmd.Env = append(os.Environ(), []string{
		"GOPATH=" + config.ProjectGopath,
		"PWD=" + dir,
	}...)

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func getShell() string {
	return "bash"
}

package util

import (
	"os"
	"os/exec"
	"strings"
	"errors"
	"github.com/jiusanzhou/tentacle/log"
)

func DoCmd(cmds string) error {
	args := strings.Split(cmds, " ")

	log.Debug("Commands: ", args)

	if len(args) < 1 {
		return errors.New("Command errors, no enough args.")
	}

	e := exec.Command(args[0], args[1:]...)
	e.Stdout = os.Stdout
	e.Stderr = os.Stderr
	e.Run()

	return nil
}

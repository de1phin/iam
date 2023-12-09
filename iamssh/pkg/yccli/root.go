package yc

import (
	"errors"
	"os/exec"
)

func ycExecute(args ...string) ([]byte, error) {
	command := exec.Command("yc", args...)
	out, err := command.Output()
	if err != nil {
		if exitErr := err.(*exec.ExitError); exitErr != nil {
			return nil, errors.New(string(exitErr.Stderr))
		}
		return nil, err
	}

	return out, nil
}

package cmd

import (
	"bytes"
	"errors"
	"os/exec"
)

// 在Shell执行命令
func ShellExec(cmd string) (*bytes.Buffer, error) {
	command := exec.Command("sh")
	in := bytes.NewBuffer(nil)
	out := bytes.NewBuffer(nil)
	errOut := bytes.NewBuffer(nil)
	command.Stdin = in
	command.Stdout = out
	command.Stderr = errOut
	in.WriteString(cmd)
	in.WriteString("\n")
	in.WriteString("exit\n")
	if err := command.Run(); err != nil {
		return nil, errors.New(errOut.String())
	}
	return out, nil
}

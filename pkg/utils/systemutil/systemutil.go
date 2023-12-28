package systemutil

import (
	"github.com/sirupsen/logrus"
	"os"
	"syscall"
)

func Restart() {
	exe, err := os.Executable()
	if err != nil {
		logrus.Fatal(err)
	}
	args := os.Args
	env := os.Environ()
	err = syscall.Exec(exe, args, env)
	if err != nil {
		logrus.Fatal(err)
	}
}

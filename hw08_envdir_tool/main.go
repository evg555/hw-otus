package main

import (
	"fmt"
	"os"
)

const ErrWrongParams = "Usage: ./<go-envdir> <path> <command> <arg1> <arg2>"

func main() {
	if len(os.Args) < 2 {
		fmt.Println(ErrWrongParams)
		os.Exit(1)
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(RunCmd(os.Args[2:], env))
}

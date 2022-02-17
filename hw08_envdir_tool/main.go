package main

import (
	"fmt"
	"os"
)

func main() {
	// 1 аргумент - сам скрипт
	// 2 аргумент - каталог
	// 3 аргумент - команда
	// 4... параметры к команде
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "not enough arguments")
		os.Exit(111)
	}
	dir := os.Args[1]

	readDir, err := ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(111)
	}
	os.Exit(RunCmd(os.Args[2:], readDir))
}

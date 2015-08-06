package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dchest/uniuri"
)

func main() {
	// custom pattern matching!
	args := os.Args[1:]
	var path string
	if len(args) != 0 {
		path = args[len(args)-1]
		if len(args) == 1 {
			args = []string{}
		} else {
			args = args[:len(args)-1]
		}
	}

	// i copied half of the go tool for that function
	resolved := importPaths([]string{path})

	// create a new tmpdir
	td, err := ioutil.TempDir("", "multicov_")
	if err != nil {
		log.Fatal(err)
	}

	// Result file
	result := []byte("mode: count\n")

	// Test each package
	for _, pkg := range resolved {
		cp := filepath.Join(td, uniuri.New())
		cmd := exec.Command(
			"go",
			append(append([]string{
				"test",
				"-covermode=count",
				"-coverprofile=" + cp,
			}, args...), pkg)...,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		if _, err := os.Stat(cp); err == nil {
			report, err := ioutil.ReadFile(cp)
			if err != nil {
				log.Fatal(err)
			}

			result = append(result, report[12:]...)
		}

	}

	// Write into coverage.out
	if err := ioutil.WriteFile("coverage.out", result, 0644); err != nil {
		log.Fatal(err)
	}
}

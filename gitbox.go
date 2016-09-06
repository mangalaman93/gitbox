package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mangalaman93/gitbox/box"
	"net/url"
)

func main() {
	wd, errWD := os.Getwd()
	if errWD != nil {
		log.Fatalf("Unable to find current WD :: %v", errWD)
	}
	baseWD := filepath.Base(wd)

	// We intercept following commands -
	//  - push, pull, fetch, clone
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "push":
			if box.IsBoxRepoCmd(baseWD) {
				box.Push(baseWD)
				return
			}
		case "pull":
			if box.IsBoxRepoCmd(baseWD) {
				box.Pull(baseWD)
				return
			}
		case "clone":
			if len(os.Args) > 2 {
				rawURL := os.Args[2]
				boxURL, errURL := url.Parse(rawURL)
				if errURL != nil {
					log.Fatalf("Invalid url %s :: %v", rawURL, errWD)
				}

				if boxURL.Scheme == "box" {
					box.Clone(baseWD)
					return
				}
			}
		}
	}

	// If the command is not intercepted, we continue and run a git command.
	out, err := exec.Command("git", os.Args[1:]...).CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}

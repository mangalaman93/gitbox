package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mangalaman93/gitbox/box"
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
			isBoxRepo, errBox := box.IsBoxRepo(baseWD)
			if errBox != nil {
				log.Fatalln("Unable to check remote url :: ", errBox)
			}
			if isBoxRepo {
				errBox := box.Push(baseWD)
				if errBox != nil {
					log.Fatalln("error in pushing :: ", errBox)
				}
				return
			}
		case "pull":
			isBoxRepo, errBox := box.IsBoxRepo(baseWD)
			if errBox != nil {
				log.Fatalln("Unable to check remote url :: ", errBox)
			}
			if isBoxRepo {
				errBox := box.Pull(baseWD)
				if errBox != nil {
					log.Fatalln("error in pulling :: ", errBox)
				}
				return
			}
		case "clone":
			if len(os.Args) > 2 {
				isBoxURL, errBox := box.IsBoxRemoteURL(os.Args[2])
				if errBox != nil {
					log.Fatalln("Invalid url to clone :: ", errBox)
				}
				if isBoxURL {
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

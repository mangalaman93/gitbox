package box

import (
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/libgit2/git2go"
)

const (
	cGitOrigin = "origin"
	cBoxScheme = "box"
)

func IsBoxRemoteURL(rawURL string) (bool, error) {
	boxURL, errURL := url.Parse(rawURL)
	if errURL != nil {
		return false, errURL
	}

	return (strings.ToLower(boxURL.Scheme) == cBoxScheme), nil
}

func IsBoxRepo(cw string) (bool, error) {
	gitRepo, errRepo := git.OpenRepository(cw)
	if errRepo != nil {
		return false, errRepo
	}

	originRemote, errRemote := gitRepo.Remotes.Lookup(cGitOrigin)
	if errRemote != nil {
		return false, errRemote
	}

	return IsBoxRemoteURL(originRemote.Url())
}

func Push(cw string) error {
	// Four step process
	//  - create temp bare repo
	//  - push to bare repo
	//  - sync to bare repo
	//  - sync files with box

	// First create a bare repo
	tempDirForBareRepo, errTemp := ioutil.TempDir(os.TempDir(), "")
	if errTemp != nil {
		return errTemp
	}
	defer os.RemoveAll(tempDirForBareRepo)

	_, errInit := git.InitRepository(tempDirForBareRepo, true)
	if errInit != nil {
		return errInit
	}

	// Second, push to bare repo from current repo
	gitRepo, errRepo := git.OpenRepository(cw)
	if errRepo != nil {
		return errRepo
	}

	bareRemoteOnGitRepo, errRemote := gitRepo.Remotes.Create("temp", tempDirForBareRepo)
	if errRemote != nil {
		return errRemote
	}
	defer gitRepo.Remotes.Delete(bareRemoteOnGitRepo.Name())

	// TODO: pick the current branch.
	bareRemoteOnGitRepo.Push([]string{"refs/heads/master"}, &git.PushOptions{})

	// Third, call sync here with box for bare repo
	// tempDirForBareRepo with (get box url from repo itself)

	// Fourth, sync files
	tempDirForFiles, errDir := ioutil.TempDir(os.TempDir(), cw)
	if errDir != nil {
		return errDir
	}
	defer os.RemoveAll(tempDirForFiles)

	_, errRepo = git.Clone(tempDirForBareRepo, tempDirForFiles, &git.CloneOptions{})
	if errRepo != nil {
		return errRepo
	}

	os.RemoveAll(filepath.Join(tempDirForFiles, ".git"))
	// sync files tempDirForFiles with box url

	return nil
}

func Pull(cw string) error {
	// Following steps -
	// - sync bare repo
	// - pull from bare repo
	tempDirForBareRepo, errDir := ioutil.TempDir(os.TempDir(), cw)
	if errDir != nil {
		return errDir
	}
	defer os.RemoveAll(tempDirForBareRepo)

	// sync bare repo from tempDirForBareRepo to box URL

	gitRepo, errRepo := git.OpenRepository(cw)
	if errRepo != nil {
		return errRepo
	}

	tempRemoteOnGitRepo, errRemote := gitRepo.Remotes.Create("temp", tempDirForBareRepo)
	if errRemote != nil {
		return errRemote
	}
	defer gitRepo.Remotes.Delete(tempRemoteOnGitRepo.Name())

	errFetch := tempRemoteOnGitRepo.Fetch([]string{"refs/heads/master"},
		&git.PushOptions{}, "")
	if errFetch != nil {
		return errFetch
	}

	remoteRef, errRef := gitRepo.References.Lookup("refs/remotes/temp/master")
	if errRef != nil {
		return errRef
	}

	mergeRemoteHead, err := gitRepo.AnnotatedCommitFromRef(remoteRef)
	if err != nil {
		return err
	}

	mergeHeads := make([]*git.AnnotatedCommit, 1)
	mergeHeads[0] = mergeRemoteHead
	errMerge := gitRepo.Merge(mergeHeads, nil, nil)
	if errMerge != nil {
		return errMerge
	}

	return nil
}

func Clone(cloneURL string) error {
	boxURL, errURL := url.Parse(cloneURL)
	if errURL != nil {
		return errURL
	}

	tempDirForBareRepo, errDir := ioutil.TempDir(os.TempDir(), "")
	if errDir != nil {
		return errDir
	}
	defer os.RemoveAll(tempDirForBareRepo)

	// sync with box

	cw := filepath.Base(boxURL.Path)
	gitRepo, errRepo := git.Clone(tempDirForBareRepo, cw, &git.CloneOptions{})
	if errRepo != nil {
		return errRepo
	}
	gitRepo.Remotes.SetPushUrl(cGitOrigin, cloneURL)

	return nil
}

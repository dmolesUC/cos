// +build mage

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	appName = "cos"
	cmdPkg = "github.com/dmolesUC3/cos/cmd"
)

var gocmd = mg.GoCmd()

var oses = []string{"darwin", "linux", "windows"}
var architectures = []string{"amd64"}

// ------------------------------------------------------------
// Targets

// Install installs cos in $GOPATH/bin.
func Install() error {
	binName := appName
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	gopath, err := sh.Output(gocmd, "env", "GOPATH")
	if err != nil {
		return fmt.Errorf("error determining GOPATH: %v", err)
	}
	binDir := filepath.Join(gopath, "bin")
	binPath := filepath.Join(binDir, binName)

	flags, err := ldFlags()
	if err != nil {
		return fmt.Errorf("error determining ldflags: %v", err)
	}
	return sh.RunV(gocmd, "build", "-o", binPath, "-ldflags", flags)
}

// Build builds a cos binary for the current platform.
func Build() error {
	binName := appName
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	flags, err := ldFlags()
	if err != nil {
		return fmt.Errorf("error determining ldflags: %v", err)
	}
	return sh.RunV(gocmd, "build", "-ldflags", flags)
}

// BuildLinux builds a linux-amd64 binary (the most common cross-compile case)
func BuildLinux() error {
	return buildFor("linux", "amd64")
}

// BuildAll builds a cos binary for each target platform.
func BuildAll() error {
	for _, os_ := range oses {
		for _, arch := range architectures {
			err := buildFor(os_, arch)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Platforms lists target platforms for buildAll.
func Platforms() {
	for _, os_ := range oses {
		for _, arch := range architectures {
			fmt.Printf("%s-%s\n", os_, arch)
		}
	}
}

// Clean removes compiled binaries from the current working directory.
func Clean() error {
	var binRe = regexp.MustCompile("^" + appName + "(-[a-zA-Z0-9]+-[a-zA-Z0-9]+)?(.exe)?$")

	rmcmd := "rm"
	if runtime.GOOS == "windows" {
		rmcmd = "del"
	}

	files, err := ioutil.ReadDir("./")
	if err != nil {
		return err
	}

	for _, f := range files {
		mode := f.Mode()
		isPlainFile := mode.IsRegular() && mode&os.ModeSymlink == 0
		isExecutable := mode&0111 != 0
		if isPlainFile && isExecutable {
			name := f.Name()
			if binRe.MatchString(name) {
				err := sh.RunV(rmcmd, name)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// ------------------------------------------------------------
// Helper functions

func ldFlags() (string, error) {
	commitHash, err := sh.Output("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return "", err
	}
	tag, err := sh.Output("git", "describe", "--tags")
	if err != nil {
		return "", err
	}
	timestamp := time.Now().Format(time.RFC3339)

	flagVals := map[string]string{
		"commitHash": commitHash,
		"tag":        tag,
		"timestamp":  timestamp,
	}

	var flags []string
	for k, v := range flagVals {
		flag := fmt.Sprintf("-X %s.%s=%s", cmdPkg, k, v)
		flags = append(flags, flag)
	}
	return strings.Join(flags, " "), nil
}

func binNameFor(os_ string, arch string) string {
	binName := appName
	binName = fmt.Sprintf("%s-%s-%s", binName, os_, arch)
	if os_ == "windows" {
		binName += ".exe"
	}
	return binName
}

func buildFor(os_, arch string) error {
	binName := binNameFor(os_, arch)

	flags, err := ldFlags()
	if err != nil {
		return fmt.Errorf("error determining ldflags: %v", err)
	}

	env := map[string]string{
		"GOOS":   os_,
		"GOARCH": arch,
	}
	return sh.RunWith(env, gocmd, "build", "-o", binName, "-ldflags", flags)
}


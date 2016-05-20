// +build task

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
)

func sourceDir() (dir string, ok bool) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "", ok
	}
	dir = filepath.Dir(file)

	// if the last segment is "task", strip that out; allows
	// segregating task files in a subdir
	parent, taskDir := filepath.Split(dir)
	if taskDir == "task" {
		dir = parent
	}

	return dir, ok
}

func goBuild(src string, action string, args ...string) error {
	cmd := exec.Command(
		"go", action,
		"-v",
	)
	cmd.Args = append(cmd.Args, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func gitVersion(dir string) (string, error) {
	const prefix = "release/"
	cmd := exec.Command(
		"git", "describe",
		"--match", prefix+"*",
		"--dirty=+edited",
	)
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	buf, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git describe: %v", err)
	}
	if !bytes.HasPrefix(buf, []byte(prefix)) {
		return "", fmt.Errorf("version tag does not begin with prefix: %q", buf)
	}
	buf = buf[len(prefix):]
	buf = bytes.TrimSuffix(buf, []byte("\n"))
	return string(buf), nil
}

// containingPackage returns the package that contains the given
// command. If the command import path is `x/cmd/foo`, returns `x`;
// otherwise, returns the command itself.
func containingPackage(cmdImportPath string) string {
	p := path.Dir(cmdImportPath)
	for {
		tmp := path.Base(p)
		p = path.Dir(p)
		if tmp == "cmd" {
			return p
		}
		if p == "." {
			break
		}
	}
	return cmdImportPath
}

var prog = filepath.Base(os.Args[0])

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", prog)
	fmt.Fprintf(os.Stderr, "  %s IMPORT_PATH\n", prog)
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("")

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	command := flag.Arg(0)
	parent := containingPackage(command)

	log.Printf("# building %v", command)

	os.Setenv("CGO_ENABLED", "0")
	os.Setenv("GOOS", "linux")

	src, ok := sourceDir()
	if !ok {
		log.Fatal("cannot determine source directory")
	}

	version, err := gitVersion(src)
	if err != nil {
		log.Fatalf("error extracting version: %v", err)
	}

	err = goBuild(src, "build", "-v",
		"-i",
		"-tags", "netgo",
		"-ldflags", "-X "+parent+"/internal/version.Version="+version+" -w",
		command,
	)
	if err != nil {
		log.Fatalf("go build: %v", err)
	}
}

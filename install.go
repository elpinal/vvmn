package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/pkg/errors"
)

var cmdInstall = &Command{
	Run:       runInstall,
	UsageLine: "install version",
	Short:     "Install Vim",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdInstall.Flag.BoolVar(&flagA, "a", false, "")
}

// runInstall executes install command and return exit code.
func runInstall(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "vvmn install: no vim version specified")
		return 1
	}
	dir := path.Join(VvmnDir, "repo")
	if _, err := os.Stat(dir); err != nil {
		_, err := exec.Command("git", "clone", "--bare", RepoURL, dir).CombinedOutput()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}

	cmd := exec.Command("git", "fetch", "--tags")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed fetch"))
		return 1
	}

	version := args[0]
	if version == "latest" {
		var err error
		version, err = latestTag()
		if err != nil {
			fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed get the latest version"))
			return 1
		}
	}
	cmd = exec.Command("git", "archive", "--prefix="+version+"/", version)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed git archive"))
		fmt.Fprintln(os.Stderr, stderr.String())
		return 1
	}

	srcDir := path.Join(VvmnDir, "src")
	if _, err := os.Stat(srcDir); err != nil {
		if err := os.MkdirAll(srcDir, 0777); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}
	cmd = exec.Command("tar", "xf", "-")
	cmd.Dir = srcDir
	cmd.Stdin = bytes.NewReader(out)
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed tar"))
		return 1
	}

	options := []string{"--prefix="+path.Join(VvmnDir, "vims/"+version)}
	if len(args) > 1 {
		options = append(options, args[1:]...)
	}
	cmd = exec.Command("./configure", options...)
	cmd.Dir = path.Join(srcDir, version)
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed configure"))
		return 1
	}

	cmd = exec.Command("make", "all", "install")
	cmd.Dir = path.Join(srcDir, version)
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed install"))
		return 1
	}
	return 0
}

func latestTag() (string, error) {
	cmd := exec.Command("git", "rev-list", "--tags", "--max-count=1")
	cmd.Dir = path.Join(VvmnDir, "repo")
	out, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "failed git rev-list")
	}
	sha := string(bytes.TrimSuffix(out, []byte("\n")))
	cmd = exec.Command("git", "describe", "--tags", sha)
	cmd.Dir = path.Join(VvmnDir, "repo")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	tag, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, stderr.String())
	}
	return string(bytes.TrimSuffix(tag, []byte("\n"))), nil
}

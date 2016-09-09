package vvmn

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

// doubleError is a type which has two error.
type doubleError struct {
	a, b error
}

func (e *doubleError) Error() string {
	if e == nil {
		return ""
	}
	if e.a == nil {
		return e.b.Error()
	}
	if e.b == nil {
		return e.a.Error()
	}
	return fmt.Sprintf("%v\n%v", e.a, e.b)
}

// build builds the specified Vim version.
func build(version string, options []string) *doubleError {
	options = append([]string{"--prefix=" + filepath.Join(vvmnrootVim, version)}, options...)
	cmd := exec.Command("./configure", options...)
	cmd.Dir = filepath.Join(vvmnrootSrc, version)
	var buf bytes.Buffer
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return &doubleError{errors.Wrap(err, "./configure failed"), fmt.Errorf(buf.String())}
	}

	cmd = exec.Command("make", "all", "install")
	cmd.Dir = filepath.Join(vvmnrootSrc, version)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return &doubleError{errors.Wrap(err, "make all install failed"), fmt.Errorf(stderr.String())}
	}

	return nil
}

// Install installs Vim version.
func Install(version string, options ...string) error {
	if exist(filepath.Join(vvmnrootVim, version, "bin", "vim")) {
		return nil
	}
	if err := build(version, options); err != nil {
		return err
	}
	return nil
}

// checkout checkouts the specified version of the Vim repository.
func checkout(version string) *doubleError {
	cmd := exec.Command("git", "archive", "--prefix="+version+"/", version)
	cmd.Dir = vvmnrootRepo
	var buf bytes.Buffer
	cmd.Stderr = &buf
	out, err := cmd.Output()
	if err != nil {
		return &doubleError{errors.Wrap(err, "git archive failed"), fmt.Errorf("%s", buf.String())}
	}

	if !exist(vvmnrootSrc) {
		if err := os.MkdirAll(vvmnrootSrc, 0777); err != nil {
			return &doubleError{err, nil}
		}
	}

	cmd = exec.Command("tar", "xf", "-")
	cmd.Dir = vvmnrootSrc
	cmd.Stdin = bytes.NewReader(out)
	output, err := cmd.CombinedOutput()
	if err != nil {
		_ = os.RemoveAll(filepath.Join(vvmnrootSrc, version))
		return &doubleError{errors.Wrap(err, "tar xf - failed"), fmt.Errorf("%s", output)}
	}

	return nil
}

// update updates the Vim repository.
func update() *doubleError {
	cmd := exec.Command("git", "fetch")
	cmd.Dir = vvmnrootRepo
	if out, err := cmd.CombinedOutput(); err != nil {
		return &doubleError{errors.Wrap(err, "failed to fetch"), fmt.Errorf("%s", out)}
	}
	return nil
}

// mirror mirrors the Vim repository.
func mirror() *doubleError {
	out, err := exec.Command("git", "clone", "--mirror", RepoURL, vvmnrootRepo).CombinedOutput()
	if err != nil {
		return &doubleError{errors.Wrap(err, "cloning repository failed"), fmt.Errorf("%s", out)}
	}
	return nil
}

// download fetches the Vim repository.
func download() error {
	if !exist(vvmnrootRepo) {
		if err := mirror(); err != nil {
			return err
		}
	}
	if err := update(); err != nil {
		return err
	}
	return nil
}

// Download fetches the Vim repository and check out version.
func Download(version string) error {
	if exist(filepath.Join(vvmnrootSrc, version)) {
		return nil
	}
	if err := download(); err != nil {
		return err
	}
	if err := checkout(version); err != nil {
		return err
	}
	return nil
}

// Get downloads and installs Vim.
func Get(version string, options ...string) error {
	if err := Download(version); err != nil {
		return err
	}
	if err := Install(version, options...); err != nil {
		return err
	}
	return nil
}

// latestTag reports the latest tag of the Vim repository.
func latestTag() (string, error) {
	cmd := exec.Command("git", "rev-list", "--tags", "--max-count=1")
	cmd.Dir = vvmnrootRepo
	out, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "git rev-list failed")
	}
	sha := string(bytes.TrimSuffix(out, []byte("\n")))
	cmd = exec.Command("git", "describe", "--tags", sha)
	cmd.Dir = vvmnrootRepo
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	tag, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, stderr.String())
	}
	return string(bytes.TrimSuffix(tag, []byte("\n"))), nil
}

// LatestTag downloads the updated Vim repository and reports
// the latest tag of it.
func LatestTag() (string, error) {
	if err := download(); err != nil {
		return "", err
	}
	return latestTag()
}

package vvmn

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// A Info represents various Vim versions.
type Info struct {
	Current    string
	Installed  []string
	Downloaded []string
}

// List returns information of Vim versions.
func List() Info {
	if !exist(vvmnrootVim) {
		return Info{}
	}
	current, _ := os.Readlink(filepath.Join(vvmnrootVim, "current"))
	currentVersion := filepath.Base(current)
	if !exist(current) {
		currentVersion = ""
	}
	versions, err := ioutil.ReadDir(vvmnrootVim)
	if err != nil {
		return Info{}
	}
	var installed, downloaded []string
	for _, version := range versions {
		ver := version.Name()
		if ver == "current" || ver == currentVersion {
			continue
		}
		if exist(filepath.Join(vvmnrootVim, version.Name(), "bin", "vim")) {
			installed = append(installed, ver)
		} else {
			downloaded = append(downloaded, ver)
		}
	}
	return Info{Current: currentVersion, Installed: installed, Downloaded: downloaded}
}

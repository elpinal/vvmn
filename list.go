package vvmn

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// A Info represents Vim name and states.
type Info struct {
	Name      string
	Current   bool
	Installed bool
}

// List returns information of installed Vim versions.
func List() []Info {
	if !exist(vvmnrootVim) {
		return nil
	}
	current, _ := os.Readlink(filepath.Join(vvmnrootVim, "current"))
	currentVersion := filepath.Base(current)
	versions, err := ioutil.ReadDir(vvmnrootVim)
	if err != nil {
		return nil
	}
	var info []Info
	for _, version := range versions {
		ver := version.Name()
		if ver == "current" {
			continue
		}
		var installed bool
		if exist(filepath.Join(vvmnrootVim, version.Name(), "bin", "vim")) {
			installed = true
		}
		info = append(info, Info{Name: ver, Current: ver == currentVersion, Installed: installed})
	}
	return info
}

package vvmn

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func mustRemoveAll(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	flag.Parse()

	tempdir, err := ioutil.TempDir("", "")
	SetRoot(tempdir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "TempDir: %v", err)
		os.Exit(2)
	}

	if !exist("testdata/vim") {
		if testing.Verbose() {
			log.SetFlags(log.Lshortfile)
			log.Print("fetching for test...")
		}
		out, err := exec.Command("git", "clone", "--depth=1", "--bare", "--branch=v7.4.2222", RepoURL, "testdata/vim").CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to fetch for test: %v\n", err)
			fmt.Fprintln(os.Stderr, string(out))
			os.Exit(2)
		}
	}
	RepoURL = "testdata/vim"

	r := m.Run()

	mustRemoveAll(tempdir)
	os.Exit(r)
}

func TestGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	if err := Get("v7.4.2222"); err != nil {
		t.Fatalf("Get: %v", err)
	}
	bin := filepath.Join(vvmnrootVim, "v7.4.2222", "bin", "vim")
	out, err := exec.Command(bin, "--version").Output()
	if err != nil {
		t.Fatalf("vim version: %v", err)
	}
	if !bytes.Contains(out, []byte("2222")) {
		t.Fatalf("%v must contain 2222", out)
	}
}

func TestDownload(t *testing.T) {
	if !testing.Short() {
		t.Skip("skipping in non-short mode")
	}
	if err := Download("v7.4.2222"); err != nil {
		t.Fatalf(`Download("v7.4.2222") failed: %v`, err)
	}
}

func TestListAndRun(t *testing.T) {
	if err := Get("v7.4.2222"); err != nil {
		t.Fatalf("GetBinary: %v", err)
	}
	if list := List(); len(list) == 0 {
		t.Error("could not list Vim versions")
	} else if list[0].Name != "v7.4.2222" {
		t.Error("could not find v7.4.2222")
	}
	if err := Run("v7.4.2222", "--version"); err != nil {
		t.Fatalf("vim version: %v", err)
	}
}

func TestLatestTag(t *testing.T) {
	if out, err := LatestTag(); err != nil {
		t.Fatalf("LatestTag: %v", err)
	} else if out != "v7.4.2222" {
		t.Fatalf("vimt %v, want v7.4.2222", out)
	}
}

func TestUse(t *testing.T) {
	if err := Use("v7.4.2222"); err != nil {
		t.Fatalf("Use: %v", err)
	}
}

func TestRemove(t *testing.T) {
	if err := Remove("v7.4.2222"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
}

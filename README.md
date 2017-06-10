# Vvmn

Vim Version Manager Next

## Install

To install, use `go get`:

```bash
$ go get github.com/elpinal/vvmn/cmd/vvmn
```

## Command vvmn

Vvmn is a tool for managing Vim versions.

Usage:

```bash
vvmn command [arguments]
```

The commands are:

```bash
get         download and install Vim
list        list installed Vim versions
use         select a Vim version to use
remove      remove Vim versions
run         execute the specified Vim version
version     print vvmn version
```

Use "vvmn help [command]" for more information about a command.

### Download and install Vim

Get specific Vim versions:

```bash
$ vvmn get v7.4.2222
```

Get the latest tagged Vim, such as v7.4.2347:

```bash
$ vvmn get latest
```

### List installed Vim versions

Know what Vim versions is installed or downloaded:

```bash
$ vvmn list
Current:

	v7.4.2347

Installed:

	v7.4.1000
	v7.4.2222
```

### Select a Vim version to use

Select a Vim version to use:

```bash
$ vvmn use v7.4.2222
```

### Remove Vim versions

Say goodbye to particular Vim versions:

```bash
$ vvmn remove v7.4.2222
```

### Execute the specified Vim version

Execute another Vim version:

```bash
$ vvmn run v7.4.1000
```

### Print vvmn version

```bash
$ vvmn version
```

## Contribution

1. Fork ([https://github.com/elpinal/vvmn/fork](https://github.com/elpinal/vvmn/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[elpinal](https://github.com/elpinal)

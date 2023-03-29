//go:build mage
// +build mage

package main

import (
	"github.com/dmvolod/mageutil/bintool"
	"github.com/dmvolod/mageutil/shellcmd"
)

var (
	linter = bintool.Must(bintool.New(
		"golangci-lint{{.BinExt}}",
		"1.52.2",
		"https://github.com/golangci/golangci-lint/releases/download/v{{.Version}}/golangci-lint-{{.Version}}-{{.GOOS}}-{{.GOARCH}}{{.ArchiveExt}}",
	))

	mage = bintool.Must(bintool.New(
		"mage{{.BinExt}}",
		"1.14.0",
		"https://github.com/magefile/mage/releases/download/v{{.Version}}/mage_{{.Version}}_{{.GOOS}}-{{.GOARCH}}{{.ArchiveExt}}",
		bintool.WithGoBinFolder(),
		bintool.WithOsSubstitution(mageOsRepl),
		bintool.WithArchSubstitution(mageArchRepl),
	))

	mageOsRepl = map[string]string{
		"darwin":  "macOS",
		"linux":   "Linux",
		"windows": "Windows",
	}
	mageArchRepl = map[string]string{
		"amd64": "64bit",
		"386":   "32bit",
		"arm":   "ARM",
		"arm64": "ARM64",
	}
)

func Mage() error {
	return mage.Ensure()
}

func Lint() error {
	if err := linter.Ensure(); err != nil {
		return err
	}

	return linter.Command("run").Run()
}

func Test() error {
	return shellcmd.Command("go test -race -coverprofile=coverage.out ./...").Run()
}

func Coverage() error {
	return shellcmd.Command("go tool cover -html=coverage.out").Run()
}

func Generate() error {
	return shellcmd.Command("go generate ./...").Run()
}

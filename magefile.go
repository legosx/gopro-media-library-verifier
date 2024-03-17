//go:build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"runtime"
)

const (
	ldflags = "-ldflags=-extldflags -static"
)

type (
	Test  mg.Namespace
	Build mg.Namespace
	Lint  mg.Namespace
)

type (
	Arch string
	OS   string
)

const (
	ArchARM64 Arch = "arm64"
	ArchAMD64 Arch = "amd64"

	OSDarwin OS = "darwin"
	OSLinux  OS = "linux"
)

type BuildEnv struct {
	Architecture Arch
	OS           OS
}

func MacOS() BuildEnv {
	return BuildEnv{
		Architecture: Arch(runtime.GOARCH),
		OS:           OSDarwin,
	}
}

var Aliases = map[string]interface{}{
	"build": Build.Verifier,
}

func (Build) Verifier() error {
	return goBuild("main.go", "bin/verifier", MacOS())
}

// Generate mocks
func Generate() error {
	return sh.RunV("go", "generate", "./...")
}

func goBuild(input, output string, buildEnv BuildEnv) error {
	env := map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        string(buildEnv.OS),
		"GOARCH":      string(buildEnv.Architecture),
	}
	a := []string{"build", "-a", ldflags}
	a = append(a, "-o", output, "-v", input)

	return sh.RunWithV(env, "go", a...)
}

func (Test) Unit() error {
	return test()
}

func (Lint) Go() error {
	return sh.RunV("bash", "-c", "golangci-lint run")
}

func test(args ...string) error {
	a := []string{"test", "./...", "-test.short", "-race"}
	a = append(a, args...)

	return sh.RunV("go", a...)
}

func (Test) Coverage() (err error) {
	if err = test("-coverprofile=coverage.tmp", "-covermode=atomic", "-coverpkg", "./..."); err != nil {
		return err
	}

	if err = sh.RunV("bash", "-c", "cat coverage.tmp > coverage"); err != nil {
		return err
	}

	if err = sh.RunV("bash", "-c", "go tool cover -func=coverage | tail -n1"); err != nil {
		return err
	}

	return sh.RunV("bash", "-c", "gocover-cobertura < coverage > coverage.xml")
}

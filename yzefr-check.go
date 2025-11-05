package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

// Version as printed with -version option
var Version = "UNKNOWN"

const (
	// Help as printed with -help option
	Help = `yzefr-check [-h] [-v] files...
Project does something:
-help       To print this help
-version    To print version
files...    Files to check`
)

type Header struct {
	Title   string `yaml:"title" validate:"required"`
	Desc    string `yaml:"desc" validate:"required"`
	Authors []struct {
		Name  string `yaml:"name" validate:"required"`
		Email string `yaml:"email"`
	}
	Date    string   `yaml:"date" validate:"required"`
	Release string   `yaml:"release"`
	Games   []string `yaml:"games" validate:"required"`
	Tags    []string `yaml:"tags" validate:"required"`
}

func main() {
	help := flag.Bool("help", false, "Print help")
	version := flag.Bool("version", false, "Print version")
	flag.Parse()
	if *help {
		fmt.Println(Help)
		os.Exit(0)
	}
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}
	if err := run(flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR %v", err)
		os.Exit(1)
	}
}

func run(files []string) error {
	validate := validator.New()
	for _, file := range files {
		text, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("reading file %s: %v", file, err)
		}
		header, err := extractHeader(text)
		if err != nil {
			return fmt.Errorf("extracting header from file %s: %v", file, err)
		}
		var info Header
		if err := yaml.UnmarshalStrict(header, &info); err != nil {
			return fmt.Errorf("checking file %s: %v", file, err)
		}
		if err := validate.Struct(info); err != nil {
			return err
		}
	}
	return nil
}

func extractHeader(text []byte) ([]byte, error) {
	r := regexp.MustCompile(`(?s)---(.+?)---`)
	match := r.FindSubmatch(text)
	r.FindSubmatch(text)
	if len(match) == 0 {
		return nil, fmt.Errorf("header not found")
	}
	if len(match) > 2 {
		return nil, fmt.Errorf("found %d headers", len(match))
	}
	return bytes.TrimSpace(match[1]), nil
}

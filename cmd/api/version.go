package main

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	version   string
	buildDate string
)

func printVersion() {
	fmt.Print(format(version, buildDate))
}

func format(version, buildDate string) string {
	version = strings.TrimPrefix(version, "v")

	var dateStr string
	if buildDate != "" {
		dateStr = fmt.Sprintf(" (%s)", buildDate)
	}
	if version == "" {
		version = "develop"
	}

	return fmt.Sprintf("wb-l0 version %s%s\n%s\n", version, dateStr, changelogURL(version))
}

func changelogURL(version string) string {
	path := "https://github.com/maypok86/wb-l0"
	r := regexp.MustCompile(`^v?\d+\.\d+\.\d+(-[\w.]+)?$`)
	if !r.MatchString(version) {
		return fmt.Sprintf("%s/releases/latest", path)
	}

	url := fmt.Sprintf("%s/releases/tag/v%s", path, strings.TrimPrefix(version, "v"))
	return url
}

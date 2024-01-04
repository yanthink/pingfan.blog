package helpers

import (
	"regexp"
	"strings"
)

func MatchImageUrls(md string) (urls []string) {
	labelPattern := `(?:[^][]+)`
	hrefPattern := `(?:[^ ()]+|\([^ )]+\))*`
	titlePattern := `"[^"]*"|'[^\']*'?`

	imagePattern := strings.NewReplacer(
		"label", labelPattern,
		"href", hrefPattern,
		"title", titlePattern,
	).Replace(`!\[(label)]\(\s*(href)(?:\s+(title))?\s*\)`)

	imageRegex := regexp.MustCompile(imagePattern)

	matches := imageRegex.FindAllStringSubmatch(md, -1)

	for _, match := range matches {
		urls = append(urls, match[2])
	}

	return
}

package main

import (
	"errors"
	"log"
	"regexp"
)

type HTMLElem struct {
	Content    string
	Attributes map[string]string
}

func (h *HTMLElem) removeRubyTags() {
	rubyRegex := regexp.MustCompile(`\</{0,1}(rb|ruby|rt)\>`)
	rubyRegex.ReplaceAllString(h.Content, "")
}

// func Get

func GetTitle(fileContent []byte) string {
	tag := "title"
	indices, err := findTagIndices(fileContent, tag)
	if err != nil {
		log.Fatalf("Unable to get title %v", err)
	}
	// indices returned are start and end indices of title element content.
	// the right side of the slice is inclusive
	return string(fileContent[indices[0]:indices[1]])
}

func findTagIndices(fileContent []byte, tag string) ([]int, error) {
	r := regexp.MustCompile(`\<` + tag)
	gt := regexp.MustCompile(`\>`)
	endTag := regexp.MustCompile(`\<\/` + tag)

	tg := r.FindIndex(fileContent) // found tag, now search past that location
	if tg == nil {
		return []int{}, errors.New("tag does not exist")
	}
	t := tg[1] + 1
	// found end of start tag - mark index of
	// start of content
	startIndex := gt.FindIndex(fileContent[t:])
	if startIndex == nil {
		return []int{}, errors.New("error parsing for end of start tag")
	}
	s := startIndex[1] + 1 // save location of first valid content char
	endIndex := endTag.FindIndex(fileContent[s:])
	if endIndex == nil {
		return []int{}, errors.New("error parsing for end of content")
	}
	e := endIndex[0] - 1 // save location last valid character (just before `<` sign)

	return []int{s, e}, nil
}

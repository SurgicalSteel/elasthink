package module

import (
	"github.com/SurgicalSteel/elasthink/util"
	"regexp"
	"strings"
)

func tokenizeIndonesianSearchTerm(s string) map[string]int {
	s = strings.ToLower(s)
	regex := regexp.MustCompile(`[^a-zA-Z0-9 ]`)
	s = regex.ReplaceAllString(s, ` `)
	rawWords := strings.Split(s, " ")

	words := make([]string, 0)

	for _, rw := range rawWords {
		rw = strings.Trim(rw, " ")
		if len(rw) > 0 {
			words = append(words, rw)
		}
	}

	wordsSet := util.CreateWordSet(words)

	return wordsSet //util.WordsSetSubtraction(wordsSet, moduleObj.StopwordSet)
}

package module

// Elasthink, An alternative to elasticsearch engine written in Go for small set of documents that uses inverted index to build the index and utilizes redis to store the indexes.
// Copyright (C) 2020 Yuwono Bangun Nagoro (a.k.a SurgicalSteel)
//
// This file is part of Elasthink
//
// Elasthink is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Elasthink is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
import (
	"github.com/SurgicalSteel/elasthink/util"
	"regexp"
	"strings"
)

//tokenize is a tokenizer function, basically to remove punctuations and (optional) to remove stopwords to a document name or a search term
func tokenize(s string) map[string]int {
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

	if moduleObj.IsUsingStopwordRemoval {
		return util.WordsSetSubtraction(wordsSet, moduleObj.StopwordSet)
	}

	return wordsSet
}

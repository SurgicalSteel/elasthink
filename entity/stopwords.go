package entity

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
import ()

//StopwordData is a struct that represent stopwords data that we have in a file for a specified language
type StopwordData struct {
	Words []string `json:"words"` //Words is a collection of stopwords for a specified language
}

var stopwords map[string]int

//InitStopword is the function to initialize a stopword set from a given stopword data
func InitStopword(stopwordData StopwordData) {
	tempStopwords := make(map[string]int)
	for _, word := range stopwordData.Words {
		tempStopwords[word] = 1
	}
	stopwords = tempStopwords
}

//GetStopword is a function to get the stopword
func GetStopword() map[string]int {
	return stopwords
}

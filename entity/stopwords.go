package entity

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

package util

import ()

//CreateWordSet create a set of words from a slice of string
func CreateWordSet(words []string) map[string]int {
	result := make(map[string]int)

	for _, vw := range words {
		if _, okr := result[vw]; !okr {
			result[vw] = 1
		}
	}

	return result
}

//WordsSetUnion do a union from two given words sets
func WordsSetUnion(ma, mb map[string]int) map[string]int {
	result := make(map[string]int)
	for ka, _ := range ma {
		if _, oka := result[ka]; !oka {
			result[ka] = 1
		}
	}

	for kb, _ := range mb {
		if _, okb := result[kb]; !okb {
			result[kb] = 1
		}
	}
	return result
}

//WordsSetSubtraction do a substraction between two sets ma - mb
func WordsSetSubtraction(ma, mb map[string]int) map[string]int {
	result := ma
	for kb, _ := range mb {
		if _, okb := result[kb]; okb {
			delete(result, kb)
		}
	}
	return result
}

//WordsSetIntersection find the intersection between two sets (ma & mb)
func WordsSetIntersection(ma, mb map[string]int) map[string]int {
	result := make(map[string]int)
	for ka, _ := range ma {
		if _, okb := mb[ka]; okb {
			result[ka] = 1
		}
	}
	return result
}

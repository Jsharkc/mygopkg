package stringutil

import (
	"strings"
	"unicode/utf8"

	"github.com/hbollon/go-edlib"
)

func Similarity(str1, str2 string) (float32, error) {
	return edlib.StringsSimilarity(str1, str2, edlib.Levenshtein)
}
func FindSimilarityIndex(target string, strList ...string) (int, float32, error) {
	var higherMatchPercent float32
	var tmpIndex int
	var sim float32
	var err error
	for index, strToCmp := range strList {
		if utf8.RuneCountInString(strToCmp) <= 2 {
			sim = 0.0
		} else if strings.Contains(strToCmp, target) {
			sim = 1.0
		} else {
			sim, err = edlib.StringsSimilarity(target, strToCmp, edlib.Jaro)
			if err != nil {
				return index, sim, err
			}
		}
		if sim == 1.0 {
			return index, sim, nil
		} else if sim > higherMatchPercent {
			higherMatchPercent = sim
			tmpIndex = index
		}
	}
	return tmpIndex, higherMatchPercent, nil
}

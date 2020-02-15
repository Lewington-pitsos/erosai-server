package scanner

import (
	"strings"

	"bitbucket.org/lewington/erosai/assist"
)

var keyTerms = []string{
	"dominated",
	"stepdad",
	"stepmom",
	"teen",
	"cum",
	"lubed",
	"big butts",
	"sex",
	"step sister",
	"fucked",
	"cock",
	"assfucked",
	"big cock",
	"slut",
	"sluts",
	"fucks",
	"fucking",
	"cumming",
	"teen girl",
	"gets fucked",
	"big tits",
	"hentai",
	"milf",
	"anal",
	"anal fuck",
	"busty",
	"sexy",
	"porn",
	"xxx",
}

type wordInspector struct {
}

func (i *wordInspector) score(page string) int {
	score := 0
	for _, word := range keyTerms {
		score += strings.Count(page, word)
	}

	return assist.Min(100, score*15)
}

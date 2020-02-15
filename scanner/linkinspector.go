package scanner

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

var source = rand.NewSource(time.Now().Unix())
var random = rand.New(source)

type linkInspector struct {
	regex          *regexp.Regexp
	remainingLinks []string
	viewedLinks    []string
}

func (l *linkInspector) score(page string) int {
	l.viewedLinks = []string{}
	l.remainingLinks = l.regex.FindAllString(page, -1)

	score := 0

	for i := 0; i < 5; {

		if len(l.remainingLinks) == 0 {
			break
		}

		index := random.Intn(len(l.remainingLinks))

		link := l.remainingLinks[index]

		l.remainingLinks = append(l.remainingLinks[:index], l.remainingLinks[index+1:]...)

		if !l.hasBeenViewed(link) {
			fmt.Println(link)
			l.viewedLinks = append(l.viewedLinks, link)
			i++
			// return highest score
		}
	}

	return score
}

func (l *linkInspector) hasBeenViewed(link string) bool {
	for _, viewedLink := range l.viewedLinks {
		if viewedLink == link {
			return true
		}
	}

	return false
}

func newLinkInspector(regex *regexp.Regexp) *linkInspector {
	return &linkInspector{
		regex,
		[]string{},
		[]string{},
	}
}

package scanner

import (
	"fmt"
	"crypto/tls"
	"net/http"
	"math/rand"
	"regexp"
	"time"
	"strconv"
	"io/ioutil"

	"bitbucket.org/lewington/erosai-server/globals"
)

var client = &http.Client{
	Timeout: time.Second * 10,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

var source = rand.NewSource(time.Now().Unix())
var random = rand.New(source)

type linkInspector struct {
	regex          *regexp.Regexp
	remainingLinks []string
	viewedLinks    []string
	saver *imageSaver
}

func (l *linkInspector) score(page string) int {
	l.viewedLinks = []string{}
	l.remainingLinks = l.regex.FindAllString(page, -1)

	score := 0

	for i := 0; i < 10; {
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
			tempScore := int(100.0 * l.scoreImage(link))
			if tempScore > score {
				score = tempScore
			}
			if score > globals.NsfwCutoff {
				break
			}
		}
	}

	fmt.Println(score)

	return score
}

func (l *linkInspector) scoreImage(URL string) float64 {
	filename, err := l.saver.save(URL)
	defer l.saver.delete()

	if err != nil {
		fmt.Println(err)
		return -1.0
	}

	request, err := http.NewRequest("GET", globals.MLServerEndpoint + "?filename=" + filename, nil)
	if err != nil {
		fmt.Println(err)
		return -1.0
	}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return -1.0
	}

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return -1.0
	}


	score, err := strconv.ParseFloat(string(responseBytes), 64)

	if err != nil {
		fmt.Println(err)
		return -1.0
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
		&imageSaver{},
	}
}

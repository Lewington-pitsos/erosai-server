package scanner

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"bitbucket.org/lewington/erosai/lg"
	"bitbucket.org/lewington/erosai/shared"
	"github.com/PuerkitoBio/goquery"
)

var client = &http.Client{
	Timeout: time.Second * 10,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

type Scanner struct {
	input chan shared.Link
	regex *regexp.Regexp
}

func (s *Scanner) work() {
	go s.scanLnks()
}

func (s *Scanner) scanLnks() {
	for link := range s.input {
		request, err := http.NewRequest("GET", link.URL, nil)
		if err != nil {
			lg.L.Debug("error creating request for URL %v", link.URL)
		}
		resp, err := client.Do(request)
		if err != nil {
			lg.L.Debug("error requesting URL %v", link.URL)
		}
		s.containsPorn(resp)
	}
}

func (s *Scanner) containsPorn(resp *http.Response) (int, error) {

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	responseString := string(responseBytes)

	if err != nil {
		return 0, err
	}

	links := s.regex.FindAllString(responseString, -1)

	for _, link := range links {
		fmt.Println(link)
	}

	return 0, nil
}

func (s *Scanner) extractImageLinksFunc(links []string) func(int, *goquery.Selection) bool {
	return func(index int, imageNode *goquery.Selection) bool {
		link, exists := imageNode.Attr("src")

		if exists {
			links[index] = link
		}

		return true
	}
}

func New(input chan shared.Link) *Scanner {
	s := &Scanner{
		input,
		regexp.MustCompile(`http\S{5,250}.(jpg|jpeg|png)`),
	}

	s.work()

	return s
}

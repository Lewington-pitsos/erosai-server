package scanner

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"bitbucket.org/lewington/erosai/lg"
	"bitbucket.org/lewington/erosai/shared"
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

	inspectors []inspector
}

func (s *Scanner) work() {
	go s.scanLnks()
}

func (s *Scanner) scanLnks() {
	for link := range s.input {
		request, err := http.NewRequest("GET", link.URL, nil)
		if err != nil {
			lg.L.Debug("error creating request for URL %v", link.URL)
			break
		}
		resp, err := client.Do(request)
		if err != nil {
			lg.L.Debug("error requesting URL %v", link.URL)
			break
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

	var score int
	for _, inspector := range s.inspectors {
		score = inspector.score(responseString)
		if score > 70 {
			return score, nil
		}
	}

	return score, nil
}

func New(input chan shared.Link) *Scanner {
	s := &Scanner{
		input,
		[]inspector{
			&wordInspector{},
			newLinkInspector(regexp.MustCompile(`http\S{5,250}.(jpg|jpeg|png)`)),
		},
	}

	s.work()

	return s
}

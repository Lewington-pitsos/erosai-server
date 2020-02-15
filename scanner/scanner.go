package scanner

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"bitbucket.org/lewington/erosai/database"
	"bitbucket.org/lewington/erosai/lg"
	"bitbucket.org/lewington/erosai/shared"
)

var pornCutoff = 70

var client = &http.Client{
	Timeout: time.Second * 10,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

type Scanner struct {
	input      chan shared.Link
	arch       database.Archivist
	inspectors []inspector
}

func (s *Scanner) work() {
	go s.scanLnks()
}

func (s *Scanner) scanLnks() {
	for link := range s.input {
		lg.L.Debug("scanning %v", link)
		request, err := http.NewRequest("GET", link.URL, nil)
		if err != nil {
			lg.L.Debug("error creating request for URL %v", link.URL)
			continue
		}
		resp, err := client.Do(request)
		if err != nil {
			lg.L.Debug("error requesting URL %v", link.URL)
			continue
		}

		score, err := s.pornScore(resp)

		if err != nil {
			lg.L.Debug("error scoring URL %v", link.URL)
			continue
		}

		link.Scanned = true
		link.Porn = score
		s.arch.UpdateLink(link)
	}
}

func (s *Scanner) pornScore(resp *http.Response) (int, error) {
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
		if score > pornCutoff {
			return score, nil
		}
	}

	return score, nil
}

func New(input chan shared.Link) *Scanner {
	s := &Scanner{
		input,
		database.NewArchivist(),
		[]inspector{
			// &wordInspector{},
			newLinkInspector(regexp.MustCompile(`http\S{5,250}.(jpg|jpeg|png)`)),
		},
	}

	s.work()

	return s
}

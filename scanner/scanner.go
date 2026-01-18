package scanner

import (
	"io/ioutil"
	"net/http"
	"regexp"

	"bitbucket.org/lewington/erosai-server/database"
	"bitbucket.org/lewington/erosai-server/lg"
	"bitbucket.org/lewington/erosai-server/shared"
)

var nsfwCutoff = 70

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

		score, err := s.nsfwScore(resp)

		if err != nil {
			lg.L.Debug("error scoring URL %v", link.URL)
			continue
		}

		link.Scanned = true
		link.Nsfw = score
		s.arch.UpdateLink(link)
	}
}

func (s *Scanner) nsfwScore(resp *http.Response) (int, error) {
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
		if score > nsfwCutoff {
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

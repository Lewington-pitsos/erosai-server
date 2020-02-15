package scanner

type inspector interface {
	score(page string) int
}

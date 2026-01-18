package shared

type Link struct {
	ID      int
	URL     string
	Scanned bool
	Nsfw    int
}

func NewUnscannedLink(ID int, URL string) Link {
	return Link{
		ID,
		URL,
		false,
		-1,
	}
}

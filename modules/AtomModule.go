package modules

type Author struct {
	Name string `xml:"name"`
	Uri  string `xml:"uri"`
}

type Entry struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	Summary string `xml:"summary"`
	Updated string `xml:"updated"`
}

type AtomModule struct {
	Title    string   `xml:"title"`
	Subtitle string   `xml:"subtitle"`
	Updated  string   `xml:"updated"`
	Entry    []*Entry `xml:"entry"`
}

// Package movie is a set of Go types on to which the IMDB movie record is
// unmarshalled.
package movie

import "strings"

// Movie parses an IMDB movie record text blob, extracting all metadata about
// a movie title.
// The core details are extracted from the MOVI entry, but also all the other
// entry types; ADPT, NOVL, CRIT, SCRP, etc.
type Movie struct {
	Title         string
	Year          int
	Month         int
	TV            bool
	SeriesName    string
	SeriesNumber  int
	EpisodeNumber int

	Adaptations         []Adaptation
	Books               []Book
	Critiques           []Critique
	Essays              []Essay
	Interviews          []Interview
	Novels              []Novel
	Others              []Other
	ProductionProtocols []ProductionProtocol
	Screenplays         []Screenplay
}

// Adaptation parses an ADPT (adapted literary source) record entry.
type Adaptation struct {
	Book
}

// Book parse a BOOK (monographic book) record entry.
type Book struct {
	Title          string
	Author         string
	Publisher      Publisher
	Date           Date
	PageCount      int
	Volume         string
	Issue          string
	ISBN           string
	FirstPublished int
	Note           string
	MiscInfo       string // from the "In:" info; usually just www links or other random text.
}

// Novel parses a NOVL (original literary source) record entry.
type Novel struct {
	Book
}

// Critique parses a CRIT (printed media reviews) record entry.
type Critique struct {
	Publication
}

// Essay parses an ESSY (printed essay) record entry.
type Essay struct {
	Publication
}

// Interview parses an IVIW (interview with cast or crew) record entry.
type Interview struct {
	Publication
}

// Other parses an OTHR (other literature) record entry.
type Other struct {
	Publication
}

// ProductionProtocol parses a PROT (production protocol) record entry.
type ProductionProtocol struct {
	Publication
}

// Screenplay parses a SCRP (published screenplay) record entry.
type Screenplay struct {
	Publication
}

// Publication is a base type used by non-book entries such as CRIT, ESSY, etc.
// generally representing a magazine.
type Publication struct {
	// core publication details.
	Name      string
	Publisher Publisher
	Date      Date
	Volume    string
	Issue     string
	ISSN      string // sometime contains an ISBN

	// details related specifically to the IMDB entries.
	ArticleAuthor string
	ArticleTitle  string
	ArticlePages  string // e.g. `1-17`, `56`, `23, 24, 66`

	// The interview subject, only used by IVIW.
	ArticleInterviewee string
}

// Publisher metadata is used by all books and publications for the publisher details.
type Publisher struct {
	Name    string
	City    string
	State   string // provence, county, state, etc.
	Country string
}

// Date is a generic type for storing dates without having to parse into time.Time objects.
type Date struct {
	Year  int
	Month int
	Day   int
}

// Unmarshall processes all record entries types.
func Unmarshall(data string, movie *Movie) {
	entry := extractEntryDataTypes(data)

	entry.movieTitleDetails(movie)
	entry.adaptations(movie)
	entry.books(movie)
	entry.novels(movie)

	entry.critiques(movie)
	entry.essays(movie)
	entry.interviews(movie)
	entry.others(movie)
	entry.productionProtocols(movie)
	entry.screenplays(movie)
}

// UnmarshallBooks processes only the record entries types that are types of books.
func UnmarshallBooks(data string, movie *Movie) {
	entry := extractEntryDataTypes(data)

	entry.movieTitleDetails(movie)
	entry.adaptations(movie)
	entry.books(movie)
	entry.novels(movie)
}

// IsAdaptation checks all book types (ADPT, BOOK, NOVL) and returns true if a
// title/author match is found.
func (m *Movie) IsAdaptation(title, author string) bool {
	for _, a := range m.Adaptations {
		if m.titleMatches(a.Title, title) && m.authorMatches(a.Author, author) {
			return true
		}
	}
	for _, b := range m.Books {
		if m.titleMatches(b.Title, title) && m.authorMatches(b.Author, author) {
			return true
		}
	}
	for _, n := range m.Novels {
		if m.titleMatches(n.Title, title) && m.authorMatches(n.Author, author) {
			return true
		}
	}

	return false
}

func (m *Movie) titleMatches(srcTitle, testTitle string) bool {
	title := strings.ToLower(srcTitle)
	testable := strings.ToLower(testTitle)

	title = strings.ReplaceAll(title, "the ", "")
	testable = strings.ReplaceAll(testable, "the ", "")

	return strings.Contains(title, testable)
}

func (m *Movie) authorMatches(srcAuthor, testAuthor string) bool {
	author := strings.ToLower(srcAuthor)
	testable := strings.ToLower(testAuthor)

	author = strings.ReplaceAll(author, ",", "")
	testable = strings.ReplaceAll(testable, ",", "")

	// perhaps splitting the name on spaces and checking each part is present
	// would be a reasonable approach:
	matching := true
	names := strings.Fields(testable)
	for _, name := range names {
		if !strings.Contains(author, name) {
			matching = false
		}
	}

	return matching
}

var monthRomanToInt = map[string]int{
	"i": 1, "ii": 2, "iii": 3, "iv": 4, "v": 5, "vi": 6,
	"vii": 7, "viii": 8, "ix": 9, "x": 10, "xi": 11, "xii": 12,
}

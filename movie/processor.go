package movie

import (
	"regexp"
	"strconv"
	"strings"
)

// compile all regular expression up front for major performance improvements
var (
	titleDetailsRegExp     = regexp.MustCompile(`\A(.*?) \(([0-9?]{4})(?:/([IVX]+))?\)`)
	authorRegExp           = regexp.MustCompile(`^(.+?)\. +"`) // NOTE: also matches the opening " of the title
	titleRegExp            = regexp.MustCompile(`^"([^"]+?)"\.? *`)
	publisherRegExp        = regexp.MustCompile(`^\(([^)]+?)\)(?:, *)?|^([^:]+?): *`)
	publisherNameRegEx     = regexp.MustCompile(`^([^,]+)(?:, *)?`)
	notesRegExp            = regexp.MustCompile(`\((.+?)\)$`)
	bracesRegExp           = regexp.MustCompile(`^\((.+?)\)$`)
	randomTextRegExp       = regexp.MustCompile(`(?i), *(?:\(BK\)|\(HB\)|\(MG\)|\(NP\)|\(Novel\)|NONE|Pg\. N/?A|\(tme\d+\))`)
	inRegExp               = regexp.MustCompile(`In: "(.+?)"(?:, *)?`)
	isbnRegExp             = regexp.MustCompile(`, *IS[BS]N(?:-\d\d)?: ([0-9X-]+)`)
	pageCountRegExp        = regexp.MustCompile(`, *Pg. *([0-9]+)`)
	pageRangeCleanupRegExp = regexp.MustCompile(`(?i), *Pg\. *(?:pg[ds]?[.;?]|pg>\.|p/ n°\.|p[a^]gs\.|Pages: *) *`)
	pageRangeRegExp        = regexp.MustCompile(`(?i), *Pg\. *((?:[a-z]?[0-9]+)(?:(?:-|\+|, *| *to *)[a-z]?[0-9]+)*)`)
	firstPublishedRegExp   = regexp.MustCompile(`(?i)First published.+?(\d\d\d\d).?`)
	publishedRegExp        = regexp.MustCompile(`\(?((?:\d{1,2} +)?(?:[JFMASOND][a-z]+ +)?\d{4})\)?`)
	volumeNumberRegExp     = regexp.MustCompile(`, *Vol. *#? *([0-9]+)`)
	issueNumberRegExp      = regexp.MustCompile(`, *Iss. *#? *([0-9]+)`)
)

type Key string

// List of all the different record entries available to a movie record.
const (
	ADPT Key = "ADPT"
	BOOK Key = "BOOK"
	CRIT Key = "CRIT"
	ESSY Key = "ESSY"
	IVIW Key = "IVIW"
	MOVI Key = "MOVI"
	NOVL Key = "NOVL"
	OTHR Key = "OTHR"
	PROT Key = "PROT"
	SCRP Key = "SCRP"
)

type textEntry map[Key][]string

func extractEntryDataTypes(data string) textEntry {
	e := textEntry{}

	lines := strings.Split(data, "\n")
	for _, line := range lines {
		l := strings.TrimSpace(line)
		if len(l) == 0 {
			continue
		}
		kv := strings.SplitN(l, ":", 2)
		k := Key(strings.TrimSpace(kv[0]))
		v := strings.TrimSpace(kv[1])

		switch k {
		case ADPT:
			e[ADPT] = append(e[ADPT], v)
		case BOOK:
			e[BOOK] = append(e[BOOK], v)
		case CRIT:
			e[CRIT] = append(e[CRIT], v)
		case ESSY:
			e[ESSY] = append(e[ESSY], v)
		case IVIW:
			e[IVIW] = append(e[IVIW], v)
		case MOVI:
			e[MOVI] = append(e[MOVI], v)
		case NOVL:
			e[NOVL] = append(e[NOVL], v)
		case OTHR:
			e[OTHR] = append(e[OTHR], v)
		case PROT:
			e[PROT] = append(e[PROT], v)
		case SCRP:
			e[SCRP] = append(e[SCRP], v)
		default:
			// ignore unknown keys
		}
	}

	return e
}

func (e textEntry) movieTitleDetails(movie *Movie) {
	if len(e[MOVI]) == 0 {
		return
	}

	m := e[MOVI][0]
	results := titleDetailsRegExp.FindStringSubmatch(m)

	if len(results) >= 1 {
		movie.Title = results[1]
	}

	if len(results) >= 2 {
		movie.Year, _ = strconv.Atoi(results[2])
	}

	if len(results) >= 3 {
		movie.Month = results[3]
	}

	if strings.Contains(m, "(TV)") {
		movie.TV = true
	}
}

func (e textEntry) adaptations(movie *Movie) {
	for _, text := range e[ADPT] {
		a := Adaptation{}
		e.bookParser(&a.Book, text)
		movie.Adaptations = append(movie.Adaptations, a)
	}
}

func (e textEntry) books(movie *Movie) {
	for _, text := range e[BOOK] {
		book := Book{}
		e.bookParser(&book, text)
		movie.Books = append(movie.Books, book)
	}
}

func (e textEntry) novels(movie *Movie) {
	for _, text := range e[NOVL] {
		novel := Novel{}
		e.bookParser(&novel.Book, text)
		movie.Novels = append(movie.Novels, novel)
	}
}

// bookParser extracts the various components from the text string.
// Anything that has a knowable marker (e.g. `Pg.`, `ISBN:`) should be
// extracted first, leaving only positional items (no knowable marker,
// such as author names), to be done last.
//
// IMPORTANT: the order of func calls is super important. Key based matching (e.g. ISBN)
// should go before positional based (e.g. Author, Title).
func (e textEntry) bookParser(book *Book, data string) {
	data = e.cleanSurroundingBraces(data)
	data = e.cleanRandomText(data)

	// TODO: should these values be used?
	_, data = e.extractIn(data)
	_, data = e.extractVolumeNumber(data)
	_, data = e.extractIssueNumber(data)

	book.ISBN, data = e.extractISBN(data)
	book.PageCount, data = e.extractPageCount(data)

	// NOTE: must be done before other date processing
	book.FirstPublished, data = e.extractFirstPublishedDate(data)

	book.Date, data = e.extractPublishedDate(data)

	// the items below are positional based

	book.Author, data = e.extractAuthor(data)
	book.Title, data = e.extractTitle(data)

	book.Publisher, data = e.extractPublisher(data)
	book.Note, data = e.extractNotes(data)
}

var months = []string{"january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december"}

func (e textEntry) monthAsNumber(month string) int {
	for i, m := range months {
		if strings.ToLower(month) == m {
			return i + 1
		}
	}
	return 0
}

func (e textEntry) critiques(movie *Movie) {
	for _, text := range e[CRIT] {
		novel := Critique{}
		e.publicationParser(&novel.publication, text)
		movie.Critiques = append(movie.Critiques, novel)
	}
}

func (e textEntry) essays(movie *Movie) {
	for _, text := range e[ESSY] {
		essay := Essay{}
		e.publicationParser(&essay.publication, text)
		movie.Essays = append(movie.Essays, essay)
	}
}

func (e textEntry) interviews(movie *Movie) {
	for _, text := range e[IVIW] {
		interview := Interview{}
		e.publicationParser(&interview.publication, text)
		movie.Interviews = append(movie.Interviews, interview)
	}
}

func (e textEntry) others(movie *Movie) {
	for _, text := range e[OTHR] {
		others := Other{}
		e.publicationParser(&others.publication, text)
		movie.Others = append(movie.Others, others)
	}
}

func (e textEntry) productionProtocols(movie *Movie) {
	for _, text := range e[PROT] {
		protocol := ProductionProtocol{}
		e.publicationParser(&protocol.publication, text)
		movie.ProductionProtocols = append(movie.ProductionProtocols, protocol)
	}
}

func (e textEntry) screenplays(movie *Movie) {
	for _, text := range e[SCRP] {
		screenplay := Screenplay{}
		e.publicationParser(&screenplay.publication, text)
		movie.Screenplays = append(movie.Screenplays, screenplay)
	}
}

// publicationParser extracts the various components from the text string.
// Anything that has a knowable marker (e.g. `Pg.`, `ISSNN:`) should be
// extracted first, leaving only positional items (no knowable marker,
// such as author names), to be done last.
func (e textEntry) publicationParser(pub *publication, data string) {
	data = e.cleanSurroundingBraces(data)
	data = e.cleanRandomText(data)

	pub.Volume, data = e.extractVolumeNumber(data)
	pub.Issue, data = e.extractIssueNumber(data)
	pub.ISSN, data = e.extractISBN(data)
	pub.ArticlePages, data = e.extractPageRange(data)

	pub.Name, data = e.extractIn(data)
	pub.Date, data = e.extractPublishedDate(data)

	//
	// the items below are positional based
	//
	pub.ArticleAuthor, data = e.extractAuthor(data)
	pub.ArticleTitle, data = e.extractTitle(data)

	pub.Publisher, data = e.extractPublisher(data)
}

// extractAuthor from the text, based on its position in the text.
func (e textEntry) extractAuthor(data string) (author string, str string) {
	results := authorRegExp.FindStringSubmatch(data)
	if authorRegExp.MatchString(data) {
		author = strings.TrimSpace(results[1])
	}
	str = authorRegExp.ReplaceAllString(data, "") // NOTE: this also replaces the opening "
	str = strings.TrimSpace(str)
	// NOTE: add the " we removed above, but sometimes there is no author, so check first
	if len(str) > 0 && str[0] != '"' {
		str = `"` + str
	}

	return
}

// extractTitle from the text, based on its position in the text.
func (e textEntry) extractTitle(data string) (title string, str string) {
	results := titleRegExp.FindStringSubmatch(data)
	if titleRegExp.MatchString(data) {
		title = strings.TrimSpace(results[1])
	}
	str = titleRegExp.ReplaceAllString(data, "")
	str = strings.TrimSpace(str)

	return
}

// extractPublisher from the text, based on its position in the text.
func (e textEntry) extractPublisher(data string) (pub Publisher, str string) {
	results := publisherRegExp.FindStringSubmatch(data)
	if publisherRegExp.MatchString(data) {
		match := results[1]
		if len(match) == 0 {
			match = results[2]
		}
		parts := strings.Split(match, ",")
		loc := make([]string, len(parts))
		for i, l := range parts {
			loc[i] = strings.TrimSpace(l)
		}
		pub.Country = loc[len(loc)-1]
		if len(loc) >= 2 {
			parts := strings.Join(loc[0:len(loc)-1], ", ")
			pub.City = parts
		}
	}
	data = publisherRegExp.ReplaceAllString(data, "")
	data = strings.TrimSpace(data)

	// now extract the publisher name
	results = publisherNameRegEx.FindStringSubmatch(data)
	if publisherNameRegEx.MatchString(data) {
		pub.Name = strings.TrimSpace(results[1])
	}
	str = publisherNameRegEx.ReplaceAllString(data, "")
	str = strings.TrimSpace(str)

	return
}

// extractNotes from the text, based on its position in the text.
// NOTE: expected to always be last item
func (e textEntry) extractNotes(data string) (note string, str string) {
	results := notesRegExp.FindStringSubmatch(data)
	if notesRegExp.MatchString(data) {
		note = strings.TrimSpace(results[1])
	}
	str = notesRegExp.ReplaceAllString(data, "")
	str = strings.TrimSpace(str)

	return
}

// Remove wrapping braces (only a _few_ entries are wrapped in braces)
func (e textEntry) cleanSurroundingBraces(data string) string {
	results := bracesRegExp.FindStringSubmatch(data)
	if !bracesRegExp.MatchString(data) {
		return data
	}
	result := strings.TrimSpace(results[1])
	data = bracesRegExp.ReplaceAllString(data, result)
	return strings.TrimSpace(data)
}

// Remove random items (e.g. `(BK)`, `NONE`) from data.
func (e textEntry) cleanRandomText(data string) string {
	data = randomTextRegExp.ReplaceAllString(data, "")
	return strings.TrimSpace(data)
}

// extractIn from the text, based on the key
func (e textEntry) extractIn(data string) (in string, str string) {
	results := inRegExp.FindStringSubmatch(data)
	if inRegExp.MatchString(data) {
		in = strings.TrimSpace(results[1])
	}
	str = inRegExp.ReplaceAllString(data, "")
	str = strings.TrimSpace(str)

	return
}

// extractISBN from the text, based on the key
func (e textEntry) extractISBN(data string) (isbn string, str string) {
	results := isbnRegExp.FindStringSubmatch(data)
	if isbnRegExp.MatchString(data) {
		isbn = strings.TrimSpace(results[1])
	}
	str = isbnRegExp.ReplaceAllString(data, "")
	str = strings.TrimSpace(str)

	return
}

// extractPageCount from the text, based on the `Pg.` key
func (e textEntry) extractPageCount(data string) (pages int, str string) {
	results := pageCountRegExp.FindStringSubmatch(data)
	if pageCountRegExp.MatchString(data) {
		count := strings.TrimSpace(results[1])
		pages, _ = strconv.Atoi(count)
	}
	str = pageCountRegExp.ReplaceAllString(data, "")
	str = strings.TrimSpace(str)

	return
}

// extractPageRange from the text, based on the `Pg.` key
func (e textEntry) extractPageRange(data string) (pages string, str string) {
	// do some clean up first
	if pageRangeCleanupRegExp.MatchString(data) {
		data = pageRangeCleanupRegExp.ReplaceAllString(data, ", Pg.")
	}

	results := pageRangeRegExp.FindStringSubmatch(data)
	if pageRangeRegExp.MatchString(data) {
		pages = strings.TrimSpace(results[1])
	}
	str = pageRangeRegExp.ReplaceAllString(data, "")
	str = strings.TrimSpace(str)

	return
}

// extractFirstPublishedDate from the text, based on a string match
// NOTE: must be done before other date processing
func (e textEntry) extractFirstPublishedDate(data string) (year int, str string) {
	results := firstPublishedRegExp.FindStringSubmatch(data)
	if firstPublishedRegExp.MatchString(data) {
		year, _ = strconv.Atoi(results[1])
	}
	str = firstPublishedRegExp.ReplaceAllString(data, "")
	str = strings.TrimSpace(str)

	return
}

// extractPublishedDate from the text, based on a string match
func (e textEntry) extractPublishedDate(data string) (date Date, str string) {
	results := publishedRegExp.FindStringSubmatch(data)
	if len(results) > 1 {
		dates := strings.Split(results[1], " ")

		// year will always be present, so just grab the last part
		year, _ := strconv.Atoi(dates[len(dates)-1])
		date.Year = year

		if len(dates) == 3 {
			day, _ := strconv.Atoi(dates[0])
			date.Day = day
			date.Month = e.monthAsNumber(dates[1])
		} else if len(dates) == 2 {
			date.Month = e.monthAsNumber(dates[0])
		}
	}
	str = publishedRegExp.ReplaceAllString(data, "")
	str = strings.TrimSpace(str)

	return
}

// extractVolumeNumber from the text, based on the `Vol.` key
func (e textEntry) extractVolumeNumber(data string) (vol string, str string) {
	results := volumeNumberRegExp.FindStringSubmatch(data)
	if volumeNumberRegExp.MatchString(data) {
		vol = strings.TrimSpace(results[1])
	}
	str = volumeNumberRegExp.ReplaceAllString(data, "")
	str = strings.TrimSpace(str)

	return
}

// extractIssueNumber from the text, based on the `Iss.` key
func (e textEntry) extractIssueNumber(data string) (issue string, str string) {
	results := issueNumberRegExp.FindStringSubmatch(data)
	if issueNumberRegExp.MatchString(data) {
		issue = strings.TrimSpace(results[1])
	}
	str = issueNumberRegExp.ReplaceAllString(data, "")
	str = strings.TrimSpace(str)

	return
}

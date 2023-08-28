package movie_test

import (
	"testing"

	"github.com/mrcook/imdblit/movie"
)

type BookTableData struct {
	text, title, author, isbn, note string
	pageCount, firstPublished       int
	volume, issue                   string
	publisher                       [4]string // name, city, country
	date                            [3]int    // year, month, day
}

type PublicationTableData struct {
	text, author, title, pages       string
	publication, volume, issue, issn string
	publisher                        [4]string // name, city, country
	date                             [3]int    // year, month, day
}

func TestMovieTitleDetails(t *testing.T) {
	mov := movie.Movie{}

	entry := `MOVI: Creature from the Black Lagoon (1954/XI) (TV)`
	movie.Unmarshall(entry, &mov)

	if mov.Title != "Creature from the Black Lagoon" {
		t.Fatalf("expected title to be extracted, got: '%s'", mov.Title)
	}
	if mov.Year != 1954 {
		t.Fatalf("expected year as 1954, got %d", mov.Year)
	}
	if mov.Month != 11 {
		t.Fatalf("expected month as 11, got %d", mov.Month)
	}
	if !mov.TV {
		t.Fatalf("expected TV to be true")
	}
}

func TestMultipleTitles(t *testing.T) {
	mov := movie.Movie{}

	// NOTE: this should never happen
	entry := `MOVI: Creature from the Black Lagoon (1954)
MOVI: The Creature from the Black Lagoon (1976)`
	movie.Unmarshall(entry, &mov)

	if mov.Title != "Creature from the Black Lagoon" {
		t.Fatalf("expected title to be extracted, got: '%s'", mov.Title)
	}
	if mov.Year != 1954 {
		t.Fatalf("expected year as 1954, got %d", mov.Year)
	}
}

func TestMovieQuotedTitles(t *testing.T) {
	mov := movie.Movie{}

	entry := `MOVI: "A Little Princess" (1973)`
	movie.Unmarshall(entry, &mov)

	if mov.Title != "A Little Princess" {
		t.Fatalf("expected title to be extracted, got: '%s'", mov.Title)
	}
}

func TestMovieSeriesInfo(t *testing.T) {
	mov := movie.Movie{}

	entry := `MOVI: "1,000 Places to See Before You Die" (2007) {Australia (#1.5)}`
	movie.Unmarshall(entry, &mov)

	if mov.Title != "1,000 Places to See Before You Die" {
		t.Fatalf("expected title to be extracted, got: '%s'", mov.Title)
	}
	if mov.Year != 2007 {
		t.Fatalf("expected year as 2007, got %d", mov.Year)
	}
	if mov.SeriesName != "Australia" {
		t.Fatalf("expected series to be Australia, got %s", mov.SeriesName)
	}
	if mov.SeriesNumber != 1 {
		t.Fatalf("expected series number of 1, got %d", mov.SeriesNumber)
	}
	if mov.EpisodeNumber != 5 {
		t.Fatalf("expected episode number of 5, got %d", mov.EpisodeNumber)
	}
}

func TestMovieSeriesName(t *testing.T) {
	mov := movie.Movie{}

	entry := `MOVI: "A Taste of Shakespeare" (1995) {King Lear}`
	movie.Unmarshall(entry, &mov)

	if mov.Title != "A Taste of Shakespeare" {
		t.Fatalf("expected title to be extracted, got: '%s'", mov.Title)
	}
	if mov.Year != 1995 {
		t.Fatalf("expected year as 1995, got %d", mov.Year)
	}
	if mov.SeriesName != "King Lear" {
		t.Fatalf("expected series to be King Lear, got %s", mov.SeriesName)
	}
}

func TestMovieSeriesEpisodes(t *testing.T) {
	mov := movie.Movie{}

	entry := `MOVI: "A Shared House" (2015) {(#1.4)}`
	movie.Unmarshall(entry, &mov)

	if mov.Title != "A Shared House" {
		t.Fatalf("expected title to be extracted, got: '%s'", mov.Title)
	}
	if mov.Year != 2015 {
		t.Fatalf("expected year as 2015, got %d", mov.Year)
	}
	if mov.SeriesNumber != 1 {
		t.Fatalf("expected series number of 1, got %d", mov.SeriesNumber)
	}
	if mov.EpisodeNumber != 4 {
		t.Fatalf("expected episode number of 4, got %d", mov.EpisodeNumber)
	}
}

func TestAdaptations(t *testing.T) {
	testItems := []BookTableData{
		{
			text:   `ADPT: Leslie Haskin. "Between Heaven & Ground Zero"`,
			author: "Leslie Haskin", title: "Between Heaven & Ground Zero",
		},
		{
			text:   `ADPT: Zona Gale. "Bill". 1927`,
			author: "Zona Gale", title: "Bill", date: [3]int{1927},
		},
		{
			text:   `ADPT: Gardner, Craig Shaw. "Back to the Future Part III (Novelization of the screenplay by Bob Gale)". (London, UK), Berkley Books, Berkley Publishing Group, 1 June 1990, Pg. 248, (BK), ISBN-10: 042512240X, (uncredited novel by co-screenwriter Doe: https: //www.example.com/doe)`,
			author: "Gardner, Craig Shaw", title: "Back to the Future Part III (Novelization of the screenplay by Bob Gale)",
			isbn: "042512240X", note: "uncredited novel by co-screenwriter Doe: https: //www.example.com/doe",
			pageCount: 248, publisher: [4]string{"Berkley Books", "London", "", "UK"}, date: [3]int{1990, 6, 1},
		},
		{
			text:   `ADPT: Siegfried Lenz. "Die Flut ist pünktlich". Hoffmann und Campe Verlag GmbH, (BK), (short story)`,
			author: "Siegfried Lenz", title: "Die Flut ist pünktlich",
			note: "short story", publisher: [4]string{"Hoffmann und Campe Verlag GmbH", "", "", ""},
		},
		{
			text:   `ADPT: Gary King. "Blind Rage: The Many Faces of Murder". In: "Amazon", Onyx Books (1995), (BK), (Novel), ISBN-13: 9780451405326`,
			author: "Gary King", title: "Blind Rage: The Many Faces of Murder",
			isbn: "9780451405326", publisher: [4]string{"Onyx Books"}, date: [3]int{1995},
		},
		{
			text:   `ADPT: Thomas Lee Howell. "Bully Boys". In: "Original novel" (USA), NONE, February 2010, Pg. 45, (BK)`,
			author: "Thomas Lee Howell", title: "Bully Boys",
			pageCount: 45, publisher: [4]string{"", "", "", "USA"}, date: [3]int{2010, 2, 0},
		},
		{
			text:   `ADPT: Francisco Ibanez Talavera,. "Clever & Smart". (Germany), ConPart-Verlag (Condor Verlag), Vol. 1, 1972, (MG)`,
			author: "Francisco Ibanez Talavera,", title: "Clever & Smart",
			publisher: [4]string{"ConPart-Verlag (Condor Verlag)", "", "", "Germany"}, date: [3]int{1972},
		},
	}

	for i, item := range testItems {
		mov := movie.Movie{}
		movie.Unmarshall(item.text, &mov)

		if len(mov.Adaptations) != 1 {
			t.Fatalf("(#%d) expected 1 adaptation to be found, got %d", i, len(mov.Adaptations))
		}
		a := mov.Adaptations[0]

		if a.Title != item.title {
			t.Errorf("(#%d) unexpected title, got '%s'", i, a.Title)
		}
		if a.Author != item.author {
			t.Errorf("(#%d) unexpected author, got '%s'", i, a.Author)
		}
		if a.PageCount != item.pageCount {
			t.Errorf("(#%d) unexpected page count, got '%d'", i, a.PageCount)
		}
		if a.ISBN != item.isbn {
			t.Errorf("(#%d) unexpected ISBN, got '%s'", i, a.ISBN)
		}
		if a.Note != item.note {
			t.Errorf("(#%d) unexpected note, got '%s'", i, a.Note)
		}
		if a.Publisher.Name != item.publisher[0] {
			t.Errorf("(#%d) unexpected publisher name, got '%s'", i, a.Publisher.Name)
		}
		if a.Publisher.City != item.publisher[1] {
			t.Errorf("(#%d) unexpected publisher city, got '%s'", i, a.Publisher.City)
		}
		if a.Publisher.State != item.publisher[2] {
			t.Errorf("(#%d) unexpected publisher state, got '%s'", i, a.Publisher.State)
		}
		if a.Publisher.Country != item.publisher[3] {
			t.Errorf("(#%d) unexpected publisher country, got '%s'", i, a.Publisher.Country)
		}
		if a.Date.Year != item.date[0] {
			t.Errorf("(#%d) unexpected date year, got '%d'", i, a.Date.Year)
		}
		if a.Date.Month != item.date[1] {
			t.Errorf("(#%d) unexpected date month, got '%d'", i, a.Date.Month)
		}
		if a.Date.Day != item.date[2] {
			t.Errorf("(#%d) unexpected date day, got '%d'", i, a.Date.Day)
		}
	}
}

func TestBooks(t *testing.T) {
	testItems := []BookTableData{
		{
			text:   `BOOK: Hickman, Roger. "Miklós Rózsa's Ben-Hur: A Film Score Guide". (Lanham, Maryland, USA), The Scarecrow Press, Inc., Vol. First, Iss. September, 2011, ISBN-13: 978-0-810-88100-4, (hb)`,
			author: "Hickman, Roger", title: "Miklós Rózsa's Ben-Hur: A Film Score Guide",
			isbn: "978-0-810-88100-4", publisher: [4]string{"The Scarecrow Press", "Lanham", "Maryland", "USA"}, date: [3]int{2011}, volume: "1", issue: "September",
		},
		{
			text:   `BOOK: Norlander, Emil. "Anderssonskans Kalle". (Stockholm), Ardor, Vol. 4th, Iss. 12, 1933, Pg. 155, (BK), (First published in 1901. Illustrated by O. A-n (sign. för Oskar Andersson.)`,
			author: "Norlander, Emil", title: "Anderssonskans Kalle",
			pageCount: 155, note: "Illustrated by O. A-n (sign. för Oskar Andersson.",
			publisher: [4]string{"Ardor", "", "", "Stockholm"}, date: [3]int{1933}, volume: "4", issue: "12", firstPublished: 1901,
		},
		{
			text:   `BOOK: Cunningham, Douglas A., editor. "The San Francisco of Alfred Hitchcock's Vertigo: Place, Pilgrimage, and Commemoration". Lanham, MD: The Scarecrow Press, 2011, ISBN-10: 0810881225`,
			author: "Cunningham, Douglas A., editor", title: "The San Francisco of Alfred Hitchcock's Vertigo: Place, Pilgrimage, and Commemoration",
			isbn: "0810881225", publisher: [4]string{"The Scarecrow Press", "Lanham", "", "MD"}, date: [3]int{2011},
		},

		// TODO: invalid parsing; inconsistent formatting compared to most other books. Is this common? Does it need handling?
		{
			text:   `BOOK: Mitchell, Charles. The Complete H.P. Lovecraft Filmography (Bibliographies and Indexes in the Performing Arts ). "New York: Greenwood Publishing Group". October 1, Iss. 26, 2001, Pg. 248, ISBN-10: 0313316414`,
			author: "Mitchell, Charles. The Complete H.P. Lovecraft Filmography (Bibliographies and Indexes in the Performing Arts )", title: "New York: Greenwood Publishing Group",
			pageCount: 248, isbn: "0313316414",
			publisher: [4]string{`October 1`}, date: [3]int{2001}, volume: "", issue: "26",
		},
	}

	for i, item := range testItems {
		mov := movie.Movie{}
		movie.Unmarshall(item.text, &mov)

		if len(mov.Books) != 1 {
			t.Fatalf("(#%d) expected 1 adaptation to be found, got %d", i, len(mov.Books))
		}
		b := mov.Books[0]

		if b.Title != item.title {
			t.Errorf("(#%d) unexpected title, got '%s'", i, b.Title)
		}
		if b.Author != item.author {
			t.Errorf("(#%d) unexpected author, got '%s'", i, b.Author)
		}
		if b.PageCount != item.pageCount {
			t.Errorf("(#%d) unexpected page count, got '%d'", i, b.PageCount)
		}
		if b.ISBN != item.isbn {
			t.Errorf("(#%d) unexpected ISBN, got '%s'", i, b.ISBN)
		}
		if b.Note != item.note {
			t.Errorf("(#%d) unexpected note, got '%s'", i, b.Note)
		}
		if b.Publisher.Name != item.publisher[0] {
			t.Errorf("(#%d) unexpected publisher name, got '%s'", i, b.Publisher.Name)
		}
		if b.Publisher.City != item.publisher[1] {
			t.Errorf("(#%d) unexpected publisher city, got '%s'", i, b.Publisher.City)
		}
		if b.Publisher.State != item.publisher[2] {
			t.Errorf("(#%d) unexpected publisher state, got '%s'", i, b.Publisher.State)
		}
		if b.Publisher.Country != item.publisher[3] {
			t.Errorf("(#%d) unexpected publisher country, got '%s'", i, b.Publisher.Country)
		}
		if b.Date.Year != item.date[0] {
			t.Errorf("(#%d) unexpected date year, got '%d'", i, b.Date.Year)
		}
		if b.Date.Month != item.date[1] {
			t.Errorf("(#%d) unexpected date month, got '%d'", i, b.Date.Month)
		}
		if b.Date.Day != item.date[2] {
			t.Errorf("(#%d) unexpected date day, got '%d'", i, b.Date.Day)
		}
		if b.Volume != item.volume {
			t.Errorf("(#%d) unexpected volume, got '%s'", i, b.Volume)
		}
		if b.Issue != item.issue {
			t.Errorf("(#%d) unexpected issue, got '%s'", i, b.Issue)
		}
		if b.FirstPublished != item.firstPublished {
			t.Errorf("(#%d) unexpected first published year, got '%d'", i, b.FirstPublished)
		}
	}
}

func TestNovels(t *testing.T) {
	testItems := []BookTableData{
		{
			text:   `NOVL: H.G. Wells. "The Food of the Gods"`,
			author: "H.G. Wells", title: "The Food of the Gods",
		},
		{
			text:   `NOVL: Wyndham, John. "Midwich Cuckoos, The". (London, England, UK), Michael Joseph Ltd., December 1957, Pg. 239, (BK), ISBN-10: 0345299116`,
			author: "Wyndham, John", title: "Midwich Cuckoos, The",
			pageCount: 239, isbn: "0345299116",
			publisher: [4]string{"Michael Joseph Ltd.", "London", "England", "UK"}, date: [3]int{1957, 12},
		},

		// How some ugly, inconsistent cases, are currently handled
		{
			text:   `NOVL: Véry, Pierre. "Goupi mains rouges, first published in 1937 by Editions Gallimard, Paris, France"`,
			author: "Véry, Pierre", title: "Goupi mains rouges, by Editions Gallimard, Paris, France",
			firstPublished: 1937,
		},
		{
			text:      `NOVL: (Etlar, Carit. Stormen på København den 11. Februar 1659 og Gøngehøvdingen)`,
			publisher: [4]string{`"Etlar`}, date: [3]int{1659},
		},
	}

	for i, item := range testItems {
		mov := movie.Movie{}
		movie.Unmarshall(item.text, &mov)

		if len(mov.Novels) != 1 {
			t.Fatalf("(#%d) expected 1 novel to be found, got %d", i, len(mov.Novels))
		}
		b := mov.Novels[0]

		if b.Title != item.title {
			t.Errorf("(#%d) unexpected title, got '%s'", i, b.Title)
		}
		if b.Author != item.author {
			t.Errorf("(#%d) unexpected author, got '%s'", i, b.Author)
		}
		if b.PageCount != item.pageCount {
			t.Errorf("(#%d) unexpected page count, got '%d'", i, b.PageCount)
		}
		if b.ISBN != item.isbn {
			t.Errorf("(#%d) unexpected ISBN, got '%s'", i, b.ISBN)
		}
		if b.Note != item.note {
			t.Errorf("(#%d) unexpected note, got '%s'", i, b.Note)
		}
		if b.Publisher.Name != item.publisher[0] {
			t.Errorf("(#%d) unexpected publisher name, got '%s'", i, b.Publisher.Name)
		}
		if b.Publisher.City != item.publisher[1] {
			t.Errorf("(#%d) unexpected publisher city, got '%s'", i, b.Publisher.City)
		}
		if b.Publisher.State != item.publisher[2] {
			t.Errorf("(#%d) unexpected publisher state, got '%s'", i, b.Publisher.State)
		}
		if b.Publisher.Country != item.publisher[3] {
			t.Errorf("(#%d) unexpected publisher country, got '%s'", i, b.Publisher.Country)
		}
		if b.Date.Year != item.date[0] {
			t.Errorf("(#%d) unexpected date year, got '%d'", i, b.Date.Year)
		}
		if b.Date.Month != item.date[1] {
			t.Errorf("(#%d) unexpected date month, got '%d'", i, b.Date.Month)
		}
		if b.Date.Day != item.date[2] {
			t.Errorf("(#%d) unexpected date day, got '%d'", i, b.Date.Day)
		}
		if b.FirstPublished != item.firstPublished {
			t.Errorf("(#%d) unexpected first published year, got '%d'", i, b.FirstPublished)
		}
	}
}

func TestCritiques(t *testing.T) {
	testItems := []PublicationTableData{
		{
			text:   `CRIT: Delmas, Jean. "Science-fiction, fantastique, cinéma d'hypothèse". In: "Jeune Cinéma" (Paris, France), Fédération Jean Vigo, Iss. # 13, March 1966, Pg. 1, (MG), ISSN: 0758-4202`,
			author: "Delmas, Jean", title: "Science-fiction, fantastique, cinéma d'hypothèse",
			publication: "Jeune Cinéma", pages: "1", issue: "13", issn: "0758-4202",
			publisher: [4]string{"Fédération Jean Vigo", "Paris", "", "France"}, date: [3]int{1966, 3},
		},
		{
			text:  `CRIT: "La chaussée des géants". In: "Cinémagazine" (Paris, France), Vol. 29, Iss. # 33, 1924, Pg. pgs. 261, (MG)`,
			title: "La chaussée des géants", publication: "Cinémagazine", volume: "29", issue: "33", pages: "261",
			publisher: [4]string{"", "Paris", "", "France"}, date: [3]int{1924},
		},
		{
			text:   `CRIT: Kochert, Mélanie. "Blanche-Neige". In: "L'Estrade" (Metz, Moselle, France), SAS Indola Presse, Iss. # 21, May 2012, Pg. 8, (MG), ISSN: 2109-4217`,
			author: "Kochert, Mélanie", title: "Blanche-Neige",
			publication: "L'Estrade", pages: "8", issue: "21", issn: "2109-4217",
			publisher: [4]string{"SAS Indola Presse", "Metz", "Moselle", "France"}, date: [3]int{2012, 5},
		},
		{text: `CRIT: "Title", Pg. 57-58`, title: "Title", pages: "57-58"},
		{text: `CRIT: "Title", Pg. c1+c10`, title: "Title", pages: "c1+c10"},
		{text: `CRIT: "Title", Pg. W20+W21`, title: "Title", pages: "W20+W21"},
		{text: `CRIT: "Title", Pg. Pages: 2-5`, title: "Title", pages: "2-5"},
		{text: `CRIT: "Title", Pg. pags. 14`, title: "Title", pages: "14"},
		{text: `CRIT: "Title", Pg. pg>. 20`, title: "Title", pages: "20"},
		{text: `CRIT: "Title", Pg. pg; 55`, title: "Title", pages: "55"},
		{text: `CRIT: "Title", Pg. pgd. 88`, title: "Title", pages: "88"},
		{text: `CRIT: "Title", Pg. pgs. 7, 22`, title: "Title", pages: "7, 22"},
		{text: `CRIT: "Title", Pg. pgs; 213 to 222, 224, 383`, title: "Title", pages: "213 to 222, 224, 383"},
	}

	for i, item := range testItems {
		mov := movie.Movie{}
		movie.Unmarshall(item.text, &mov)

		if len(mov.Critiques) != 1 {
			t.Fatalf("(#%d) expected 1 critique to be found, got %d", i, len(mov.Critiques))
		}
		b := mov.Critiques[0]

		if b.ArticleAuthor != item.author {
			t.Errorf("(#%d) unexpected author, got '%s'", i, b.ArticleAuthor)
		}
		if b.ArticleTitle != item.title {
			t.Errorf("(#%d) unexpected title, got '%s'", i, b.ArticleTitle)
		}
		if b.ArticlePages != item.pages {
			t.Errorf("(#%d) unexpected page range, got '%s'", i, b.ArticlePages)
		}
		if b.Publisher.Name != item.publisher[0] {
			t.Errorf("(#%d) unexpected publisher name, got '%s'", i, b.Publisher.Name)
		}
		if b.Publisher.City != item.publisher[1] {
			t.Errorf("(#%d) unexpected publisher city, got '%s'", i, b.Publisher.City)
		}
		if b.Publisher.State != item.publisher[2] {
			t.Errorf("(#%d) unexpected publisher state, got '%s'", i, b.Publisher.State)
		}
		if b.Publisher.Country != item.publisher[3] {
			t.Errorf("(#%d) unexpected publisher country, got '%s'", i, b.Publisher.Country)
		}
		if b.Date.Year != item.date[0] {
			t.Errorf("(#%d) unexpected date year, got '%d'", i, b.Date.Year)
		}
		if b.Date.Month != item.date[1] {
			t.Errorf("(#%d) unexpected date month, got '%d'", i, b.Date.Month)
		}
		if b.Date.Day != item.date[2] {
			t.Errorf("(#%d) unexpected date day, got '%d'", i, b.Date.Day)
		}
		if b.Volume != item.volume {
			t.Errorf("(#%d) unexpected volume, got '%s'", i, b.Volume)
		}
		if b.Issue != item.issue {
			t.Errorf("(#%d) unexpected issue, got '%s'", i, b.Issue)
		}
		if b.ISSN != item.issn {
			t.Errorf("(#%d) unexpected ISSN, got '%s'", i, b.ISSN)
		}
	}
}

// ALl entries are either Movie, Book, or Publication. Other tests in this file test for the parsing of each of those,
// therefore, this test is to check that all types are found an processed.
func TestExtractingAllEntries(t *testing.T) {
	// A real example movie entry. The only fake data is for the SCRP.
	var fullExample = `MOVI: Creature from the Black Lagoon (1954)

ADPT: Fearn, John Russell writing as Vargo Statten. "Creature from the Black Lagoon". (London, UK), Dragon Books, 1954, (BK)

BOOK: Tom Weaver. "The Creature Chronicles: Exploring the Black Lagoon Trilogy". (Jefferson NC), McFarland & Co., 2014, (BK), ISBN-13: 978-0-7864-9418-7
BOOK: Weaver, Tom. ""Creature from the Black Lagoon" Absecon NJ: MagicImage Books, ("making-of" book complete with script)". 1992, ISBN-10: 1882127307

CRIT: D.. "La mujer y el monstruo". In: "ABC" (Madrid), 8 June 1955, Pg. 60-61, (NP)
CRIT: Walker, John. In: "Total Film" (UK), July 1999, Pg. 106, (MG)

ESSY: "The Creature from the Black Lagoon". In: "The Economist" (UK), Vol. 8570, 8 March 2008, Pg. 89, (MG)

IVIW: "Fangoria" by: Tom Weaver, "The Creator from the Black Lagoon" (interview with screenwriter Harry Essex). (USA), Iss. 68, 1987
IVIW: "Starlog" (US), by: Tom Weaver, "The Producer from the Black Lagoon" (interview with producer William Alland). Iss. 218, September 1995
IVIW: "Starlog" by: Tom Weaver, "Creature Hunter" (interview with Richard Denning). (USA), Iss. 164, March 1991
IVIW: "Starlog" by: Tom Weaver, "Creature King" (interview with "the land creature" Ben Chapman). (USA), Iss. 180, July 1992
IVIW: "Starlog" by: Tom Weaver, "Creature Love" (interview with leading lady Julie Adams). (USA), Iss. 167, June 1991
IVIW: Robert Nott. "The Monster Mash-up: A New Book Dishes on the 'Creature' Films". In: "Pasatiempo" (Santa Fe NM), 31 October 2014, Pg. 34-36, (NP), (article on the "Creature" movies and the book "The Creature Chronicles")
IVIW: Tom Weaver. "Anatomy of a Mermaid: Ginger Stanley". In: "Classic Images" (Muscatine IA), Iss. 463, January 2014, Pg. 6-15, 70-81, (MG), (interview with Julie Adams' stunt double)
IVIW: Weaver, Tom. "Science Fiction and Fantasy Film Flashbacks". (Jefferson NC), McFarland & Co., 1998, Pg. 288-94, (BK), ISBN-10: 0786405643

NOVL: Dreadstone, Carl. "Creature from the Black Lagoon". (New York City, New York, USA), Berkley Medallion Books, Berkley Publishing Group, 27 June 1977, Pg. 194, (BK), ISBN-10: 042503464X, (Adapted from the Screenplay by Arthur Ross and Barry Essex)

OTHR: "Das Ungeheuer der schwarzen Lagune". In: "Illustrierter Film-Kurier" (Vienna, Austria), N. Freund - Metropol-Verlag, Iss. 2020, December 1954, Pg. 4

PROT: Weaver, Tom. "Behind the Black Lagoon". In: "Fangoria" (USA), Iss. 120, 1993, Pg. 14-19, 58-59

SCRP: Doe, John. "A Fake Screenplay.". In: "A Fake Screenplay. Screenplay by John Doe, Jane Doe, and Jenny Moe. Original story by John Doe and Jane Doe; © 1937" (New York City, New York, USA), Viking Press/Collection: A Viking Film Book, 27 July 1972, Pg. 276, (BK), ISBN-10: 0670258873

`

	mov := movie.Movie{}

	movie.Unmarshall(fullExample, &mov)

	if mov.Title != "Creature from the Black Lagoon" {
		t.Errorf("unexpected title extracted: '%s'", mov.Title)
	}
	if len(mov.Adaptations) != 1 {
		t.Errorf("expected 1 adaptation to be found, got %d", len(mov.Adaptations))
	}
	if len(mov.Books) != 2 {
		t.Errorf("expected 2 books to be found, got %d", len(mov.Books))
	}
	if len(mov.Critiques) != 2 {
		t.Errorf("expected 2 critiques to be found, got %d", len(mov.Critiques))
	}
	if len(mov.Essays) != 1 {
		t.Errorf("expected 1 essay to be found, got %d", len(mov.Essays))
	}
	if len(mov.Interviews) != 8 {
		t.Errorf("expected 8 interviews to be found, got %d", len(mov.Interviews))
	}
	if len(mov.Novels) != 1 {
		t.Errorf("expected 1 novel to be found, got %d", len(mov.Novels))
	}
	if len(mov.Others) != 1 {
		t.Errorf("expected 1 other to be found, got %d", len(mov.Others))
	}
	if len(mov.ProductionProtocols) != 1 {
		t.Errorf("expected 1 protocol to be found, got %d", len(mov.ProductionProtocols))
	}
	if len(mov.Screenplays) != 1 {
		t.Errorf("expected 1 screenplay to be found, got %d", len(mov.Screenplays))
	}
}

func TestIsAdaptation(t *testing.T) {
	testItems := [][]string{
		{`ADPT: H.P. Lovecraft. "The Thing on the Doorstep"`, "The Thing on the Doorstep", "Lovecraft, H.P."},
		{`BOOK: Dickens, Charles. "A Christmas Carol"`, "A Christmas Carol", "Dickens, Charles"},
		{`NOVL: Wells, H.G.. "The Food of the Gods"`, "Food of the Gods", "H.G. Wells"},
	}

	for i, item := range testItems {
		mov := movie.Movie{}
		movie.Unmarshall(item[0], &mov)

		IsAdaptation := mov.IsAdaptation(item[1], item[2])
		if !IsAdaptation {
			t.Errorf("(#%d) expected movie to be an adaptation", i)
		}
	}
}

func TestIsNotAdaptation(t *testing.T) {
	testItems := [][]string{
		{`ADPT: P.H. Lovecraft. "Nyarlathotep"`, "Nyarlathotep", "Lovecraft, H.P."},
		{`BOOK: Charles Dickens. "Oliver Twist"`, "Ollies Twist", "Dickens, Charles"},
		{`NOVL: H.G. Wells. "The War of the Words"`, "War of the Worlds", "H.G. Wells"},
	}

	for i, item := range testItems {
		mov := movie.Movie{}
		movie.Unmarshall(item[0], &mov)

		IsAdaptation := mov.IsAdaptation(item[1], item[2])
		if IsAdaptation {
			t.Errorf("(#%d) expected movie to not be an adaptation", i)
		}
	}
}

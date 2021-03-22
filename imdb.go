// Package imdblit provides the ability to search an IMDB literature.list database
// file, decoding each record to a set of Go structs.
package imdblit

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/mrcook/imdblit/movie"
)

// An IMDB reads and processes movie records values from an input stream.
type IMDB struct {
	r io.Reader

	createdOn    time.Time
	totalRecords int
}

// NewIMDB returns a new IMDB that reads from r.
func NewIMDB(r io.Reader) *IMDB {
	return &IMDB{r: bufio.NewReader(r)}
}

// DatabaseCreatedOn will contain the datetime that the DB file was generated,
// once the .list has been parsed.
func (db IMDB) DatabaseCreatedOn() time.Time {
	return db.createdOn
}

// TotalRecordCount will contain the total number of records found in the file,
// once the .list has been parsed.
func (db IMDB) TotalRecordCount() int {
	return db.totalRecords
}

// FindMovieAdaptations processes the DB and returns movies that are
// adaptations of the given book title/author.
//
// The search will only parse the book types: ADPT, BOOK, and NOVL, which will
// speed up the processing considerably.
func (db *IMDB) FindMovieAdaptations(title, author string) []movie.Movie {
	const recordDivider = "-------------------------------------------------------------------------------"

	var movies []movie.Movie

	scanner := bufio.NewScanner(db.r)
	if err := db.readDBHeader(scanner); err != nil {
		return movies
	}

	recordText := ""

	done := false
	for !done {
		done = !scanner.Scan()
		line := scanner.Text()
		if line == recordDivider {
			db.totalRecords++
		}

		// if this is the very first divider in the list, ignore it and read the next line
		if line == recordDivider && db.totalRecords == 1 {
			continue
		}

		// if a divider for the next movie entry is reached, check if the
		// current movie entry is valid, and add to list if so.
		if done || line == recordDivider {
			mov := movie.Movie{}
			movie.UnmarshallBooks(recordText, &mov)

			if mov.IsAdaptation(title, author) {
				movies = append(movies, mov)
			}

			recordText = ""
		} else {
			recordText += line + "\n"
		}
	}

	// quick reverse sort
	sort.Slice(movies, func(i, j int) bool {
		return movies[i].Year > movies[j].Year
	})

	return movies
}

// Reads the header section of the database file, reading the created on datetime,
// and setting the scanner pointer position to the start of the record entries.
func (db *IMDB) readDBHeader(scanner *bufio.Scanner) error {
	for {
		if !scanner.Scan() {
			return fmt.Errorf("scanner unexectedly stopped")
		}

		line := scanner.Text()

		if strings.Contains(line, " Date: ") {
			parts := strings.Split(line, " Date: ")
			date := strings.TrimSpace(parts[len(parts)-1])
			db.createdOn, _ = time.Parse("Mon Jan 2 15:04:05 2006", date)
		}

		if line == "LITERATURE LIST" {
			// read the next line, which is lots of ====, before returning
			if !scanner.Scan() {
				return fmt.Errorf("scanner unexectedly stopped")
			}
			return nil
		}
	}
}

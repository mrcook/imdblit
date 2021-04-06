package imdblit_test

import (
	"bytes"
	"testing"

	imdb "github.com/mrcook/imdblit"
)

var imdbText = `CRC: 0x527C5E79  File: literature.list  Date: Fri Dec 22 00:00:00 2017

Copyright 1991-2017 The Internet Movie Database Ltd. All rights reserved.

http://www.imdb.com

literature.list

2017-12-19

-----------------------------------------------------------------------------

LITERATURE LIST
===============
-------------------------------------------------------------------------------
MOVI: Dissonances (2003)

NOVL: Dixon, Stephen. "Interstate". (BK)

-------------------------------------------------------------------------------
MOVI: "The Last of the Mohicans" (1971)

NOVL: Cooper, James Fenimore. "Last of the Mohicans, The". Bantam New York, 1982, ISBN-10: 0553213296, (originally published 1826)

-------------------------------------------------------------------------------
MOVI: Mansfield Park (1983)

NOVL: Austen, Jane. "Mansfield Park"

-------------------------------------------------------------------------------
MOVI: Mansfield Park (2007) (TV)

NOVL: Austen, Jane. "Mansfield Park"

PROT: Glendinning, Lee. "New Generation Of Teenagers Prepare To Be Seduced With Rebirth Of Austen". In: "The Independent" (UK), Independent News & Media Ltd, Vol. 6345, 16 February 2007, Pg. 3, (NP)

-------------------------------------------------------------------------------
MOVI: Mansion of the Doomed (1976)

CRIT: Relizzo, Donald. In: "Demonique" (Los Angeles, California, USA), FantaCo Enterprises Inc., Vol. 4, 1983, Pg. 20, (MG)
`

func TestBasicProcessing(t *testing.T) {
	file := bytes.NewBuffer([]byte(imdbText)) // Fake a file read
	db := imdb.NewIMDB(file)
	_ = db.FindMovieAdaptations("Mansfield Park", "Jane Austen")

	if db.TotalRecordCount() != 5 {
		t.Errorf("expected a total of 5 movie records to have been processed, got %d", db.TotalRecordCount())
	}

	date := db.DatabaseCreatedOn().String()
	if date != "2017-12-22 00:00:00 +0000 UTC" {
		t.Errorf("unexpected datetime the DB was created, got %s", date)
	}
}

func TestMovieAdaptations(t *testing.T) {
	file := bytes.NewBuffer([]byte(imdbText)) // Fake a file read
	db := imdb.NewIMDB(file)

	movies := db.FindMovieAdaptations("Mansfield Park", "Jane Austen")

	if len(movies) != 2 {
		t.Fatalf("expected 2 movies to be found, got %d", len(movies))
	}

	mov := movies[0]

	if mov.Title != "Mansfield Park" {
		t.Errorf("unexpected movie title, got: %s", mov.Title)
	}
	if mov.Year != 2007 {
		t.Errorf("unexpected movie year, got: %d", mov.Year)
	}
}

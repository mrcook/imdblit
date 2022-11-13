# IMDB literature.list wrapper

A small library for processing an exported IMDB `literature.list` database file.

This library was created for my own ebook toolchain, so it contains features
only for that purpose; fetching movies that are adaptations of novels.


## Usage

A basic example might be:

```go
package main

import (
	"fmt"
	"os"

	"github.com/mrcook/imdblit"
)

func main() {
	imdbFile, err := os.Open("./literature.list")
	if err != nil {
		return
	}
	db := imdblit.NewIMDB(imdbFile)

	movies := db.FindMovieAdaptations("A Christmas Carol", "Charles Dickens")
	for _, movie := range movies {
		fmt.Printf("%d: %s\n", movie.Year, movie.Title)
	}
}
```


## LICENSE

Copyright (c) 2018-2022 Michael R. Cook. All rights reserved.

This work is licensed under the terms of the MIT license.
For a copy, see <https://opensource.org/licenses/MIT>.

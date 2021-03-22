# IMDB literature.list wrapper

A small library for processing an exported IMDB `literature.list` database file.

Currently, this library is useful for fetching movies that are adaptations of novels.

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
	db := imdblit.New(imdbFile)

	movies := db.MovieAdaptations("A Christmas Carol", "Charles Dickens")

	for _, movie := range movies {
		fmt.Printf("%d: %s\n", movie.Year, movie.Title)
	}
}
```


## LICENSE

Copyright (c) 2018-2021 Michael R. Cook. All rights reserved.

This work is licensed under the terms of the MIT license.
For a copy, see <https://opensource.org/licenses/MIT>.

# IMDBlit Changelog


## 0.9.0 (2023-08-28)

* Extract more data from movies (series/episode info), and books (volume/issue).
* Change the movie `Month` to be an `int`, and convert Roman numerals during processing.
* Add an `ExtractAll` function that returns all movies from the database file.
* `FindMovieAdaptations` now also returns an error. 


## 0.8.0 (2022-11-13)

The official IMDB `lterature.list` file is encoded as Windows 1252 so each
line in the file is now decoded before being processed.

That means there is now a dependency on the `golang.org/x/text` package.


## 0.7.0 (2021-04-05)

Major performance improvement during FindMovieAdaptations: ~3.5x faster!
On my laptop this drops from 11.5s to 3.2s. This was achieved by performing
all `regexp.MustCompile` once only rather than on each func call.


## 0.6.0 (2021-03-22)

Initial release (extracted from existing codebase).

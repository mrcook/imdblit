# IMDBlit Changelog


## 0.7.0 (2021-04-05)

Major performance improvement during FindMovieAdaptations: ~3.5x faster!
On my laptop this drops from 11.5s to 3.2s. This was achieved by performing
all `regexp.MustCompile` once only rather than on each func call.


## 0.6.0 (2021-03-22)

Initial release (extracted from existing codebase).

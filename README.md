ladder
======

Efficient solution of the all pairs shortest paths probem in Go. 
(http://en.wikipedia.org/wiki/Shortest_path_problem) 
At present the program is
not configured to print or find solutions but rather to sum the length of the shortest path
between every pair of connected words (nodes). This is a "fingerprint" of the APSP distance
matrix and serves to show that multiple solutions are deriving the same underlying result.

The code here is associated with 
the following golang-nuts discussion about approach, parallelism, efficiency, and timing.

https://groups.google.com/forum/#!topic/golang-nuts/ScFRRxqHTkY

The ladder program presumes "/usr/share/dict/words" exists, or, you can specify any number of word files on the command line. You can use it in these ways:

`go build`

_answer the question posted to golang-nuts about 4-letter words_

`./ladder -n 4`

_get detailed timing information_

`./ladder -t -n 4`

_get a variety of interesting facts by raising the verbosity level_


```
./ladder -v 1 -n 4
./ladder -v 2 -n 4
./ladder -v 3 -n 4
```

To go fast and benefit from parallelism, be sure to set GOMAXPROCS.

_use 4 cores and 4 SMT phantom CPUs to process all of /usr/share/dict/words_

```
export GOMAXPROCS=8
./ladder -v 1
```

Watch the resource usage graphs if you have tools to visuaize them.

There are many tests and benchmarks. To test:

```
go test -v
```

As it turns out, the tests were the most interesting part of this whole effort. Since I could not find any references to "known analytic values of sum of ASSP for paramateratized graphs" then I had to choose some simple graphs and derive the analytic solutions by hand. It was not too difficult generally, but some were tricky. Most challenging was the 2D lattice or grid graph, where an m x n lattice of nodes represnts an (m-1) x (n-1) grid of cells, like graph paper. The shortest paths between any pair of nodes on such a grid is one problem, and then the sum of all such paths is the second task. For the case of n x n square lattices, the sum of the lengths of the shortest paths between every pair of nodes is:

((2 * n * n) * (n - 1) * n * (n + 1)) / 3

there is a solution like this in `lattice_test.go` for a growing collection of simple graph types. _(Update: Now that I did all this, I see that the Weiner Index is essentially half the value I was looking for, and that the analytic solutions were already available. Ref:  http://mathworld.wolfram.com/WienerIndex.html)_

Benchmarks come in two varieties, V1 and V2. V1 are single-threaded while V2 use every processor allowed by GOMAXPROCS.

_choose your parallelism level_

```
export GOMAXPROCS=8
go test -v -bench=.
```

When I compile, I usually do it this way:

`go build -gcflags="-l -l -l -l -l -l -l -l -l -l"`

but I also compare with -B mode to measure bounds checking cost

`go build -gcflags="-l -l -l -l -l -l -l -l -l -l -B"`

You'll also want some sample word lists. The tests expect a subdirectory named words with files 
named webster-1, webster-2, ..., webster-9 representing 1-9 character words abstracted from the 
system word list stored (on my Mac) in `/usr/share/dict/words.` These files are here in the `words` 
subdirectory.

I was lazy when I built my files, using this bash script:
```
mtj$ cat BUILD 
grep "^.$" webster | tr "[A-Z]" "[a-z]" | sort -u > webster-1
grep "^..$" webster | tr "[A-Z]" "[a-z]" | sort -u > webster-2
grep "^...$" webster | tr "[A-Z]" "[a-z]" | sort -u > webster-3
grep "^....$" webster | tr "[A-Z]" "[a-z]" | sort -u > webster-4
grep "^.....$" webster | tr "[A-Z]" "[a-z]" | sort -u > webster-5
grep "^......$" webster | tr "[A-Z]" "[a-z]" | sort -u > webster-6
grep "^.......$" webster | tr "[A-Z]" "[a-z]" | sort -u > webster-7
grep "^........$" webster | tr "[A-Z]" "[a-z]" | sort -u > webster-8
grep "^.........$" webster | tr "[A-Z]" "[a-z]" | sort -u > webster-9
grep "^..........$" webster | tr "[A-Z]" "[a-z]" | sort -u > webster-10
grep "^...........$" webster | tr "[A-Z]" "[a-z]" | sort -u > webster-11
grep "^............$" webster | tr "[A-Z]" "[a-z]" | sort -u > webster-12
grep "^.............$" webster | tr "[A-Z]" "[a-z]" | sort -u > webster-13
grep "^..............$" webster | tr "[A-Z]" "[a-z]" | sort -u > webster-14
grep "^...............$" webster | tr "[A-Z]" "[a-z]" | sort -u > webster-15
grep "^................$" webster | tr "[A-Z]" "[a-z]" | sort -u > webster-16
grep "^.................$" webster | tr "[A-Z]" "[a-z]" | sort -u > webster-17
grep "^..................$" webster | tr "[A-Z]" "[a-z]" | sort -u > webster-18
grep "^...................$" webster | tr "[A-Z]" "[a-z]" | sort -u > webster-19
grep "^....................$" webster | tr "[A-Z]" "[a-z]" | sort -u > webster-20
```

however, now that the program exists, it would be fine to do the following to create these files:

```
go build 
mkdir words
./ladder -o words/webster-1 -n 1
./ladder -o words/webster-2 -n 2
./ladder -o words/webster-3 -n 3
:
```

...using whatever looping/shell structure makes sense to you. Easier is to have the program read all of /usr/share/dict/words and pull out only the n-letter words. This is specified with the `-n` option. 
For 7-letter words:

```
./ladder -n 7
```

Note that the file reading code is carefully written so that it can extract words from arbitrary text, such as Project Gutenberg books (http://www.gutenberg.org/) and UTF-8 encoded files in languages like Greek and Chinese. Chinese, Japanese, and Korean are interesting in the context of Doublets because they have a much higher word-to-word edge density than languages like English and Latin. This density is reported when the verbosity is raised above 1. (try `-v 2`)

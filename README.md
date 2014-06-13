ladder
======

Efficient solution of an all source shortest paths probem in Go. The code here is associated with 
the following golang-nuts discussion about approach, parallelism, efficiency, timing, and so on.

https://groups.google.com/forum/#!topic/golang-nuts/ScFRRxqHTkY

The ladder program presumes `/usr/share/dict/words`, so you can use it in these ways:

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

To go fast and understand the parallelism approach, be sure to set GOMAXPROCS.

_use 4 cores and 4 SMT phantom CPUs to process all of /usr/share/dict/words_
```
export GOMAXPROCS=8
./ladder -v 1
```

Be sure to watch the resource usage graphs if you have that.

There are benchmarks, but no tests (I've not been able to find any other data to test against that 
says what the sum of lengths of all shortest paths is for some word set.) To test:

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
system word list stored (on my Mac) in `/usr/share/dict/words.`

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

however, now that the program exists, it would be fine to do the following:
```
go build 
mkdir words
./ladder -o words/webster-1 -n 1
./ladder -o words/webster-2 -n 2
./ladder -o words/webster-3 -n 3
:
```

...using whatever looping/shell structure makes sense to you.

The files are there now.

package main

/*
 * ladder.go -- count shortest paths among words in Doublet puzzle
 */

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"
)

const WIDEST = 20

type Index uint32 // array index type for node and edge lists
const INFINITY = 1 << 30

type Indexes []Index

func (a Indexes) Len() int           { return len(a) }
func (a Indexes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Indexes) Less(i, j int) bool { return a[i] < a[j] }

// flag processor global variables
var wordsize, verbose int
var timing bool
var output string

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	flag.IntVar(&wordsize, "n", 0, "number of letters (zero means any)")
	flag.BoolVar(&timing, "t", false, "time and report progress")
	flag.IntVar(&verbose, "v", 0, "verbosity level")
	flag.StringVar(&output, "o", "", "output wordset to file")
}

var MaxProcs = runtime.GOMAXPROCS(0)
var NumCPU = runtime.NumCPU()

func main() {
	start := time.Now()
	meter := NewMeter()

	if verbose >= 1 {
		log.Printf("execution begins")
	}

	flag.Parse()

	// Read words from files named on the command line, or if none is given,
	// from "/usr/share/dict/words". Each word will be a node in our graph.
	filenames := flag.Args()
	if len(filenames) == 0 { // set default file name
		filenames = []string{"/usr/share/dict/words"}
	}

	word := readWords(filenames, wordsize)
	if timing {
		meter.SetWork(float64(len(word))) // unique words/sec
		log.Printf("%v read %v words", meter, len(word))
	}

	if output != "" {
		writeWords(word, output)
		if timing {
			meter.SetWork(0)
			log.Printf("%v wrote %v words", meter, len(word))
		}
	}

	// Determine which word-to-word transformations are allowed by the rules
	// of Lewis Carroll's Doublets puzzle. These are the graph's edges.
	pair := findPairs(word)
	if timing {
		sum := 0
		for _, p := range pair {
			sum += len(p)
		}
		sum /= 2
		meter.SetWork(float64(sum)) // pairs/sec
		log.Printf("%v find %v pairs", meter, sum)
	}

	// Determine graph's connected components. Each component is disconnected
	// from the others so searching and counting are independent sub-problems.
	component := findComponents(word, pair)
	if timing {
		meter.SetWork(float64(len(component))) // connected components/sec
		log.Printf("%v find %v components", meter, len(component))
	}

	// count one shortest length path between each word pair in each component
	var count, total int
	if false {
		count, total = sumAllSourcesShortestPathsV1(word, pair, component)
		if timing {
			meter.SetWork(float64(count)) // paths/sec
			log.Printf("%v find %v paths", meter, count)
		}
		fmt.Printf("%12d word pairs\n", count)
		fmt.Printf("%12d summed lengths of one shortest path per pair\n", total)
	}

	// count one shortest length path between each word pair in each component
	if true {
		count, total = sumAllSourcesShortestPathsV2(word, pair, component)
		if timing {
			meter.SetWork(float64(count)) // paths/sec
			log.Printf("%v find %v paths", meter, count)
		}
		fmt.Printf("%12d word pairs\n", count)
		fmt.Printf("%12d summed lengths of one shortest path per pair\n", total)
	}

	elapsed := float64(time.Now().Sub(start)) / 1e9
	if verbose >= 1 {
		log.Printf("execution ends, elapsed time = %.6f seconds", elapsed)
	}
	if timing {
		fmt.Printf("# %12.6f %12d %12d %2d %6d %v\n", elapsed, count, total, wordsize, len(word), filenames)
	}
}

// Read words from files and return a clean, ordered word list
func readWords(name []string, length int) []string {
	// interpret word length parameter
	var minLength, maxLength int
	switch {
	case length == 0:
		minLength = 1
		maxLength = WIDEST
	case 1 <= length && length <= WIDEST:
		minLength = length
		maxLength = length
	default:
		log.Fatalf("error: wordlength (%d) must be in 1..WIDEST (1..%d)", length, WIDEST)
	}

	names := len(name)
	if verbose >= 1 {
		switch {
		case length == 0:
			log.Printf("selecting words of length %d..%d from %d file%s\n", minLength, maxLength, names, plural(names))
		default:
			log.Printf("selecting %d-letter words from %d file%s\n", minLength, names, plural(names))
		}
	}

	// gather words from files using a map
	unique := make(map[string]struct{})
	var totalAdded, totalLong, totalRead int

	// custom splitter for scanner using "all non-letters except apostrophe"
	splitter := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		return scanCutset(data, atEOF, " \t\n\r0123456789`~!@#$%^&*()-—_=+[{]}\\|;:\",<.>/?") // allow apostrophe
	}

	for _, n := range name {
		var wordsAdded, wordsLong, wordsRead int

		// access named file
		file, err := os.Open(n)
		if err != nil {
			log.Printf("%v: %v", n, err)
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(splitter)
		for scanner.Scan() {
			word := scanner.Text()
			word = strings.Replace(word, "'", "", -1) // remove apostrophes ("o'clock" ==> "oclock")
			word = strings.Replace(word, "’", "", -1) // remove apostrophes ("o'clock" ==> "oclock")
			word = strings.ToLower(word)

			switch l := utf8.RuneCountInString(word); {
			case minLength <= l && l <= maxLength:
				unique[word] = struct{}{}
				wordsAdded++
			case l > WIDEST:
				wordsLong++
			}
			wordsRead++
		}
		totalAdded += wordsAdded
		totalLong += wordsLong
		totalRead += wordsRead

		if wordsLong > 0 {
			log.Printf("warning: skipped %6d words longer than WIDEST=%d runes in %s", wordsLong, WIDEST, n)
		}
		if verbose >= 1 {
			log.Printf("  added %7d of %7d words from file %s", wordsAdded, wordsRead, n)
		}
	}
	if len(unique) < 1 {
		log.Fatal("error: no words found")
	}

	// sort words into alphabetic order
	word := make([]string, len(unique))
	words := 0
	for s := range unique {
		word[words] = s
		words++
	}
	sort.Strings(word)

	if totalLong > 0 {
		log.Printf("skipped total of %6d words longer than WIDEST=%d runes", totalLong, WIDEST)
	}
	if verbose >= 1 {
		log.Printf("read total of %d unique words (skipped %d repeated words)", words, totalAdded-words)
	}
	if verbose >= 2 {
		fmt.Println()
		fmt.Printf("words:\n")
		printWords(word)
	}
	return word
}

// scanCutset is a version of strings.ScanWords that represents a split
// function for a Scanner that returns each non-cutset-separated string
// of text, with surrounding cutset characters deleted. It will never
// return an empty string. It cannot be used directly with scanner, but
// must be wrapped in another function to supply the cutset.
func scanCutset(data []byte, atEOF bool, cutset string) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !strings.ContainsRune(cutset, r) {
			break
		}
	}
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	// Scan until space, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if strings.ContainsRune(cutset, r) {
			return i + width, data[start:i], nil
		}
	}
	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data.
	return 0, nil, nil
}

func findPairs(word []string) []Indexes {
	widest := widestString(word)
	if widest > WIDEST {
		log.Fatalf("constant 'WIDEST=%v' is too small, must be >= %v for chosen words", WIDEST, widest)
	}

	// make a list of words sharing each "change one letter" word variation
	link := make(map[[WIDEST]rune]Indexes, 5*len(word))
	var key [WIDEST]rune
	for wn, w := range word {
		runes := []rune(w)
		for i, r := range runes {
			key[i] = r
		}
		for i, r := range runes {
			key[i] = '?'
			list := link[key]
			link[key] = append(list, Index(wn))
			key[i] = r
		}
		for i := range runes {
			key[i] = 0
		}
	}

	pair := make([]Indexes, len(word))
	for _, list := range link {
		for _, wn1 := range list {
			for _, wn2 := range list {
				if wn1 != wn2 {
					pair[wn1] = append(pair[wn1], wn2)
				}
			}
		}
	}
	for _, p := range pair {
		sort.Sort(p) // keep ordered by word number
	}

	if verbose >= 1 {
		total := 0
		for _, v := range pair {
			total += len(v)
		}
		total /= 2 // undirected edges go both ways
		log.Printf("found %d edge%s between words\n", total, plural(total))
	}
	if verbose >= 2 {
		fmt.Printf("linked words:\n")
		for wn, w := range word {
			if len(pair[wn]) > 0 { // not "aloof" as DEK would say
				fmt.Printf("%5d: %-6s -> ", wn, w)
				fmt.Printf("%s", word[pair[wn][0]])
				for i := 1; i < len(pair[wn]); i++ {
					fmt.Printf(", %s", word[pair[wn][i]])
				}
				fmt.Printf("\n")
			}
		}
		fmt.Println()
	}
	return pair

}

type Component struct {
	word  Indexes
	words int
}
type Components []Component

func (c Components) Len() int      { return len(c) }
func (c Components) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c Components) Less(i, j int) bool {
	// non-increasing order (largest first, so that initial allocations will generally suffice)
	return c[i].words > c[j].words || (c[i].words == c[j].words && c[i].word[0] < c[j].word[0])
}

// find connected components
func findComponents(word []string, pair []Indexes) Components {
	// every node has a corresponding component id
	id := make([]Index, len(word))
	for i := range id {
		id[i] = INFINITY // means "not part of any component"
	}
	ids := Index(0)
	queue := make([]Index, len(word)) // queue for BFS
	m := make(Indexes, 0, len(word))
	var component Components
	for w := range pair {
		// find every active word reachable from this one
		if id[w] == INFINITY {
			switch {
			case len(pair[w]) == 0: // disconnected word
				id[w] = ids
				m = append(m, Index(w))
			case len(pair[w]) == 1 && len(pair[pair[w][0]]) == 1: // disconnected pair
				w2 := pair[w][0]
				id[w] = ids
				id[w2] = ids
				m = append(m, Index(w), Index(w2))
			default:
				var head, tail int
				id[w] = ids            // id of current component
				queue[tail] = Index(w) // push starting word onto queue
				tail++
				m = append(m, Index(w))

				for head < tail { // breadth first traversal from w
					n := queue[head]
					head++
					for _, wn := range pair[n] {
						if id[wn] == INFINITY {
							id[wn] = ids
							queue[tail] = wn
							tail++
							m = append(m, wn)
						}
					}
				}
				sort.Sort(m) // order by word number
			}
			members := make(Indexes, len(m))
			copy(members, m)
			component = append(component, Component{members, len(members)})
			m = m[:0]
			ids++
		}
	}
	sort.Sort(component) // keep ordered by component size
	components := len(component)

	if verbose >= 1 {
		log.Printf("found %d connected component%s\n", components, plural(components))
	}
	if verbose >= 2 {
		log.Printf("component sizes:\n")
		for i := 0; i < components; {
			n := component[i].words
			j := i
			for j < components && component[j].words == n {
				j++
			}
			log.Printf("%6d with %4d word%s\n", j-i, n, plural(n))
			i = j
		}
	}
	if verbose >= 2 {
		log.Printf("component list:\n")
		for cn, c := range component {
			switch {
			case c.words < 13:
				fmt.Printf("%4d: size = %4d, words: ", cn, c.words)
				for _, wn := range c.word {
					fmt.Printf("%s ", word[wn])
				}
				fmt.Println()
				continue
			default:
				fmt.Printf("%4d: size = %4d, words:\n", cn, c.words)
				var list []string
				for _, wn := range c.word {
					list = append(list, word[wn])
				}
				printWords(list)
			}
		}
	}
	return component
}

func sumAllSourcesShortestPathsV1(word []string, pair []Indexes, component []Component) (int, int) {
	var totalPairs, totalPaths int
	if len(component) > 0 {
		distance := make([]Index, len(word))           // shortest distance to every node
		queue := make([]Index, len(component[0].word)) // queue of newly processed fringe nodes
		done := make([]bool, len(word))                // state: has node been processed?
		// note: special cases for 1 and 2 words are optional speedups
		for _, c := range component {
			switch c.words {
			case 1:
			case 2:
				totalPairs += 2
				totalPaths += 2
			default:
				totalPairs += c.words * (c.words - 1)
				for _, w := range c.word {
					totalPaths += ssspBFS(c.word, pair, Index(w), distance, queue, done)
				}
			}
		}
	}
	return totalPairs, totalPaths
}

// Compute Single Source Shortest Paths (SSSP) between a single source node and all
// other nodes of the connected component. Return distance array in d and the sum
// of the lengths of the shortest paths. Uses simple Breadth First Search which is
// an optimal foundation for SSSP/ASSP in unweighted adjacency-list graphs. BFS is
// friendly to parallelism since is has no impediment to concurrent evaluation.
func ssspBFS(word []Index, pair []Indexes, w Index, distance, queue []Index, done []bool) int {
	for _, wn := range word {
		distance[wn] = INFINITY // not known to be rechable (can be graph or component)
		done[wn] = false        // not yet discovered and processed by traversal
	}
	distance[w] = 0
	done[w] = true

	// push starting word onto queue
	var head, tail int
	queue[tail] = w
	tail++

	// breadth first traversal of graph rooted at w
	total := 0
	for head < tail { // while queue is not empty
		n := queue[head]
		head++
		d := distance[n] + 1
		for _, wn := range pair[n] {
			if !done[wn] {
				done[wn] = true
				distance[wn] = d
				queue[tail] = wn
				tail++
				total += int(d)
			}
		}
	}
	return total
}

// ideally the breakpoint would be determined by a test
const BREAKPOINT = 32 // switch from internal to external parallelism

func sumAllSourcesShortestPathsV2(word []string, pair []Indexes, component []Component) (int, int) {
	var i, j, totalPairs, totalPaths int
	components := len(component)

	// optimization -- skip parallel framework overhead when nothing but simple tasks
	if true {
		if components > 0 && component[0].words <= 16 {
			return sumAllSourcesShortestPathsV1(word, pair, component)
		}
	}

	// PHASE 1 -- solve large problems sequentially, using parallel workers fo each
	for ; i < components && component[i].words >= BREAKPOINT; i++ {
		c := component[i]
		totalPairs += c.words * (c.words - 1)
		totalPaths += ssspWordsParallel(word, pair, c)
	}

	// PHASE 2 -- solve medium problems in parallel, using a single worker for each
	for j = i; j < components && component[j].words > 2; j++ {
	}
	if i < j {
		pairCount, pathCount := ssspComponentsParallel(word, pair, component[i:j])
		totalPairs += pairCount
		totalPaths += pathCount
		i = j
	}

	// PHASE 3 -- solve small (nodes <= 2) problems directly
	for ; i < len(component); i++ {
		c := component[i]
		switch {
		case c.words == 1: // single aloof word with no solutions
		case c.words == 2: // single pair of words with two solutions (a->b and b->a) of length 1
			totalPairs += 2
			totalPaths += 2
		default:
			log.Fatal("internal error: phase 3 problem with more than 2 nodes")
		}
	}
	return totalPairs, totalPaths
}

func ssspComponentsParallel(word []string, pair []Indexes, component []Component) (int, int) {
	var totalPairs, totalPaths int
	tasks := make(chan Component) //, 1024)
	results := make(chan int)     //, 1024)

	// start workers
	workers := MaxProcs
	for k := 0; k < workers; k++ {
		go func(id int, in chan Component, out chan int, word []string, pair []Indexes) {
			distance := make(Indexes, len(word))
			done := make([]bool, len(word))
			var queue Indexes

			for c := range in {
				if cap(queue) < c.words {
					queue = make(Indexes, c.words) // should happen once in each worker
				} else {
					queue = queue[:c.words]
				}

				out <- ssspWordsSerial(word, pair, c, distance, queue, done)
			}
		}(k, tasks, results, word, pair)
	}

	// dispatch tasks to workers
	go func(out chan Component, component []Component) {
		for _, c := range component {
			out <- c
		}
		close(out)
	}(tasks, component)

	// harvest results from workers
	for _ = range component {
		totalPaths += <-results
	}
	close(results)

	// determine number of pairs for these components
	for _, c := range component {
		totalPairs += c.words * (c.words - 1)
	}
	return totalPairs, totalPaths
}

func ssspWordsParallel(word []string, pair []Indexes, c Component) int {
	tasks := make(chan Index) //, 1024)
	results := make(chan int) //, 1024)

	// start workers
	workers := minInt(c.words, MaxProcs)
	for i := 0; i < workers; i++ {
		go func(id int, in chan Index, out chan int, word []string, pair []Indexes, c Component) {
			distance := make(Indexes, len(word))
			queue := make(Indexes, len(c.word))
			done := make([]bool, len(word))
			for w := range in {

				// if cap(queue) < c.words {
				// 	queue = make(Indexes, c.words)
				// } else {
				// 	queue = queue[:c.words]
				// }

				out <- ssspBFS(c.word, pair, Index(w), distance, queue, done)
			}
		}(i, tasks, results, word, pair, c)
	}

	// start dispatcher
	go func(out chan Index, c Component) {
		for _, w := range c.word {
			out <- Index(w)
		}
		close(out)
	}(tasks, c)

	// harvest results from workers
	total := 0
	for _ = range c.word {
		total += <-results
	}
	close(results)
	return total
}

func ssspWordsSerial(word []string, pair []Indexes, c Component, distance, queue []Index, done []bool) int {
	total := 0
	for _, w := range c.word {
		total += ssspBFS(c.word, pair, w, distance, queue, done)
	}
	return total
}

//
// utility functions
//

func printWords(word []string) {
	printWidth := 120
	words := len(word)
	if words > 0 {
		k := 0
		widest := widestString(word)
		cols := printWidth / (widest + 1)
		for i := 0; i < (words+cols-1)/cols; i++ {
			for j := 0; j < cols && k < words; j++ {
				fmt.Printf("%-*s ", widest, word[k])
				k++
			}
			fmt.Println()
		}
		fmt.Println()
	}
}

func widestString(word []string) int {
	widest := 0
	for _, w := range word {
		widest = maxInt(widest, utf8.RuneCountInString(w))
	}
	return widest
}

func writeWords(word []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, s := range word {
		w.WriteString(s)
		w.WriteString("\n")
	}
	w.Flush()

	if verbose >= 1 {
		log.Printf("wrote %v words to file %v", len(word), filename)
	}
	return nil
}

func plural(n int) string {
	if n == 1 {
		return " "
	}
	return "s"
}

func maxInt(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func minIndex(a, b Index) Index {
	if a <= b {
		return a
	}
	return b
}

//
// simple performance meter: instantiate one and print it first on output lines.
// set work to display activities per second (bytes, pages, queries, words etc.)
// run this program with "-t" to see it in action. It is helpful to see elapsed
// times and also effective degree of parallelism across various parts of code.
//

type Meter struct {
	now     time.Time
	user    float64
	system  float64
	memory  uint64
	elapsed float64
	dUser   float64
	dSystem float64
	dMemory uint64
	work    float64
}

func NewMeter() *Meter {
	m := &Meter{now: time.Now()}
	m.user, m.system, m.memory = ProcessTimes()
	return m
}

func ProcessTimes() (user, system float64, size uint64) {
	var usage syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &usage); err != nil {
		fmt.Printf("Error: unable to gather resource usage data: %v\n", err)
	}
	user = float64(usage.Utime.Sec) + float64(usage.Utime.Usec)/1e6
	system = float64(usage.Stime.Sec) + float64(usage.Stime.Usec)/1e6
	size = uint64(uint32(usage.Maxrss))
	return
}

func (m *Meter) SetWork(work float64) {
	m.work = work
}

func (m *Meter) String() string {
	now := time.Now()
	user, system, memory := ProcessTimes()

	elapsed := float64((now.Sub(m.now))) / 1e9
	dUser := user - m.user
	dSystem := system - m.system
	dMemory := memory - m.memory

	var s string
	if m.work > 0 && dUser >= 0.0001 {
		s = fmt.Sprintf("%12.6f (%10.3f+%9.3f) %7.3f%% %9.3f MiB %9.0f/sec (%11.0f/sec)",
			elapsed, dUser, dSystem, 100*(dUser+dSystem)/elapsed, float64(dMemory)/(1024.0*1024.0), m.work/dUser, m.work/elapsed)
	} else {
		s = fmt.Sprintf("%12.6f (%10.3f+%9.3f) %7.3f%% %9.3f MiB                                ",
			elapsed, dUser, dSystem, 100*(dUser+dSystem)/elapsed, float64(dMemory)/(1024.0*1024.0))
	}

	m.now = now
	m.user = user
	m.system = system
	m.memory = memory
	m.elapsed = elapsed
	m.dUser = dUser
	m.dSystem = dSystem
	m.dMemory = dMemory

	return s
}

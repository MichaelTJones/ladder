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
	count, total := sumAllSourcesShortestPaths(word, pair, component)
	if timing {
		meter.SetWork(float64(count)) // paths/sec
		log.Printf("%v find %v paths", meter, count)
	}
	fmt.Printf("%12d word pairs\n", count)
	fmt.Printf("%12d summed lengths of one shortest path per pair\n", total)

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
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			word := scanner.Text()
			wordsRead++

			if true { // sanitize (half-hearted job does not look inside strings)
				word = strings.Trim(word, " \t0123456789.,?;:'\"[]{}=+~!@#$%^&*()\\|-")
				word = strings.ToLower(word)
			}

			switch l := utf8.RuneCountInString(word); {
			case minLength <= l && l <= maxLength:
				unique[word] = struct{}{}
				wordsAdded++
			case l > WIDEST:
				wordsLong++
			}
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

func findPairs(word []string) []Indexes {
	widest := widestString(word)
	if widest > WIDEST {
		log.Fatalf("constant 'WIDEST=%v' is too small, must be >= %v for chosen words", WIDEST, widest)
	}

	// create a map entry for unique "change one letter" word variants
	links := 0
	link := make(map[[WIDEST]rune]int, 128)
	for _, w := range word {
		var key [WIDEST]rune
		for i, c := range w {
			key[i] = c
		}
		for i := range w {
			t := key[i]
			key[i] = '?'
			if _, found := link[key]; !found {
				link[key] = links
				links++
			}
			key[i] = t
		}
	}
	if verbose >= 2 {
		log.Printf("%v links\n", len(link))
	}
	if verbose >= 3 {
		fmt.Printf("links:\n")
		ls := make([]string, len(link))
		ln := 0
		for l := range link {
			r := WIDEST - 1
			for r > 0 && l[r] == 0 {
				r--
			}
			ls[ln] = string(l[:r+1])
			ln++
		}
		sort.Strings(ls)
		printWords(ls)
	}

	// create lists of linked words with a common index in link map
	list := make([]Indexes, links)
	for wn, w := range word {
		var key [WIDEST]rune
		for i, c := range w {
			key[i] = c
		}
		for i := range w {
			t := key[i]
			key[i] = '?'
			j := link[key]
			list[j] = append(list[j], Index(wn))
			key[i] = t
		}
	}

	// create a customized list of linked words for each word
	pair := make([]Indexes, len(word))
	for wn, w := range word {
		var key [WIDEST]rune
		for i, c := range w {
			key[i] = c
		}
		for i := range w {
			t := key[i]
			key[i] = '?'
			j := link[key]
			for _, k := range list[j] {
				if k != Index(wn) { // do not include self in list
					pair[wn] = append(pair[wn], k)
				}
			}
			key[i] = t
		}
		sort.Sort(pair[wn]) // keep ordered by word number (optional)
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

func findComponents(word []string, pair []Indexes) Components {
	var component Components

	// every node starts in the "active and ready to join" state
	active := make([]bool, len(word))
	for i := range active {
		active[i] = true
	}

	// special cases for smallest (trivial) components (optional, but faster)
	if true {
		// find disconnected words as a special case
		for wn, a := range active {
			if a && len(pair[wn]) == 0 {
				component = append(component,
					Component{Indexes{Index(wn)}, 1})
				active[wn] = false
			}
		}

		// find disconnected two-word pairs as a special case
		for wn1, a := range active {
			if a && len(pair[wn1]) == 1 {
				wn2 := pair[wn1][0]
				if len(pair[wn2]) == 1 {
					if word[wn1] < word[wn2] {
						component = append(component,
							Component{Indexes{Index(wn1), Index(wn2)}, 2})
					}
					active[wn1] = false
					active[wn2] = false
				}
			}
		}
	}

	// general case: find remaining connected components
	last := make(Indexes, 0, 3*1024) // 11% faster when preallocated
	this := make(Indexes, 0, 3*1024)
	for wn, a := range active {
		// find every active word reachable from this one
		if a {
			var reach Indexes // grow slowly as it will persist
			reach = append(reach, Index(wn))
			last = append(last, Index(wn))
			active[wn] = false
			for len(last) > 0 {
				for _, l := range last {
					for _, p := range pair[l] {
						if active[p] {
							reach = append(reach, p)
							this = append(this, p)
							active[p] = false
						}
					}
				}
				last, this = this, last[:0]
			}
			sort.Sort(reach) // keep ordered by word number (optional, 1/2 speed)
			component = append(component, Component{reach, len(reach)})
			last = last[:0]
			this = this[:0]
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

// ideally the breakpoint would be determined by a test
const BREAKPOINT = 32 // switch from internal to external parallelism

func sumAllSourcesShortestPaths(word []string, pair []Indexes, component []Component) (int, int) {
	var i, j, totalPairs, totalPaths int
	components := len(component)

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
			d := make(Indexes, len(word))
			m := make([]uint8, len(word))
			var q Indexes

			for c := range in {
				if cap(q) < c.words {
					q = make(Indexes, c.words) // should happen once in each worker
				} else {
					q = q[:c.words]
				}
				out <- ssspWordsSerial(word, pair, c, d, q, m)
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
			d := make(Indexes, len(word))
			mode := make([]uint8, len(word))
			queue := make(Indexes, len(c.word))
			for w := range in {

				// if cap(queue) < c.words {
				// 	queue = make(Indexes, c.words)
				// } else {
				// 	queue = queue[:c.words]
				// }

				out <- ssspBFS(c.word, pair, Index(w), d, queue, mode)
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

func ssspWordsSerial(word []string, pair []Indexes, c Component, d Indexes, q Indexes, m []uint8) int {
	total := 0
	for _, w := range c.word {
		total += ssspBFS(c.word, pair, w, d, q, m)
	}
	return total
}

const (
	Undiscovered = iota
	// Discovered
	Processed
)

// Compute the Single Source Shortest Paths between one single source node and all
// other nodes in the connected component. Return distance array in d and the sum
// of the lengths of the shortest paths. Uses simple Breadth First Search which is
// an optimal foundation for SSSP/ASSP in unweighted adjacency-list graphs. BFS is
// friendly to parallelism since is has no impediment to concurrent evaluation.
func ssspBFS(word Indexes, pair []Indexes, w Index, d Indexes, queue Indexes, mode []uint8) int {
	for _, wn := range word {
		d[wn] = INFINITY        // not known to be rechable (only set distances in component)
		mode[wn] = Undiscovered // not yet been discovered
	}
	d[w] = 0
	mode[w] = Processed

	// push starting word onto queue
	var head, tail int
	queue[tail] = w
	tail++

	// breadth first traversal of graph rooted at w
	total := 0
	for distance := Index(1); head < tail; distance++ { // while queue is not empty
		last := tail
		for head < last { // for all nodes at this same distance
			n := queue[head]
			head++
			for _, wn := range pair[n] {
				if mode[wn] == Undiscovered {
					d[wn] = distance
					total += int(distance) // build sum incrementally
					mode[wn] = Processed
					queue[tail] = wn
					tail++
				}
			}
		}
		head = last
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
package main

import "testing"

// NOTE: the files mentioned below are 1 through 9 letter extracts from /usr/share/dict/words

//
// TESTS
//

// build linear graphs
//
//              sum of
//       pairs   paths   illustration
//    1:     0       0   O
//    2:     2       2   O --- O
//    3:     6       8   O --- O --- O
//    4:    12      20   O --- O --- O --- O
//    5:    20      40   O --- O --- O --- O --- O
//    6:    30      70   O --- O --- O --- O --- O --- O
//    7:    42     112   O --- O --- O --- O --- O --- O --- O
//                  :
//  200: 39800 2666600
//
func buildLinearGraph(n int) ([]string, []Indexes, []Component) {
	node := make([]string, n)
	// for i := range node {
	// 	node[i] = fmt.Sprintf("n%d", i)
	// }

	edge := make([]Indexes, n)
	for i := range edge {
		switch {
		case i == 0:
			edge[i] = []Index{Index(i + 1)}
		case i < n-1:
			edge[i] = []Index{Index(i - 1), Index(i + 1)}
		default:
			edge[i] = []Index{Index(i - 1)}
		}
	}

	member := make([]Index, n)
	for i := range member {
		member[i] = Index(i)
	}
	component := []Component{{member, n}}

	return node, edge, component
}

func TestLinearGraphsV1(t *testing.T) {
	for n := 1; n <= 200; n++ {
		pair1 := n * (n - 1)
		path1 := ((n - 1) * n * (n + 1)) / 3

		node, edge, component := buildLinearGraph(n)
		pair2, path2 := sumAllSourcesShortestPathsV1(node, edge, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

func TestLinearGraphsV2(t *testing.T) {
	for n := 1; n <= 200; n++ {
		pair1 := n * (n - 1)
		path1 := ((n - 1) * n * (n + 1)) / 3

		node, edge, component := buildLinearGraph(n)
		pair2, path2 := sumAllSourcesShortestPathsV2(node, edge, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

// fully-connected 4-node graph
//
//  O ----- O
//  | \   / |
//  |  \ /  |
//  |   X   |
//  |  / \  |
//  | /   \ |
//  O ----- O

// every node is connected to every other node, all shortest paths have length one
func buildFullyConnectedGraph(n int) ([]string, []Indexes, []Component) {
	node := make([]string, n)
	// for i := range node {
	// 	node[i] = fmt.Sprintf("n%d", i)
	// }

	edge := make([]Indexes, n)
	for i := range edge {
		edge[i] = make([]Index, n-1)
		for j := 0; j < i; j++ {
			edge[i][j] = Index(j) // link to every lower-numbered node
		}
		for j := i + 1; j < n; j++ {
			edge[i][j-1] = Index(j) // link to every higher-numbered node
		}
	}

	member := make([]Index, n)
	for i := range member {
		member[i] = Index(i)
	}
	component := []Component{{member, n}}

	return node, edge, component
}

func TestFullyConnectedGraphsV1(t *testing.T) {
	for n := 1; n <= 200; n++ {
		pair1 := n * (n - 1)
		path1 := n * (n - 1)

		node, edge, component := buildFullyConnectedGraph(n)
		pair2, path2 := sumAllSourcesShortestPathsV1(node, edge, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

func TestFullyConnectedGraphsV2(t *testing.T) {
	for n := 1; n <= 200; n++ {
		pair1 := n * (n - 1)
		path1 := n * (n - 1)

		node, edge, component := buildFullyConnectedGraph(n)
		pair2, path2 := sumAllSourcesShortestPathsV2(node, edge, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

// build lattice graph
//
//  3x4:
//    O --- O --- O --- O
//    |     |     |     |
//    O --- O --- O --- O
//    |     |     |     |
//    O --- O --- O --- O

func buildLatticeGraph(rows, cols int) ([]string, []Indexes, []Component) {
	n := rows * cols
	node := make([]string, n)
	// for r := 0; r < rows; r++ {
	// 	for c := 0; c < cols; c++ {
	// 		node[row*cols+col] = fmt.Sprintf("n[%d][%d]", row, col)
	// 	}
	// }

	edge := make([]Indexes, n)

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			i := r*cols + c

			if r > 0 {
				edge[i] = append(edge[i], Index((r-1)*cols+c)) // UP
			}
			if c > 0 {
				edge[i] = append(edge[i], Index(i-1)) // LEFT: i-1 == r*cols+(c-1)
			}
			if c < cols-1 {
				edge[i] = append(edge[i], Index(i+1)) // RIGHT: i+1 == r*cols+(c+1)
			}
			if r < rows-1 {
				edge[i] = append(edge[i], Index((r+1)*cols+c)) // DOWN
			}

		}
	}

	member := make([]Index, n)
	for i := range member {
		member[i] = Index(i)
	}
	component := []Component{{member, n}}

	return node, edge, component
}

//                 sum of
//       pairs      paths
//  2:      12         16
//  3:      72        144
//  4:     240        640
//  5:     600       2000
//  6:    1260       5040
//  7:    2352      10976
//  8:    4032      21504
//  :       :          :
// 37: 1872792   46195536
// 38: 2083692   52786864
// 39: 2311920   60109920
// 40: 2558400   68224000

func TestLatticeGraphsV1(t *testing.T) {
	// test paths on square n x n lattices
	for n := 2; n <= 40; n++ {
		pair1 := (n * n) * (n*n - 1)
		path1 := (2 * (n - 1) * n * n * n * (n + 1)) / 3

		node, edge, component := buildLatticeGraph(n, n)
		pair2, path2 := sumAllSourcesShortestPathsV1(node, edge, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

func TestLatticeGraphsV2(t *testing.T) {
	// test paths on square n x n lattices
	for n := 2; n <= 40; n++ {
		pair1 := (n * n) * (n*n - 1)
		path1 := (2 * (n - 1) * n * n * n * (n + 1)) / 3

		node, edge, component := buildLatticeGraph(n, n) // square n x n lattice
		pair2, path2 := sumAllSourcesShortestPathsV2(node, edge, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

//
// BENCHMARKS
//

func benchmarkReadWords(b *testing.B, f string, length int) {
	slice := []string{f}
	word := readWords(slice, length)
	if len(word) < 1 {
		b.Logf("no words in input file")
	}
	b.ResetTimer()

	for BN := 0; BN < b.N; BN++ {
		readWords(slice, length)
	}
}

func BenchmarkReadWords_webster1(b *testing.B) { benchmarkReadWords(b, "words/webster-1", 1) }
func BenchmarkReadWords_webster2(b *testing.B) { benchmarkReadWords(b, "words/webster-2", 2) }
func BenchmarkReadWords_webster3(b *testing.B) { benchmarkReadWords(b, "words/webster-3", 3) }
func BenchmarkReadWords_webster4(b *testing.B) { benchmarkReadWords(b, "words/webster-4", 4) }
func BenchmarkReadWords_webster5(b *testing.B) { benchmarkReadWords(b, "words/webster-5", 5) }
func BenchmarkReadWords_webster6(b *testing.B) { benchmarkReadWords(b, "words/webster-6", 6) }
func BenchmarkReadWords_webster7(b *testing.B) { benchmarkReadWords(b, "words/webster-7", 7) }
func BenchmarkReadWords_webster8(b *testing.B) { benchmarkReadWords(b, "words/webster-8", 8) }
func BenchmarkReadWords_webster9(b *testing.B) { benchmarkReadWords(b, "words/webster-9", 9) }

func benchmarkFindPairs(b *testing.B, f string, length int) {
	word := readWords([]string{f}, length)
	if len(word) < 1 {
		b.Logf("no words in input file")
	}
	b.ResetTimer()

	for BN := 0; BN < b.N; BN++ {
		findPairs(word)
	}
}

func BenchmarkFindPairs_webster1(b *testing.B) { benchmarkFindPairs(b, "words/webster-1", 1) }
func BenchmarkFindPairs_webster2(b *testing.B) { benchmarkFindPairs(b, "words/webster-2", 2) }
func BenchmarkFindPairs_webster3(b *testing.B) { benchmarkFindPairs(b, "words/webster-3", 3) }
func BenchmarkFindPairs_webster4(b *testing.B) { benchmarkFindPairs(b, "words/webster-4", 4) }
func BenchmarkFindPairs_webster5(b *testing.B) { benchmarkFindPairs(b, "words/webster-5", 5) }
func BenchmarkFindPairs_webster6(b *testing.B) { benchmarkFindPairs(b, "words/webster-6", 6) }
func BenchmarkFindPairs_webster7(b *testing.B) { benchmarkFindPairs(b, "words/webster-7", 7) }
func BenchmarkFindPairs_webster8(b *testing.B) { benchmarkFindPairs(b, "words/webster-8", 8) }
func BenchmarkFindPairs_webster9(b *testing.B) { benchmarkFindPairs(b, "words/webster-9", 9) }

func benchmarkFindComponents(b *testing.B, f string, length int) {
	word := readWords([]string{f}, length)
	if len(word) < 1 {
		b.Logf("no words in input file")
	}
	pair := findPairs(word)
	b.ResetTimer()

	for BN := 0; BN < b.N; BN++ {
		findComponents(word, pair)
	}
}

func BenchmarkFindComponents_webster1(b *testing.B) { benchmarkFindComponents(b, "words/webster-1", 1) }
func BenchmarkFindComponents_webster2(b *testing.B) { benchmarkFindComponents(b, "words/webster-2", 2) }
func BenchmarkFindComponents_webster3(b *testing.B) { benchmarkFindComponents(b, "words/webster-3", 3) }
func BenchmarkFindComponents_webster4(b *testing.B) { benchmarkFindComponents(b, "words/webster-4", 4) }
func BenchmarkFindComponents_webster5(b *testing.B) { benchmarkFindComponents(b, "words/webster-5", 5) }
func BenchmarkFindComponents_webster6(b *testing.B) { benchmarkFindComponents(b, "words/webster-6", 6) }
func BenchmarkFindComponents_webster7(b *testing.B) { benchmarkFindComponents(b, "words/webster-7", 7) }
func BenchmarkFindComponents_webster8(b *testing.B) { benchmarkFindComponents(b, "words/webster-8", 8) }
func BenchmarkFindComponents_webster9(b *testing.B) { benchmarkFindComponents(b, "words/webster-9", 9) }

func benchmarkSumASSPV1(b *testing.B, f string, length int) {
	word := readWords([]string{f}, length)
	if len(word) < 1 {
		b.Logf("no words in input file")
	}
	pair := findPairs(word)
	component := findComponents(word, pair)
	b.ResetTimer()

	for BN := 0; BN < b.N; BN++ {
		sumAllSourcesShortestPathsV1(word, pair, component)
	}
}

func BenchmarkSumASSPV1_webster1(b *testing.B) { benchmarkSumASSPV1(b, "words/webster-1", 1) }
func BenchmarkSumASSPV1_webster2(b *testing.B) { benchmarkSumASSPV1(b, "words/webster-2", 2) }
func BenchmarkSumASSPV1_webster3(b *testing.B) { benchmarkSumASSPV1(b, "words/webster-3", 3) }
func BenchmarkSumASSPV1_webster4(b *testing.B) { benchmarkSumASSPV1(b, "words/webster-4", 4) }
func BenchmarkSumASSPV1_webster5(b *testing.B) { benchmarkSumASSPV1(b, "words/webster-5", 5) }
func BenchmarkSumASSPV1_webster6(b *testing.B) { benchmarkSumASSPV1(b, "words/webster-6", 6) }
func BenchmarkSumASSPV1_webster7(b *testing.B) { benchmarkSumASSPV1(b, "words/webster-7", 7) }
func BenchmarkSumASSPV1_webster8(b *testing.B) { benchmarkSumASSPV1(b, "words/webster-8", 8) }
func BenchmarkSumASSPV1_webster9(b *testing.B) { benchmarkSumASSPV1(b, "words/webster-9", 9) }

// func BenchmarkSumASSPV1_webster10(b *testing.B) { benchmarkSumASSPV1(b, "words/webster-10", 10) }
// func BenchmarkSumASSPV1_webster11(b *testing.B) { benchmarkSumASSPV1(b, "words/webster-11", 11) }
// func BenchmarkSumASSPV1_webster12(b *testing.B) { benchmarkSumASSPV1(b, "words/webster-12", 12) }
// func BenchmarkSumASSPV1_webster13(b *testing.B) { benchmarkSumASSPV1(b, "words/webster-13", 13) }
// func BenchmarkSumASSPV1_webster14(b *testing.B) { benchmarkSumASSPV1(b, "words/webster-14", 14) }
// func BenchmarkSumASSPV1_webster15(b *testing.B) { benchmarkSumASSPV1(b, "words/webster-15", 15) }
// func BenchmarkSumASSPV1_webster16(b *testing.B) { benchmarkSumASSPV1(b, "words/webster-16", 16) }
// func BenchmarkSumASSPV1_webster17(b *testing.B) { benchmarkSumASSPV1(b, "words/webster-17", 17) }
// func BenchmarkSumASSPV1_webster18(b *testing.B) { benchmarkSumASSPV1(b, "words/webster-18", 18) }
// func BenchmarkSumASSPV1_webster19(b *testing.B) { benchmarkSumASSPV1(b, "words/webster-19", 19) }
// func BenchmarkSumASSPV1_webster20(b *testing.B) { benchmarkSumASSPV1(b, "words/webster-20", 20) }

func benchmarkSumASSPV2(b *testing.B, f string, length int) {
	word := readWords([]string{f}, length)
	if len(word) < 1 {
		b.Logf("no words in input file")
	}
	pair := findPairs(word)
	component := findComponents(word, pair)
	b.ResetTimer()

	for BN := 0; BN < b.N; BN++ {
		sumAllSourcesShortestPathsV2(word, pair, component)
	}
}

func BenchmarkSumASSPV2_webster1(b *testing.B) { benchmarkSumASSPV2(b, "words/webster-1", 1) }
func BenchmarkSumASSPV2_webster2(b *testing.B) { benchmarkSumASSPV2(b, "words/webster-2", 2) }
func BenchmarkSumASSPV2_webster3(b *testing.B) { benchmarkSumASSPV2(b, "words/webster-3", 3) }
func BenchmarkSumASSPV2_webster4(b *testing.B) { benchmarkSumASSPV2(b, "words/webster-4", 4) }
func BenchmarkSumASSPV2_webster5(b *testing.B) { benchmarkSumASSPV2(b, "words/webster-5", 5) }
func BenchmarkSumASSPV2_webster6(b *testing.B) { benchmarkSumASSPV2(b, "words/webster-6", 6) }
func BenchmarkSumASSPV2_webster7(b *testing.B) { benchmarkSumASSPV2(b, "words/webster-7", 7) }
func BenchmarkSumASSPV2_webster8(b *testing.B) { benchmarkSumASSPV2(b, "words/webster-8", 8) }
func BenchmarkSumASSPV2_webster9(b *testing.B) { benchmarkSumASSPV2(b, "words/webster-9", 9) }

// func BenchmarkSumASSPV2_webster10(b *testing.B) { benchmarkSumASSPV2(b, "words/webster-10", 10) }
// func BenchmarkSumASSPV2_webster11(b *testing.B) { benchmarkSumASSPV2(b, "words/webster-11", 11) }
// func BenchmarkSumASSPV2_webster12(b *testing.B) { benchmarkSumASSPV2(b, "words/webster-12", 12) }
// func BenchmarkSumASSPV2_webster13(b *testing.B) { benchmarkSumASSPV2(b, "words/webster-13", 13) }
// func BenchmarkSumASSPV2_webster14(b *testing.B) { benchmarkSumASSPV2(b, "words/webster-14", 14) }
// func BenchmarkSumASSPV2_webster15(b *testing.B) { benchmarkSumASSPV2(b, "words/webster-15", 15) }
// func BenchmarkSumASSPV2_webster16(b *testing.B) { benchmarkSumASSPV2(b, "words/webster-16", 16) }
// func BenchmarkSumASSPV2_webster17(b *testing.B) { benchmarkSumASSPV2(b, "words/webster-17", 17) }
// func BenchmarkSumASSPV2_webster18(b *testing.B) { benchmarkSumASSPV2(b, "words/webster-18", 18) }
// func BenchmarkSumASSPV2_webster19(b *testing.B) { benchmarkSumASSPV2(b, "words/webster-19", 19) }
// func BenchmarkSumASSPV2_webster20(b *testing.B) { benchmarkSumASSPV2(b, "words/webster-20", 20) }

package main

import "testing"

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

// 3 x 4 lattice graph
//
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
	// 		node[r*cols+c] = fmt.Sprintf("n[%d][%d]", r, c)
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
		path1 := ((2 * n * n) * (n - 1) * n * (n + 1)) / 3 // this was tricky to derive

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
		path1 := ((2 * n * n) * (n - 1) * n * (n + 1)) / 3 // this was tricky to derive

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

//   NOTE: Word iles mentioned below are 1 through 9 letter extracts from /usr/share/dict/words.
//         They are also in the "words" subdirectory.

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

func benchmarkSumLinearV1(b *testing.B, n int) {
	node, edge, component := buildLinearGraph(n)
	b.ResetTimer()

	for BN := 0; BN < b.N; BN++ {
		sumAllSourcesShortestPathsV1(node, edge, component)
	}
}

func BenchmarkSumLinearV1_2000(b *testing.B)  { benchmarkSumLinearV1(b, 2000) }
func BenchmarkSumLinearV1_4000(b *testing.B)  { benchmarkSumLinearV1(b, 4000) }
func BenchmarkSumLinearV1_6000(b *testing.B)  { benchmarkSumLinearV1(b, 6000) }
func BenchmarkSumLinearV1_8000(b *testing.B)  { benchmarkSumLinearV1(b, 8000) }
func BenchmarkSumLinearV1_10000(b *testing.B) { benchmarkSumLinearV1(b, 10000) }
func BenchmarkSumLinearV1_12000(b *testing.B) { benchmarkSumLinearV1(b, 12000) }
func BenchmarkSumLinearV1_14000(b *testing.B) { benchmarkSumLinearV1(b, 14000) }
func BenchmarkSumLinearV1_16000(b *testing.B) { benchmarkSumLinearV1(b, 16000) }
func BenchmarkSumLinearV1_18000(b *testing.B) { benchmarkSumLinearV1(b, 18000) }

func benchmarkSumLinearV2(b *testing.B, n int) {
	node, edge, component := buildLinearGraph(n)
	b.ResetTimer()

	for BN := 0; BN < b.N; BN++ {
		sumAllSourcesShortestPathsV2(node, edge, component)
	}
}

func BenchmarkSumLinearV2_2000(b *testing.B)  { benchmarkSumLinearV2(b, 2000) }
func BenchmarkSumLinearV2_4000(b *testing.B)  { benchmarkSumLinearV2(b, 4000) }
func BenchmarkSumLinearV2_6000(b *testing.B)  { benchmarkSumLinearV2(b, 6000) }
func BenchmarkSumLinearV2_8000(b *testing.B)  { benchmarkSumLinearV2(b, 8000) }
func BenchmarkSumLinearV2_10000(b *testing.B) { benchmarkSumLinearV2(b, 10000) }
func BenchmarkSumLinearV2_12000(b *testing.B) { benchmarkSumLinearV2(b, 12000) }
func BenchmarkSumLinearV2_14000(b *testing.B) { benchmarkSumLinearV2(b, 14000) }
func BenchmarkSumLinearV2_16000(b *testing.B) { benchmarkSumLinearV2(b, 16000) }
func BenchmarkSumLinearV2_18000(b *testing.B) { benchmarkSumLinearV2(b, 18000) }

func benchmarkSumFullyConnectedV1(b *testing.B, n int) {
	node, edge, component := buildFullyConnectedGraph(n)
	b.ResetTimer()

	for BN := 0; BN < b.N; BN++ {
		sumAllSourcesShortestPathsV1(node, edge, component)
	}
}

func BenchmarkSumFullyConnectedV1_100(b *testing.B)  { benchmarkSumFullyConnectedV1(b, 100) }
func BenchmarkSumFullyConnectedV1_200(b *testing.B)  { benchmarkSumFullyConnectedV1(b, 200) }
func BenchmarkSumFullyConnectedV1_300(b *testing.B)  { benchmarkSumFullyConnectedV1(b, 300) }
func BenchmarkSumFullyConnectedV1_400(b *testing.B)  { benchmarkSumFullyConnectedV1(b, 400) }
func BenchmarkSumFullyConnectedV1_500(b *testing.B)  { benchmarkSumFullyConnectedV1(b, 500) }
func BenchmarkSumFullyConnectedV1_600(b *testing.B)  { benchmarkSumFullyConnectedV1(b, 600) }
func BenchmarkSumFullyConnectedV1_700(b *testing.B)  { benchmarkSumFullyConnectedV1(b, 700) }
func BenchmarkSumFullyConnectedV1_800(b *testing.B)  { benchmarkSumFullyConnectedV1(b, 800) }
func BenchmarkSumFullyConnectedV1_900(b *testing.B)  { benchmarkSumFullyConnectedV1(b, 900) }
func BenchmarkSumFullyConnectedV1_1000(b *testing.B) { benchmarkSumFullyConnectedV1(b, 1000) }
func BenchmarkSumFullyConnectedV1_1100(b *testing.B) { benchmarkSumFullyConnectedV1(b, 1100) }

func benchmarkSumFullyConnectedV2(b *testing.B, n int) {
	node, edge, component := buildFullyConnectedGraph(n)
	b.ResetTimer()

	for BN := 0; BN < b.N; BN++ {
		sumAllSourcesShortestPathsV2(node, edge, component)
	}
}

func BenchmarkSumFullyConnectedV2_100(b *testing.B)  { benchmarkSumFullyConnectedV2(b, 100) }
func BenchmarkSumFullyConnectedV2_200(b *testing.B)  { benchmarkSumFullyConnectedV2(b, 200) }
func BenchmarkSumFullyConnectedV2_300(b *testing.B)  { benchmarkSumFullyConnectedV2(b, 300) }
func BenchmarkSumFullyConnectedV2_400(b *testing.B)  { benchmarkSumFullyConnectedV2(b, 400) }
func BenchmarkSumFullyConnectedV2_500(b *testing.B)  { benchmarkSumFullyConnectedV2(b, 500) }
func BenchmarkSumFullyConnectedV2_600(b *testing.B)  { benchmarkSumFullyConnectedV2(b, 600) }
func BenchmarkSumFullyConnectedV2_700(b *testing.B)  { benchmarkSumFullyConnectedV2(b, 700) }
func BenchmarkSumFullyConnectedV2_800(b *testing.B)  { benchmarkSumFullyConnectedV2(b, 800) }
func BenchmarkSumFullyConnectedV2_900(b *testing.B)  { benchmarkSumFullyConnectedV2(b, 900) }
func BenchmarkSumFullyConnectedV2_1000(b *testing.B) { benchmarkSumFullyConnectedV2(b, 1000) }
func BenchmarkSumFullyConnectedV2_1100(b *testing.B) { benchmarkSumFullyConnectedV2(b, 1100) }

func benchmarkSumLatticeV1(b *testing.B, n int) {
	node, edge, component := buildLatticeGraph(n, n)
	b.ResetTimer()

	for BN := 0; BN < b.N; BN++ {
		sumAllSourcesShortestPathsV1(node, edge, component)
	}
}

func BenchmarkSumLatticeV1_10(b *testing.B)  { benchmarkSumLatticeV1(b, 10) }
func BenchmarkSumLatticeV1_20(b *testing.B)  { benchmarkSumLatticeV1(b, 20) }
func BenchmarkSumLatticeV1_30(b *testing.B)  { benchmarkSumLatticeV1(b, 30) }
func BenchmarkSumLatticeV1_40(b *testing.B)  { benchmarkSumLatticeV1(b, 40) }
func BenchmarkSumLatticeV1_50(b *testing.B)  { benchmarkSumLatticeV1(b, 50) }
func BenchmarkSumLatticeV1_60(b *testing.B)  { benchmarkSumLatticeV1(b, 60) }
func BenchmarkSumLatticeV1_70(b *testing.B)  { benchmarkSumLatticeV1(b, 70) }
func BenchmarkSumLatticeV1_80(b *testing.B)  { benchmarkSumLatticeV1(b, 80) }
func BenchmarkSumLatticeV1_90(b *testing.B)  { benchmarkSumLatticeV1(b, 90) }
func BenchmarkSumLatticeV1_100(b *testing.B) { benchmarkSumLatticeV1(b, 100) }

func benchmarkSumLatticeV2(b *testing.B, n int) {
	node, edge, component := buildLatticeGraph(n, n)
	b.ResetTimer()

	for BN := 0; BN < b.N; BN++ {
		sumAllSourcesShortestPathsV2(node, edge, component)
	}
}

func BenchmarkSumLatticeV2_10(b *testing.B)  { benchmarkSumLatticeV2(b, 10) }
func BenchmarkSumLatticeV2_20(b *testing.B)  { benchmarkSumLatticeV2(b, 20) }
func BenchmarkSumLatticeV2_30(b *testing.B)  { benchmarkSumLatticeV2(b, 30) }
func BenchmarkSumLatticeV2_40(b *testing.B)  { benchmarkSumLatticeV2(b, 40) }
func BenchmarkSumLatticeV2_50(b *testing.B)  { benchmarkSumLatticeV2(b, 50) }
func BenchmarkSumLatticeV2_60(b *testing.B)  { benchmarkSumLatticeV2(b, 60) }
func BenchmarkSumLatticeV2_70(b *testing.B)  { benchmarkSumLatticeV2(b, 70) }
func BenchmarkSumLatticeV2_80(b *testing.B)  { benchmarkSumLatticeV2(b, 80) }
func BenchmarkSumLatticeV2_90(b *testing.B)  { benchmarkSumLatticeV2(b, 90) }
func BenchmarkSumLatticeV2_100(b *testing.B) { benchmarkSumLatticeV2(b, 100) }

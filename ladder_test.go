package main

import "testing"

func buildGraph(n int) ([]string, []Component) {
	node := make([]string, n)
	// for i := range node {
	// 	node[i] = fmt.Sprintf("n%d", i)
	// }

	element := make([]Index, n)
	for i := range element {
		element[i] = Index(i)
	}
	component := []Component{{element, len(element)}}

	return node, component
}

// Path graph, P_n
// http://en.wikipedia.org/wiki/Path_graph
//
//    P_1:    O
//    P_2:    O --- O
//    P_3:    O --- O --- O
//    P_4:    O --- O --- O --- O
//    P_5:    O --- O --- O --- O --- O

func buildPathGraph(n int) ([]string, []Indexes, []Component) {
	node, component := buildGraph(n)

	// build path graph adjacency lists
	a := make([]Indexes, n)
	for i := range a {
		switch {
		case i == 0:
			a[i] = []Index{Index(i + 1)}
		case i < n-1:
			a[i] = []Index{Index(i - 1), Index(i + 1)}
		case i == n-1:
			a[i] = []Index{Index(i - 1)}
		}
	}

	return node, a, component
}

// Path graph counts
//
//                sum of
//      n  pairs   paths
//      1:     0       0
//      2:     2       2
//      3:     6       8
//      4:    12      20
//      5:    20      40
//      6:    30      70
//      7:    42     112
//      8:    56     168
//             ⋮
//    197: 38612 2548392
//    198: 39006 2587398
//    199: 39402 2626800
//    200: 39800 2666600

func TestPathGraphV1(t *testing.T) {
	for n := 1; n <= 200; n++ {
		pair1 := n * (n - 1)         // number of pairs
		path1 := (n * (n*n - 1)) / 3 // sum of all pairs shortest path lengths

		node, a, component := buildPathGraph(n)
		pair2, path2 := sumAllSourcesShortestPathsV1(node, a, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

func TestPathGraphV2(t *testing.T) {
	for n := 1; n <= 200; n++ {
		pair1 := n * (n - 1)
		path1 := (n * (n*n - 1)) / 3 // sum of all pairs shortest path lengths

		node, a, component := buildPathGraph(n)
		pair2, path2 := sumAllSourcesShortestPathsV2(node, a, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

// Cycle graph, C_n
// http://en.wikipedia.org/wiki/Cycle_graph
//
// 4 node cycle graph
//
//    O ---- O
//    |      |
//    |      |
//    O ---- O

func buildCycleGraph(n int) ([]string, []Indexes, []Component) {
	node, component := buildGraph(n)

	// build cycle graph adjacency lists
	a := make([]Indexes, n)
	for i := range a {
		a[i] = []Index{
			Index((i - 1 + n) % n),
			Index((i + 1) % n),
		}
	}

	return node, a, component
}

// Cycle graph counts
//
//                sum of
//      n  pairs   paths
//      3:     6       6
//      4:    12      16
//      5:    20      30
//      6:    30      54
//      7:    42      84
//      8:    56     128
//             ⋮
//    197: 38612 1911294
//    198: 39006 1940598
//    199: 39402 1970100
//    200: 39800 2000000

func TestCycleGraphV1(t *testing.T) {
	for n := 3; n <= 200; n++ {
		pair1 := n * (n - 1)

		var path1 int // sum of all pairs shortest path lengths
		switch n & 1 {
		case 0:
			path1 = (n * n * n) / 4
		case 1:
			path1 = (n * (n*n - 1)) / 4
		}

		node, a, component := buildCycleGraph(n)
		pair2, path2 := sumAllSourcesShortestPathsV1(node, a, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

func TestCycleGraphV2(t *testing.T) {
	for n := 3; n <= 200; n++ {
		pair1 := n * (n - 1)

		var path1 int // sum of all pairs shortest path lengths
		switch n & 1 {
		case 0:
			path1 = (n * n * n) / 4
		case 1:
			path1 = (n * (n*n - 1)) / 4
		}

		node, a, component := buildCycleGraph(n)
		pair2, path2 := sumAllSourcesShortestPathsV2(node, a, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

// Complete graph, K_n
// http://en.wikipedia.org/wiki/Complete_graph
//
// 4 node complete graph
//
//  O ----- O
//  | \   / |
//  |  \ /  |
//  |   X   |
//  |  / \  |
//  | /   \ |
//  O ----- O

// every node is connected to every other node, all shortest paths have length one
func buildCompleteGraph(n int) ([]string, []Indexes, []Component) {
	node, component := buildGraph(n)

	// build complete graph adjacency lists
	a := make([]Indexes, n)
	for i := range a {
		a[i] = make([]Index, n-1)
		for j := 0; j < i; j++ {
			a[i][j] = Index(j) // link to every lower-numbered node
		}
		for j := i + 1; j < n; j++ {
			a[i][j-1] = Index(j) // link to every higher-numbered node
		}
	}

	return node, a, component
}

// Complete graph counts
//
//                sum of
//      n  pairs   paths
//      1:     0       0
//      2:     2       2
//      3:     6       6
//      4:    12      12
//      5:    20      20
//      6:    30      30
//      7:    42      42
//      8:    56      56
//             ⋮
//     97:  9312    9312
//     98:  9506    9506
//     99:  9702    9702
//    100:  9900    9900

func TestCompleteGraphV1(t *testing.T) {
	for n := 1; n <= 100; n++ {
		pair1 := n * (n - 1)
		path1 := n * (n - 1) // sum of all pairs shortest path lengths

		node, a, component := buildCompleteGraph(n)
		pair2, path2 := sumAllSourcesShortestPathsV1(node, a, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

func TestCompleteGraphV2(t *testing.T) {
	for n := 1; n <= 100; n++ {
		pair1 := n * (n - 1)
		path1 := n * (n - 1) // sum of all pairs shortest path lengths

		node, a, component := buildCompleteGraph(n)
		pair2, path2 := sumAllSourcesShortestPathsV2(node, a, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

// Star graph, S_n
// http://en.wikipedia.org/wiki/Star_graph
//
// 5 node star graph
//
//          O
//          |
//    O --- O --- O
//          |
//          O

func buildStarGraph(n int) ([]string, []Indexes, []Component) {
	node, component := buildGraph(n)

	// build star graph adjacency lists
	a := make([]Indexes, n)
	for i := range a {
		switch i {
		case n - 1:
			a[i] = make(Indexes, n-1)
			for j := 0; j < n-1; j++ {
				a[i][j] = Index(j)
			}
		default:
			a[i] = []Index{Index(n - 1)}
		}
	}

	return node, a, component
}

// Star graph counts
//
//                sum of
//      n  pairs   paths
//      1:     0       0
//      2:     2       2
//      3:     6       8
//      4:    12      18
//      5:    20      32
//      6:    30      50
//      7:    42      72
//      8:    56      98
//             ⋮
//    197: 38612   76832
//    198: 39006   77618
//    199: 39402   78408
//    200: 39800   79202

func TestStarGraphV1(t *testing.T) {
	for n := 1; n <= 200; n++ {
		pair1 := n * (n - 1)
		path1 := 2 * (n - 1) * (n - 1) // sum of all pairs shortest path lengths

		node, a, component := buildStarGraph(n)
		pair2, path2 := sumAllSourcesShortestPathsV1(node, a, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

func TestStarGraphV2(t *testing.T) {
	for n := 1; n <= 200; n++ {
		pair1 := n * (n - 1)
		path1 := 2 * (n - 1) * (n - 1) // sum of all pairs shortest path lengths

		node, a, component := buildStarGraph(n)
		pair2, path2 := sumAllSourcesShortestPathsV2(node, a, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

// Wheel graph, W_n
// http://en.wikipedia.org/wiki/Wheel_graph
//
// 5 node wheel graph
//
//  O ----- O
//  | \   / |
//  |  \ /  |
//  |   O   |
//  |  / \  |
//  | /   \ |
//  O ----- O

func buildWheelGraph(n int) ([]string, []Indexes, []Component) {
	node, component := buildGraph(n)

	// build wheel graph adjacency lists
	a := make([]Indexes, n)
	for i := range a {
		switch i {
		case n - 1:
			a[i] = make(Indexes, n-1)
			for j := 0; j < n-1; j++ {
				a[i][j] = Index(j)
			}
		default:
			m := n - 1
			a[i] = []Index{
				Index((i - 1 + m) % m),
				Index((i + 1) % m),
				Index(n - 1)}
		}
	}

	return node, a, component
}

// Wheel graph counts
//
//                sum of
//      n  pairs   paths
//      4:    12      12
//      5:    20      24
//      6:    30      40
//      7:    42      60
//      8:    56      84
//      9:    72     112
//     10:    90     144
//     11:   110     180
//             ⋮
//    297: 87912  174640
//    298: 88506  175824
//    299: 89102  177012
//    300: 89700  178204

func TestWheelGraphV1(t *testing.T) {
	for n := 4; n <= 300; n++ {
		pair1 := n * (n - 1)
		path1 := 2 * (n - 1) * (n - 2) // sum of all pairs shortest path lengths

		node, a, component := buildWheelGraph(n)
		pair2, path2 := sumAllSourcesShortestPathsV1(node, a, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

func TestWheelGraphV2(t *testing.T) {
	for n := 4; n <= 300; n++ {
		pair1 := n * (n - 1)
		path1 := 2 * (n - 1) * (n - 2) // sum of all pairs shortest path lengths

		node, a, component := buildWheelGraph(n)
		pair2, path2 := sumAllSourcesShortestPathsV2(node, a, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

// 2DGrid graph (grid graph, square grid graph)
// http://en.wikipedia.org/wiki/2DGrid_graph#Square_grid_graph
//
// 3 x 4 grid graph
//
//    O --- O --- O --- O
//    |     |     |     |
//    O --- O --- O --- O
//    |     |     |     |
//    O --- O --- O --- O

func build2DGridGraph(nx, ny int) ([]string, []Indexes, []Component) {
	n := nx * ny
	node, component := buildGraph(n)

	a := make([]Indexes, n)
	for y := 0; y < ny; y++ {
		for x := 0; x < nx; x++ {
			i := y*nx + x
			if x > 0 {
				a[i] = append(a[i], Index(i-1)) // -X
			}
			if x < nx-1 {
				a[i] = append(a[i], Index(i+1)) // +X
			}
			if y > 0 {
				a[i] = append(a[i], Index((y-1)*nx+x)) // -Y
			}
			if y < ny-1 {
				a[i] = append(a[i], Index((y+1)*nx+x)) // +Y
			}

		}
	}

	return node, a, component
}

// 2DGrid graph counts
//
//                   sum of
//           pairs    paths
//      1:       0        0
//      2:      12       16
//      3:      72      144
//      4:     240      640
//      5:     600     2000
//      6:    1260     5040
//      7:    2352    10976
//      8:    4032    21504
//             ⋮
//     37: 1872792 46195536
//     38: 2083692 52786864
//     39: 2311920 60109920
//     40: 2558400 68224000

func Test2DGridGraphV1(t *testing.T) {
	for n := 1; n <= 20; n++ {
		pair1 := (n * n) * (n*n - 1)
		// sum of all pairs shortest path lengths (tricky to derive with pencil and paper)
		path1 := (2 * n * n * n * (n*n - 1)) / 3

		node, a, component := build2DGridGraph(n, n)
		pair2, path2 := sumAllSourcesShortestPathsV1(node, a, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

func Test2DGridGraphV2(t *testing.T) {
	for n := 2; n <= 20; n++ {
		pair1 := (n * n) * (n*n - 1)
		path1 := (2 * n * n * n * (n*n - 1)) / 3

		node, a, component := build2DGridGraph(n, n) // square n x n lattice
		pair2, path2 := sumAllSourcesShortestPathsV2(node, a, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

// 3DGrid graph (grid graph, square grid graph)
// http://en.wikipedia.org/wiki/3DGrid_graph#Square_grid_graph

func build3DGridGraph(nx, ny, nz int) ([]string, []Indexes, []Component) {
	n := nx * ny * nz
	node, component := buildGraph(n)

	a := make([]Indexes, n)
	for z := 0; z < nz; z++ {
		for y := 0; y < ny; y++ {
			for x := 0; x < nx; x++ {
				i := z*nx*ny + y*nx + x
				if x > 0 {
					a[i] = append(a[i], Index(z*nx*ny+y*nx+(x-1))) // -X
				}
				if x < nx-1 {
					a[i] = append(a[i], Index(z*nx*ny+y*nx+(x+1))) // +X
				}
				if y > 0 {
					a[i] = append(a[i], Index(z*nx*ny+(y-1)*nx+x)) // -Y
				}
				if y < ny-1 {
					a[i] = append(a[i], Index(z*nx*ny+(y+1)*nx+x)) // +Y
				}
				if z > 0 {
					a[i] = append(a[i], Index((z-1)*nx*ny+y*nx+x)) // -Z
				}
				if z < nz-1 {
					a[i] = append(a[i], Index((z+1)*nx*ny+y*nx+x)) // +Z
				}
			}
		}
	}

	return node, a, component
}

// 3DGrid graph counts
//
//                   sum of
//           pairs    paths
//      1:       0        0
//      2:      12       16
//      3:      72      144
//      4:     240      640
//      5:     600     2000
//      6:    1260     5040
//      7:    2352    10976
//      8:    4032    21504
//             ⋮
//     37: 1872792 46195536
//     38: 2083692 52786864
//     39: 2311920 60109920
//     40: 2558400 68224000

func Test3DGridGraphV1(t *testing.T) {
	// test paths on cubic n x n x n lattices
	for n := 1; n <= 14; n++ {
		pair1 := (n * n * n) * (n*n*n - 1)
		path1 := n * n * n * n * n * (n*n - 1)

		node, a, component := build3DGridGraph(n, n, n)
		pair2, path2 := sumAllSourcesShortestPathsV1(node, a, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

func Test3DGridGraphV2(t *testing.T) {
	// test paths on cubic n x n x n lattices
	for n := 1; n <= 14; n++ {
		pair1 := (n * n * n) * (n*n*n - 1)
		path1 := n * n * n * n * n * (n*n - 1)

		node, a, component := build3DGridGraph(n, n, n)
		pair2, path2 := sumAllSourcesShortestPathsV2(node, a, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

// Complete bipartite graph, K_{m,n}
// http://en.wikipedia.org/wiki/Complete_bipartite_graph
//
// K_{5,3}:
//
//    O     O     O     O     O  5 in M group (not connected to each other)
//
//      \.................../    edge between each M and every N (M*N)
//       \                 /     edge between each N and every M (M*N)
//
//          O     O     O        3 in N group (not connected to each other)

func buildCompleteBipartiteGraph(m, n int) ([]string, []Indexes, []Component) {
	node, component := buildGraph(m + n)

	// build complete bipartite graph adjacency lists
	ms := make(Indexes, m)
	for i := range ms { // all the M nodes
		ms[i] = Index(i)
	}
	ns := make(Indexes, n)
	for i := range ns { // all the N nodes
		ns[i] = Index(m + i)
	}
	a := make([]Indexes, m+n)
	for i := 0; i < m; i++ {
		a[i] = ns
	}
	for i := m; i < m+n; i++ {
		a[i] = ms
	}

	return node, a, component
}

// Complete bipartate graph counts
//
//                  sum of
//  m  n    pairs    paths
// 20, 1:     420      800
// 20, 2:     462      844
// 20, 3:     506      892
// 20, 4:     552      944
// 20, 5:     600     1000
// 20, 6:     650     1060
// 20, 7:     702     1124
// 20, 8:     756     1192
// 20, 9:     812     1264
// 20,10:     870     1340

func TestBipartiteGraphV1(t *testing.T) {
	for m := 1; m <= 30; m++ {
		for n := 1; n <= 30; n++ {
			pair1 := (m + n) * (m + n - 1)
			path1 := 2 * (m*(m-1) + m*n + n*(n-1)) // sum of all pairs shortest path lengths

			node, a, component := buildCompleteBipartiteGraph(m, n)
			pair2, path2 := sumAllSourcesShortestPathsV1(node, a, component)

			if pair1 != pair2 || path1 != path2 {
				t.Errorf("%dx%d: expect (%d, %d) computed (%d, %d)", m, n, pair1, path1, pair2, path2)
			}
		}
	}
}

func TestBipartiteGraphV2(t *testing.T) {
	for m := 1; m <= 30; m++ {
		for n := 1; n <= 30; n++ {
			pair1 := (m + n) * (m + n - 1)
			path1 := 2 * (m*(m-1) + m*n + n*(n-1)) // sum of all pairs shortest path lengths

			node, a, component := buildCompleteBipartiteGraph(m, n)
			pair2, path2 := sumAllSourcesShortestPathsV2(node, a, component)

			if pair1 != pair2 || path1 != path2 {
				t.Errorf("%dx%d: expect (%d, %d) computed (%d, %d)", m, n, pair1, path1, pair2, path2)
			}
		}
	}
}

// Complete binary tree, T_n (full, same height, all leaves full)
// http://en.wikipedia.org/wiki/Binary_tree
func buildCompleteBinaryTree(height int) ([]string, []Indexes, []Component) {
	n := (1 << uint(height+1)) - 1 // 2**(height+1) - 1
	node, component := buildGraph(n)

	// build path graph adjacency lists
	a := make([]Indexes, n)
	for i := range a {
		switch {
		case i == 0: // root
			a[i] = []Index{
				Index(2*i + 1), // left child
				Index(2*i + 2), // right child
			}
		case i < n/2: // internal
			a[i] = []Index{
				Index(2*i + 1),     // left child
				Index(2*i + 2),     // right child
				Index((i - 1) / 2), // parent
			}
		default: // leaf
			a[i] = []Index{
				Index((i - 1) / 2), // parent
			}
		}
	}

	return node, a, component
}

// Complete binary tree counts
//
//                    sum of
// height   pairs      paths
//      1:       6         8
//      2:      42        96
//      3:     210       736
//      4:     930      4608
//      5:    3906     25728
//      6:   16002    133632
//      7:   64770    660992
//      8:  260610   3158016
//      9: 1045506  14706688
//     10: 4188162  67166208

func TestBinaryTreeV1(t *testing.T) {
	// test paths on square n x n lattic11es
	for n := 1; n <= 10; n++ {
		p := 1 << uint(n+1) // 2**(n+1)
		pair1 := (p - 1) * (p - 2)
		path1 := 2 * p * ((n-2)*(p+1) + 6)

		node, a, component := buildCompleteBinaryTree(n)
		pair2, path2 := sumAllSourcesShortestPathsV1(node, a, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

func TestBinaryTreeV2(t *testing.T) {
	// test paths on square n x n lattic11es
	for n := 1; n <= 10; n++ {
		p := 1 << uint(n+1) // 2**(n+1)
		pair1 := (p - 1) * (p - 2)
		path1 := 2 * p * ((n-2)*(p+1) + 6)

		node, a, component := buildCompleteBinaryTree(n)
		pair2, path2 := sumAllSourcesShortestPathsV2(node, a, component)

		if pair1 != pair2 || path1 != path2 {
			t.Errorf("%2d: expect (%d, %d) computed (%d, %d)", n, pair1, path1, pair2, path2)
		}
	}
}

//
// BENCHMARKS
//

//   NOTE: Word files mentioned below are 1 through 9 letter extracts from /usr/share/dict/words.
//         They are also in the "words" subdirectory.

func benchmarkReadWords(b *testing.B, f string, length int) {
	slice := []string{f}
	word, _ := readWords(slice, length)
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
	word, runes := readWords([]string{f}, length)
	if len(word) < 1 {
		b.Logf("no words in input file")
	}
	b.ResetTimer()

	for BN := 0; BN < b.N; BN++ {
		findPairs(word, runes)
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
	word, runes := readWords([]string{f}, length)
	if len(word) < 1 {
		b.Logf("no words in input file")
	}
	pair := findPairs(word, runes)
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
	word, runes := readWords([]string{f}, length)
	if len(word) < 1 {
		b.Logf("no words in input file")
	}
	pair := findPairs(word, runes)
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
	word, runes := readWords([]string{f}, length)
	if len(word) < 1 {
		b.Logf("no words in input file")
	}
	pair := findPairs(word, runes)
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

//
// Benchmark the solving of various simple graph types
//

func benchmarkSumPathV1(b *testing.B, n int) {
	node, a, component := buildPathGraph(n)
	b.ResetTimer()

	for BN := 0; BN < b.N; BN++ {
		sumAllSourcesShortestPathsV1(node, a, component)
	}
}

func BenchmarkSumPathV1_2000(b *testing.B)  { benchmarkSumPathV1(b, 2000) }
func BenchmarkSumPathV1_4000(b *testing.B)  { benchmarkSumPathV1(b, 4000) }
func BenchmarkSumPathV1_6000(b *testing.B)  { benchmarkSumPathV1(b, 6000) }
func BenchmarkSumPathV1_8000(b *testing.B)  { benchmarkSumPathV1(b, 8000) }
func BenchmarkSumPathV1_10000(b *testing.B) { benchmarkSumPathV1(b, 10000) }
func BenchmarkSumPathV1_12000(b *testing.B) { benchmarkSumPathV1(b, 12000) }
func BenchmarkSumPathV1_14000(b *testing.B) { benchmarkSumPathV1(b, 14000) }
func BenchmarkSumPathV1_16000(b *testing.B) { benchmarkSumPathV1(b, 16000) }
func BenchmarkSumPathV1_18000(b *testing.B) { benchmarkSumPathV1(b, 18000) }

func benchmarkSumPathV2(b *testing.B, n int) {
	node, a, component := buildPathGraph(n)
	b.ResetTimer()

	for BN := 0; BN < b.N; BN++ {
		sumAllSourcesShortestPathsV2(node, a, component)
	}
}

func BenchmarkSumPathV2_2000(b *testing.B)  { benchmarkSumPathV2(b, 2000) }
func BenchmarkSumPathV2_4000(b *testing.B)  { benchmarkSumPathV2(b, 4000) }
func BenchmarkSumPathV2_6000(b *testing.B)  { benchmarkSumPathV2(b, 6000) }
func BenchmarkSumPathV2_8000(b *testing.B)  { benchmarkSumPathV2(b, 8000) }
func BenchmarkSumPathV2_10000(b *testing.B) { benchmarkSumPathV2(b, 10000) }
func BenchmarkSumPathV2_12000(b *testing.B) { benchmarkSumPathV2(b, 12000) }
func BenchmarkSumPathV2_14000(b *testing.B) { benchmarkSumPathV2(b, 14000) }
func BenchmarkSumPathV2_16000(b *testing.B) { benchmarkSumPathV2(b, 16000) }
func BenchmarkSumPathV2_18000(b *testing.B) { benchmarkSumPathV2(b, 18000) }

// func benchmarkSumCycleV1(b *testing.B, n int) {
// 	node, a, component := buildCycleGraph(n)
// 	b.ResetTimer()

// 	for BN := 0; BN < b.N; BN++ {
// 		sumAllSourcesShortestPathsV1(node, a, component)
// 	}
// }

// func BenchmarkSumCycleV1_2000(b *testing.B)  { benchmarkSumCycleV1(b, 2000) }
// func BenchmarkSumCycleV1_4000(b *testing.B)  { benchmarkSumCycleV1(b, 4000) }
// func BenchmarkSumCycleV1_6000(b *testing.B)  { benchmarkSumCycleV1(b, 6000) }
// func BenchmarkSumCycleV1_8000(b *testing.B)  { benchmarkSumCycleV1(b, 8000) }
// func BenchmarkSumCycleV1_10000(b *testing.B) { benchmarkSumCycleV1(b, 10000) }
// func BenchmarkSumCycleV1_12000(b *testing.B) { benchmarkSumCycleV1(b, 12000) }
// func BenchmarkSumCycleV1_14000(b *testing.B) { benchmarkSumCycleV1(b, 14000) }
// func BenchmarkSumCycleV1_16000(b *testing.B) { benchmarkSumCycleV1(b, 16000) }
// func BenchmarkSumCycleV1_18000(b *testing.B) { benchmarkSumCycleV1(b, 18000) }

// func benchmarkSumCycleV2(b *testing.B, n int) {
// 	node, a, component := buildCycleGraph(n)
// 	b.ResetTimer()

// 	for BN := 0; BN < b.N; BN++ {
// 		sumAllSourcesShortestPathsV2(node, a, component)
// 	}
// }

// func BenchmarkSumCycleV2_2000(b *testing.B)  { benchmarkSumCycleV2(b, 2000) }
// func BenchmarkSumCycleV2_4000(b *testing.B)  { benchmarkSumCycleV2(b, 4000) }
// func BenchmarkSumCycleV2_6000(b *testing.B)  { benchmarkSumCycleV2(b, 6000) }
// func BenchmarkSumCycleV2_8000(b *testing.B)  { benchmarkSumCycleV2(b, 8000) }
// func BenchmarkSumCycleV2_10000(b *testing.B) { benchmarkSumCycleV2(b, 10000) }
// func BenchmarkSumCycleV2_12000(b *testing.B) { benchmarkSumCycleV2(b, 12000) }
// func BenchmarkSumCycleV2_14000(b *testing.B) { benchmarkSumCycleV2(b, 14000) }
// func BenchmarkSumCycleV2_16000(b *testing.B) { benchmarkSumCycleV2(b, 16000) }
// func BenchmarkSumCycleV2_18000(b *testing.B) { benchmarkSumCycleV2(b, 18000) }

// func benchmarkSumCompleteV1(b *testing.B, n int) {
// 	node, a, component := buildCompleteGraph(n)
// 	b.ResetTimer()

// 	for BN := 0; BN < b.N; BN++ {
// 		sumAllSourcesShortestPathsV1(node, a, component)
// 	}
// }

// func BenchmarkSumCompleteV1_100(b *testing.B)  { benchmarkSumCompleteV1(b, 100) }
// func BenchmarkSumCompleteV1_200(b *testing.B)  { benchmarkSumCompleteV1(b, 200) }
// func BenchmarkSumCompleteV1_300(b *testing.B)  { benchmarkSumCompleteV1(b, 300) }
// func BenchmarkSumCompleteV1_400(b *testing.B)  { benchmarkSumCompleteV1(b, 400) }
// func BenchmarkSumCompleteV1_500(b *testing.B)  { benchmarkSumCompleteV1(b, 500) }
// func BenchmarkSumCompleteV1_600(b *testing.B)  { benchmarkSumCompleteV1(b, 600) }
// func BenchmarkSumCompleteV1_700(b *testing.B)  { benchmarkSumCompleteV1(b, 700) }
// func BenchmarkSumCompleteV1_800(b *testing.B)  { benchmarkSumCompleteV1(b, 800) }
// func BenchmarkSumCompleteV1_900(b *testing.B)  { benchmarkSumCompleteV1(b, 900) }
// func BenchmarkSumCompleteV1_1000(b *testing.B) { benchmarkSumCompleteV1(b, 1000) }
// func BenchmarkSumCompleteV1_1100(b *testing.B) { benchmarkSumCompleteV1(b, 1100) }

// func benchmarkSumCompleteV2(b *testing.B, n int) {
// 	node, a, component := buildCompleteGraph(n)
// 	b.ResetTimer()

// 	for BN := 0; BN < b.N; BN++ {
// 		sumAllSourcesShortestPathsV2(node, a, component)
// 	}
// }

// func BenchmarkSumCompleteV2_100(b *testing.B)  { benchmarkSumCompleteV2(b, 100) }
// func BenchmarkSumCompleteV2_200(b *testing.B)  { benchmarkSumCompleteV2(b, 200) }
// func BenchmarkSumCompleteV2_300(b *testing.B)  { benchmarkSumCompleteV2(b, 300) }
// func BenchmarkSumCompleteV2_400(b *testing.B)  { benchmarkSumCompleteV2(b, 400) }
// func BenchmarkSumCompleteV2_500(b *testing.B)  { benchmarkSumCompleteV2(b, 500) }
// func BenchmarkSumCompleteV2_600(b *testing.B)  { benchmarkSumCompleteV2(b, 600) }
// func BenchmarkSumCompleteV2_700(b *testing.B)  { benchmarkSumCompleteV2(b, 700) }
// func BenchmarkSumCompleteV2_800(b *testing.B)  { benchmarkSumCompleteV2(b, 800) }
// func BenchmarkSumCompleteV2_900(b *testing.B)  { benchmarkSumCompleteV2(b, 900) }
// func BenchmarkSumCompleteV2_1000(b *testing.B) { benchmarkSumCompleteV2(b, 1000) }
// func BenchmarkSumCompleteV2_1100(b *testing.B) { benchmarkSumCompleteV2(b, 1100) }

// func benchmarkSumStarV1(b *testing.B, n int) {
// 	node, a, component := buildStarGraph(n)
// 	b.ResetTimer()

// 	for BN := 0; BN < b.N; BN++ {
// 		sumAllSourcesShortestPathsV1(node, a, component)
// 	}
// }

// func BenchmarkSumStarV1_2000(b *testing.B)  { benchmarkSumStarV1(b, 2000) }
// func BenchmarkSumStarV1_4000(b *testing.B)  { benchmarkSumStarV1(b, 4000) }
// func BenchmarkSumStarV1_6000(b *testing.B)  { benchmarkSumStarV1(b, 6000) }
// func BenchmarkSumStarV1_8000(b *testing.B)  { benchmarkSumStarV1(b, 8000) }
// func BenchmarkSumStarV1_10000(b *testing.B) { benchmarkSumStarV1(b, 10000) }
// func BenchmarkSumStarV1_12000(b *testing.B) { benchmarkSumStarV1(b, 12000) }
// func BenchmarkSumStarV1_14000(b *testing.B) { benchmarkSumStarV1(b, 14000) }
// func BenchmarkSumStarV1_16000(b *testing.B) { benchmarkSumStarV1(b, 16000) }
// func BenchmarkSumStarV1_18000(b *testing.B) { benchmarkSumStarV1(b, 18000) }

// func benchmarkSumStarV2(b *testing.B, n int) {
// 	node, a, component := buildStarGraph(n)
// 	b.ResetTimer()

// 	for BN := 0; BN < b.N; BN++ {
// 		sumAllSourcesShortestPathsV2(node, a, component)
// 	}
// }

// func BenchmarkSumStarV2_2000(b *testing.B)  { benchmarkSumStarV2(b, 2000) }
// func BenchmarkSumStarV2_4000(b *testing.B)  { benchmarkSumStarV2(b, 4000) }
// func BenchmarkSumStarV2_6000(b *testing.B)  { benchmarkSumStarV2(b, 6000) }
// func BenchmarkSumStarV2_8000(b *testing.B)  { benchmarkSumStarV2(b, 8000) }
// func BenchmarkSumStarV2_10000(b *testing.B) { benchmarkSumStarV2(b, 10000) }
// func BenchmarkSumStarV2_12000(b *testing.B) { benchmarkSumStarV2(b, 12000) }
// func BenchmarkSumStarV2_14000(b *testing.B) { benchmarkSumStarV2(b, 14000) }
// func BenchmarkSumStarV2_16000(b *testing.B) { benchmarkSumStarV2(b, 16000) }
// func BenchmarkSumStarV2_18000(b *testing.B) { benchmarkSumStarV2(b, 18000) }

// func benchmarkSumWheelV1(b *testing.B, n int) {
// 	node, a, component := buildWheelGraph(n)
// 	b.ResetTimer()

// 	for BN := 0; BN < b.N; BN++ {
// 		sumAllSourcesShortestPathsV1(node, a, component)
// 	}
// }

// func BenchmarkSumWheelV1_2000(b *testing.B)  { benchmarkSumWheelV1(b, 2000) }
// func BenchmarkSumWheelV1_4000(b *testing.B)  { benchmarkSumWheelV1(b, 4000) }
// func BenchmarkSumWheelV1_6000(b *testing.B)  { benchmarkSumWheelV1(b, 6000) }
// func BenchmarkSumWheelV1_8000(b *testing.B)  { benchmarkSumWheelV1(b, 8000) }
// func BenchmarkSumWheelV1_10000(b *testing.B) { benchmarkSumWheelV1(b, 10000) }
// func BenchmarkSumWheelV1_12000(b *testing.B) { benchmarkSumWheelV1(b, 12000) }
// func BenchmarkSumWheelV1_14000(b *testing.B) { benchmarkSumWheelV1(b, 14000) }
// func BenchmarkSumWheelV1_16000(b *testing.B) { benchmarkSumWheelV1(b, 16000) }
// func BenchmarkSumWheelV1_18000(b *testing.B) { benchmarkSumWheelV1(b, 18000) }

// func benchmarkSumWheelV2(b *testing.B, n int) {
// 	node, a, component := buildWheelGraph(n)
// 	b.ResetTimer()

// 	for BN := 0; BN < b.N; BN++ {
// 		sumAllSourcesShortestPathsV2(node, a, component)
// 	}
// }

// func BenchmarkSumWheelV2_2000(b *testing.B)  { benchmarkSumWheelV2(b, 2000) }
// func BenchmarkSumWheelV2_4000(b *testing.B)  { benchmarkSumWheelV2(b, 4000) }
// func BenchmarkSumWheelV2_6000(b *testing.B)  { benchmarkSumWheelV2(b, 6000) }
// func BenchmarkSumWheelV2_8000(b *testing.B)  { benchmarkSumWheelV2(b, 8000) }
// func BenchmarkSumWheelV2_10000(b *testing.B) { benchmarkSumWheelV2(b, 10000) }
// func BenchmarkSumWheelV2_12000(b *testing.B) { benchmarkSumWheelV2(b, 12000) }
// func BenchmarkSumWheelV2_14000(b *testing.B) { benchmarkSumWheelV2(b, 14000) }
// func BenchmarkSumWheelV2_16000(b *testing.B) { benchmarkSumWheelV2(b, 16000) }
// func BenchmarkSumWheelV2_18000(b *testing.B) { benchmarkSumWheelV2(b, 18000) }

func benchmarkSum2DGridV1(b *testing.B, n int) {
	node, a, component := build2DGridGraph(n, n)
	b.ResetTimer()

	for BN := 0; BN < b.N; BN++ {
		sumAllSourcesShortestPathsV1(node, a, component)
	}
}

func BenchmarkSum2DGridV1_10(b *testing.B)  { benchmarkSum2DGridV1(b, 10) }
func BenchmarkSum2DGridV1_20(b *testing.B)  { benchmarkSum2DGridV1(b, 20) }
func BenchmarkSum2DGridV1_30(b *testing.B)  { benchmarkSum2DGridV1(b, 30) }
func BenchmarkSum2DGridV1_40(b *testing.B)  { benchmarkSum2DGridV1(b, 40) }
func BenchmarkSum2DGridV1_50(b *testing.B)  { benchmarkSum2DGridV1(b, 50) }
func BenchmarkSum2DGridV1_60(b *testing.B)  { benchmarkSum2DGridV1(b, 60) }
func BenchmarkSum2DGridV1_70(b *testing.B)  { benchmarkSum2DGridV1(b, 70) }
func BenchmarkSum2DGridV1_80(b *testing.B)  { benchmarkSum2DGridV1(b, 80) }
func BenchmarkSum2DGridV1_90(b *testing.B)  { benchmarkSum2DGridV1(b, 90) }
func BenchmarkSum2DGridV1_100(b *testing.B) { benchmarkSum2DGridV1(b, 100) }

func benchmarkSum2DGridV2(b *testing.B, n int) {
	node, a, component := build2DGridGraph(n, n)
	b.ResetTimer()

	for BN := 0; BN < b.N; BN++ {
		sumAllSourcesShortestPathsV2(node, a, component)
	}
}

func BenchmarkSum2DGridV2_10(b *testing.B)  { benchmarkSum2DGridV2(b, 10) }
func BenchmarkSum2DGridV2_20(b *testing.B)  { benchmarkSum2DGridV2(b, 20) }
func BenchmarkSum2DGridV2_30(b *testing.B)  { benchmarkSum2DGridV2(b, 30) }
func BenchmarkSum2DGridV2_40(b *testing.B)  { benchmarkSum2DGridV2(b, 40) }
func BenchmarkSum2DGridV2_50(b *testing.B)  { benchmarkSum2DGridV2(b, 50) }
func BenchmarkSum2DGridV2_60(b *testing.B)  { benchmarkSum2DGridV2(b, 60) }
func BenchmarkSum2DGridV2_70(b *testing.B)  { benchmarkSum2DGridV2(b, 70) }
func BenchmarkSum2DGridV2_80(b *testing.B)  { benchmarkSum2DGridV2(b, 80) }
func BenchmarkSum2DGridV2_90(b *testing.B)  { benchmarkSum2DGridV2(b, 90) }
func BenchmarkSum2DGridV2_100(b *testing.B) { benchmarkSum2DGridV2(b, 100) }

// func benchmarkSum3DGridV1(b *testing.B, n int) {
// 	node, a, component := build3DGridGraph(n, n, n)
// 	b.ResetTimer()

// 	for BN := 0; BN < b.N; BN++ {
// 		sumAllSourcesShortestPathsV1(node, a, component)
// 	}
// }

// func BenchmarkSum3DGridV1_2(b *testing.B)  { benchmarkSum3DGridV1(b, 2) }
// func BenchmarkSum3DGridV1_4(b *testing.B)  { benchmarkSum3DGridV1(b, 4) }
// func BenchmarkSum3DGridV1_6(b *testing.B)  { benchmarkSum3DGridV1(b, 6) }
// func BenchmarkSum3DGridV1_8(b *testing.B)  { benchmarkSum3DGridV1(b, 8) }
// func BenchmarkSum3DGridV1_10(b *testing.B) { benchmarkSum3DGridV1(b, 10) }
// func BenchmarkSum3DGridV1_12(b *testing.B) { benchmarkSum3DGridV1(b, 12) }
// func BenchmarkSum3DGridV1_14(b *testing.B) { benchmarkSum3DGridV1(b, 14) }
// func BenchmarkSum3DGridV1_16(b *testing.B) { benchmarkSum3DGridV1(b, 16) }
// func BenchmarkSum3DGridV1_18(b *testing.B) { benchmarkSum3DGridV1(b, 18) }
// func BenchmarkSum3DGridV1_20(b *testing.B) { benchmarkSum3DGridV1(b, 20) }

// func benchmarkSum3DGridV2(b *testing.B, n int) {
// 	node, a, component := build3DGridGraph(n, n, n)
// 	b.ResetTimer()

// 	for BN := 0; BN < b.N; BN++ {
// 		sumAllSourcesShortestPathsV2(node, a, component)
// 	}
// }

// func BenchmarkSum3DGridV2_2(b *testing.B)  { benchmarkSum3DGridV2(b, 2) }
// func BenchmarkSum3DGridV2_4(b *testing.B)  { benchmarkSum3DGridV2(b, 4) }
// func BenchmarkSum3DGridV2_6(b *testing.B)  { benchmarkSum3DGridV2(b, 6) }
// func BenchmarkSum3DGridV2_8(b *testing.B)  { benchmarkSum3DGridV2(b, 8) }
// func BenchmarkSum3DGridV2_10(b *testing.B) { benchmarkSum3DGridV2(b, 10) }
// func BenchmarkSum3DGridV2_12(b *testing.B) { benchmarkSum3DGridV2(b, 12) }
// func BenchmarkSum3DGridV2_14(b *testing.B) { benchmarkSum3DGridV2(b, 14) }
// func BenchmarkSum3DGridV2_16(b *testing.B) { benchmarkSum3DGridV2(b, 16) }
// func BenchmarkSum3DGridV2_18(b *testing.B) { benchmarkSum3DGridV2(b, 18) }
// func BenchmarkSum3DGridV2_20(b *testing.B) { benchmarkSum3DGridV2(b, 20) }

// func benchmarkSumBipartiteV1(b *testing.B, m int) {
// 	node, a, component := buildCompleteBipartiteGraph(m, m/2)
// 	b.ResetTimer()

// 	for BN := 0; BN < b.N; BN++ {
// 		sumAllSourcesShortestPathsV1(node, a, component)
// 	}
// }

// func BenchmarkSumBipartiteV1_100(b *testing.B)  { benchmarkSumBipartiteV1(b, 100) }
// func BenchmarkSumBipartiteV1_200(b *testing.B)  { benchmarkSumBipartiteV1(b, 200) }
// func BenchmarkSumBipartiteV1_300(b *testing.B)  { benchmarkSumBipartiteV1(b, 300) }
// func BenchmarkSumBipartiteV1_400(b *testing.B)  { benchmarkSumBipartiteV1(b, 400) }
// func BenchmarkSumBipartiteV1_500(b *testing.B)  { benchmarkSumBipartiteV1(b, 500) }
// func BenchmarkSumBipartiteV1_600(b *testing.B)  { benchmarkSumBipartiteV1(b, 600) }
// func BenchmarkSumBipartiteV1_700(b *testing.B)  { benchmarkSumBipartiteV1(b, 700) }
// func BenchmarkSumBipartiteV1_800(b *testing.B)  { benchmarkSumBipartiteV1(b, 800) }
// func BenchmarkSumBipartiteV1_900(b *testing.B)  { benchmarkSumBipartiteV1(b, 900) }
// func BenchmarkSumBipartiteV1_1000(b *testing.B) { benchmarkSumBipartiteV1(b, 1000) }

// func benchmarkSumBipartiteV2(b *testing.B, m int) {
// 	node, a, component := buildCompleteBipartiteGraph(m, m/2)
// 	b.ResetTimer()

// 	for BN := 0; BN < b.N; BN++ {
// 		sumAllSourcesShortestPathsV2(node, a, component)
// 	}
// }

// func BenchmarkSumBipartiteV2_100(b *testing.B)  { benchmarkSumBipartiteV2(b, 100) }
// func BenchmarkSumBipartiteV2_200(b *testing.B)  { benchmarkSumBipartiteV2(b, 200) }
// func BenchmarkSumBipartiteV2_300(b *testing.B)  { benchmarkSumBipartiteV2(b, 300) }
// func BenchmarkSumBipartiteV2_400(b *testing.B)  { benchmarkSumBipartiteV2(b, 400) }
// func BenchmarkSumBipartiteV2_500(b *testing.B)  { benchmarkSumBipartiteV2(b, 500) }
// func BenchmarkSumBipartiteV2_600(b *testing.B)  { benchmarkSumBipartiteV2(b, 600) }
// func BenchmarkSumBipartiteV2_700(b *testing.B)  { benchmarkSumBipartiteV2(b, 700) }
// func BenchmarkSumBipartiteV2_800(b *testing.B)  { benchmarkSumBipartiteV2(b, 800) }
// func BenchmarkSumBipartiteV2_900(b *testing.B)  { benchmarkSumBipartiteV2(b, 900) }
// func BenchmarkSumBipartiteV2_1000(b *testing.B) { benchmarkSumBipartiteV2(b, 1000) }

// func benchmarkSumTreeV1(b *testing.B, n int) {
// 	node, a, component := buildCompleteBinaryTree(n)
// 	b.ResetTimer()

// 	for BN := 0; BN < b.N; BN++ {
// 		sumAllSourcesShortestPathsV1(node, a, component)
// 	}
// }

// func BenchmarkSumTreeV1_2(b *testing.B)  { benchmarkSumTreeV1(b, 2) }
// func BenchmarkSumTreeV1_4(b *testing.B)  { benchmarkSumTreeV1(b, 4) }
// func BenchmarkSumTreeV1_6(b *testing.B)  { benchmarkSumTreeV1(b, 6) }
// func BenchmarkSumTreeV1_8(b *testing.B)  { benchmarkSumTreeV1(b, 8) }
// func BenchmarkSumTreeV1_10(b *testing.B) { benchmarkSumTreeV1(b, 10) }
// func BenchmarkSumTreeV1_12(b *testing.B) { benchmarkSumTreeV1(b, 12) }
// func BenchmarkSumTreeV1_14(b *testing.B) { benchmarkSumTreeV1(b, 14) }

// func benchmarkSumTreeV2(b *testing.B, n int) {
// 	node, a, component := buildCompleteBinaryTree(n)
// 	b.ResetTimer()

// 	for BN := 0; BN < b.N; BN++ {
// 		sumAllSourcesShortestPathsV2(node, a, component)
// 	}
// }

// func BenchmarkSumTreeV2_2(b *testing.B)  { benchmarkSumTreeV2(b, 2) }
// func BenchmarkSumTreeV2_4(b *testing.B)  { benchmarkSumTreeV2(b, 4) }
// func BenchmarkSumTreeV2_6(b *testing.B)  { benchmarkSumTreeV2(b, 6) }
// func BenchmarkSumTreeV2_8(b *testing.B)  { benchmarkSumTreeV2(b, 8) }
// func BenchmarkSumTreeV2_10(b *testing.B) { benchmarkSumTreeV2(b, 10) }
// func BenchmarkSumTreeV2_12(b *testing.B) { benchmarkSumTreeV2(b, 12) }
// func BenchmarkSumTreeV2_14(b *testing.B) { benchmarkSumTreeV2(b, 14) }

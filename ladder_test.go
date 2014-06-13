package main

import "testing"

// NOTE: the files mentioned below are 1 through 9 letter extracts from /usr/share/dict/words

//
// TESTS
//

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

func benchmarkSumASSP(b *testing.B, f string, length int) {
	word := readWords([]string{f}, length)
	if len(word) < 1 {
		b.Logf("no words in input file")
	}
	pair := findPairs(word)
	component := findComponents(word, pair)
	b.ResetTimer()

	for BN := 0; BN < b.N; BN++ {
		sumAllSourcesShortestPaths(word, pair, component)
	}
}

func BenchmarkSumASSP_webster1(b *testing.B) { benchmarkSumASSP(b, "words/webster-1", 1) }
func BenchmarkSumASSP_webster2(b *testing.B) { benchmarkSumASSP(b, "words/webster-2", 2) }
func BenchmarkSumASSP_webster3(b *testing.B) { benchmarkSumASSP(b, "words/webster-3", 3) }
func BenchmarkSumASSP_webster4(b *testing.B) { benchmarkSumASSP(b, "words/webster-4", 4) }
func BenchmarkSumASSP_webster5(b *testing.B) { benchmarkSumASSP(b, "words/webster-5", 5) }
func BenchmarkSumASSP_webster6(b *testing.B) { benchmarkSumASSP(b, "words/webster-6", 6) }
func BenchmarkSumASSP_webster7(b *testing.B) { benchmarkSumASSP(b, "words/webster-7", 7) }
func BenchmarkSumASSP_webster8(b *testing.B) { benchmarkSumASSP(b, "words/webster-8", 8) }
func BenchmarkSumASSP_webster9(b *testing.B) { benchmarkSumASSP(b, "words/webster-9", 9) }

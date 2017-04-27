package main

import "testing"

// func BenchmarkDict1(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		process("dict1", "out1")
// 	}
// }
//
// func BenchmarkDict2(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		process("dict2", "out2")
// 	}
// }

func TestParseIndex(t *testing.T) {
	if parseIndex("1") != 1 {
		t.Error("failed parsing 1")
	}
	if parseIndex("2") != 2 {
		t.Error("failed parsing 2")
	}
	if parseIndex("200000") != 200000 {
		t.Error("failed parsing 200000", parseIndex("200000"))
	}
}

package main

import "testing"

func TestGeneratePasswordHash(t *testing.T) {
	result, _ := GeneratePasswordHash("a@2.a", "secret")
	t.Log(result)
}

func BenchmarkGeneratePasswordHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GeneratePasswordHash("a@2.a", "secret")
	}
}

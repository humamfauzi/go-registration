package main

import "testing"

func TestGeneratePasswordHash(t *testing.T) {
	result, _ := GeneratePasswordHash("a@2.a", "secret")
	t.Log(result)
}

func TestGenerateWebToken(t *testing.T) {
	Id := "DummyId"
	result, err := GenerateWebToken(Id)
	if err != nil {
		t.Error("should not error", err)
		return
	}
	t.Log(string(result))
}

func TestValidateWebToken(t *testing.T) {
	Id := "DummyId"
	result, err := GenerateWebToken(Id)
	if err != nil {
		t.Error("should not error", err)
		return
	}
	ok := ValidateWebToken(result)
	if !ok {
		t.Error("should not error", err)
		return
	}
	return
}

func BenchmarkGeneratePasswordHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GeneratePasswordHash("a@2.a", "secret")
	}
}

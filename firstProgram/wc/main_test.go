package main

import (
	"bytes"
	"testing"
)

// TestCountLines test the count function set to count lines
func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("word1, word2, word3, word4\nline2\nline3 word1")
	exp := 3
	res := count(b, true, false)
	if res != exp {
		t.Errorf("Expected %d, got %d instead.\n", exp, res)
	}
}

// TestCountBytes test the count funtion set to count bytes
func TestCountBytes(t *testing.T) {
	b := bytes.NewBufferString("word1")
	exp := 5
	res := count(b, false, true)
	if res != exp {
		t.Errorf("Expected %d, got %d instead.\n", exp, res)
	}
}

// TestCountWords test the count function set to count words
func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("word1, word2, word3, word4\n")
	exp := 4
	res := count(b, false, false)
	if res != exp {
		t.Errorf("Expected %d, got %d instead.\n", exp, res)
	}
}

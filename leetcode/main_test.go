package main_test

import (
	"testing"

	main "leetcode"
)

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		testName string
		number   int
		want     bool
	}{
		{
			testName: "'10' is not a Palindrome",
			number:   10,
			want:     false,
		},
		{
			testName: "'121' is a Palindrome",
			number:   121,
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			isPalindrome := main.IsPalindrome(tt.number)
			if tt.want != isPalindrome {
				t.Errorf("got %v, want %v", isPalindrome, tt.want)
			}
		})
	}
}

func BenchmarkIsPalindrome(b *testing.B) {
	for b.Loop() {
		main.IsPalindrome(121)
	}
}

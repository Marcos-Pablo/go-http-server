package main

import "testing"

func Test_cleanBody(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		s    string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanBody(tt.s)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("cleanBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cleanBody(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		s        string
		badWords map[string]struct{}
		want     string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanBody(tt.s, tt.badWords)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("cleanBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

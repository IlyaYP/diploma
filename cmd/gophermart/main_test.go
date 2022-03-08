package main

import (
	"testing"
)

func TestGeneralTests(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GeneralTests(); (err != nil) != tt.wantErr {
				t.Errorf("GeneralTests() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

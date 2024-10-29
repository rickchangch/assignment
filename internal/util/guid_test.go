package util

import (
	"fmt"
	"testing"
)

func TestGenXid(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "length",
			want: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenXid(); len(got) != tt.want {
				fmt.Print()
				t.Errorf("GenXid() = %v, want length %v", got, tt.want)
			}
		})
	}
}

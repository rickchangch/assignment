package util

import (
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	type In struct {
		field1 string
	}
	type Out struct {
		field1 string
	}
	type args struct {
		in []In
		fn func(In) Out
	}

	tests := []struct {
		name string
		args args
		want []Out
	}{
		{
			name: "positve",
			args: args{
				in: []In{
					{"test1"},
					{"test2"},
				},
				fn: func(in In) Out {
					return Out{in.field1}
				},
			},
			want: []Out{
				{"test1"},
				{"test2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.in, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

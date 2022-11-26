package main

import (
	"reflect"
	"testing"
)

func Test_cutAndShiftCards(t *testing.T) {
	type args struct {
		cards []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "simple shuffle",
			args: args{
				cards: []int{3, 4, 6},
			},
			want: []int{4, 6, 3},
		},
		{
			name: "one entry",
			args: args{
				cards: []int{3},
			},
			want: []int{3},
		},
		{
			name: "nil",
			args: args{
				cards: nil,
			},
			want: nil,
		},
		{
			name: "no entries",
			args: args{
				cards: []int{},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cutAndShiftCards(tt.args.cards); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cutAndShiftCards() = %v, want %v", got, tt.want)
			}
		})
	}
}

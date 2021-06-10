package util

import (
	"reflect"
	"testing"
)

func TestStringSliceContains(t *testing.T) {
	type args struct {
		slice []string
		val   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "expect true",
			args: args{
				slice: []string{"tes", "Test"},
				val:   "test",
			},
			want: true,
		},
		{
			name: "expect false",
			args: args{
				slice: []string{"tes", "Test"},
				val:   "testing",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSliceContains(tt.args.slice, tt.args.val); got != tt.want {
				t.Errorf("StringSliceContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSliceMergeUnique(t *testing.T) {
	type args struct {
		existingSlice []string
		mergeSlice    []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "no change",
			args: args{
				existingSlice: []string{"test", "test2"},
				mergeSlice:    []string{"test", "Test2"},
			},
			want: []string{"test", "test2"},
		},
		{
			name: "no change empty",
			args: args{
				existingSlice: []string{},
				mergeSlice:    []string{},
			},
			want: []string{},
		},
		{
			name: "with change",
			args: args{
				existingSlice: []string{"test", "test2"},
				mergeSlice:    []string{"test", "Test2", "test3"},
			},
			want: []string{"test", "test2", "test3"},
		},
		{
			name: "with change no empty",
			args: args{
				existingSlice: []string{"", "en"},
				mergeSlice:    []string{"fr", "", "en"},
			},
			want: []string{"en", "fr"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSliceMergeUnique(tt.args.existingSlice, tt.args.mergeSlice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSliceMergeUnique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSliceContainsAny(t *testing.T) {
	type args struct {
		slice []string
		vals  []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "expect true",
			args: args{
				slice: []string{"tes", "Test"},
				vals:  []string{"nope", "test"},
			},
			want: true,
		},
		{
			name: "expect false",
			args: args{
				slice: []string{"tes", "Test"},
				vals:  []string{"testing"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSliceContainsAny(tt.args.slice, tt.args.vals); got != tt.want {
				t.Errorf("StringSliceContainsAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

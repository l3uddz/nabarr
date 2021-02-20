package util

import "testing"

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

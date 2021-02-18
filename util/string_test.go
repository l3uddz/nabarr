package util

import "testing"

func TestStripNonAlphaNumeric(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "remove trailing slash",
			args: args{
				value: "tt1234567/",
			},
			want: "tt1234567",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StripNonAlphaNumeric(tt.args.value); got != tt.want {
				t.Errorf("StripNonAlphaNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

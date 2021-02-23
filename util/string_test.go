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

func TestStripNonNumeric(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "remove non numeric",
			args: args{
				value: "10, 560",
			},
			want: "10560",
		},
		{
			name: "remove nothing",
			args: args{
				value: "100",
			},
			want: "100",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StripNonNumeric(tt.args.value); got != tt.want {
				t.Errorf("StripNonNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

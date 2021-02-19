package util

import "testing"

func TestAtof64(t *testing.T) {
	type args struct {
		val        string
		defaultVal float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "expect whole value",
			args: args{
				val:        "1",
				defaultVal: 0,
			},
			want: 1.0,
		},
		{
			name: "expect non whole value",
			args: args{
				val:        "1.5",
				defaultVal: 0,
			},
			want: 1.5,
		},
		{
			name: "expect default",
			args: args{
				val:        "invalid",
				defaultVal: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Atof64(tt.args.val, tt.args.defaultVal); got != tt.want {
				t.Errorf("Atof64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAtoi(t *testing.T) {
	type args struct {
		val        string
		defaultVal int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "expect value",
			args: args{
				val:        "5",
				defaultVal: 0,
			},
			want: 5,
		},
		{
			name: "expect default",
			args: args{
				val:        "invalid",
				defaultVal: 2,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Atoi(tt.args.val, tt.args.defaultVal); got != tt.want {
				t.Errorf("Atoi() = %v, want %v", got, tt.want)
			}
		})
	}
}

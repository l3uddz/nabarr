package util

import "testing"

func TestContainsMovieCategory(t *testing.T) {
	type args struct {
		cats []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "want true",
			args: args{
				cats: []string{
					"2000",
					"5000",
				},
			},
			want: true,
		}, {
			name: "want false",
			args: args{
				cats: []string{
					"3000",
					"4000",
				},
			},
			want: false,
		}, {
			name: "want false",
			args: args{
				cats: []string{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsMovieCategory(tt.args.cats); got != tt.want {
				t.Errorf("ContainsMovieCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsTvCategory(t *testing.T) {
	type args struct {
		cats []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "want true",
			args: args{
				cats: []string{
					"5000",
					"4000",
					"2000",
				},
			},
			want: true,
		}, {
			name: "want false",
			args: args{
				cats: []string{
					"3000",
					"4000",
				},
			},
			want: false,
		}, {
			name: "want false",
			args: args{
				cats: []string{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsTvCategory(tt.args.cats); got != tt.want {
				t.Errorf("ContainsTvCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

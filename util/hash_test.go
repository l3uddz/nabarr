package util

import "testing"

func TestAsSHA256(t *testing.T) {
	type args struct {
		o interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "hash struct",
			args: args{
				o: struct {
					Name    string
					Surname string
					Age     int
				}{
					Name:    "John",
					Surname: "Smith",
					Age:     18,
				},
			},
			want: "f20fe06d96e179073fc3eebac62d7a2edf3164f0c50524d82c0c6390013bbc4a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AsSHA256(tt.args.o); got != tt.want {
				t.Errorf("AsSHA256() = %v, want %v", got, tt.want)
			}
		})
	}
}

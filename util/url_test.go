package util

import (
	"net/url"
	"testing"
)

func TestJoinURL(t *testing.T) {
	type args struct {
		base  string
		paths []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "single path",
			args: args{
				base:  "https://www.google.co.uk/",
				paths: []string{"search"},
			},
			want: "https://www.google.co.uk/search",
		},
		{
			name: "multiple path",
			args: args{
				base:  "https://www.google.co.uk",
				paths: []string{"search", "string"},
			},
			want: "https://www.google.co.uk/search/string",
		},
		{
			name: "multiple path with slashes",
			args: args{
				base:  "https://www.google.co.uk/",
				paths: []string{"/search/", "/string/"},
			},
			want: "https://www.google.co.uk/search/string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinURL(tt.args.base, tt.args.paths...); got != tt.want {
				t.Errorf("JoinURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURLWithQuery(t *testing.T) {
	type args struct {
		base string
		q    url.Values
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "url with values",
			args: args{
				base: JoinURL("https://api.trakt.tv", "search", "tvdb", "12345"),
				q: url.Values{
					"extended": []string{"full"},
					"type":     []string{"show"},
				},
			},
			want:    "https://api.trakt.tv/search/tvdb/12345?extended=full&type=show",
			wantErr: false,
		},
		{
			name: "trakt url without values",
			args: args{
				base: JoinURL("https://api.trakt.tv", "search", "tvdb", "12345"),
				q:    nil,
			},
			want:    "https://api.trakt.tv/search/tvdb/12345",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := URLWithQuery(tt.args.base, tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("URLWithQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("URLWithQuery() got = %v, want %v", got, tt.want)
			}
		})
	}
}

package cmd

import (
	"testing"

	"golang.org/x/exp/slices"

	"git.thrls.net/thiagorls/gosos/storage"
)

func TestParseRemoveArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{"Valid URL", []string{"http://example.com"}, "http://example.com", false},
		{"No arguments", []string{}, "", true},
		{"Multiple arguments", []string{"http://example.com", "extra"}, "http://example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRemoveArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRemoveArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseRemoveArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolveTarget(t *testing.T) {
	urls := []string{
		"http://example.com",
		"http://test.com",
		"https://third.com",
	}

	tests := []struct {
		name    string
		target  string
		want    string
		wantErr bool
	}{
		{"Literal URL passes through", "http://example.com", "http://example.com", false},
		{"Unknown URL passes through", "http://not-in-list.com", "http://not-in-list.com", false},
		{"Index 0", "0", "http://example.com", false},
		{"Index 1", "1", "http://test.com", false},
		{"Last index", "2", "https://third.com", false},
		{"Index out of range", "3", "", true},
		{"Negative index", "-1", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveTarget(tt.target, urls)
			if (err != nil) != tt.wantErr {
				t.Fatalf("resolveTarget() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("resolveTarget() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolveTargetEmptyList(t *testing.T) {
	if _, err := resolveTarget("0", nil); err == nil {
		t.Errorf("expected error when resolving index on empty list, got nil")
	}
}

func TestRemoveURLFromList(t *testing.T) {
	tests := []struct {
		name    string
		urlList *storage.URLList
		url     string
		want    *storage.URLList
		wantErr bool
	}{
		{
			name:    "Remove existing URL",
			urlList: &storage.URLList{URLs: []string{"http://example.com", "http://test.com"}},
			url:     "http://example.com",
			want:    &storage.URLList{URLs: []string{"http://test.com"}},
			wantErr: false,
		},
		{
			name:    "Remove non-existing URL",
			urlList: &storage.URLList{URLs: []string{"http://example.com", "http://test.com"}},
			url:     "http://nonexistent.com",
			want:    &storage.URLList{URLs: []string{"http://example.com", "http://test.com"}},
			wantErr: true,
		},
		{
			name:    "Remove from empty list",
			urlList: &storage.URLList{URLs: []string{}},
			url:     "http://example.com",
			want:    &storage.URLList{URLs: []string{}},
			wantErr: true,
		},
		{
			name:    "Remove last URL in list",
			urlList: &storage.URLList{URLs: []string{"http://example.com"}},
			url:     "http://example.com",
			want:    &storage.URLList{URLs: []string{}},
			wantErr: false,
		},
		{
			name:    "Remove URL with different scheme",
			urlList: &storage.URLList{URLs: []string{"http://example.com", "https://example.com"}},
			url:     "https://example.com",
			want:    &storage.URLList{URLs: []string{"http://example.com"}},
			wantErr: false,
		},
		{
			name:    "Remove duplicate URL",
			urlList: &storage.URLList{URLs: []string{"http://example.com", "http://test.com", "http://example.com"}},
			url:     "http://example.com",
			want:    &storage.URLList{URLs: []string{"http://test.com", "http://example.com"}},
			wantErr: false,
		},
		{
			name:    "Case sensitivity check",
			urlList: &storage.URLList{URLs: []string{"http://EXAMPLE.com", "http://test.com"}},
			url:     "http://example.com",
			want:    &storage.URLList{URLs: []string{"http://EXAMPLE.com", "http://test.com"}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := removeURLFromList(tt.urlList, tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("removeURLFromList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !slices.Equal(tt.urlList.URLs, tt.want.URLs) {
				t.Errorf("removeURLFromList() = %v, want %v", tt.urlList, tt.want)
			}
		})
	}
}

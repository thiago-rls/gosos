package cmd

import (
	"flag"
	"git.thrls.net/thiagorls/gosos/storage"
	"testing"
)

func TestValidateArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"Valid args", []string{"http://example.com"}, false},
		{"No args", []string{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := flag.NewFlagSet("test", flag.ContinueOnError)
			cmd.Parse(tt.args)
			err := validateArgs(cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		urlStr  string
		wantErr bool
	}{
		{"Valid URL", "http://example.com", false},
		{"Invalid URL", "not-a-url", true},
		{"Missing scheme", "example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURL(tt.urlStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddURLToList(t *testing.T) {
	urlList := &storage.URLList{URLs: []string{"http://existing.com"}}
	newURL := "http://new-example.com"

	err := addURLToList(urlList, newURL)
	if err != nil {
		t.Errorf("addURLToList() error = %v", err)
	}

	if len(urlList.URLs) != 2 || urlList.URLs[1] != newURL {
		t.Errorf("addURLToList() did not add URL correctly")
	}
}

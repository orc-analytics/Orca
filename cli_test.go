package main

import (
	"testing"

	dlyr "github.com/predixus/orca/internal/datalayers"
)

func TestDatalayerConsistency(t *testing.T) {
	// Test 1: Verify all datalayer suggestions can be cast to dlyr.Platform
	for _, suggestion := range datalayerSuggestions {
		// This will panic if the cast is invalid
		_ = dlyr.Platform(suggestion)
	}

	// Test 2: Verify datalayerSuggestions and connectionTemplates have exact same entries
	if len(datalayerSuggestions) != len(connectionTemplates) {
		t.Errorf("datalayerSuggestions length (%d) does not match connectionTemplates length (%d)",
			len(datalayerSuggestions), len(connectionTemplates))
	}

	for _, suggestion := range datalayerSuggestions {
		if _, exists := connectionTemplates[suggestion]; !exists {
			t.Errorf(
				"datalayer suggestion %q exists but has no corresponding connection template",
				suggestion,
			)
		}
	}

	for templateKey := range connectionTemplates {
		found := false
		for _, suggestion := range datalayerSuggestions {
			if templateKey == suggestion {
				found = true
				break
			}
		}
		if !found {
			t.Errorf(
				"connection template %q exists but has no corresponding datalayer suggestion",
				templateKey,
			)
		}
	}
}

func TestGetFullConnStr(t *testing.T) {
	tests := []struct {
		name     string
		template connStringTemplate
		want     string
	}{
		{
			name: "PostgreSQL template",
			template: connStringTemplate{
				prefix:     "postgresql://",
				components: []string{"user", "password", "host", "port", "dbname"},
				separators: []string{":", "@", ":", "/"},
			},
			want: "postgresql://<user>:<password>@<host>:<port>/<dbname>",
		},
		{
			name: "Simple template",
			template: connStringTemplate{
				prefix:     "simple://",
				components: []string{"first", "second"},
				separators: []string{":"},
			},
			want: "simple://<first>:<second>",
		},
		{
			name: "Single component template",
			template: connStringTemplate{
				prefix:     "basic://",
				components: []string{"only"},
				separators: []string{},
			},
			want: "basic://<only>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.template.getFullConnStr()
			if got != tt.want {
				t.Errorf("getFullConnStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

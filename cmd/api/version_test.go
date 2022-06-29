package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	expects := "wb-l0 version 1.4.0 (2022-06-21)\nhttps://github.com/maypok86/wb-l0/releases/tag/v1.4.0\n"
	got := format("1.4.0", "2022-06-21")
	require.Equal(t, expects, got)
}

func TestChangelogURL(t *testing.T) {
	tests := []struct {
		name string
		tag  string
		url  string
	}{
		{
			name: "Tag 0.3.2",
			tag:  "0.3.2",
			url:  "https://github.com/maypok86/wb-l0/releases/tag/v0.3.2",
		},
		{
			name: "Tag v0.3.2",
			tag:  "v0.3.2",
			url:  "https://github.com/maypok86/wb-l0/releases/tag/v0.3.2",
		},
		{
			name: "Tag 0.3.2-pre.1",
			tag:  "0.3.2-pre.1",
			url:  "https://github.com/maypok86/wb-l0/releases/tag/v0.3.2-pre.1",
		},
		{
			name: "Tag 0.3.5-90-gdd3f0e0",
			tag:  "0.3.5-90-gdd3f0e0",
			url:  "https://github.com/maypok86/wb-l0/releases/latest",
		},
		{
			name: "Tag deadbeef",
			tag:  "deadbeef",
			url:  "https://github.com/maypok86/wb-l0/releases/latest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := changelogURL(tt.tag)
			require.Equal(t, tt.url, got)
		})
	}
}

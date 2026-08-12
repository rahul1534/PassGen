package wordlist

import (
	"embed"
	"strings"
)

//go:embed words.txt
var file embed.FS

// Words returns the bundled passphrase word list.
func Words() ([]string, error) {
	data, err := file.ReadFile("words.txt")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	words := make([]string, 0, len(lines))
	for _, line := range lines {
		word := strings.TrimSpace(line)
		if word != "" {
			words = append(words, word)
		}
	}
	return words, nil
}

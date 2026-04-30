package theme

import (
	"embed"
	"strings"
)

//go:embed themes/*.json
var themesFS embed.FS

var builtinThemes = map[string][]byte{}

func init() {
	entries, err := themesFS.ReadDir("themes")
	if err != nil {
		return
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".json")
		data, err := themesFS.ReadFile("themes/" + e.Name())
		if err != nil {
			continue
		}
		builtinThemes[name] = data
	}
}

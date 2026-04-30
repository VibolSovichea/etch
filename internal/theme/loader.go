package theme

import (
	"encoding/json"
	"fmt"
	"os"
)

func Load(nameOrPath string) (*Theme, error) {
	if nameOrPath == "" || nameOrPath == "default" {
		return Default(), nil
	}

	if data, ok := builtinThemes[nameOrPath]; ok {
		return parse(data)
	}

	data, err := os.ReadFile(nameOrPath)
	if err != nil {
		return nil, fmt.Errorf("theme: cannot read %q: %w", nameOrPath, err)
	}
	return parse(data)
}

func parse(data []byte) (*Theme, error) {
	t := Default()
	if err := json.Unmarshal(data, t); err != nil {
		return nil, fmt.Errorf("theme: invalid JSON: %w", err)
	}
	if err := validate(t); err != nil {
		return nil, err
	}
	return t, nil
}

func validate(t *Theme) error {
	required := []struct {
		name  string
		value string
	}{
		{"palette.background.primary", t.Palette.Background.Primary},
		{"palette.text.primary", t.Palette.Text.Primary},
		{"palette.accent.primary", t.Palette.Accent.Primary},
	}
	for _, r := range required {
		if r.value == "" {
			return fmt.Errorf("theme: missing required field %q", r.name)
		}
	}
	return nil
}

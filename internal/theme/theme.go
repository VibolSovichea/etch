package theme

import (
	"sort"
	"strings"
)

// Theme defines the complete visual configuration for the TUI.
type Theme struct {
	Meta       Meta       `json:"meta"`
	Palette    Palette    `json:"palette"`
	Border     Border     `json:"border"`
	Components Components `json:"components"`
	Effects    Effects    `json:"effects"`
	Behavior   Behavior   `json:"behavior"`
}

type Meta struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Mode    string `json:"mode"`
}

type Palette struct {
	Background BgColors     `json:"background"`
	Text       TextColors   `json:"text"`
	Accent     AccentColors `json:"accent"`
	State      StateColors  `json:"state"`
}

type BgColors struct {
	Primary  string `json:"primary"`
	Panel    string `json:"panel"`
	Elevated string `json:"elevated"`
}

type TextColors struct {
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
	Muted     string `json:"muted"`
	Dim       string `json:"dim"`
	Inverse   string `json:"inverse"`
}

type AccentColors struct {
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
	Tertiary  string `json:"tertiary"`
}

type StateColors struct {
	Success string `json:"success"`
	Error   string `json:"error"`
	Warning string `json:"warning"`
	Info    string `json:"info"`
}

type Border struct {
	Style      string `json:"style"`
	Foreground string `json:"foreground"`
}

type Components struct {
	Panel  PanelComponent  `json:"panel"`
	Text   TextComponent   `json:"text"`
	Input  InputComponent  `json:"input"`
	Status StatusComponent `json:"status"`
	Graph  GraphComponent  `json:"graph"`
}

type PanelComponent struct {
	BorderForeground string `json:"border_foreground"`
	TitleForeground  string `json:"title_foreground"`
}

type TextComponent struct {
	Heading    string `json:"heading"`
	Body       string `json:"body"`
	Code       string `json:"code"`
	Link       string `json:"link"`
	Blockquote string `json:"blockquote"`
	ListMarker string `json:"list_marker"`
}

type InputComponent struct {
	Text        string `json:"text"`
	Placeholder string `json:"placeholder"`
	Cursor      string `json:"cursor"`
	Label       string `json:"label"`
}

type StatusComponent struct {
	NormalBg  string `json:"normal_bg"`
	NormalFg  string `json:"normal_fg"`
	InsertBg  string `json:"insert_bg"`
	InsertFg  string `json:"insert_fg"`
	CommandBg string `json:"command_bg"`
	CommandFg string `json:"command_fg"`
	BarBg     string `json:"bar_bg"`
	BarFg     string `json:"bar_fg"`
}

type GraphComponent struct {
	Bar   string `json:"bar"`
	Label string `json:"label"`
	Grid  string `json:"grid"`
}

type Effects struct {
	Glow     bool `json:"glow"`
	Scanline bool `json:"scanline"`
	Noise    bool `json:"noise"`
}

type Behavior struct {
	Density string `json:"density"`
}

// Resolve looks up a token path (e.g. "accent.primary") and returns
// the color value. If the input is already a hex color, it is returned
// as-is. Unknown paths return a visible fallback (#FF00FF).
func (t *Theme) Resolve(path string) string {
	if strings.HasPrefix(path, "#") {
		return path
	}
	if val, ok := t.tokens()[path]; ok && val != "" {
		return val
	}
	return "#FF00FF"
}

func (t *Theme) tokens() map[string]string {
	return map[string]string{
		"background.primary":  t.Palette.Background.Primary,
		"background.panel":    t.Palette.Background.Panel,
		"background.elevated": t.Palette.Background.Elevated,

		"text.primary":   t.Palette.Text.Primary,
		"text.secondary": t.Palette.Text.Secondary,
		"text.muted":     t.Palette.Text.Muted,
		"text.dim":       t.Palette.Text.Dim,
		"text.inverse":   t.Palette.Text.Inverse,

		"accent.primary":   t.Palette.Accent.Primary,
		"accent.secondary": t.Palette.Accent.Secondary,
		"accent.tertiary":  t.Palette.Accent.Tertiary,

		"state.success": t.Palette.State.Success,
		"state.error":   t.Palette.State.Error,
		"state.warning": t.Palette.State.Warning,
		"state.info":    t.Palette.State.Info,
	}
}

// ResolveComponent resolves a component token reference against the palette.
// Component fields store token paths like "accent.primary" which get
// resolved to actual color values.
func (t *Theme) ResolveComponent(ref string) string {
	return t.Resolve(ref)
}

// Default returns the default theme matching etch's original warm parchment palette.
func Default() *Theme {
	return &Theme{
		Meta: Meta{Name: "default", Version: "1.0", Mode: "dark"},
		Palette: Palette{
			Background: BgColors{
				Primary:  "#1A1816",
				Panel:    "#2A2825",
				Elevated: "#3A3835",
			},
			Text: TextColors{
				Primary:   "#FAF3E8",
				Secondary: "#E8D5B7",
				Muted:     "#8B8680",
				Dim:       "#5C5955",
				Inverse:   "#1A1816",
			},
			Accent: AccentColors{
				Primary:   "#D4A843",
				Secondary: "#C4873B",
				Tertiary:  "#7D9B76",
			},
			State: StateColors{
				Success: "#7D9B76",
				Error:   "#A0522D",
				Warning: "#C4873B",
				Info:    "#8B8680",
			},
		},
		Border: Border{
			Style:      "rounded",
			Foreground: "text.dim",
		},
		Components: Components{
			Panel: PanelComponent{
				BorderForeground: "text.dim",
				TitleForeground:  "accent.primary",
			},
			Text: TextComponent{
				Heading:    "accent.primary",
				Body:       "text.primary",
				Code:       "accent.secondary",
				Link:       "accent.tertiary",
				Blockquote: "text.muted",
				ListMarker: "accent.secondary",
			},
			Input: InputComponent{
				Text:        "text.secondary",
				Placeholder: "text.dim",
				Cursor:      "accent.primary",
				Label:       "accent.secondary",
			},
			Status: StatusComponent{
				NormalBg:  "accent.primary",
				NormalFg:  "text.inverse",
				InsertBg:  "accent.tertiary",
				InsertFg:  "text.inverse",
				CommandBg: "accent.secondary",
				CommandFg: "text.inverse",
				BarBg:     "background.panel",
				BarFg:     "text.muted",
			},
			Graph: GraphComponent{
				Bar:   "accent.primary",
				Label: "text.muted",
				Grid:  "text.dim",
			},
		},
		Effects:  Effects{},
		Behavior: Behavior{Density: "normal"},
	}
}

// BuiltinNames returns the names of all built-in themes, sorted alphabetically.
func BuiltinNames() []string {
	names := make([]string, 0, len(builtinThemes))
	for name := range builtinThemes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

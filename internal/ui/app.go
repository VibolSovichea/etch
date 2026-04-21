package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/vibolsovichea/scripture/internal/asset"
	"github.com/vibolsovichea/scripture/internal/config"
	"github.com/vibolsovichea/scripture/internal/note"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const maxRecent = 5

// --- View enum ---

type view int

const (
	viewDashboard view = iota
	viewFinder
)

// --- Finder sub-modes ---

type finderMode int

const (
	finderBrowse finderMode = iota
	finderSearch
	finderCreate
	finderCreateTags
	finderDelete
)

// --- Dashboard actions ---

type dashAction int

const (
	dashActionFind dashAction = iota
	dashActionNew
	dashActionQuit
)

var dashActions = []struct {
	key  string
	label string
	action dashAction
}{
	{"f", "Find Note", dashActionFind},
	{"n", "New Note", dashActionNew},
	{"q", "Quit", dashActionQuit},
}

// --- App Model ---

type AppModel struct {
	cfg    *config.Config
	ascii  string
	notes  []*note.Note
	recent []*note.Note
	width  int
	height int

	// Dashboard state
	view       view
	dashCursor int // cursor over actions + recent files
	// actions are indices 0..len(dashActions)-1
	// recent files are indices len(dashActions)..len(dashActions)+len(recent)-1

	// Finder state
	finderMode   finderMode
	filtered     []*note.Note
	finderCursor int
	finderScroll int
	input        textinput.Model
	newTitle     string

	err error
}

func NewAppModel(cfg *config.Config) AppModel {
	ti := textinput.New()
	ti.CharLimit = 256
	ti.Width = 40
	ti.PromptStyle = inputLabelStyle
	ti.TextStyle = lipgloss.NewStyle().Foreground(sand)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(gold)

	ascii := strings.TrimRight(asset.ASCIIArt, "\n")

	m := AppModel{
		cfg:   cfg,
		ascii: ascii,
		input: ti,
		view:  viewDashboard,
	}
	m.loadNotes()
	return m
}

func (m *AppModel) loadNotes() {
	notes, err := note.ListAll(m.cfg.VaultPath)
	if err != nil {
		m.err = err
		return
	}
	m.notes = notes
	m.filtered = notes

	// Sort by modified time descending for recents
	sorted := make([]*note.Note, len(notes))
	copy(sorted, notes)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Modified.After(sorted[j].Modified)
	})
	if len(sorted) > maxRecent {
		sorted = sorted[:maxRecent]
	}
	m.recent = sorted

	if m.finderCursor >= len(m.filtered) {
		m.finderCursor = max(0, len(m.filtered)-1)
	}
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

// --- Update ---

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case editorFinishedMsg:
		m.loadNotes()
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.view == viewDashboard {
			return m.updateDashboard(msg)
		}
		return m.updateFinder(msg)
	}
	return m, nil
}

// --- Dashboard Update ---

func (m AppModel) updateDashboard(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	totalItems := len(dashActions) + len(m.recent)

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "j", "down":
		if m.dashCursor < totalItems-1 {
			m.dashCursor++
		}
	case "k", "up":
		if m.dashCursor > 0 {
			m.dashCursor--
		}
	case "f":
		m.openFinder(finderSearch)
		return m, textinput.Blink
	case "n":
		m.openFinder(finderCreate)
		m.input.Placeholder = "Note title"
		return m, textinput.Blink
	case "q":
		return m, tea.Quit
	case "enter":
		return m.handleDashEnter()
	}
	return m, nil
}

func (m AppModel) handleDashEnter() (tea.Model, tea.Cmd) {
	if m.dashCursor < len(dashActions) {
		switch dashActions[m.dashCursor].action {
		case dashActionFind:
			m.openFinder(finderSearch)
			return m, textinput.Blink
		case dashActionNew:
			m.openFinder(finderCreate)
			m.input.Placeholder = "Note title"
			return m, textinput.Blink
		case dashActionQuit:
			return m, tea.Quit
		}
	}

	// Recent file selected
	recentIdx := m.dashCursor - len(dashActions)
	if recentIdx >= 0 && recentIdx < len(m.recent) {
		n := m.recent[recentIdx]
		c := makeEditorCmd(n.Path)
		return m, tea.ExecProcess(c, func(err error) tea.Msg {
			return editorFinishedMsg{err}
		})
	}
	return m, nil
}

func (m *AppModel) openFinder(mode finderMode) {
	m.view = viewFinder
	m.finderMode = mode
	m.finderCursor = 0
	m.finderScroll = 0
	m.filtered = m.notes
	m.input.SetValue("")
	if mode == finderSearch {
		m.input.Placeholder = "Search notes..."
	}
	m.input.Focus()
}

// --- Finder Update ---

func (m AppModel) updateFinder(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.finderMode {
	case finderDelete:
		return m.handleFinderDelete(msg)
	case finderCreate, finderCreateTags:
		return m.handleFinderCreate(msg)
	default:
		return m.handleFinderNormal(msg)
	}
}

func (m AppModel) handleFinderNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.view = viewDashboard
		m.input.Blur()
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	case "down", "ctrl+n":
		if m.finderCursor < len(m.filtered)-1 {
			m.finderCursor++
			m.clampFinderScroll()
		}
	case "up", "ctrl+p":
		if m.finderCursor > 0 {
			m.finderCursor--
			m.clampFinderScroll()
		}
	case "enter":
		if len(m.filtered) > 0 {
			n := m.filtered[m.finderCursor]
			c := makeEditorCmd(n.Path)
			return m, tea.ExecProcess(c, func(err error) tea.Msg {
				return editorFinishedMsg{err}
			})
		}
	case "ctrl+d":
		if len(m.filtered) > 0 {
			m.finderMode = finderDelete
		}
		return m, nil
	default:
		// Pass to text input for searching
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)

		query := m.input.Value()
		if query == "" {
			m.filtered = m.notes
		} else {
			m.filtered = filterNotes(m.notes, query)
		}
		m.finderCursor = 0
		m.finderScroll = 0
		return m, cmd
	}
	return m, nil
}

func (m AppModel) handleFinderCreate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		if m.finderMode == finderCreateTags {
			// Go back to title
			m.finderMode = finderCreate
			m.input.Placeholder = "Note title"
			m.input.SetValue(m.newTitle)
			return m, nil
		}
		m.view = viewDashboard
		m.input.Blur()
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	case "enter":
		if m.finderMode == finderCreate {
			m.newTitle = m.input.Value()
			if m.newTitle == "" {
				m.view = viewDashboard
				m.input.Blur()
				return m, nil
			}
			m.finderMode = finderCreateTags
			m.input.Placeholder = "Tags (comma separated, optional)"
			m.input.SetValue("")
			return m, nil
		}
		// finderCreateTags
		tags := parseTags(m.input.Value())
		dir := filepath.Join(m.cfg.VaultPath, "notes")
		n, err := note.Create(dir, m.newTitle, tags)
		if err != nil {
			m.err = err
			m.view = viewDashboard
			m.input.Blur()
			return m, nil
		}
		m.view = viewDashboard
		m.input.Blur()
		m.loadNotes()
		c := makeEditorCmd(n.Path)
		return m, tea.ExecProcess(c, func(err error) tea.Msg {
			return editorFinishedMsg{err}
		})
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m AppModel) handleFinderDelete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		if len(m.filtered) > 0 {
			n := m.filtered[m.finderCursor]
			trashDir := filepath.Join(m.cfg.VaultPath, ".scripture", "trash")
			n.Delete(trashDir)
			m.loadNotes()
			m.filtered = m.notes
		}
		m.finderMode = finderSearch
	default:
		m.finderMode = finderSearch
	}
	return m, nil
}

// --- Views ---

func (m AppModel) View() string {
	if m.width == 0 {
		return ""
	}

	switch m.view {
	case viewFinder:
		return m.viewFinder()
	default:
		return m.viewDashboard()
	}
}

// --- Dashboard View ---

func (m AppModel) viewDashboard() string {
	var sections []string

	// ASCII art
	art := asciiStyle.Render(m.ascii)
	sections = append(sections, art)
	sections = append(sections, "")

	// Actions
	for i, a := range dashActions {
		key := dashActionKeyStyle.Render(a.key)
		label := dashActionStyle.Render(a.label)
		cursor := "  "
		if i == m.dashCursor {
			cursor = dashActionSelectedStyle.Render("> ")
			label = dashActionSelectedStyle.Render(a.label)
		}
		sections = append(sections, fmt.Sprintf("  %s %s  %s", cursor, key, label))
	}

	// Recent files
	if len(m.recent) > 0 {
		sections = append(sections, "")
		sections = append(sections, "  "+dashRecentTitleStyle.Render("  Recent Notes"))
		sections = append(sections, "")

		for i, n := range m.recent {
			globalIdx := len(dashActions) + i
			cursor := "  "
			title := dashRecentItemStyle.Render(truncate(n.Title, m.width-20))
			date := dashRecentDateStyle.Render(relativeTime(n.Modified))

			if globalIdx == m.dashCursor {
				cursor = dashRecentSelectedStyle.Render("> ")
				title = dashRecentSelectedStyle.Render(truncate(n.Title, m.width-20))
			}

			sections = append(sections, fmt.Sprintf("  %s %s  %s", cursor, title, date))
		}
	}

	// Footer
	sections = append(sections, "")
	sections = append(sections, dashFooterStyle.Render(
		fmt.Sprintf("    Scripture  —  %d notes in vault", len(m.notes)),
	))

	content := strings.Join(sections, "\n")

	// Center everything
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

// --- Finder View (Telescope-style) ---

func (m AppModel) viewFinder() string {
	// Modal dimensions
	modalW := min(m.width-4, 100)
	modalH := min(m.height-4, 30)
	if modalW < 40 {
		modalW = 40
	}
	if modalH < 10 {
		modalH = 10
	}

	listW := modalW * 2 / 5
	previewW := modalW - listW

	// Title bar
	title := m.finderTitle()
	titleBar := finderTitleStyle.Render(title)

	// Input bar
	inputBar := m.finderInputBar(modalW)

	// Content area height
	contentH := modalH - 4 // title + input + border overhead

	// Left: file list
	left := m.finderList(listW, contentH)

	// Right: preview
	right := m.finderPreview(previewW, contentH)

	// Separator
	sep := finderPreviewDivStyle.Render(
		strings.Repeat("│\n", contentH),
	)

	content := lipgloss.JoinHorizontal(lipgloss.Top, left, sep, right)

	// Bottom help
	help := m.finderHelp()

	// Compose modal
	modal := lipgloss.JoinVertical(lipgloss.Left,
		titleBar,
		inputBar,
		content,
		help,
	)

	framed := finderBorderStyle.
		Width(modalW).
		Render(modal)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		framed,
	)
}

func (m AppModel) finderTitle() string {
	switch m.finderMode {
	case finderCreate:
		return " New Note"
	case finderCreateTags:
		return " New Note — Tags"
	case finderDelete:
		return " Delete Note"
	default:
		return " Find Notes"
	}
}

func (m AppModel) finderInputBar(w int) string {
	if m.finderMode == finderDelete {
		n := m.filtered[m.finderCursor]
		return deleteWarnStyle.Render(
			fmt.Sprintf(" Delete \"%s\"? (y/n)", truncate(n.Title, w-25)),
		)
	}

	prompt := inputLabelStyle.Render(" > ")
	return prompt + m.input.View()
}

func (m AppModel) finderList(w, h int) string {
	var lines []string

	if len(m.filtered) == 0 {
		empty := finderPreviewEmptyStyle.Render("  No matches")
		lines = append(lines, empty)
	} else {
		visibleCount := h
		for i, n := range m.filtered {
			if i < m.finderScroll {
				continue
			}
			if len(lines) >= visibleCount {
				break
			}

			cursor := "  "
			if i == m.finderCursor {
				cursor = finderCursorStyle.Render("> ")
			}

			maxW := w - 6
			title := truncate(n.Title, maxW)
			if i == m.finderCursor {
				title = finderItemSelectedStyle.Render(title)
			} else {
				title = finderItemStyle.Render(title)
			}

			lines = append(lines, cursor+title)
		}
	}

	// Pad to fill height
	for len(lines) < h {
		lines = append(lines, "")
	}

	// Count indicator
	countStr := finderCountStyle.Render(
		fmt.Sprintf("  %d/%d", len(m.filtered), len(m.notes)),
	)
	if len(lines) > 0 {
		lines[len(lines)-1] = countStr
	}

	return lipgloss.NewStyle().
		Width(w).
		Height(h).
		Render(strings.Join(lines, "\n"))
}

func (m AppModel) finderPreview(w, h int) string {
	if len(m.filtered) == 0 {
		empty := finderPreviewEmptyStyle.Render("No note selected")
		return lipgloss.NewStyle().
			Width(w).
			Height(h).
			Padding(0, 1).
			Render(lipgloss.Place(w-2, h, lipgloss.Center, lipgloss.Center, empty))
	}

	n := m.filtered[m.finderCursor]
	var b strings.Builder

	// Title
	b.WriteString(finderPreviewTitleStyle.Render(n.Title))
	b.WriteString("\n")

	// Tags
	if len(n.Tags) > 0 {
		tags := make([]string, len(n.Tags))
		for i, t := range n.Tags {
			tags[i] = "#" + t
		}
		b.WriteString(finderPreviewTagStyle.Render(strings.Join(tags, "  ")))
		b.WriteString("\n")
	}

	// Dates
	b.WriteString(finderPreviewMetaStyle.Render(
		fmt.Sprintf("Created %s  Modified %s",
			n.Created.Format("Jan 02, 2006"),
			n.Modified.Format("Jan 02, 2006"))))
	b.WriteString("\n")

	// Divider
	divW := w - 4
	if divW < 5 {
		divW = 5
	}
	b.WriteString(finderPreviewDivStyle.Render(strings.Repeat("─", divW)))
	b.WriteString("\n")

	// Body
	if n.Body == "" {
		b.WriteString(finderPreviewEmptyStyle.Render("Empty note"))
	} else {
		// Truncate body to fit
		bodyLines := strings.Split(n.Body, "\n")
		maxLines := h - 6
		if maxLines < 1 {
			maxLines = 1
		}
		if len(bodyLines) > maxLines {
			bodyLines = bodyLines[:maxLines]
			bodyLines = append(bodyLines, finderPreviewEmptyStyle.Render("..."))
		}
		b.WriteString(finderPreviewBodyStyle.Render(strings.Join(bodyLines, "\n")))
	}

	return lipgloss.NewStyle().
		Width(w).
		Height(h).
		Padding(0, 1).
		Render(b.String())
}

func (m AppModel) finderHelp() string {
	if m.finderMode == finderDelete {
		return ""
	}
	if m.finderMode == finderCreate || m.finderMode == finderCreateTags {
		return helpDescStyle.Render(" Enter confirm  Esc back")
	}

	parts := []string{
		helpKeyStyle.Render("Enter") + helpDescStyle.Render(" open"),
		helpKeyStyle.Render("Ctrl+d") + helpDescStyle.Render(" delete"),
		helpKeyStyle.Render("Esc") + helpDescStyle.Render(" back"),
	}
	return " " + strings.Join(parts, "  ")
}

// --- Helpers ---

func (m *AppModel) clampFinderScroll() {
	visibleH := m.finderVisibleCount()
	if m.finderCursor < m.finderScroll {
		m.finderScroll = m.finderCursor
	}
	if m.finderCursor >= m.finderScroll+visibleH {
		m.finderScroll = m.finderCursor - visibleH + 1
	}
	if m.finderScroll < 0 {
		m.finderScroll = 0
	}
}

func (m AppModel) finderVisibleCount() int {
	modalH := min(m.height-4, 30)
	h := modalH - 6
	if h < 3 {
		h = 3
	}
	return h
}

func makeEditorCmd(path string) *exec.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	return exec.Command(editor, path)
}

type editorFinishedMsg struct{ err error }

func filterNotes(notes []*note.Note, query string) []*note.Note {
	q := strings.ToLower(query)
	var result []*note.Note
	for _, n := range notes {
		if strings.Contains(strings.ToLower(n.Title), q) ||
			strings.Contains(strings.ToLower(n.Body), q) ||
			containsTag(n.Tags, q) {
			result = append(result, n)
		}
	}
	return result
}

func containsTag(tags []string, query string) bool {
	for _, t := range tags {
		if strings.Contains(strings.ToLower(t), query) {
			return true
		}
	}
	return false
}

func parseTags(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	var tags []string
	for _, t := range strings.Split(s, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}

func truncate(s string, maxW int) string {
	if maxW <= 0 {
		return ""
	}
	if len(s) <= maxW {
		return s
	}
	if maxW <= 3 {
		return s[:maxW]
	}
	return s[:maxW-3] + "..."
}

func relativeTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		m := int(diff.Minutes())
		if m == 1 {
			return "1 min ago"
		}
		return fmt.Sprintf("%d mins ago", m)
	case diff < 24*time.Hour:
		h := int(diff.Hours())
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", h)
	case diff < 7*24*time.Hour:
		d := int(diff.Hours() / 24)
		if d == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%d days ago", d)
	default:
		return t.Format("Jan 02, 2006")
	}
}

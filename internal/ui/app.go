package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/vibolsovichea/scripture/internal/config"
	"github.com/vibolsovichea/scripture/internal/note"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// View modes
type mode int

const (
	modeNormal mode = iota
	modeSearch
	modeCreate
	modeCreateTags
	modeDelete
)

var (
	sidebarStyle = lipgloss.NewStyle().
			Width(30).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderRight(true).
			Padding(1)

	previewStyle = lipgloss.NewStyle().
			Padding(1, 2)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	tagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)
)

type AppModel struct {
	cfg      *config.Config
	notes    []*note.Note
	filtered []*note.Note
	cursor   int
	width    int
	height   int
	mode     mode
	input    textinput.Model
	newTitle string // used during create flow
	err      error
}

func NewAppModel(cfg *config.Config) AppModel {
	ti := textinput.New()
	ti.CharLimit = 256
	ti.Width = 40

	m := AppModel{
		cfg:   cfg,
		input: ti,
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
	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

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
		// Handle input modes first
		if m.mode == modeSearch || m.mode == modeCreate || m.mode == modeCreateTags {
			return m.handleInputMode(msg)
		}
		if m.mode == modeDelete {
			return m.handleDeleteMode(msg)
		}
		return m.handleNormalMode(msg)
	}
	return m, nil
}

func (m AppModel) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "j", "down":
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
		}
	case "k", "up":
		if m.cursor > 0 {
			m.cursor--
		}
	case "n":
		m.mode = modeCreate
		m.input.Placeholder = "Note title"
		m.input.SetValue("")
		m.input.Focus()
		return m, textinput.Blink
	case "/":
		m.mode = modeSearch
		m.input.Placeholder = "Search..."
		m.input.SetValue("")
		m.input.Focus()
		return m, textinput.Blink
	case "e", "enter":
		if len(m.filtered) > 0 {
			n := m.filtered[m.cursor]
			c := makeEditorCmd(n.Path)
			return m, tea.ExecProcess(c, func(err error) tea.Msg {
				return editorFinishedMsg{err}
			})
		}
	case "d":
		if len(m.filtered) > 0 {
			m.mode = modeDelete
		}
	case "r":
		m.loadNotes()
	}
	return m, nil
}

func (m AppModel) handleInputMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeNormal
		m.input.Blur()
		if m.mode != modeSearch {
			m.filtered = m.notes
		}
		return m, nil
	case "enter":
		switch m.mode {
		case modeSearch:
			query := m.input.Value()
			m.filtered = filterNotes(m.notes, query)
			m.cursor = 0
			m.mode = modeNormal
			m.input.Blur()
		case modeCreate:
			m.newTitle = m.input.Value()
			if m.newTitle == "" {
				m.mode = modeNormal
				m.input.Blur()
				return m, nil
			}
			m.mode = modeCreateTags
			m.input.Placeholder = "Tags (comma separated, optional)"
			m.input.SetValue("")
			return m, nil
		case modeCreateTags:
			tags := parseTags(m.input.Value())
			dir := filepath.Join(m.cfg.VaultPath, "notes")
			n, err := note.Create(dir, m.newTitle, tags)
			if err != nil {
				m.err = err
				m.mode = modeNormal
				m.input.Blur()
				return m, nil
			}
			m.mode = modeNormal
			m.input.Blur()
			m.loadNotes()
			c := makeEditorCmd(n.Path)
			return m, tea.ExecProcess(c, func(err error) tea.Msg {
				return editorFinishedMsg{err}
			})
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	// Live search filtering
	if m.mode == modeSearch {
		query := m.input.Value()
		if query == "" {
			m.filtered = m.notes
		} else {
			m.filtered = filterNotes(m.notes, query)
		}
		m.cursor = 0
	}

	return m, cmd
}

func (m AppModel) handleDeleteMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		if len(m.filtered) > 0 {
			n := m.filtered[m.cursor]
			trashDir := filepath.Join(m.cfg.VaultPath, ".scripture", "trash")
			n.Delete(trashDir)
			m.loadNotes()
		}
		m.mode = modeNormal
	default:
		m.mode = modeNormal
	}
	return m, nil
}

func (m AppModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	sidebar := m.renderSidebar()
	preview := m.renderPreview()

	sideW := 32
	previewW := m.width - sideW - 4
	if previewW < 20 {
		previewW = 20
	}

	contentH := m.height - 3 // leave room for help bar

	left := sidebarStyle.Width(sideW).Height(contentH).Render(sidebar)
	right := previewStyle.Width(previewW).Height(contentH).Render(preview)

	main := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	help := m.renderHelp()

	return lipgloss.JoinVertical(lipgloss.Left, main, help)
}

func (m AppModel) renderSidebar() string {
	var b strings.Builder

	b.WriteString(headerStyle.Render("Scripture"))
	b.WriteString("\n\n")

	if m.mode == modeSearch || m.mode == modeCreate || m.mode == modeCreateTags {
		label := "Search"
		if m.mode == modeCreate {
			label = "New Note"
		} else if m.mode == modeCreateTags {
			label = "Tags"
		}
		b.WriteString(fmt.Sprintf("%s: %s\n\n", label, m.input.View()))
	}

	if len(m.filtered) == 0 {
		b.WriteString(dimStyle.Render("No notes found.\nPress 'n' to create one."))
		return b.String()
	}

	for i, n := range m.filtered {
		prefix := "  "
		if i == m.cursor {
			prefix = "> "
		}

		title := n.Title
		if i == m.cursor {
			title = selectedStyle.Render(title)
		}

		tags := ""
		if len(n.Tags) > 0 {
			tags = " " + tagStyle.Render("["+strings.Join(n.Tags, ", ")+"]")
		}

		b.WriteString(fmt.Sprintf("%s%s%s\n", prefix, title, tags))
	}

	return b.String()
}

func (m AppModel) renderPreview() string {
	if len(m.filtered) == 0 {
		return dimStyle.Render("No note selected")
	}

	n := m.filtered[m.cursor]
	var b strings.Builder

	b.WriteString(headerStyle.Render(n.Title))
	b.WriteString("\n")

	if len(n.Tags) > 0 {
		b.WriteString(tagStyle.Render("Tags: " + strings.Join(n.Tags, ", ")))
		b.WriteString("\n")
	}

	b.WriteString(dimStyle.Render(fmt.Sprintf("Created: %s  Modified: %s",
		n.Created.Format("2006-01-02"),
		n.Modified.Format("2006-01-02"))))
	b.WriteString("\n\n")

	b.WriteString(n.Body)

	return b.String()
}

func (m AppModel) renderHelp() string {
	switch m.mode {
	case modeDelete:
		return helpStyle.Render("Delete this note? (y/n)")
	case modeSearch, modeCreate, modeCreateTags:
		return helpStyle.Render("Enter to confirm • Esc to cancel")
	default:
		return helpStyle.Render("n:new  e/↵:edit  d:delete  /:search  r:refresh  q:quit")
	}
}

// makeEditorCmd creates an exec.Cmd for the user's editor.
func makeEditorCmd(path string) *exec.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	c := exec.Command(editor, path)
	return c
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

func getEnv(key string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return ""
}

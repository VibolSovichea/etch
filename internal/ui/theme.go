package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/vibolsovichea/etch/internal/theme"
)

var activeTheme *theme.Theme

var (
	borderStyle lipgloss.Border

	asciiStyle              lipgloss.Style
	dashActionStyle         lipgloss.Style
	dashActionKeyStyle      lipgloss.Style
	dashActionSelectedStyle lipgloss.Style
	dashRecentTitleStyle    lipgloss.Style
	dashRecentItemStyle     lipgloss.Style
	dashRecentSelectedStyle lipgloss.Style
	dashRecentDateStyle     lipgloss.Style
	dashFooterStyle         lipgloss.Style

	finderBorderStyle       lipgloss.Style
	finderTitleStyle        lipgloss.Style
	finderInputStyle        lipgloss.Style
	finderItemStyle         lipgloss.Style
	finderItemSelectedStyle lipgloss.Style
	finderCursorStyle       lipgloss.Style
	finderTagStyle          lipgloss.Style
	finderDateStyle         lipgloss.Style
	finderPreviewTitleStyle lipgloss.Style
	finderPreviewMetaStyle  lipgloss.Style
	finderPreviewTagStyle   lipgloss.Style
	finderPreviewBodyStyle  lipgloss.Style
	finderPreviewDivStyle   lipgloss.Style
	finderPreviewEmptyStyle lipgloss.Style
	finderCountStyle        lipgloss.Style

	inputLabelStyle        lipgloss.Style
	deleteWarnStyle        lipgloss.Style
	helpKeyStyle           lipgloss.Style
	helpDescStyle          lipgloss.Style
	createBorderStyle      lipgloss.Style
	createTitleStyle       lipgloss.Style
	createLabelStyle       lipgloss.Style
	createStepActiveStyle  lipgloss.Style
	createStepDoneStyle    lipgloss.Style
	createStepPendingStyle lipgloss.Style
	createValueStyle       lipgloss.Style

	edGutterStyle       lipgloss.Style
	edGutterActiveStyle lipgloss.Style
	edGutterSepStyle    lipgloss.Style
	edTextStyle         lipgloss.Style
	edCursorBlockStyle     lipgloss.Style
	edCursorUnderlineStyle lipgloss.Style
	edStatusBarStyle    lipgloss.Style
	edModeNormalStyle   lipgloss.Style
	edModeInsertStyle   lipgloss.Style
	edModeCommandStyle  lipgloss.Style
	edFileNameStyle     lipgloss.Style
	edModifiedStyle     lipgloss.Style
	edPosStyle          lipgloss.Style
	edCmdLineStyle      lipgloss.Style

	mdHeadingStyle    lipgloss.Style
	mdBoldStyle       lipgloss.Style
	mdItalicStyle     lipgloss.Style
	mdCodeInlineStyle lipgloss.Style
	mdCodeBlockStyle  lipgloss.Style
	mdLinkStyle       lipgloss.Style
	mdBlockquoteStyle lipgloss.Style
	mdListMarkerStyle lipgloss.Style
	mdHrStyle         lipgloss.Style

	inputTextStyle        lipgloss.Style
	inputCursorStyle      lipgloss.Style
	inputPlaceholderStyle lipgloss.Style
	editorBodyStyle       lipgloss.Style
)

func init() {
	InitTheme(theme.Default())
}

func InitTheme(t *theme.Theme) {
	activeTheme = t
	buildStyles(t)
}

func ThemeColor(path string) lipgloss.Color {
	return lipgloss.Color(activeTheme.Resolve(path))
}

func CurrentTheme() *theme.Theme {
	return activeTheme
}

func buildStyles(t *theme.Theme) {
	accentPrimary := lipgloss.Color(t.Resolve("accent.primary"))
	accentSecondary := lipgloss.Color(t.Resolve("accent.secondary"))
	accentTertiary := lipgloss.Color(t.Resolve("accent.tertiary"))
	textPrimary := lipgloss.Color(t.Resolve("text.primary"))
	textSecondary := lipgloss.Color(t.Resolve("text.secondary"))
	textMuted := lipgloss.Color(t.Resolve("text.muted"))
	textDim := lipgloss.Color(t.Resolve("text.dim"))
	stateError := lipgloss.Color(t.Resolve("state.error"))
	stateSuccess := lipgloss.Color(t.Resolve("state.success"))

	compPanelBorder := lipgloss.Color(t.ResolveComponent(t.Components.Panel.BorderForeground))
	compPanelTitle := lipgloss.Color(t.ResolveComponent(t.Components.Panel.TitleForeground))
	compTextHeading := lipgloss.Color(t.ResolveComponent(t.Components.Text.Heading))
	compTextBody := lipgloss.Color(t.ResolveComponent(t.Components.Text.Body))
	compTextCode := lipgloss.Color(t.ResolveComponent(t.Components.Text.Code))
	compTextLink := lipgloss.Color(t.ResolveComponent(t.Components.Text.Link))
	compTextBlockquote := lipgloss.Color(t.ResolveComponent(t.Components.Text.Blockquote))
	compTextListMarker := lipgloss.Color(t.ResolveComponent(t.Components.Text.ListMarker))
	compInputText := lipgloss.Color(t.ResolveComponent(t.Components.Input.Text))
	compInputPlaceholder := lipgloss.Color(t.ResolveComponent(t.Components.Input.Placeholder))
	compInputCursor := lipgloss.Color(t.ResolveComponent(t.Components.Input.Cursor))
	compInputLabel := lipgloss.Color(t.ResolveComponent(t.Components.Input.Label))
	compStatusNormalBg := lipgloss.Color(t.ResolveComponent(t.Components.Status.NormalBg))
	compStatusNormalFg := lipgloss.Color(t.ResolveComponent(t.Components.Status.NormalFg))
	compStatusInsertBg := lipgloss.Color(t.ResolveComponent(t.Components.Status.InsertBg))
	compStatusInsertFg := lipgloss.Color(t.ResolveComponent(t.Components.Status.InsertFg))
	compStatusCommandBg := lipgloss.Color(t.ResolveComponent(t.Components.Status.CommandBg))
	compStatusCommandFg := lipgloss.Color(t.ResolveComponent(t.Components.Status.CommandFg))
	compStatusBarBg := lipgloss.Color(t.ResolveComponent(t.Components.Status.BarBg))
	compStatusBarFg := lipgloss.Color(t.ResolveComponent(t.Components.Status.BarFg))

	glowBold := t.Effects.Glow

	switch t.Border.Style {
	case "normal":
		borderStyle = lipgloss.NormalBorder()
	case "thick":
		borderStyle = lipgloss.ThickBorder()
	case "double":
		borderStyle = lipgloss.DoubleBorder()
	default:
		borderStyle = lipgloss.RoundedBorder()
	}

	asciiStyle = lipgloss.NewStyle().
		Foreground(accentPrimary).
		Bold(true)

	dashActionStyle = lipgloss.NewStyle().
		Foreground(textSecondary)

	dashActionKeyStyle = lipgloss.NewStyle().
		Foreground(accentSecondary).
		Bold(true)

	dashActionSelectedStyle = lipgloss.NewStyle().
		Foreground(accentPrimary).
		Bold(true)

	dashRecentTitleStyle = lipgloss.NewStyle().
		Foreground(textMuted).
		Bold(true)

	dashRecentItemStyle = lipgloss.NewStyle().
		Foreground(textSecondary)

	dashRecentSelectedStyle = lipgloss.NewStyle().
		Foreground(accentPrimary).
		Bold(true)

	dashRecentDateStyle = lipgloss.NewStyle().
		Foreground(textDim)

	dashFooterStyle = lipgloss.NewStyle().
		Foreground(textDim).
		Italic(true)

	finderBorderStyle = lipgloss.NewStyle().
		BorderStyle(borderStyle).
		BorderForeground(compPanelBorder)

	finderTitleStyle = lipgloss.NewStyle().
		Foreground(compPanelTitle).
		Bold(true)

	finderInputStyle = lipgloss.NewStyle().
		Foreground(textSecondary)

	finderItemStyle = lipgloss.NewStyle().
		Foreground(textSecondary)

	finderItemSelectedStyle = lipgloss.NewStyle().
		Foreground(accentPrimary).
		Bold(true)

	finderCursorStyle = lipgloss.NewStyle().
		Foreground(accentSecondary).
		Bold(true)

	finderTagStyle = lipgloss.NewStyle().
		Foreground(accentTertiary)

	finderDateStyle = lipgloss.NewStyle().
		Foreground(textDim)

	finderPreviewTitleStyle = lipgloss.NewStyle().
		Foreground(accentPrimary).
		Bold(true)

	finderPreviewMetaStyle = lipgloss.NewStyle().
		Foreground(textMuted)

	finderPreviewTagStyle = lipgloss.NewStyle().
		Foreground(accentTertiary).
		Italic(true)

	finderPreviewBodyStyle = lipgloss.NewStyle().
		Foreground(textPrimary)

	finderPreviewDivStyle = lipgloss.NewStyle().
		Foreground(textDim)

	finderPreviewEmptyStyle = lipgloss.NewStyle().
		Foreground(textDim).
		Italic(true)

	finderCountStyle = lipgloss.NewStyle().
		Foreground(textMuted)

	inputLabelStyle = lipgloss.NewStyle().
		Foreground(compInputLabel).
		Bold(true)

	deleteWarnStyle = lipgloss.NewStyle().
		Foreground(stateError).
		Bold(true)

	helpKeyStyle = lipgloss.NewStyle().
		Foreground(accentSecondary).
		Bold(true)

	helpDescStyle = lipgloss.NewStyle().
		Foreground(textMuted)

	createBorderStyle = lipgloss.NewStyle().
		BorderStyle(borderStyle).
		BorderForeground(accentPrimary)

	createTitleStyle = lipgloss.NewStyle().
		Foreground(accentPrimary).
		Bold(true)

	createLabelStyle = lipgloss.NewStyle().
		Foreground(textMuted).
		Bold(true)

	createStepActiveStyle = lipgloss.NewStyle().
		Foreground(accentPrimary).
		Bold(true)

	createStepDoneStyle = lipgloss.NewStyle().
		Foreground(stateSuccess)

	createStepPendingStyle = lipgloss.NewStyle().
		Foreground(textDim)

	createValueStyle = lipgloss.NewStyle().
		Foreground(textSecondary)

	edGutterStyle = lipgloss.NewStyle().
		Foreground(textDim).
		Width(5).
		Align(lipgloss.Right).
		PaddingRight(1)

	edGutterActiveStyle = lipgloss.NewStyle().
		Foreground(accentSecondary).
		Bold(true).
		Width(5).
		Align(lipgloss.Right).
		PaddingRight(1)

	edGutterSepStyle = lipgloss.NewStyle().
		Foreground(textDim)

	edTextStyle = lipgloss.NewStyle().
		Foreground(textPrimary)

	edCursorBlockStyle = lipgloss.NewStyle().
		Reverse(true)

	edCursorUnderlineStyle = lipgloss.NewStyle().
		Underline(true)

	edStatusBarStyle = lipgloss.NewStyle().
		Background(compStatusBarBg).
		Foreground(compStatusBarFg)

	edModeNormalStyle = lipgloss.NewStyle().
		Background(compStatusNormalBg).
		Foreground(compStatusNormalFg).
		Bold(true).
		Padding(0, 1)

	edModeInsertStyle = lipgloss.NewStyle().
		Background(compStatusInsertBg).
		Foreground(compStatusInsertFg).
		Bold(true).
		Padding(0, 1)

	edModeCommandStyle = lipgloss.NewStyle().
		Background(compStatusCommandBg).
		Foreground(compStatusCommandFg).
		Bold(true).
		Padding(0, 1)

	edFileNameStyle = lipgloss.NewStyle().
		Background(compStatusBarBg).
		Foreground(textSecondary).
		Padding(0, 1)

	edModifiedStyle = lipgloss.NewStyle().
		Background(compStatusBarBg).
		Foreground(stateError).
		Bold(true).
		Padding(0, 1)

	edPosStyle = lipgloss.NewStyle().
		Background(compStatusBarBg).
		Foreground(textMuted).
		Padding(0, 1)

	edCmdLineStyle = lipgloss.NewStyle().
		Foreground(textSecondary)

	mdHeadingStyle = lipgloss.NewStyle().
		Foreground(compTextHeading).
		Bold(true)

	mdBoldStyle = lipgloss.NewStyle().
		Foreground(textSecondary).
		Bold(true)

	mdItalicStyle = lipgloss.NewStyle().
		Foreground(textSecondary).
		Italic(true)

	mdCodeInlineStyle = lipgloss.NewStyle().
		Foreground(compTextCode)

	mdCodeBlockStyle = lipgloss.NewStyle().
		Foreground(compTextCode)

	mdLinkStyle = lipgloss.NewStyle().
		Foreground(compTextLink)

	mdBlockquoteStyle = lipgloss.NewStyle().
		Foreground(compTextBlockquote).
		Italic(true)

	mdListMarkerStyle = lipgloss.NewStyle().
		Foreground(compTextListMarker)

	mdHrStyle = lipgloss.NewStyle().
		Foreground(textDim)

	inputTextStyle = lipgloss.NewStyle().Foreground(compInputText)
	inputCursorStyle = lipgloss.NewStyle().Foreground(compInputCursor)
	inputPlaceholderStyle = lipgloss.NewStyle().Foreground(compInputPlaceholder)
	editorBodyStyle = lipgloss.NewStyle().Foreground(compTextBody)

	if glowBold {
		asciiStyle = asciiStyle.Bold(true)
		finderCursorStyle = finderCursorStyle.Bold(true)
		finderPreviewTitleStyle = finderPreviewTitleStyle.Bold(true)
		mdHeadingStyle = mdHeadingStyle.Bold(true)
		mdCodeInlineStyle = mdCodeInlineStyle.Bold(true)
		mdCodeBlockStyle = mdCodeBlockStyle.Bold(true)
		mdLinkStyle = mdLinkStyle.Bold(true)
		mdListMarkerStyle = mdListMarkerStyle.Bold(true)
	}
}

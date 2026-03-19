package output

import (
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

type Alignment int

const (
	Left Alignment = iota
	Center
	Right
)

type TableStyle int

const (
	Rounded TableStyle = iota
	Sharp
	Minimal
	None
)

type Column struct {
	Name      string
	Priority  int
	MinWidth  int
	MaxWidth  int
	Align     Alignment
	ColorFunc func(value string) lipgloss.Style
	origIndex int // set by adaptColumns to map back to row data index
}

type Table struct {
	Columns    []Column
	Rows       [][]string
	Style      TableStyle
	SortCol    int
	SortDesc   bool
	ShowFooter bool
	FooterRow  []string
	mu         sync.RWMutex
}

func NewTable() *Table {
	return &Table{
		Columns:    []Column{},
		Rows:       [][]string{},
		Style:      Rounded,
		SortCol:    -1,
		SortDesc:   false,
		ShowFooter: false,
	}
}

func (t *Table) AddColumn(col Column) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if col.Priority == 0 {
		col.Priority = 5
	}
	if col.MinWidth == 0 {
		col.MinWidth = 10
	}
	if col.MaxWidth == 0 {
		col.MaxWidth = 50
	}
	if col.ColorFunc == nil {
		col.ColorFunc = func(s string) lipgloss.Style { return lipgloss.NewStyle() }
	}

	col.origIndex = len(t.Columns)
	t.Columns = append(t.Columns, col)
}

func (t *Table) AddRow(row []string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Rows = append(t.Rows, row)
}

func (t *Table) SetFooter(row []string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.FooterRow = row
	t.ShowFooter = true
}

func (t *Table) Render(width int) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if len(t.Columns) == 0 {
		return ""
	}

	if width <= 0 {
		width = TermWidth()
	}

	visibleCols := adaptColumns(t.Columns, width)
	colWidths := calculateColumnWidths(visibleCols, t.Rows, width)

	var result strings.Builder

	switch t.Style {
	case Rounded:
		result.WriteString(t.renderRounded(visibleCols, colWidths))
	case Sharp:
		result.WriteString(t.renderSharp(visibleCols, colWidths))
	case Minimal:
		result.WriteString(t.renderMinimal(visibleCols, colWidths))
	case None:
		result.WriteString(t.renderNone(visibleCols, colWidths))
	default:
		result.WriteString(t.renderRounded(visibleCols, colWidths))
	}

	return result.String()
}

func (t *Table) renderRounded(cols []Column, widths []int) string {
	var result strings.Builder
	theme := GetTheme()

	result.WriteString(theme.Border.Render("╭"))
	for i := range cols {
		result.WriteString(strings.Repeat("─", widths[i]+2))
		if i < len(cols)-1 {
			result.WriteString(theme.Border.Render("┬"))
		}
	}
	result.WriteString(theme.Border.Render("╮"))
	result.WriteString("\n")

	result.WriteString(theme.Border.Render("│"))
	for i, col := range cols {
		headerText := padText(col.Name, widths[i], col.Align)
		result.WriteString(" " + theme.Header.Render(headerText) + " ")
		result.WriteString(theme.Border.Render("│"))
	}
	result.WriteString("\n")

	result.WriteString(theme.Border.Render("├"))
	for i := range cols {
		result.WriteString(strings.Repeat("─", widths[i]+2))
		if i < len(cols)-1 {
			result.WriteString(theme.Border.Render("┼"))
		}
	}
	result.WriteString(theme.Border.Render("┤"))
	result.WriteString("\n")

	for _, row := range t.Rows {
		result.WriteString(theme.Border.Render("│"))
		for i, col := range cols {
			cell := rowCell(row, col)
			if cell != "" || col.origIndex < len(row) {
				cellText := truncateText(cell, widths[i])
				cellText = padText(cellText, widths[i], col.Align)
				styled := col.ColorFunc(cell)
				result.WriteString(" " + styled.Render(cellText) + " ")
			} else {
				result.WriteString(" " + strings.Repeat(" ", widths[i]) + " ")
			}
			result.WriteString(theme.Border.Render("│"))
		}
		result.WriteString("\n")
	}

	if t.ShowFooter && len(t.FooterRow) > 0 {
		result.WriteString(theme.Border.Render("├"))
		for i := range cols {
			result.WriteString(strings.Repeat("─", widths[i]+2))
			if i < len(cols)-1 {
				result.WriteString(theme.Border.Render("┼"))
			}
		}
		result.WriteString(theme.Border.Render("┤"))
		result.WriteString("\n")

		result.WriteString(theme.Border.Render("│"))
		for i, col := range cols {
			cell := rowCell(t.FooterRow, col)
			if cell != "" {
				cellText := truncateText(cell, widths[i])
				cellText = padText(cellText, widths[i], col.Align)
				result.WriteString(" " + theme.Muted.Render(cellText) + " ")
			} else {
				result.WriteString(" " + strings.Repeat(" ", widths[i]) + " ")
			}
			result.WriteString(theme.Border.Render("│"))
		}
		result.WriteString("\n")
	}

	result.WriteString(theme.Border.Render("╰"))
	for i := range cols {
		result.WriteString(strings.Repeat("─", widths[i]+2))
		if i < len(cols)-1 {
			result.WriteString(theme.Border.Render("┴"))
		}
	}
	result.WriteString(theme.Border.Render("╯"))
	result.WriteString("\n")

	return result.String()
}

func (t *Table) renderSharp(cols []Column, widths []int) string {
	var result strings.Builder
	theme := GetTheme()

	result.WriteString(theme.Border.Render("┌"))
	for i := range cols {
		result.WriteString(strings.Repeat("─", widths[i]+2))
		if i < len(cols)-1 {
			result.WriteString(theme.Border.Render("┬"))
		}
	}
	result.WriteString(theme.Border.Render("┐"))
	result.WriteString("\n")

	result.WriteString(theme.Border.Render("│"))
	for i, col := range cols {
		headerText := padText(col.Name, widths[i], col.Align)
		result.WriteString(" " + theme.Header.Render(headerText) + " ")
		result.WriteString(theme.Border.Render("│"))
	}
	result.WriteString("\n")

	result.WriteString(theme.Border.Render("├"))
	for i := range cols {
		result.WriteString(strings.Repeat("─", widths[i]+2))
		if i < len(cols)-1 {
			result.WriteString(theme.Border.Render("┼"))
		}
	}
	result.WriteString(theme.Border.Render("┤"))
	result.WriteString("\n")

	for _, row := range t.Rows {
		result.WriteString(theme.Border.Render("│"))
		for i, col := range cols {
			cell := rowCell(row, col)
			cellText := truncateText(cell, widths[i])
			cellText = padText(cellText, widths[i], col.Align)
			styled := col.ColorFunc(cell)
			result.WriteString(" " + styled.Render(cellText) + " ")
			result.WriteString(theme.Border.Render("│"))
		}
		result.WriteString("\n")
	}

	if t.ShowFooter && len(t.FooterRow) > 0 {
		result.WriteString(theme.Border.Render("├"))
		for i := range cols {
			result.WriteString(strings.Repeat("─", widths[i]+2))
			if i < len(cols)-1 {
				result.WriteString(theme.Border.Render("┼"))
			}
		}
		result.WriteString(theme.Border.Render("┤"))
		result.WriteString("\n")

		result.WriteString(theme.Border.Render("│"))
		for i, col := range cols {
			cell := rowCell(t.FooterRow, col)
			cellText := truncateText(cell, widths[i])
			cellText = padText(cellText, widths[i], col.Align)
			result.WriteString(" " + theme.Muted.Render(cellText) + " ")
			result.WriteString(theme.Border.Render("│"))
		}
		result.WriteString("\n")
	}

	result.WriteString(theme.Border.Render("└"))
	for i := range cols {
		result.WriteString(strings.Repeat("─", widths[i]+2))
		if i < len(cols)-1 {
			result.WriteString(theme.Border.Render("┴"))
		}
	}
	result.WriteString(theme.Border.Render("┘"))
	result.WriteString("\n")

	return result.String()
}

func (t *Table) renderMinimal(cols []Column, widths []int) string {
	var result strings.Builder
	theme := GetTheme()

	for i, col := range cols {
		headerText := padText(col.Name, widths[i], col.Align)
		result.WriteString(theme.Header.Render(headerText))
		if i < len(cols)-1 {
			result.WriteString("  ")
		}
	}
	result.WriteString("\n")

	for i := range cols {
		result.WriteString(strings.Repeat("─", widths[i]))
		if i < len(cols)-1 {
			result.WriteString("  ")
		}
	}
	result.WriteString("\n")

	for _, row := range t.Rows {
		for i, col := range cols {
			cell := rowCell(row, col)
			cellText := truncateText(cell, widths[i])
			cellText = padText(cellText, widths[i], col.Align)
			styled := col.ColorFunc(cell)
			result.WriteString(styled.Render(cellText))
			if i < len(cols)-1 {
				result.WriteString("  ")
			}
		}
		result.WriteString("\n")
	}

	if t.ShowFooter && len(t.FooterRow) > 0 {
		for i := range cols {
			result.WriteString(strings.Repeat("─", widths[i]))
			if i < len(cols)-1 {
				result.WriteString("  ")
			}
		}
		result.WriteString("\n")

		for i, col := range cols {
			cell := rowCell(t.FooterRow, col)
			cellText := truncateText(cell, widths[i])
			cellText = padText(cellText, widths[i], col.Align)
			result.WriteString(theme.Muted.Render(cellText))
			if i < len(cols)-1 {
				result.WriteString("  ")
			}
		}
		result.WriteString("\n")
	}

	return result.String()
}

func (t *Table) renderNone(cols []Column, widths []int) string {
	var result strings.Builder

	for i, col := range cols {
		headerText := padText(col.Name, widths[i], col.Align)
		result.WriteString(headerText)
		if i < len(cols)-1 {
			result.WriteString("  ")
		}
	}
	result.WriteString("\n")

	for _, row := range t.Rows {
		for i, col := range cols {
			cell := rowCell(row, col)
			cellText := truncateText(cell, widths[i])
			cellText = padText(cellText, widths[i], col.Align)
			result.WriteString(cellText)
			if i < len(cols)-1 {
				result.WriteString("  ")
			}
		}
		result.WriteString("\n")
	}

	if t.ShowFooter && len(t.FooterRow) > 0 {
		for i, col := range cols {
			cell := rowCell(t.FooterRow, col)
			cellText := truncateText(cell, widths[i])
			cellText = padText(cellText, widths[i], col.Align)
			result.WriteString(cellText)
			if i < len(cols)-1 {
				result.WriteString("  ")
			}
		}
		result.WriteString("\n")
	}

	return result.String()
}

// rowCell returns the cell value for a visible column, using origIndex
// to map back to the row's data slice.
func rowCell(row []string, col Column) string {
	idx := col.origIndex
	if idx < len(row) {
		return row[idx]
	}
	return ""
}

func calculateColumnWidths(cols []Column, rows [][]string, termWidth int) []int {
	widths := make([]int, len(cols))

	for i, col := range cols {
		widths[i] = len(col.Name)

		for _, row := range rows {
			cell := rowCell(row, col)
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}

		if widths[i] < col.MinWidth {
			widths[i] = col.MinWidth
		}
		if col.MaxWidth > 0 && widths[i] > col.MaxWidth {
			widths[i] = col.MaxWidth
		}
	}

	totalWidth := 0
	for _, w := range widths {
		totalWidth += w + 3
	}

	if totalWidth > termWidth {
		availableWidth := termWidth - (3 * len(cols))
		for availableWidth > 0 && totalWidth > termWidth {
			maxReducibleIdx := -1
			maxReducibleAmount := 0

			for i, col := range cols {
				if widths[i] > col.MinWidth {
					reducibleAmount := widths[i] - col.MinWidth
					if reducibleAmount > maxReducibleAmount {
						maxReducibleAmount = reducibleAmount
						maxReducibleIdx = i
					}
				}
			}

			if maxReducibleIdx == -1 {
				break
			}

			widths[maxReducibleIdx]--
			totalWidth--
		}
	}

	return widths
}

func adaptColumns(cols []Column, width int) []Column {
	if len(cols) == 0 {
		return cols
	}

	if width <= 0 {
		width = TermWidth()
	}

	breakpoint := DetectBreakpoint(width)
	minBorderWidth := 3

	// First pass: determine which columns are eligible based on breakpoint
	eligible := make([]bool, len(cols))
	for i, col := range cols {
		if breakpoint == XS && col.Priority > 3 {
			eligible[i] = false
		} else if breakpoint == SM && col.Priority > 4 {
			eligible[i] = false
		} else {
			eligible[i] = true
		}
	}

	// Second pass: from eligible columns (in original order), include those
	// that fit within the terminal width. Always include Priority 1 columns.
	// Track original index so the renderer can map row data correctly.
	var visible []Column
	totalWidth := 0

	for i, col := range cols {
		if !eligible[i] {
			continue
		}
		colTotalWidth := col.MinWidth + minBorderWidth

		if totalWidth+colTotalWidth <= width || col.Priority == 1 {
			col.origIndex = i
			visible = append(visible, col)
			totalWidth += colTotalWidth
		}
	}

	if len(visible) == 0 && len(cols) > 0 {
		visible = []Column{cols[0]}
	}

	return visible
}

func padText(text string, width int, align Alignment) string {
	if len(text) >= width {
		return text
	}

	padding := width - len(text)

	switch align {
	case Right:
		return strings.Repeat(" ", padding) + text
	case Center:
		leftPad := padding / 2
		rightPad := padding - leftPad
		return strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad)
	default:
		return text + strings.Repeat(" ", padding)
	}
}

func truncateText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return text
	}

	if len(text) <= maxWidth {
		return text
	}

	if maxWidth <= 3 {
		return strings.Repeat(".", maxWidth)
	}

	return text[:maxWidth-1] + "…"
}

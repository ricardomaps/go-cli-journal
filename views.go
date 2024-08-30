package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/viewport"
	"strings"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1).BorderForeground(lipgloss.Color("62"))
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b).BorderForeground(lipgloss.Color("62"))
	}()
)

var listStyle = func() lipgloss.Style {
	b := lipgloss.DoubleBorder()
	style := lipgloss.NewStyle().BorderStyle(b).BorderForeground(lipgloss.Color("62")).PaddingTop(1).MarginRight(2)
	return style
}()

func (m *model) listView() string {
	return lipgloss.PlaceVertical(m.screenHeight, lipgloss.Bottom, listStyle.Render(m.entries.View()))

}

func (m model) headerView() string {
	header := titleStyle.Render(m.currentSelected.title)
	lineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("62"))
	line := lineStyle.Render(strings.Repeat("─", max(0, m.entryWindow.Width-lipgloss.Width(header))))
	return lipgloss.JoinHorizontal(lipgloss.Center, header, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.entryWindow.ScrollPercent()*100))
	lineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("62"))
	line := lineStyle.Render(strings.Repeat("─", max(0, m.entryWindow.Width-lipgloss.Width(info)-m.entries.Width()+15)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m *model) windowView() string {
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight
	m.entryWindow = viewport.New(m.screenWidth-m.entries.Width(), m.screenHeight-verticalMarginHeight)
  m.entryWindow.SetContent(m.currentSelected.text)
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.entryWindow.View(), m.footerView())
}

func (m *model) summaryView() string {
	border := lipgloss.DoubleBorder()
  headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#04B575"))
	style := lipgloss.NewStyle().BorderStyle(border).BorderForeground(lipgloss.Color("62")).Width(35)
  content := m.entries.SelectedItem().(Entry).Summary()
	return lipgloss.Place(m.screenWidth - lipgloss.Width(m.listView()),
    m.screenHeight, 
    lipgloss.Center, 
    lipgloss.Center, 
    style.Render(lipgloss.JoinVertical(lipgloss.Center, headerStyle.Render("Summary\n"), content)),)
}

func(m *model) ConfirmDeleteView() string{
	border := lipgloss.DoubleBorder()
	style := lipgloss.NewStyle().BorderStyle(border).BorderForeground(lipgloss.Color("62")).Width(35).Height(10)
  content := "Are you sure you want to delete this item?"
	return lipgloss.Place(m.screenWidth - lipgloss.Width(m.listView()),
    m.screenHeight, 
    lipgloss.Center, 
    lipgloss.Center, 
    style.Render(content),
    )
}

func (m *model) customizationView() string {
	border := lipgloss.DoubleBorder()
	style := lipgloss.NewStyle().BorderStyle(border).BorderForeground(lipgloss.Color("62"))
	return lipgloss.Place(m.screenWidth, m.screenHeight, lipgloss.Center, lipgloss.Center, style.Render("Customization Options"))
}

func (m *model) entryAddingView() string {
  border := lipgloss.DoubleBorder()
  styleBorder := lipgloss.NewStyle().BorderStyle(border).BorderForeground(lipgloss.Color("62")).Padding(0, 1)

  lineStyle := lipgloss.NewStyle().MarginBottom(1).Faint(true)
  line := lineStyle.Render(strings.Repeat("-", 30))
  
  m.addingWindow.title.Placeholder = "Choose a title here"
  m.addingWindow.title.Prompt = ""
  m.addingWindow.title.CharLimit = 20
  m.addingWindow.title.Width = 20
  m.addingWindow.summary.Placeholder = "Add a short summary here"
  m.addingWindow.summary.CharLimit = 80
  m.addingWindow.summary.Width = 20
  m.addingWindow.summary.Prompt = ""
	m.addingWindow.text.Placeholder = "Write whatever you want here"
  m.addingWindow.text.CharLimit = 0

	return lipgloss.Place(m.screenWidth,
		m.screenHeight,
		lipgloss.Center,
		lipgloss.Center,
    styleBorder.Render(
		lipgloss.JoinVertical(lipgloss.Center, 
        m.addingWindow.title.View(),
        line,
        m.addingWindow.summary.View(),
        line,
        m.addingWindow.text.View())))
}

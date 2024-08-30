package main

import tea "github.com/charmbracelet/bubbletea"

const(
  title = iota
  summary
  text
)
type index struct{
  value int
}

func(m *model) removerUpdate(msg tea.Msg) (tea.Model, tea.Cmd){
  switch msg := msg.(type){
  case tea.KeyMsg:
  switch msg.String(){
    case "enter":
      m.RemoveEntry()
      m.removing = false
    case "esc", "q", "Q":
    m.removing = false
  }
  }
  return m, nil
}

func(m *model) indexFocusNext(){
  switch m.addingWindow.index{
  case 0:
    m.addingWindow.title.Blur()
    m.addingWindow.summary.Focus()
    m.addingWindow.index++
  case 1:
    m.addingWindow.summary.Blur()
    m.addingWindow.text.Focus()
    m.addingWindow.index++
  case 2:
    m.addingWindow.text.Blur()
    m.addEntry(m.addingWindow.title.Value(), m.addingWindow.summary.Value(), m.addingWindow.text.Value())
    m.addingWindow.index = 0
    m.addingEntry = false
}
}

func(m *model) indexFocusPrev(){

  switch m.addingWindow.index{
  case title:
    m.addingWindow.title.Blur()
    m.addingEntry = false 
  case summary:
    m.addingWindow.summary.Blur()
    m.addingWindow.title.Focus()
    m.addingWindow.index--
  case text:
    m.addingWindow.text.Blur()
    m.addingWindow.summary.Focus()
    m.addingWindow.index--
}
}

func(m *model) updateFields(msg tea.Msg)(tea.Model, tea.Cmd){
  var cmd tea.Cmd
  if m.addingWindow.title.Focused(){
    m.addingWindow.title, cmd = m.addingWindow.title.Update(msg)
  } else if m.addingWindow.summary.Focused(){
    m.addingWindow.summary, cmd = m.addingWindow.summary.Update(msg)
  } else if m.addingWindow.text.Focused(){
    m.addingWindow.text, cmd = m.addingWindow.text.Update(msg)
  }
  return m, cmd
}

func(m *model) adderUpdate(msg tea.Msg)(tea.Model, tea.Cmd){
  if m.addingWindow.index == title{
    m.addingWindow.title.Focus()
  }
  var cmd tea.Cmd
  temp1, temp2 := m.updateFields(msg)
  m = temp1.(*model)
  cmd = temp2
  switch msg := msg.(type){
  case tea.KeyMsg:
    switch msg.String(){
    case "esc":
        m.indexFocusPrev()
    case "enter":
      m.indexFocusNext()
  }
}
  return m, cmd
}

func(m *model) customizingUpdate(msg string)(tea.Model, tea.Cmd){
  var cmd tea.Cmd
  switch msg{
  case "q", "esc":
    m.customizing = false
  }
  return m, cmd
}

package main
import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/help"
)

type Entry struct{
  title string
  summary string
  text string
  date time.Time
}

type model struct {
  entries  list.Model
  selected bool
  currentSelected Entry
  listVisible bool
  entryWindow viewport.Model
  ready bool
  viewingEntry bool
  quitting bool
  removing bool
  addingWindow adder
  screenWidth int
  screenHeight int
  addingEntry bool
  customizing bool
  mainPageHelp help.Model
  colorPalette []string
}

type adder struct{
  title textinput.Model
  summary textinput.Model
  text textarea.Model
  index int
}

func max(a, b int) int{
  if a>b{
    return a
  }else{
    return b
  }
}

func New()* model{
  return &model{mainPageHelp: help.New()}
}

func newAdder() *adder{
  return &adder{}
}
// THIS SECTION ABOUT THE LIST OF ENTRIES
func (e Entry) FilterValue() string { return e.title }

func (e Entry) Title() string{
  return e.title
}
func (e Entry) Description() string{
  return e.date.Format("Mon, 2006/01/02")
}

func(e Entry) Text() string{
  return e.text
}

func(e Entry) Summary() string{
  return e.summary
}

type entryDelegate struct{}

type Keymap struct{
  AddEntry key.Binding
  ToggleListVisible key.Binding
  SelectEntry key.Binding
  CustomizationOptions key.Binding
  Quit key.Binding
  RemoveEntry key.Binding
  ShowMoreHelp key.Binding
}

var Keys = Keymap{
    AddEntry: key.NewBinding(
      key.WithKeys("a", "A"),
      key.WithHelp("a/A", "add new entry"),
      ),
    ToggleListVisible: key.NewBinding(
      key.WithKeys("h", "H"),
      key.WithHelp("h/H", "hide/show list"),
      ),
    SelectEntry: key.NewBinding(
      key.WithKeys("enter"),
      key.WithHelp("enter", "view entry/select other entry"),
      ),
    CustomizationOptions: key.NewBinding(
      key.WithKeys("c", "C"),
      key.WithHelp("c/C", "open customization menu"),
      ),
    Quit: key.NewBinding(
    key.WithKeys("q", "Q", "esc"),
    key.WithHelp("q/esc", "quit"),
    ),
    RemoveEntry: key.NewBinding(
    key.WithKeys("r", "R"),
    key.WithHelp("r", "remove entry"),
    ),
  ShowMoreHelp: key.NewBinding(
    key.WithKeys("?"),
    key.WithHelp("?", "show help"),
    ),

  }

func(k Keymap) FullHelp() [][]key.Binding{
  return [][]key.Binding{
    {k.SelectEntry, k.AddEntry, k.Quit, k.ShowMoreHelp},
    {k.RemoveEntry, k.ToggleListVisible, k.CustomizationOptions},
  }
}

func(k Keymap) ShortHelp() []key.Binding{
  return []key.Binding{k.SelectEntry, k.AddEntry, k.ShowMoreHelp}

}

func (m model) Init() tea.Cmd {
    return nil
}

func (m *model) initialList(width int, height int){
  m.listVisible = true
  listDelegate := list.NewDefaultDelegate()
    m.entries = list.New([]list.Item{
    Entry{title: "first entry", date: time.Now()},
    Entry{title: "second entry", date: time.Now()},
    Entry{title: "third entry", date: time.Now()},
    Entry{title: "fourth entry", date: time.Now()},
  }, 
  listDelegate,
  width,
  height)

  m.entries.SetShowStatusBar(false)
  m.entries.Paginator.SetTotalPages(7)
  m.entries.Title = "Journal Entries"
  m.entries.SetShowHelp(false)
  m.entries.Help.Ellipsis = ""
}

func(m *model) addEntry(title, summary, text string){
  m.entries.InsertItem(0, Entry{
    title: title,
    summary: summary,
    text: text,
    date: time.Now(),
  })
}

func(m *model) RemoveEntry(){
  m.entries.RemoveItem(m.entries.Index())
}

//THIS SECTION ABOUT THE VIEWPORT WITH THE ENTRY'S TEXT

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd){
  switch msg := msg.(type){
  case tea.WindowSizeMsg:
    if !m.ready{
      m.initialList(20, msg.Height-3)
      m.screenWidth = msg.Width
      m.screenHeight = msg.Height
      m.ready = true
    }
  case tea.KeyMsg:
    if m.customizing{
      return m.customizingUpdate(msg.String())
    }
    if m.addingEntry{
      return m.adderUpdate(msg)
    }
    if m.removing{
      return m.removerUpdate(msg)
    }
    switch msg.String(){
    case "q", "ctrl+c":
      m.quitting = true
      return m, tea.Quit
    case "c", "C":
      m.customizing = true
    case "a", "A":
      m.addingWindow.text = textarea.New()
      m.addingWindow.summary = textinput.New()
      m.addingWindow.title = textinput.New()
      m.addingEntry = true 
    case "enter":
      if !m.selected{
        m.selected = true
        m.currentSelected = m.entries.SelectedItem().(Entry)
        m.viewingEntry = true
      }else{
        m.selected = false
        m.viewingEntry = false
      }
    case "h", "H":
      m.listVisible = !m.listVisible
    case "r", "R":
      m.removing = true
  }
}
    var cmd tea.Cmd
    if !m.selected{
      m.entries, cmd = m.entries.Update(msg)
    }else{
      m.entryWindow, cmd = m.entryWindow.Update(msg)
    }
    return m, cmd
}

func (m model) View() string{
  if m.quitting{
    return ""
  }
  if m.ready{
    if m.removing{
      return lipgloss.JoinHorizontal(lipgloss.Center, 
        m.listView(), 
        m.ConfirmDeleteView(),
        )
    }
    if m.customizing{
      return m.customizationView()
    }
    if m.addingEntry{
      return m.entryAddingView()
    }
    if m.viewingEntry && m.listVisible{
      return lipgloss.JoinHorizontal(lipgloss.Center, 
        m.listView(), 
        m.windowView(),
        )
    }
    if !m.viewingEntry{
      return lipgloss.JoinHorizontal(lipgloss.Left, m.listView(), m.summaryView())+
      lipgloss.Place(lipgloss.Width(m.mainPageHelp.View(Keys)), lipgloss.Height(m.mainPageHelp.View(Keys)), lipgloss.Bottom, lipgloss.Right, m.mainPageHelp.View(Keys))
    }
    return lipgloss.PlaceHorizontal(m.screenWidth, 
      lipgloss.Center, 
      m.windowView(),)
  }
  return ""
}

func main(){ 
  m := New()
  p := tea.NewProgram(m, tea.WithAltScreen())
  if _, err := p.Run(); err != nil {
    fmt.Printf("Alas, there's been an error: %v", err)
    os.Exit(1)
  }
}


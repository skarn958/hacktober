package list

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/secfault-org/hacktober/internal/ui/common"
	"github.com/secfault-org/hacktober/internal/ui/components/selector"
	"github.com/secfault-org/hacktober/internal/ui/styles"
)

type KeyMap struct {
	UpDown key.Binding
	Select key.Binding
}

func defaultKeyMap() KeyMap {
	return KeyMap{
		UpDown: key.NewBinding(
			key.WithKeys("up", "down", "k", "j"),
			key.WithHelp("↑↓ j/k", "navigate"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

type Page struct {
	screen   styles.Themed
	keyMap   KeyMap
	items    []selector.IdentifiableItem
	selector *selector.Selector
	onSelect func(selector.IdentifiableItem) tea.Cmd
}

func New(scr styles.Themed, title string, items []selector.IdentifiableItem, delegate selector.ItemDelegate, onSelect func(item selector.IdentifiableItem) tea.Cmd) *Page {
	km := defaultKeyMap()

	sel := selector.New(
		scr,
		km.Select,
		[]selector.IdentifiableItem{},
		delegate,
	)

	sel.Title = title
	sel.SetShowTitle(title != "")
	sel.SetShowHelp(false)
	sel.SetShowStatusBar(false)
	sel.DisableQuitKeybindings()

	return &Page{
		screen:   scr,
		keyMap:   km,
		items:    items,
		selector: sel,
		onSelect: onSelect,
	}
}

func (p *Page) Init() tea.Cmd {
	return tea.Batch(
		p.selector.Init(),
		p.selector.SetItems(p.items),
	)
}

func (p *Page) Update(msg tea.Msg) (common.Page, tea.Cmd) {
	if sel, ok := msg.(selector.SelectMsg); ok && p.onSelect != nil {
		if cmd := p.onSelect(sel.IdentifiableItem); cmd != nil {
			return p, cmd
		}
	}
	s, cmd := p.selector.Update(msg)
	p.selector = s
	return p, cmd
}

func (p *Page) View() string {
	return p.screen.Theme().ListPage.Render(p.selector.View())
}

func (p *Page) SetSize(w, h int) {
	s := p.screen.Theme().ListPage
	p.selector.SetSize(w-s.GetHorizontalFrameSize(), h-s.GetVerticalFrameSize())
}

func (p *Page) ShortHelp() []key.Binding {
	return []key.Binding{
		p.keyMap.UpDown,
		p.keyMap.Select,
	}
}

func (p *Page) FullHelp() [][]key.Binding {
	return [][]key.Binding{{
		p.selector.KeyMap.CursorUp,
		p.selector.KeyMap.CursorDown,
		p.keyMap.Select,
	}}
}

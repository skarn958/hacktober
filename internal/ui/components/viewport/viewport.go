package viewport

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

type Viewport struct {
	*viewport.Model

	top, bottom key.Binding
}

func New(top, bottom key.Binding) *Viewport {
	vp := viewport.New()
	vp.MouseWheelEnabled = true
	return &Viewport{
		Model:  &vp,
		top:    top,
		bottom: bottom,
	}
}

func (v *Viewport) SetSize(width, height int) {
	v.Model.SetWidth(width)
	v.Model.SetHeight(height)
}

func (v *Viewport) Init() tea.Cmd {
	return v.Model.Init()
}

func (v *Viewport) Update(msg tea.Msg) (*Viewport, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, v.top):
			v.GotoTop()
		case key.Matches(msg, v.bottom):
			v.GotoBottom()
		}
	}
	vp, cmd := v.Model.Update(msg)
	v.Model = &vp
	return v, cmd
}

func (v *Viewport) View() string {
	return v.Model.View()
}

func (v *Viewport) SetContent(content string) {
	v.Model.SetContent(content)
}

func (v *Viewport) GotoTop() {
	v.Model.GotoTop()
}

func (v *Viewport) GotoBottom() {
	v.Model.GotoBottom()
}

func (v *Viewport) ScrollPercent() float64 {
	return v.Model.ScrollPercent()
}

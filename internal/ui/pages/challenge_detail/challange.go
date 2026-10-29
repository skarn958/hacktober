package challenge_detail

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
	"github.com/secfault-org/hacktober/internal/model/challenge"
	"github.com/secfault-org/hacktober/internal/session"
	"github.com/secfault-org/hacktober/internal/ui/commands"
	"github.com/secfault-org/hacktober/internal/ui/common"
	"github.com/secfault-org/hacktober/internal/ui/components/viewport"
	"github.com/secfault-org/hacktober/internal/ui/styles"
)

type ChallengeDetailPage struct {
	width, height     int
	screen            styles.Themed
	keyMap            KeyMap
	viewport          *viewport.Viewport
	selectedChallenge *challenge.Challenge
	session           *session.Session
}

func New(scr styles.Themed, session *session.Session, chall *challenge.Challenge) *ChallengeDetailPage {
	km := defaultKeyMap()

	page := &ChallengeDetailPage{
		screen:            scr,
		keyMap:            km,
		viewport:          viewport.New(km.Top, km.Bottom),
		selectedChallenge: chall,
		session:           session,
	}
	return page
}

func (c *ChallengeDetailPage) Init() tea.Cmd {
	width := min(c.width, 120)

	r, err := glamour.NewTermRenderer(
		glamour.WithStylePath("tokyo-night"),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return nil
	}
	rendered, err := r.Render(c.selectedChallenge.ChallengeMarkdown)
	if err != nil {
		return nil
	}
	c.viewport.Model.SetContent(rendered)

	return nil
}

func (c *ChallengeDetailPage) Update(msg tea.Msg) (common.Page, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, c.keyMap.Spawn) && c.session.CanStartChallenge():
			cmds = append(cmds, commands.StartChallenge(c.session, c.selectedChallenge))
		}
	case tea.WindowSizeMsg:
		c.SetSize(msg.Width, msg.Height)
	}
	v, cmd := c.viewport.Update(msg)
	c.viewport = v
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return c, tea.Batch(cmds...)
}

func (c *ChallengeDetailPage) View() string {
	_, hm := c.getMargins()
	mainStyle := c.screen.Theme().ChallengeDetail.Body.
		Height(c.height - hm)

	return mainStyle.Render(c.viewport.View())
}

func (c *ChallengeDetailPage) getMargins() (int, int) {
	hm := c.screen.Theme().ChallengeDetail.Body.GetVerticalFrameSize()
	return 0, hm
}

func (c *ChallengeDetailPage) SetSize(width, height int) {
	c.width, c.height = width, height
	_, hm := c.getMargins()
	c.viewport.SetSize(width, height-hm)
}

func (c *ChallengeDetailPage) commonHelp() []key.Binding {
	b := make([]key.Binding, 0)
	if c.session.CanStartChallenge() {
		b = append(b, c.keyMap.Spawn)
	}

	return b
}

func (c *ChallengeDetailPage) ShortHelp() []key.Binding {
	return c.commonHelp()
}

func (c *ChallengeDetailPage) FullHelp() [][]key.Binding {
	b := make([][]key.Binding, 0)
	b = append(b, c.commonHelp())
	return b
}

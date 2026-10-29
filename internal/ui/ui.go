package ui

import (
	"context"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/secfault-org/hacktober/internal/application"
	"github.com/secfault-org/hacktober/internal/session"
	"github.com/secfault-org/hacktober/internal/ui/commands"
	"github.com/secfault-org/hacktober/internal/ui/common"
	"github.com/secfault-org/hacktober/internal/ui/components/confetti"
	"github.com/secfault-org/hacktober/internal/ui/components/footer"
	"github.com/secfault-org/hacktober/internal/ui/keymap"
	"github.com/secfault-org/hacktober/internal/ui/pages/tracks"
	"github.com/secfault-org/hacktober/internal/ui/screen"
)

type Ui struct {
	ctx          context.Context
	app          *application.App
	session      *session.Session
	global       keymap.GlobalKeyMap
	screen       *screen.Screen
	navStack     []common.Page
	footer       *footer.Footer
	confetti     confetti.Model
	showConfetti bool
	inputFocused bool
}

func New(ctx context.Context, a *application.App, sess *session.Session, w, h int) *Ui {
	scr := screen.New(w, h)
	global := keymap.DefaultGlobalKeyMap()

	timeout := time.Duration(a.Config().Container.TimeoutInMinutes) * time.Minute

	ui := &Ui{
		ctx:          ctx,
		app:          a,
		session:      sess,
		screen:       scr,
		global:       global,
		navStack:     nil,
		showConfetti: false,
		confetti:     confetti.InitialModel(scr),
	}
	ui.footer = footer.New(scr, ui, timeout)
	return ui
}

func (ui *Ui) SetSize(width, height int) {
	ui.screen.Width, ui.screen.Height = width, height
}

func (ui *Ui) layout() {
	w, _ := ui.usableSize()
	ui.footer.SetWidth(w)

	// a lower footer width can influence the height so we have to calc it again
	_, h := ui.usableSize()
	for _, page := range ui.navStack {
		page.SetSize(w, h)
	}
}

func (ui *Ui) Init() tea.Cmd {
	allTracks, err := ui.app.Repo().GetAllTracks(ui.ctx)

	if err != nil {
		return common.ErrorCmd(err)
	}

	ui.layout()

	return tea.Batch(
		ui.navPush(tracks.New(ui.screen, ui.session, allTracks)),
		ui.confetti.Init(),
	)
}

func (ui *Ui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ui.screen.Width, ui.screen.Height = msg.Width, msg.Height
	case ShowConfettiMsg:
		ui.confetti.Reset()
		ui.showConfetti = true
	case HideConfettiMsg:
		ui.showConfetti = false
		cmds = append(cmds, tea.ClearScreen)
	case commands.ChallengeSolvedMsg:
		cmds = append(cmds,
			showConfettiCmd(),
			hideConfettiCmd(2*time.Second),
			commands.StopChallenge(ui.session, ui.ctx),
		)
	case commands.FlagEnteredMsg:
		ui.inputFocused = false
		cmds = append(cmds, commands.SubmitFlag(ui.session, msg.Flag))
	case tea.KeyMsg:
		if !ui.inputFocused {
			switch {
			case key.Matches(msg, ui.global.EnterFlag) && ui.session.IsChallengeRunning():
				ui.inputFocused = true
				cmds = append(cmds, commands.EnteringFlagCmd)
			case key.Matches(msg, ui.global.Quit):
				return ui, tea.Quit
			case key.Matches(msg, ui.global.Back) && ui.canGoBack():
				ui.navPop()
			case key.Matches(msg, ui.global.Help):
				cmds = append(cmds, footer.ToggleHelpCmd)
			case key.Matches(msg, ui.global.StopContainer):
				cmds = append(cmds, commands.StopChallenge(ui.session, ui.ctx))
			}
		} else {
			switch {
			case key.Matches(msg, ui.global.Back):
				ui.inputFocused = false
				cmds = append(cmds, commands.EnteringFlagCanceled)
			}
		}
	case common.PushPageMsg:
		cmds = append(cmds, ui.navPush(msg.Page))
	case commands.ChallengeStaringMsg:
		cmds = append(cmds, commands.SpawnContainer(ui.session, ui.ctx, msg))
	case commands.ChallengeStartedMsg:
		ui.app.Logger().Debugf("Container for %s started", msg.Challenge.Name)
	case commands.ChallengeStoppedMsg:
		ui.app.Logger().Debugf("Container for %s stopped", msg.Name)
	case commands.ContainerErrorMsg:
		ui.app.Logger().Error(msg.Error())
	case commands.ChallengeTimedOutMsg:
		if ui.session.CanStopChallenge() {
			cmds = append(cmds, commands.StopChallenge(ui.session, ui.ctx))
		}
	}

	var cmd tea.Cmd
	con, cmd := ui.confetti.Update(msg)
	ui.confetti = con
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	f, cmd := ui.footer.Update(msg)
	ui.footer = f
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	if !ui.inputFocused {
		m, cmd := ui.currentPage().Update(msg)
		ui.navStack[len(ui.navStack)-1] = m
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	switch msg.(type) {
	case tea.WindowSizeMsg, footer.ToggleHelpMsg:
		ui.layout()
	}
	return ui, tea.Batch(cmds...)
}

type ShowConfettiMsg struct{}
type HideConfettiMsg struct{}

func showConfettiCmd() tea.Cmd {
	return func() tea.Msg {
		return ShowConfettiMsg{}
	}
}

func hideConfettiCmd(d time.Duration) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(d)
		return HideConfettiMsg{}
	}
}

func (ui *Ui) content() string {
	if ui.showConfetti {
		return ui.confetti.View()
	}

	appStyle := ui.screen.Theme().App
	usableHeight := ui.screen.Height - appStyle.GetVerticalFrameSize()

	mainHeight := usableHeight - ui.footer.Height()

	mainView := lipgloss.NewStyle().
		Height(mainHeight).
		MaxHeight(mainHeight).
		Render(ui.currentPage().View())

	view := lipgloss.JoinVertical(lipgloss.Left,
		mainView,
		ui.footer.View(ui.session),
	)

	return ui.screen.Theme().App.Render(view)
}

func (ui *Ui) View() tea.View {
	view := tea.NewView(ui.content())
	view.AltScreen = true
	return view
}

func (ui *Ui) usableSize() (int, int) {
	appStyle := ui.screen.Theme().App
	w := ui.screen.Width - appStyle.GetHorizontalFrameSize()
	h := ui.screen.Height - appStyle.GetVerticalFrameSize() - ui.footer.Height()
	return w, h
}

func (ui *Ui) currentPage() common.Page {
	return ui.navStack[len(ui.navStack)-1]
}

func (ui *Ui) navPush(p common.Page) tea.Cmd {
	w, h := ui.usableSize()
	p.SetSize(w, h)
	ui.navStack = append(ui.navStack, p)
	return p.Init()
}

func (ui *Ui) navPop() {
	if ui.canGoBack() {
		ui.navStack = ui.navStack[:len(ui.navStack)-1]
	}
}

func (ui *Ui) canGoBack() bool { return len(ui.navStack) > 1 }

func (ui *Ui) ShortHelp() []key.Binding {
	bindings := make([]key.Binding, 0)

	bindings = append(bindings, ui.currentPage().ShortHelp()...)
	if ui.session.CanStopChallenge() {
		bindings = append(bindings, ui.global.StopContainer, ui.global.EnterFlag)
	}
	if ui.canGoBack() {
		bindings = append(bindings, ui.global.Back)
	}
	bindings = append(bindings, ui.global.Help, ui.global.Quit)

	return bindings
}

func (ui *Ui) FullHelp() [][]key.Binding {
	bindings := make([][]key.Binding, 0)

	bindings = append(bindings, ui.currentPage().FullHelp()...)
	if ui.session.CanStopChallenge() {
		bindings = append(bindings, []key.Binding{ui.global.StopContainer}, []key.Binding{ui.global.EnterFlag})
	}
	if ui.canGoBack() {
		bindings = append(bindings, []key.Binding{ui.global.Back})
	}
	bindings = append(bindings,
		[]key.Binding{
			ui.global.Help,
			ui.global.Quit,
		},
	)

	return bindings
}

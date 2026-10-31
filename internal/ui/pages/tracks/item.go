package tracks

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
	"github.com/secfault-org/hacktober/internal/model/track"
	"github.com/secfault-org/hacktober/internal/session"
	"github.com/secfault-org/hacktober/internal/ui/styles"
)

type Item struct {
	Track *track.Track
}

func (i Item) ID() string          { return i.Track.Id }
func (i Item) Title() string       { return i.Track.Title }
func (i Item) Description() string { return i.Track.Description }
func (i Item) FilterValue() string { return i.Track.Title }

func (i Item) challengeCount(solved, total int) string {
	noun := "Challenges"
	if total == 1 {
		noun = "Challenge"
	}
	return fmt.Sprintf("%d/%d %s", solved, total, noun)
}

type ItemDelegate struct {
	screen  styles.Themed
	session *session.Session
}

func NewItemDelegate(screen styles.Themed, sess *session.Session) *ItemDelegate {
	return &ItemDelegate{screen, sess}
}

func (d ItemDelegate) progress(t *track.Track) (solved, total int) {
	total = len(t.Challenges)
	for i := range t.Challenges {
		if d.session.IsSolved(t.Challenges[i].Id) {
			solved++
		}
	}
	return
}

func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	listItem := item.(Item)
	stringBuilder := strings.Builder{}

	isSelected := index == m.Index()

	itemStyle := d.screen.Theme().TrackListItem.Normal
	if isSelected {
		itemStyle = d.screen.Theme().TrackListItem.Active
	}

	solved, total := d.progress(listItem.Track)
	complete := total > 0 && solved == total

	title := listItem.Title()
	countStr := listItem.challengeCount(solved, total)
	countStyle := itemStyle.Count
	if complete {
		countStr += " 🏆"
		countStyle = d.screen.Theme().TrackListItem.Complete.Count
	}

	if m.Width()-itemStyle.Base.GetHorizontalFrameSize()-lipgloss.Width(countStr) <= 0 {
		countStr = ""
	}
	countStyle = countStyle.
		Align(lipgloss.Right).
		Width(m.Width() - itemStyle.Base.GetHorizontalFrameSize() - lipgloss.Width(title))

	count := countStyle.Render(countStr)
	title = itemStyle.Title.Render(title)
	desc := itemStyle.Desc.Render(listItem.Description())

	stringBuilder.WriteString(lipgloss.JoinHorizontal(lipgloss.Bottom, title, count))
	stringBuilder.WriteRune('\n')
	stringBuilder.WriteString(desc)
	if _, err := fmt.Fprint(w, itemStyle.Base.Render(stringBuilder.String())); err != nil {
		log.Error("Failed to render track list item", "error", err)
	}
}

func (d ItemDelegate) Height() int {
	return d.screen.Theme().TrackListItemContainer.GetVerticalFrameSize() +
		d.screen.Theme().TrackListItemContainer.GetHeight()
}

func (d ItemDelegate) Spacing() int {
	return 1
}

func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

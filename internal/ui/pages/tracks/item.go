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
	"github.com/secfault-org/hacktober/internal/ui/styles"
)

type Item struct {
	Track *track.Track
}

func (i Item) ID() string          { return i.Track.Id }
func (i Item) Title() string       { return i.Track.Title }
func (i Item) Description() string { return i.Track.Description }
func (i Item) FilterValue() string { return i.Track.Title }

func (i Item) challengeCount() string {
	n := len(i.Track.Challenges)
	if n == 1 {
		return "1 Challenge"
	}
	return fmt.Sprintf("%d Challenges", n)
}

type ItemDelegate struct {
	screen styles.Themed
}

func NewItemDelegate(screen styles.Themed) *ItemDelegate {
	return &ItemDelegate{screen}
}

func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	listItem := item.(Item)
	stringBuilder := strings.Builder{}

	isSelected := index == m.Index()

	itemStyle := d.screen.Theme().TrackListItem.Normal
	if isSelected {
		itemStyle = d.screen.Theme().TrackListItem.Active
	}

	title := listItem.Title()
	countStr := listItem.challengeCount()
	if m.Width()-itemStyle.Base.GetHorizontalFrameSize()-lipgloss.Width(countStr) <= 0 {
		countStr = ""
	}
	countStyle := itemStyle.Count.
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

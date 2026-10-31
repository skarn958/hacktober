package tracks

import (
	tea "charm.land/bubbletea/v2"
	"github.com/secfault-org/hacktober/internal/model/track"
	"github.com/secfault-org/hacktober/internal/session"
	"github.com/secfault-org/hacktober/internal/ui/common"
	"github.com/secfault-org/hacktober/internal/ui/components/selector"
	"github.com/secfault-org/hacktober/internal/ui/pages/challenges"
	"github.com/secfault-org/hacktober/internal/ui/pages/list"
	"github.com/secfault-org/hacktober/internal/ui/styles"
)

func New(scr styles.Themed, sess *session.Session, tracks []track.Track) common.Page {
	items := make([]selector.IdentifiableItem, len(tracks))

	for i := range tracks {
		items[i] = Item{Track: &tracks[i]}
	}

	return list.New(scr, "Tracks", items, NewItemDelegate(scr, sess),
		func(item selector.IdentifiableItem) tea.Cmd {
			t := item.(Item)
			return common.PushPage(challenges.New(scr, sess, t.Track.Title, t.Track.Challenges))
		})
}

package challenges

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/secfault-org/hacktober/internal/model/challenge"
	"github.com/secfault-org/hacktober/internal/session"
	"github.com/secfault-org/hacktober/internal/ui/common"
	"github.com/secfault-org/hacktober/internal/ui/components/selector"
	"github.com/secfault-org/hacktober/internal/ui/pages/challenge_detail"
	"github.com/secfault-org/hacktober/internal/ui/pages/list"
	"github.com/secfault-org/hacktober/internal/ui/styles"
)

func New(scr styles.Themed, sess *session.Session, track string, challs []challenge.Challenge) common.Page {
	items := make([]selector.IdentifiableItem, len(challs))
	for i := range challs {
		items[i] = Item{Challenge: &challs[i]}
	}

	return list.New(scr, fmt.Sprintf("%s Challenges", track), items, NewItemDelegate(scr, sess),
		func(it selector.IdentifiableItem) tea.Cmd {
			c := it.(Item)
			if c.Challenge.Locked() {
				return nil
			}
			return common.PushPage(challenge_detail.New(scr, sess, c.Challenge))
		})
}

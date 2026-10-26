package challenges

import (
	"github.com/secfault-org/hacktober/internal/model/challenge"
	"github.com/secfault-org/hacktober/internal/ui/common"
	"github.com/secfault-org/hacktober/internal/ui/components/selector"
	"github.com/secfault-org/hacktober/internal/ui/pages/list"
	"github.com/secfault-org/hacktober/internal/ui/styles"
)

func New(scr styles.Themed, challs []challenge.Challenge) common.Page {
	items := make([]selector.IdentifiableItem, len(challs))
	for i := range challs {
		items[i] = Item{Challenge: &challs[i]}
	}

	return list.New(scr, "Challenges", items, NewItemDelegate(scr))
}

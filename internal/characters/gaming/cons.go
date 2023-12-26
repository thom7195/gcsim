package gaming

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}

	c.Core.Player.Heal(player.HealInfo{
		Caller:  c.Index,
		Target:  c.Core.Player.Active(),
		Message: "Bringer of Blessing (C1)",
		Src:     c.MaxHP() * 0.25,
		Bonus:   c.Stat(attributes.Heal),
	})
}

func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}

	c.Core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		hi := args[0].(*player.HealInfo)
		overheal := args[3].(float64)

		if overheal <= 0 {
			return false
		}

		if hi.Target != c.Core.Player.Active() && hi.Target != -1 {
			return false
		}

		c2M := make([]float64, attributes.EndStatType)
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("gaming-c2", 5*60),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				c2M[attributes.ATKP] = 0.2
				return c2M, true
			},
		})

		return false
	}, "gaming-c2")
}

func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}
	c.AddEnergy("gaming-c4", 2)
}

func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}

	c6Buff := make([]float64, attributes.EndStatType)
	c6Buff[attributes.CR] = 0.2
	c6Buff[attributes.CD] = 0.4
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("gaming-c6", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagPlunge {
				return nil, false
			}
			if !strings.Contains(atk.Info.Abil, ePlungeKey) {
				return nil, false
			}
			return c6Buff, true
		},
	})
}

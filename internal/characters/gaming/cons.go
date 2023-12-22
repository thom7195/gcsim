package gaming

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
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

	// cr and cd separately to avoid stack overflow due to NoStat attribute
	mCR := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("gaming-c6-cr", -1),
		AffectedStat: attributes.CR,
		Extra:        true,
		Amount: func() ([]float64, bool) {
			mCR[attributes.CR] = 0.2
			return mCR, true
		},
	})

	mCD := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("gaming-c6-cd", -1),
		AffectedStat: attributes.CD,
		Extra:        true,
		Amount: func() ([]float64, bool) {
			mCD[attributes.CD] = 0.4
			return mCD, true
		},
	})
}

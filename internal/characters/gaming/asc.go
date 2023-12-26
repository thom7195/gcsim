package gaming

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

/*
For 1s after hitting an opponent with Bestial Ascent's Plunging Attack: Charmed Cloudstrider,

	Gaming will recover 10% of his HP.
*/
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}

	c.QueueCharTask(func() {
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Index,
			Message: "Horned Lion's Gilded Dance Healing",
			Src:     c.MaxHP() * 0.3,
			Bonus:   c.Stat(attributes.Heal),
		})
	}, 60)
}

/*
When Gaming has less than 50% HP, he will receive a 20% Incoming Healing Bonus.
When Gaming has 50% HP or more, he will gain a 20% Pyro DMG Bonus.
*/
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	// Healing part
	mHeal := make([]float64, attributes.EndStatType)
	mHeal[attributes.Heal] = 0.2
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("gaming-a4-heal-bonus", -1),
		AffectedStat: attributes.Heal,
		Amount: func() ([]float64, bool) {
			active := c.Core.Player.ActiveChar()
			if active.CurrentHPRatio() < 0.5 {
				return mHeal, true
			}
			return nil, false
		},
	})

	a4Buff := make([]float64, attributes.EndStatType)
	a4Buff[attributes.PyroP] = 0.2
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("gaming-a4-dmg-bonus", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagPlunge {
				return nil, false
			}
			if c.CurrentHPRatio() < 0.5 {
				return nil, false
			}
			if !strings.Contains(atk.Info.Abil, ePlungeKey) {
				return nil, false
			}
			return a4Buff, true
		},
	})
}

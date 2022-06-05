package gorou

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

//When characters (other than Gorou) within the AoE of Gorou's General's War Banner
//or General's Glory deal Geo DMG to opponents, the CD of Gorou's Inuzaka All-Round Defense
//is decreased by 2s. This effect can occur once every 10s.
func (c *char) c1() {
	icd := -1
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		if c.Core.Status.Duration(generalGloryKey) == 0 && c.Core.Status.Duration(generalWarBannerKey) == 0 {
			return false
		}
		if icd > c.Core.F {
			return false
		}
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex == c.Index {
			return false
		}
		if atk.Info.Element != attributes.Geo {
			return false
		}
		icd = c.Core.F + 600
		c.ReduceActionCooldown(action.ActionSkill, 120)
		return false
	}, "gorou-c1")
}

//While General's Glory is in effect, its duration is extended by 1s when a nearby
//active character obtains an Elemental Shard from a Crystallize reaction.
//This effect can occur once every 0.1s. Max extension is 3s.
func (c *char) c2() {
	c.Core.Events.Subscribe(event.OnShielded, func(args ...interface{}) bool {
		if c.Core.Status.Duration(generalGloryKey) <= 0 {
			return false
		}
		if c.c2Extension >= 3 {
			return false
		}
		c.c2Extension++
		c.Core.Status.Add(generalGloryKey, 60)
		return false
	}, "gorou-c2")
}

func (c *char) c6() {
	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod(c6key, 720, func(ae *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if ae.Info.Element != attributes.Geo {
				return nil, false
			}
			return c.c6buff, true
		})
	}
}

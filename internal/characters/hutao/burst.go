package hutao

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var burstFrames []int

const burstHitmark = 66

func init() {
	burstFrames = frames.InitAbilSlice(100)
	burstFrames[action.ActionAttack] = 98
	burstFrames[action.ActionSkill] = 97
	burstFrames[action.ActionDash] = 98
	burstFrames[action.ActionSwap] = 95
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	low := (c.HPCurrent / c.MaxHP()) <= 0.5
	mult := burst[c.TalentLvlBurst()]
	regen := regen[c.TalentLvlBurst()]
	if low {
		mult = burstLow[c.TalentLvlBurst()]
		regen = regenLow[c.TalentLvlBurst()]
	}
	targets := p["targets"]
	//regen for p+1 targets, max at 5; if not specified then p = 1
	count := 1
	if targets > 0 {
		count = targets
	}
	if count > 5 {
		count = 5
	}
	c.Core.Player.Heal(player.HealInfo{
		Caller:  c.Index,
		Target:  c.Index,
		Message: "Spirit Soother",
		Src:     c.MaxHP() * float64(count) * regen,
		Bonus:   c.Stat(attributes.Heal),
	})

	//[2:28 PM] Aluminum | Harbinger of Jank: I think the idea is that PP won't fall off before dmg hits, but other buffs aren't snapshot
	//[2:29 PM] Isu: yes, what Aluminum said. PP can't expire during the burst animation, but any other buff can
	if burstHitmark > c.Core.Status.Duration("paramita") && c.Core.Status.Duration("paramita") > 0 {
		c.Core.Status.Add("paramita", burstHitmark) //extend this to barely cover the burst
		c.Core.Log.NewEvent("Paramita status extension for burst", glog.LogCharacterEvent, c.Index).
			Write("new_duration", c.Core.Status.Duration("paramita"))
	}

	if c.Core.Status.Duration("paramita") > 0 && c.Base.Cons >= 2 {
		c.applyBB()
	}

	//TODO: apparently damage is based on stats on contact, not at cast
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spirit Soother",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 50,
		Mult:       mult,
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(5, false, combat.TargettableEnemy), burstHitmark, burstHitmark)

	c.ConsumeEnergy(68)
	c.SetCDWithDelay(action.ActionBurst, 900, 62)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

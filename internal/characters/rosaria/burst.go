package rosaria

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(70)
	burstFrames[action.ActionDash] = 57
	burstFrames[action.ActionJump] = 59
	burstFrames[action.ActionSwap] = 69
}

// Burst attack damage queue generator
// Rosaria swings her weapon to slash surrounding opponents, then she summons a frigid Ice Lance that strikes the ground. Both actions deal Cryo DMG.
// While active, the Ice Lance periodically releases a blast of cold air, dealing Cryo DMG to surrounding opponents.
// Also includes the following effects: A4, C6
func (c *char) Burst(p map[string]int) action.ActionInfo {
	// Note - if a more advanced targeting system is added in the future
	// hit 1 is technically only on surrounding enemies, hits 2 and dot are on the lance
	// For now assume that everything hits all targets
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Rites of Termination (Hit 1)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       burst[0][c.TalentLvlBurst()],
	}
	x, y := c.Core.Combat.Target(0).Pos()
	// Hit 1 comes out on frame 15
	// 2nd hit comes after lance drop animation finishes
	c.Core.QueueAttack(ai, combat.NewCircleHit(x, y, 1, false, combat.TargettableEnemy), 15, 15, c.c6)

	ai.Abil = "Rites of Termination (Hit 2)"
	ai.Mult = burst[1][c.TalentLvlBurst()]

	//lance lands at 60f/1s
	c.Core.QueueAttack(ai, combat.NewCircleHit(0, 0, 2, false, combat.TargettableEnemy), 60, 60, c.c6)

	//duration is 8 second (extended by c2 by 4s), + 0.5
	dur := 510
	if c.Base.Cons >= 2 {
		dur += 240
	}

	// Burst is snapshot when the lance lands (when the 2nd damage proc hits)
	ai = combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Rites of Termination (DoT)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       burstDot[c.TalentLvlBurst()],
	}

	c.Core.Tasks.Add(func() {
		// dot every 2 second after lance lands
		for i := 120; i < dur; i += 120 {
			c.Core.QueueAttack(ai, combat.NewCircleHit(0, 0, 2, false, combat.TargettableEnemy), 0, i, c.c6)
		}
	}, 60)

	c.Core.Status.Add("rosariaburst", dur)

	// Handle A4
	// Casting Rites of Termination increases CRIT RATE of all nearby party members, excluding Rosaria herself, by 15% of Rosaria's CRIT RATE for 10s. CRIT RATE bonus gained this way cannot exceed 15%.
	// Uses the snapshot generated by hit #1 to ensure all mods are accounted for.
	// Confirmed via testing that mods like Rosaria A1 are accounted for, and Blizzard Strayer modifications are not
	crit_share := 0.15 * c.Stat(attributes.CR)
	if crit_share > 0.15 {
		crit_share = 0.15
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = crit_share
	for i, char := range c.Core.Player.Chars() {
		// skip Rosaria
		if i == c.Index {
			continue
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("rosaria-a4", 600),
			AffectedStat: attributes.CR,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
	c.Core.Log.NewEvent("Rosaria A4 activation", glog.LogCharacterEvent, c.Index).
		Write("ends_on", c.Core.F+600).
		Write("crit_share", crit_share)

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(6)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}

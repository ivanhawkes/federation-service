package Skills

import (
	"time"
)

//
// Combat skills.
//

// DamgeType descrribes the general form of the damage to be taken and allows specific
// situational resists to be applied.
const (
	DamageTypeAcid = iota
	DamageTypeCold
	DamageTypeCrushing
	DamageTypeDecay
	DamageTypeDisease
	DamageTypeElectricity
	DamageTypeEntropy
	DamageTypeFire
	DamageTypeHoly
	DamageTypeIce
	DamageTypeNature
	DamageTypePiercing
	DamageTypePoison
	DamageTypePlasma
	DamageTypeRadiation
	DamageTypeSlashing
	DamageTypeTearing
	DamageTypeUnholy
)

// Different types of healing which might be used in games.
const (
	HealTypeQi = iota
	HealTypePhysical
	HealTypeHoly
	HealTypeNature
	HealTypeSpirit
)

// The area of affect for the skill.
const (
	AreaOfEffectSingleTarget = iota
	AreaOfEffectCone
	AreaOfEffectColumn
	AreaOfEffectChain
	AreaOfEffectPBAOE
	AreaOfEffectTBAOE
	AreaOfEffectGTAOE
)

// Might not use this idea grabbed from TSW.
const (
	AttackTypeBasic = iota
	AttackTypeStrike
	AttackTypeBurst
	AttackTypeChannel
	AttackTypeDoT
)

// A player based area of effect that provides a buff for friendlies.
type AuraBuff struct {
	BuffId   int64
	Range    float32
	Duration float32
}

// A player based area of effect that provides a debuff on enemies.
type AuraDebuff struct {
	DebuffId int64
	Range    float32
	Duration float32
}

// Drops a banner which is stationary on the ground and provides a friendly buff. This can also
// be used for consecrate and other ground effects if made flexible enough.
type BannerBuff struct {
	BuffId   int64
	Range    float32
	Duration float32
}

// Drops a banner which is stationary on the ground and provides an enemy debuff. This can also
// be used for consecrate and other ground effects if made flexible enough.
type BannerDebuff struct {
	DebuffId int64
	Range    float32
	Duration float32
}

// Target is unable to see for duration of the skill.
type Blind struct {
	Duration float32
}

// Player breaks out of crowd control.
type BreakCC struct {
	BreakCharm      bool
	BreakFear       bool
	BreakFlee       bool
	BreakPolymorph  bool
	BreakPossession bool
	BreakBlind      bool
	BreakSilence    bool
	BreakSlow       bool
	BreakSnare      bool
	BreakStun       bool
	BreakTaunt      bool
}

// Combat charge e.g. warriors
type Charge struct {
	Distance float32
}

// Target is charmed and unable to attack the player for the duration.
type Charm struct {
	Duration float32
}

// Removes all harmful effects.
type Cleanse struct {
}

// Removes 1 or 2 harmful effects in the given DamageType.
type Cure struct {
	DamageType uint8
	MaxEffects uint8
}

// Removes 1 or 2 harmful effects of any DamageType.
type CureAny struct {
	MaxEffects uint8
}

// Straight up damage given in DPS.
type Damage struct {
	DamageType uint8
	DPS        float32
}

// Removes 1 buff from the target.
type Dispel struct {
}

type DoT struct {
	DamageType uint8
	DPS        float32
	Duration   float32

	// Number of times it can stack.
	Stack uint8
}

// Target is disarmed for the duration.
type Disarm struct {
	Duration float32
}

// Target quakes in fear for the duration.
type Fear struct {
	Duration float32
}

// Target flees in fear for the duration.
type Flee struct {
	Duration float32
}

// Straight up heal given in HPS.
type Damage struct {
	HPS float32
}

// Heal over time.
type HoT struct {
	HealType uint8
	DPS      float32
	Duration float32

	// Number of times it can stack.
	Stack uint8
}

type Knockback struct {
	Force    float32
	Duration float32
}

type Knockdown struct {
	Force    float32
	Duration float32
}

type Possession struct {
	Duration float32
}

type Polymorph struct {
	Duration float32
}

type Pull struct {
	Force float32
}

// Target has their Qi reduced by the Amount given, every second for the duration of the skill.
// Amount is Qi / second.
type QiDrain struct {
	Amount   float32
	Duration float32
}

// Restore amount of Qi to the player / party / friendlies.
// Amount is Qi / second.
type QiRestore struct {
	Amount   float32
	Duration float32
}

// Provides a buff to friendlies.
type ShoutBuff struct {
	BuffId   int64
	Range    float32
	Duration float32
}

// Places a debuff on enemies.
type ShoutDebuff struct {
	DebuffId int64
	Range    float32
	Duration float32
}

// Silence all types of spoken / Qi damage. Damage types that do not rely on a spoken component
// should be unaffected.
type Silence struct {
	Duration float32
}

// Silence just one type of damage.
type SilenceDamageType struct {
	DamageType uint8
	Duration   float32
}

// Slow the target down by the given percentage.
type Slow struct {
	Percent  float32
	Duration float32
}

// Target is snared in place for the duration.
type Snare struct {
	Duration float32

	// True if the player can spin on the spot while snared.
	CanTurn bool
}

// Player has been stunned and is unable to react for the duration.
type Stun struct {
	Duration float32
}

// Target is forced to attack the caster for the duration. In some games this might redirect
// the damage to another player and even reduce it e.g. PVP in SWTOR. For PVE - maybe the same,
// or maybe the standard model.
type Taunt struct {
	Duration float32
}

// Apply a vulnerability to the damage type for the duration.
type Vulnerability struct {
	Amount     float32
	DamageType uint8

	// Number of times it can stack.
	Stack uint8
}

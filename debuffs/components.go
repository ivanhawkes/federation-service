package Debuffs

import (
	"time"
)

// TODO: Implement a cost system...resource points, currency, XP?, skill point?

type DPS struct {
	Amount float32
}

type Qi struct {
	Amount float32
}

// Stats - hand, foot, etc

// Special abilities e.g. curse of agony, shadow word pain

type Skill struct {
	SkillId  int64
	Duration int32
}

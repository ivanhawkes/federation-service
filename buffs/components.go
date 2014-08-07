package buffs

import (
	"appengine/datastore"
	// 	"time"
)

// TODO: Implement a cost system...resource points, currency, XP?, skill point?

type DPS struct {
	Amount float32
}

type Qi struct {
	Amount float32
}

// Stats - hand, foot, etc

// Special abilities e.g. shield, power infusion.
type Skill struct {
	SkillId  datastore.Key
	Duration int32
}

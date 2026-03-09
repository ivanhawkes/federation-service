package missions

import (
	"google.golang.org/appengine/datastore"
)

// Gain XP to get new skill points. XP levels are flat, like in TSW. This is simply a way to evenly
// distribute skill points over a constant amount of time / effort.

type RewardXp struct {
	Xp int `json:"xp" xml:"XP"`
}

// Types include gold, PVP badges, dungeon badges, etc

type RewardCurrency struct {
	Type   datastore.Key `json:"type" xml:"type"`
	Amount int           `json:"amount" xml:"amount"`
}

type RewardItem struct {
	ItemId datastore.Key `json:"item_id" xml:"item-id"`
	Amount int           `json:"amount" xml:"amount"`
}

type RewardTitle struct {
	TitleId datastore.Key `json:"title_id" xml:"title-id"`
}

//
// Skills: points, categories of skills? e.g. cooking, combat, gathering? With categories we can make them
// perform a mission to unlock new skill categories.

type RewardSkillPoints struct {
	Amount int `json:"amount" xml:"amount"`
}

type RewardSkillCategory struct {
	SkillCategoryId datastore.Key `json:"skill_category_id" xml:"skill-category-id"`
}

type RewardFaction struct {
	FactionId datastore.Key `json:"faction_id" xml:"faction-id"`
	Amount    int           `json:"amount" xml:"amount"`
}

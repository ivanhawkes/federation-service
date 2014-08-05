package missions

// import (
// 	"appengine"
// 	"appengine/datastore"
// )

// Gain XP to get new skill points. XP levels are flat, like in TSW. This is simply a way to evenly
// distribute skill points over a constant amount of time / effort.

type RewardXp struct {
	Xp int `json:"xp" xml:"XP"`
}

// Types include gold, PVP badges, dungeon badges, etc

type RewardCurrency struct {
	Type   int64 `json:"type" xml:"Type"`
	Amount int   `json:"amount" xml:"Amount"`
}

type RewardItem struct {
	ItemId int64 `json:"item-id" xml:"ItemId"`
	Amount int   `json:"amount" xml:"Amount"`
}

type RewardTitle struct {
	TitleId int64 `json:"title-id" xml:"TitleId"`
}

//
// Skills: points, categories of skills? e.g. cooking, combat, gathering? With categories we can make them
// perform a mission to unlock new skill categories.

type RewardSkillPoints struct {
	Amount int `json:"amount" xml:"Amount"`
}

type RewardSkillCategory struct {
	SkillCategoryId int64 `json:"skill-category-id" xml:"SkillCategoryId"`
}

type RewardFaction struct {
	FactionId int64 `json:"faction-id" xml:"FactionId"`
	Amount    int   `json:"amount" xml:"Amount"`
}

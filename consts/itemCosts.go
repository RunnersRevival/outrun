package consts

import "github.com/RunnersRevival/outrun/enums"

// Seemingly not reflected in game? Game is just using OG Runners prices currently.
var ItemPrices = map[string]int64{
	enums.ItemIDStrInvincible: 2000,
	enums.ItemIDStrBarrier:    3000,
	enums.ItemIDStrMagnet:     3000,
	enums.ItemIDStrTrampoline: 2000,
	enums.ItemIDStrCombo:      5000,
	enums.ItemIDStrLaser:      2000,
	enums.ItemIDStrDrill:      2000,
	enums.ItemIDStrAsteroid:   1000,

	enums.ItemIDStrBoostScore:      6000,
	enums.ItemIDStrBoostTrampoline: 1000,
	enums.ItemIDStrBoostSubChara:   2000,
}

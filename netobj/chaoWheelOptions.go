package netobj

import (
	"github.com/RunnersRevival/outrun/consts"
	"github.com/RunnersRevival/outrun/enums"
	"github.com/RunnersRevival/outrun/obj"
	"github.com/jinzhu/now"
)

type ChaoWheelOptions struct {
	Rarity               []int64        `json:"rarity" db:"rarity"`
	ItemWeight           []int64        `json:"itemWeight" db:"item_weight"`
	CampaignList         []obj.Campaign `json:"campaignList" db:"campaign_list"`
	SpinCost             int64          `json:"spinCost" db:"spin_cost"`
	ChaoRouletteType     int64          `json:"chaoRouletteType" db:"chao_roulette_type"` // value from enums.ChaoWheelType*
	NumSpecialEgg        int64          `json:"numSpecialEgg" db:"num_special_egg"`
	RouletteAvailable    int64          `json:"rouletteAvailable" db:"roulette_available"`         // flag
	NumChaoRouletteToken int64          `json:"numChaoRouletteToken" db:"num_chao_roulette_token"` // number of premium roulette tickets
	NumChaoRoulette      int64          `json:"numChaoRoulette" db:"num_chao_roulette"`            // == 0 --> chaoWheelOptions.IsTutorial
	StartTime            int64          `json:"startTime" db:"start_time"`                         // TODO: Is this needed?
	EndTime              int64          `json:"endTime" db:"end_time"`                             // TODO: Is this needed?
}

func NewChaoWheelOptions(rarity, itemWeight []int64, campaignList []obj.Campaign, spinCost, chaoRouletteType, numSpecialEgg, rouletteAvailable, numChaoRouletteToken, numChaoRoulette, startTime, endTime int64) ChaoWheelOptions {
	return ChaoWheelOptions{
		rarity,
		itemWeight,
		campaignList,
		spinCost,
		chaoRouletteType,
		numSpecialEgg,
		rouletteAvailable,
		numChaoRouletteToken,
		numChaoRoulette,
		startTime,
		endTime,
	}
}

func DefaultChaoWheelOptions(playerState PlayerState) ChaoWheelOptions {
	rarity := []int64{2, 1, 100, 0, 2, 0, 100, 1}
	//rarity := []int64{0, 1, 2, 100, 0, 1, 2, 100}      // TODO: REMOVE!
	itemWeight := []int64{6, 17, 5, 17, 16, 17, 5, 17} // Could possibly fake these, but the logic shouldn't allow it to happen
	campaignList := []obj.Campaign{}
	chaoRouletteType := enums.ChaoWheelTypeNormal
	numSpecialEgg := playerState.ChaoEggs
	if numSpecialEgg >= 10 {
		chaoRouletteType = enums.ChaoWheelTypeSpecial
		rarity = []int64{2, 1, 100, 2, 1, 2, 100, 1}
		itemWeight = []int64{1, 1, 1, 1, 1, 1, 1, 1}
	}
	rouletteAvailable := int64(1)
	numChaoRouletteToken := playerState.NumChaoRouletteTicket
	spinCost := consts.ChaoRouletteTicketCost
	if numChaoRouletteToken <= 0 { // if out of chao roulette tickets
		// use higher value for the red rings
		spinCost = consts.ChaoRouletteRedRingCost
	}
	numChaoRoulette := int64(1)
	startTime := now.BeginningOfDay().UTC().Unix() + 32400 // 12 AM + 9 hours = 9 AM
	endTime := startTime + 86399                           // 23:59:59 later
	return NewChaoWheelOptions(rarity, itemWeight, campaignList, spinCost, chaoRouletteType, numSpecialEgg, rouletteAvailable, numChaoRouletteToken, numChaoRoulette, startTime, endTime)
}

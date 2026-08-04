package requests

import "github.com/RunnersRevival/outrun/obj"

type LoginRequest struct {
	Version      string `json:"version"`
	Device       string `json:"device"`
	Seq          int64  `json:"seq,string"`
	Platform     int64  `json:"platform,string"`
	Language     int64  `json:"language,string"`
	SalesLocate  int64  `json:"salesLocate,string"`
	StoreID      int64  `json:"storeId,string"`
	PlatformSNS  int64  `json:"platform_sns,string"`
	obj.LineAuth `json:"lineAuth"`
}

type LoginBonusSelectRequest struct {
	RewardID          int64 `json:"rewardId,string"`
	RewardDays        int64 `json:"rewardDays,string"`
	RewardSelect      int64 `json:"rewardSelect,string"`
	FirstRewardDays   int64 `json:"firstRewardDays,string"`
	FirstRewardSelect int64 `json:"firstRewardSelect,string"`
}

type GetMigrationPasswordRequest struct {
	UserPassword string `json:"userPassword"`
}

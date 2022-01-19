package muxhandlers

import (
	"github.com/Mtbcooler/outrun/db"
	"github.com/Mtbcooler/outrun/emess"
	"github.com/Mtbcooler/outrun/enums"
	"github.com/Mtbcooler/outrun/helper"
	"github.com/Mtbcooler/outrun/netobj"
	"github.com/Mtbcooler/outrun/responses"
	"github.com/Mtbcooler/outrun/status"
)

func GetItemStockNum(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	// TODO: Flesh out properly! The game responds with
	// [IDRouletteTicketPremium, IDRouletteTicketItem, IDSpecialEgg]
	// for item IDs, along with an event ID, likely for event characters.
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultItemStockNum(baseInfo)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err := helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

// Raid boss roulette endpoints

func GetRaidbossWheelOptions(helper *helper.Helper) {
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	wheelOptions := netobj.DefaultRaidbossWheelOptions(0, player.PlayerState.ChaoEggs, 0, enums.WheelRankNormal, 0)
	response := responses.RaidbossWheelOptions(baseInfo, wheelOptions)
	err = helper.SendInsecureResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetPrizeRaidbossWheelSpin(helper *helper.Helper) {
	// agnostic
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultPrizeRaidbossWheel(baseInfo)
	err := helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

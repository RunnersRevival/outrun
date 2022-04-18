package muxhandlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/Mtbcooler/outrun/analytics"
	"github.com/Mtbcooler/outrun/analytics/factors"
	"github.com/Mtbcooler/outrun/config"
	"github.com/Mtbcooler/outrun/config/campaignconf"
	"github.com/Mtbcooler/outrun/config/gameconf"
	"github.com/Mtbcooler/outrun/consts"
	"github.com/Mtbcooler/outrun/db"
	"github.com/Mtbcooler/outrun/db/dbaccess"
	"github.com/Mtbcooler/outrun/emess"
	"github.com/Mtbcooler/outrun/enums"
	"github.com/Mtbcooler/outrun/helper"
	"github.com/Mtbcooler/outrun/logic/campaign"
	"github.com/Mtbcooler/outrun/logic/conversion"
	"github.com/Mtbcooler/outrun/logic/gameplay"
	"github.com/Mtbcooler/outrun/netobj"
	"github.com/Mtbcooler/outrun/obj"
	"github.com/Mtbcooler/outrun/obj/constobjs"
	"github.com/Mtbcooler/outrun/requests"
	"github.com/Mtbcooler/outrun/responses"
	"github.com/Mtbcooler/outrun/status"
	"github.com/jinzhu/now"
)

func GetDailyChallengeData(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DailyChallengeData(baseInfo, player.PlayerState.NumDailyChallenge, player.PlayerState.NextNumDailyChallenge)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetCostList(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	// no player, agonstic
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefinedCostList(baseInfo, gameconf.CFile.ContinueCost)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err := helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetMileageData(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultMileageData(baseInfo, player)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetCampaignList(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	campaignList := []obj.Campaign{}
	if campaignconf.CFile.AllowCampaigns {
		for _, confCampaign := range campaignconf.CFile.CurrentCampaigns {
			newCampaign := conversion.ConfiguredCampaignToCampaign(confCampaign)
			campaignList = append(campaignList, newCampaign)
		}
	}
	helper.DebugOut("Campaign list: %v", campaignList)
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.CampaignList(baseInfo, campaignList)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err := helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func QuickActStart(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	recv := helper.GetGameRequest()
	var request requests.QuickActStartRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	responseStatus := status.OK
	campaignList := []obj.Campaign{}
	if campaignconf.CFile.AllowCampaigns {
		for _, confCampaign := range campaignconf.CFile.CurrentCampaigns {
			newCampaign := conversion.ConfiguredCampaignToCampaign(confCampaign)
			campaignList = append(campaignList, newCampaign)
		}
	}
	helper.DebugOut("Campaign list: %v", campaignList)
	// consume items
	modToStringSlice := func(ns []int64) []string {
		result := []string{}
		for _, n := range ns {
			result = append(result, fmt.Sprintf("%v", n))
		}
		return result
	}

	//update energy counter
	for time.Now().UTC().Unix() >= player.PlayerState.EnergyRenewsAt && player.PlayerState.Energy < player.PlayerVarious.EnergyRecoveryMax {
		player.PlayerState.Energy++
		player.PlayerState.EnergyRenewsAt += player.PlayerVarious.EnergyRecoveryTime
	}
	if player.PlayerState.Energy > 0 {
		if gameconf.CFile.EnableEnergyConsumption {
			if player.PlayerState.Energy >= player.PlayerVarious.EnergyRecoveryMax {
				player.PlayerState.EnergyRenewsAt = time.Now().UTC().Unix() + player.PlayerVarious.EnergyRecoveryTime
			}
			player.PlayerState.Energy--
		}
		player.PlayerState.NumPlaying++
		if !gameconf.CFile.AllItemsFree {
			consumedItems := modToStringSlice(request.Modifier)
			consumedRings := gameplay.GetRequiredItemPayment(consumedItems, player)
			for _, citemID := range consumedItems {
				if citemID[:2] == "11" { // boosts, not items
					continue
				}
				index := player.IndexOfItem(citemID)
				if index == -1 {
					helper.Uncatchable(fmt.Sprintf("Player sent bad item ID '%s', cannot continue", citemID))
					helper.InvalidRequest()
					return
				}
				if player.PlayerState.Items[index].Amount >= 1 { // can use item
					player.PlayerState.Items[index].Amount--
				}
			}
			if player.PlayerState.NumRings < consumedRings { // not enough rings
				responseStatus = status.NotEnoughRings
			}
			player.PlayerState.NumRings -= consumedRings
		}
		helper.DebugOut(fmt.Sprintf("%v", player.PlayerState.Items))
	} else {
		responseStatus = status.NotEnoughEnergy
	}
	baseInfo := helper.BaseInfo(emess.OK, responseStatus)
	response := responses.DefaultQuickActStart(baseInfo, player, campaignList)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
	err = db.SavePlayer(player)
	if err != nil {
		helper.InternalErr("Error saving player", err)
		return
	}
	_, err = analytics.Store(player.ID, factors.AnalyticTypeTimedStarts)
	if err != nil {
		helper.WarnErr("Error storing analytics (AnalyticTypeTimedStarts)", err)
	}
}

func ActStart(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	recv := helper.GetGameRequest()
	var request requests.ActStartRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	helper.DebugOut(fmt.Sprintf("%v", player.PlayerState.Items))
	responseStatus := status.OK
	campaignList := []obj.Campaign{}
	if campaignconf.CFile.AllowCampaigns {
		for _, confCampaign := range campaignconf.CFile.CurrentCampaigns {
			newCampaign := conversion.ConfiguredCampaignToCampaign(confCampaign)
			campaignList = append(campaignList, newCampaign)
		}
	}
	helper.DebugOut("Campaign list: %v", campaignList)
	// consume items
	modToStringSlice := func(ns []int64) []string {
		result := []string{}
		for _, n := range ns {
			result = append(result, fmt.Sprintf("%v", n))
		}
		return result
	}

	//update energy counter
	for time.Now().UTC().Unix() >= player.PlayerState.EnergyRenewsAt && player.PlayerState.Energy < player.PlayerVarious.EnergyRecoveryMax {
		player.PlayerState.Energy++
		player.PlayerState.EnergyRenewsAt += player.PlayerVarious.EnergyRecoveryTime
	}
	if player.PlayerState.Energy > 0 {
		if gameconf.CFile.EnableEnergyConsumption {
			if player.PlayerState.Energy >= player.PlayerVarious.EnergyRecoveryMax {
				player.PlayerState.EnergyRenewsAt = time.Now().UTC().Unix() + player.PlayerVarious.EnergyRecoveryTime
			}
			player.PlayerState.Energy--
		}
		player.PlayerState.NumPlaying++
		if !gameconf.CFile.AllItemsFree {
			consumedItems := modToStringSlice(request.Modifier)
			consumedRings := gameplay.GetRequiredItemPayment(consumedItems, player)
			for _, citemID := range consumedItems {
				if citemID[:2] == "11" { // boosts, not items
					continue
				}
				index := player.IndexOfItem(citemID)
				if index == -1 {
					helper.Uncatchable(fmt.Sprintf("Player sent bad item ID '%s', cannot continue", citemID))
					helper.InvalidRequest()
					return
				}
				if player.PlayerState.Items[index].Amount >= 1 { // can use item
					player.PlayerState.Items[index].Amount--
				}
			}
			if player.PlayerState.NumRings < consumedRings { // not enough rings
				responseStatus = status.NotEnoughRings
			}
			player.PlayerState.NumRings -= consumedRings
		}
		helper.DebugOut(fmt.Sprintf("%v", player.PlayerState.Items))
	} else {
		responseStatus = status.NotEnoughEnergy
	}
	baseInfo := helper.BaseInfo(emess.OK, responseStatus)
	response := responses.DefaultActStart(baseInfo, player, campaignList)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
	err = db.SavePlayer(player)
	if err != nil {
		helper.InternalErr("Error saving player", err)
		return
	}
	_, err = analytics.Store(player.ID, factors.AnalyticTypeStoryStarts)
	if err != nil {
		helper.WarnErr("Error storing analytics (AnalyticTypeStoryStarts)", err)
	}
}

func ActRetry(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	redRingContinuePrice := gameconf.CFile.ContinueCost
	// TODO: Add campaign support
	responseStatus := status.OK
	if player.PlayerState.NumRedRings >= int64(redRingContinuePrice) { //does the player actually have enough red rings?
		//if so, subtract 5 red rings and respond with an OK status
		player.PlayerState.NumRedRings -= int64(redRingContinuePrice)
		err = db.SavePlayer(player)
		if err != nil {
			helper.InternalErr("Error saving player", err)
			return
		}
	} else {
		//if not, respond with NotEnoughRedRings
		responseStatus = status.NotEnoughRedRings
	}
	baseInfo := helper.BaseInfo(emess.OK, responseStatus)
	response := responses.NewBaseResponse(baseInfo)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
	if responseStatus == status.OK {
		_, err = analytics.Store(player.ID, factors.AnalyticTypeRevives)
		if err != nil {
			helper.WarnErr("Error storing analytics (AnalyticTypeRevives)", err)
		}
	}
}

func ActRetryFree(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	recv := helper.GetGameRequest()
	var request requests.ActRetryFreeRequest
	err = json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	responseStatus := status.OK
	if request.FreeReviveType == 1 {
		responseStatus = status.AlreadyUsedDailyRevive
	} else {
		// TODO: only allow a certain amount of ad revives
	}
	baseInfo := helper.BaseInfo(emess.OK, responseStatus)
	response := responses.NewBaseResponse(baseInfo)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
	if responseStatus == status.OK {
		_, err = analytics.Store(player.ID, factors.AnalyticTypeRevives)
		if err != nil {
			helper.WarnErr("Error storing analytics (AnalyticTypeRevives)", err)
		}
	}
}

func QuickPostGameResults(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	recv := helper.GetGameRequest()
	var request requests.QuickPostGameResultsRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}

	//update energy counter
	for time.Now().UTC().Unix() >= player.PlayerState.EnergyRenewsAt && player.PlayerState.Energy < player.PlayerVarious.EnergyRecoveryMax {
		player.PlayerState.Energy++
		player.PlayerState.EnergyRenewsAt += player.PlayerVarious.EnergyRecoveryTime
	}

	hasSubCharacter := player.PlayerState.SubCharaID != "-1"
	var subC netobj.Character
	mainC, err := player.GetMainChara()
	if err != nil {
		helper.InternalErr("Error getting main character", err)
		return
	}
	playCharacters := []netobj.Character{ // assume only main character active right now
		mainC,
	}
	lvupCharacters := []netobj.Character{
		mainC,
	}
	copy(lvupCharacters, playCharacters)
	if hasSubCharacter {
		subC, err = player.GetSubChara()
		if err != nil {
			helper.InternalErr("Error getting sub character", err)
			return
		}
		playCharacters = []netobj.Character{ // add sub character to playCharacters
			mainC,
			subC,
		}
		lvupCharacters = append(lvupCharacters, subC)
	}
	if request.Closed == 0 { // If the game wasn't exited out of
		player.PlayerState.NumRings += request.Rings
		player.PlayerState.NumRedRings += request.RedRings
		player.PlayerState.NumRouletteTicket += request.RedRings // TODO: This was meant as a debugging function. Perhaps this should be removed at some point?
		player.PlayerState.Animals += request.Animals
		player.OptionUserResult.NumTakeAllRings += request.Rings
		player.OptionUserResult.NumTakeAllRedRings += request.RedRings
		playerTimedHighScore := player.PlayerState.TimedHighScore
		if request.Score > playerTimedHighScore {
			player.PlayerState.TimedHighScore = request.Score
		}
		_, leagueresettime, err := dbaccess.GetStartAndEndTimesForLeague(player.PlayerState.RankingLeague, player.PlayerState.RankingLeagueGroup)
		if err != nil {
			helper.InternalErr("Error getting league reset time (try restarting Outrun or sending the ResetRankingData RPC command)", err)
			return
		}
		if time.Now().UTC().Unix() < leagueresettime {
			playerTimedLeagueHighScore := player.PlayerState.QuickLeagueHighScore
			if request.Score > playerTimedLeagueHighScore {
				player.PlayerState.QuickLeagueHighScore = request.Score
			}
			player.PlayerState.TimedTotalScore += request.Score
			if player.PlayerState.TimedTotalScore > player.PlayerState.TimedHighTotalScore {
				player.PlayerState.TimedHighTotalScore = player.PlayerState.TimedTotalScore
				player.OptionUserResult.QuickTotalSumHighScore = player.PlayerState.TimedHighTotalScore
			}
		}
		helper.DebugOut("request.DailyChallengeValue: %v", request.DailyChallengeValue)
		helper.DebugOut("request.DailyChallengeComplete: %v", request.DailyChallengeComplete)
		if player.PlayerState.DailyChallengeComplete == 0 && request.DailyChallengeComplete == 1 {
			if player.PlayerState.NextNumDailyChallenge <= 0 {
				player.PlayerState.NumDailyChallenge = int64(0)
				player.PlayerState.NextNumDailyChallenge = int64(1)
			}
			player.AddOperatorMessage(
				"A Daily Challenge reward.",
				obj.NewMessageItem(
					consts.DailyMissionRewards[player.PlayerState.NextNumDailyChallenge-1],
					consts.DailyMissionRewardCounts[player.PlayerState.NextNumDailyChallenge-1],
					0,
					0,
				),
				2592000,
			)
			player.PlayerState.NumDailyChallenge = player.PlayerState.NextNumDailyChallenge
		}
		if player.PlayerState.DailyChallengeComplete == 0 {
			player.PlayerState.DailyChallengeComplete = request.DailyChallengeComplete
		}
		player.PlayerState.DailyChallengeValue = request.DailyChallengeValue
		if time.Now().UTC().Unix() >= player.PlayerState.DailyMissionEndTime {
			if player.PlayerState.DailyChallengeComplete == 1 && player.PlayerState.DailyChalSetNum < 10 {
				helper.DebugOut("Advancing to next daily mission...")
				player.PlayerState.DailyChalSetNum++
			} else {
				player.PlayerState.DailyChalCatNum = int64(rand.Intn(5))
				player.PlayerState.DailyChalSetNum = int64(0)
			}
			if player.PlayerState.DailyChallengeComplete == 0 {
				player.PlayerState.NumDailyChallenge = int64(0)
				player.PlayerState.NextNumDailyChallenge = int64(1)
			} else {
				player.PlayerState.NextNumDailyChallenge++
				if int(player.PlayerState.NextNumDailyChallenge) > len(consts.DailyMissionRewards) {
					player.PlayerState.NumDailyChallenge = int64(0)
					player.PlayerState.NextNumDailyChallenge = int64(1) //restart from beginning
					player.PlayerState.DailyChalCatNum = int64(rand.Intn(5))
					player.PlayerState.DailyChalSetNum = int64(0)
				}
			}
			player.PlayerState.DailyChalPosNum = int64(1 + rand.Intn(2))
			player.PlayerState.DailyMissionID = int64((player.PlayerState.DailyChalCatNum * 33) + (player.PlayerState.DailyChalSetNum * 3) + player.PlayerState.DailyChalPosNum)
			player.PlayerState.DailyChallengeValue = int64(0)
			player.PlayerState.DailyChallengeComplete = int64(0)
			player.PlayerState.DailyMissionEndTime = now.EndOfDay().UTC().Unix() + 1
			helper.DebugOut("New daily mission ID: %v", player.PlayerState.DailyMissionID)
		}
		//player.PlayerState.TotalDistance += request.Distance  // We don't do this in timed mode!

		// increase character(s)'s experience
		expIncrease := request.Rings + request.FailureRings // all rings collected
		// check that increases exist
		_, ok := consts.UpgradeIncreases[mainC.ID]
		if !ok {
			helper.InternalErr("Error getting upgrade increase", fmt.Errorf("no key '%s' in consts.UpgradeIncreases", mainC.ID))
			return
		}
		if hasSubCharacter {
			_, ok = consts.UpgradeIncreases[subC.ID]
			if !ok {
				helper.InternalErr("Error getting upgrade increase for sub character", fmt.Errorf("no key '%s' in consts.UpgradeIncreases", subC.ID))
				return
			}
		}

		// Recalculate level
		actualLevel := int64(0)
		index := 0
		for index < len(lvupCharacters[0].AbilityLevel) {
			actualLevel += lvupCharacters[0].AbilityLevel[index]
			index++
		}
		if lvupCharacters[0].Level != actualLevel {
			helper.DebugOut("Main character level has DESYNCED!\nStored: %v\nActual: %v", lvupCharacters[0].Level, actualLevel)
			lvupCharacters[0].Level = actualLevel
		}

		playCharacters[0].AbilityLevelUp = []int64{}
		playCharacters[0].AbilityLevelUpExp = []int64{}
		lvupCharacters[0].AbilityLevelUp = []int64{}
		lvupCharacters[0].AbilityLevelUpExp = []int64{}
		if lvupCharacters[0].Level < 100 {
			helper.DebugOut("Level of main character: %v", lvupCharacters[0].Level)
			abilityIndex := 1
			for abilityIndex == 1 || mainC.AbilityLevel[abilityIndex] >= 10 { // unused ability is at index 1
				abilityIndex = rand.Intn(len(mainC.AbilityLevel))
			}
			playCharacters[0].Exp += expIncrease
			lvupCharacters[0].Exp += expIncrease
			for lvupCharacters[0].Exp >= lvupCharacters[0].Cost && lvupCharacters[0].Status != enums.CharacterStatusMaxLevel {
				// more exp than cost = level up
				// FIXME: The level up logic seems to be not working correctly.
				if lvupCharacters[0].Level < 100 {
					lvupCharacters[0].Level++                                               // increase level
					lvupCharacters[0].AbilityLevel[abilityIndex]++                          // increase ability level
					lvupCharacters[0].Exp -= lvupCharacters[0].Cost                         // remove cost from exp
					lvupCharacters[0].Cost += consts.UpgradeIncreases[lvupCharacters[0].ID] // increase cost
					//playCharacters[0].AbilityLevelUp = append(playCharacters[0].AbilityLevelUp, int64(abilityIndex-1))
					//playCharacters[0].AbilityLevelUpExp = append(playCharacters[0].AbilityLevelUpExp, lvupCharacters[0].Cost)
					lvupCharacters[0].AbilityLevelUp = append(lvupCharacters[0].AbilityLevelUp, int64(abilityIndex-1))
					lvupCharacters[0].AbilityLevelUpExp = append(lvupCharacters[0].AbilityLevelUpExp, lvupCharacters[0].Cost)
				} else {
					lvupCharacters[0].Status = enums.CharacterStatusMaxLevel
					lvupCharacters[0].Exp = 0
				}
				abilityIndex = 1
				for abilityIndex == 1 || lvupCharacters[0].AbilityLevel[abilityIndex] >= 10 { // unused ability is at index 1
					abilityIndex = rand.Intn(len(mainC.AbilityLevel))
				}
			}
		} else {
			lvupCharacters[0].Status = enums.CharacterStatusMaxLevel
		}
		if hasSubCharacter {

			// Recalculate level
			actualLevel := int64(0)
			index := 0
			for index < len(lvupCharacters[1].AbilityLevel) {
				actualLevel += lvupCharacters[1].AbilityLevel[index]
				index++
			}
			if lvupCharacters[1].Level != actualLevel {
				helper.DebugOut("Sub character level has DESYNCED!\nStored: %v\nActual: %v", lvupCharacters[1].Level, actualLevel)
				lvupCharacters[1].Level = actualLevel
			}

			playCharacters[1].AbilityLevelUp = []int64{}
			playCharacters[1].AbilityLevelUpExp = []int64{}
			lvupCharacters[1].AbilityLevelUp = []int64{}
			lvupCharacters[1].AbilityLevelUpExp = []int64{}
			if lvupCharacters[1].Level < 100 {
				abilityIndex := 1
				for abilityIndex == 1 || subC.AbilityLevel[abilityIndex] >= 10 { // unused ability is at index 1
					abilityIndex = rand.Intn(len(subC.AbilityLevel))
				}
				playCharacters[1].Exp += expIncrease
				lvupCharacters[1].Exp += expIncrease
				for lvupCharacters[1].Exp >= lvupCharacters[1].Cost && lvupCharacters[1].Status != enums.CharacterStatusMaxLevel {
					// more exp than cost = level up
					// FIXME: The level up logic seems to be not working correctly. CHECK THE GAME CODE TO SEE IF THIS IS WHAT IT EXPECTS
					if lvupCharacters[1].Level < 100 {
						lvupCharacters[1].Level++                                               // increase level
						lvupCharacters[1].AbilityLevel[abilityIndex]++                          // increase ability level
						lvupCharacters[1].Exp -= lvupCharacters[1].Cost                         // remove cost from exp
						lvupCharacters[1].Cost += consts.UpgradeIncreases[lvupCharacters[1].ID] // increase cost
						//playCharacters[1].AbilityLevelUp = append(playCharacters[1].AbilityLevelUp, int64(abilityIndex-1))
						//playCharacters[1].AbilityLevelUpExp = append(playCharacters[1].AbilityLevelUpExp, lvupCharacters[1].Cost)
						lvupCharacters[1].AbilityLevelUp = append(lvupCharacters[1].AbilityLevelUp, int64(abilityIndex-1))
						lvupCharacters[1].AbilityLevelUpExp = append(lvupCharacters[1].AbilityLevelUpExp, lvupCharacters[1].Cost)
					} else {
						lvupCharacters[1].Status = enums.CharacterStatusMaxLevel
						lvupCharacters[1].Exp = 0
					}
					abilityIndex = 1
					for abilityIndex == 1 || lvupCharacters[1].AbilityLevel[abilityIndex] >= 10 { // unused ability is at index 1
						abilityIndex = rand.Intn(len(subC.AbilityLevel))
					}
				}
			} else {
				lvupCharacters[1].Status = enums.CharacterStatusMaxLevel
			}
		}

		helper.DebugOut("Old mainC Exp: %v / %v", mainC.Exp, mainC.Cost)
		helper.DebugOut("Old mainC Level: %v", mainC.Level)
		if hasSubCharacter {
			helper.DebugOut("Old subC Exp: %v / %v", subC.Exp, subC.Cost)
			helper.DebugOut("Old subC Level: %v", subC.Level)
		}
		helper.DebugOut("New mainC Exp: %v / %v", lvupCharacters[0].Exp, lvupCharacters[0].Cost)
		helper.DebugOut("New mainC Level: %v", lvupCharacters[0].Level)
		if hasSubCharacter {
			helper.DebugOut("New subC Exp: %v / %v", lvupCharacters[1].Exp, lvupCharacters[1].Cost)
			helper.DebugOut("New subC Level: %v", lvupCharacters[1].Level)
		}
		helper.DebugOut("Sent playCharacters[0] Exp: %v / %v", playCharacters[0].Exp, playCharacters[0].Cost)
		helper.DebugOut("Sent playCharacters[0] Level: %v", playCharacters[0].Level)
		helper.DebugOut("Sent playCharacters[0] Ability Levels: %v", playCharacters[0].AbilityLevel)
		helper.DebugOut("Sent playCharacters[0] Ability Level ups: %v", playCharacters[0].AbilityLevelUp)
		helper.DebugOut("Sent playCharacters[0] Ability Level up costs: %v", playCharacters[0].AbilityLevelUpExp)
		if hasSubCharacter {
			helper.DebugOut("Sent playCharacters[1] Exp: %v / %v", playCharacters[1].Exp, playCharacters[1].Cost)
			helper.DebugOut("Sent playCharacters[1] Level: %v", playCharacters[1].Level)
			helper.DebugOut("Sent playCharacters[1] Ability Levels: %v", playCharacters[1].AbilityLevel)
			helper.DebugOut("Sent playCharacters[1] Ability Level ups: %v", playCharacters[1].AbilityLevelUp)
			helper.DebugOut("Sent playCharacters[1] Ability Level up costs: %v", playCharacters[1].AbilityLevelUpExp)
		}
	}

	mainCIndex := player.IndexOfChara(mainC.ID) // TODO: check if -1
	subCIndex := -1
	if hasSubCharacter {
		subCIndex = player.IndexOfChara(subC.ID) // TODO: check if -1
	}

	//helper.DebugOut("CheatResult: %s", request.CheatResult)
	if request.CheatResult[0] != '0' { // unknown trigger
		helper.DebugOut("(CheatResult) flag 1 set!!!")
	}
	if request.CheatResult[1] != '0' { // unknown trigger
		helper.DebugOut("(CheatResult) flag 2 set!!!")
	}
	if request.CheatResult[2] != '0' { // unknown trigger
		helper.DebugOut("(CheatResult) flag 3 set!!!")
	}
	if request.CheatResult[3] != '0' { // unknown trigger
		helper.DebugOut("(CheatResult) flag 4 set!!!")
	}
	if request.CheatResult[4] != '0' { // unknown trigger
		helper.DebugOut("(CheatResult) flag 5 set!!!")
	}
	if request.CheatResult[5] != '0' { // unknown trigger
		helper.DebugOut("(CheatResult) flag 6 set!!!")
	}
	if request.CheatResult[6] != '0' { // unknown trigger
		helper.DebugOut("(CheatResult) flag 7 set!!!")
	}
	if request.CheatResult[7] != '0' { // unknown trigger
		helper.DebugOut("(CheatResult) flag 8 set!!!")
	}

	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultQuickPostGameResults(baseInfo, player, playCharacters)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
	// apply the save after the response so that we don't break the leveling
	lvupCharacters[0].AbilityLevelUp = []int64{}
	lvupCharacters[0].AbilityLevelUpExp = []int64{}
	mainC = lvupCharacters[0]
	if hasSubCharacter {
		lvupCharacters[1].AbilityLevelUp = []int64{}
		lvupCharacters[1].AbilityLevelUpExp = []int64{}
		subC = lvupCharacters[1]
	}
	player.CharacterState[mainCIndex] = mainC
	if hasSubCharacter {
		player.CharacterState[subCIndex] = subC
	}
	err = db.SavePlayer(player)
	if err != nil {
		helper.InternalErr("Error saving player", err)
		return
	}
	helper.DebugOut(fmt.Sprintf("%v", player.PlayerState.Items))

	_, err = analytics.Store(player.ID, factors.AnalyticTypeTimedEnds)
	if err != nil {
		helper.WarnErr("Error storing analytics (AnalyticTypeTimedEnds)", err)
	}
}

func PostGameResults(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	recv := helper.GetGameRequest()
	var request requests.PostGameResultsRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}

	//update energy counter
	for time.Now().UTC().Unix() >= player.PlayerState.EnergyRenewsAt && player.PlayerState.Energy < player.PlayerVarious.EnergyRecoveryMax {
		player.PlayerState.Energy++
		player.PlayerState.EnergyRenewsAt += player.PlayerVarious.EnergyRecoveryTime
	}

	hasSubCharacter := player.PlayerState.SubCharaID != "-1"
	var subC netobj.Character
	mainC, err := player.GetMainChara()
	if err != nil {
		helper.InternalErr("Error getting main character", err)
		return
	}
	playCharacters := []netobj.Character{ // assume only main character active right now
		mainC,
	}
	lvupCharacters := []netobj.Character{
		mainC,
	}
	copy(lvupCharacters, playCharacters)
	if hasSubCharacter {
		subC, err = player.GetSubChara()
		if err != nil {
			helper.InternalErr("Error getting sub character", err)
			return
		}
		playCharacters = []netobj.Character{ // add sub character to playCharacters
			mainC,
			subC,
		}
		lvupCharacters = append(lvupCharacters, subC)
	}
	helper.DebugOut("Pre-function")
	helper.DebugOut("Chapter: %v", player.MileageMapState.Chapter)
	helper.DebugOut("Episode: %v", player.MileageMapState.Episode)
	helper.DebugOut("StageTotalScore: %v", player.MileageMapState.StageTotalScore)
	helper.DebugOut("Point: %v", player.MileageMapState.Point)
	helper.DebugOut("request.Score: %v", request.Score)

	incentives := constobjs.GetMileageIncentives(player.MileageMapState.Episode, player.MileageMapState.Chapter) // Game wants incentives in _current_ episode-chapter
	sentIncentives := []obj.MileageIncentive{}                                                                   // The game appears to behave in unexpected ways if you send the incentives for points which haven't been reached yet.
	var oldRewardEpisode, newRewardEpisode int64
	var oldRewardChapter, newRewardChapter int64
	var oldRewardPoint, newRewardPoint int64

	if request.Closed == 0 { // If the game wasn't exited out of
		oldRewardEpisode = player.MileageMapState.Episode
		oldRewardChapter = player.MileageMapState.Chapter
		oldRewardPoint = player.MileageMapState.Point
		player.PlayerState.NumRings += request.Rings
		player.PlayerState.NumRedRings += request.RedRings
		player.PlayerState.NumRouletteTicket += request.RedRings // TODO: URGENT! Remove as soon as possible!
		player.PlayerState.Animals += request.Animals
		player.OptionUserResult.NumTakeAllRings += request.Rings
		player.OptionUserResult.NumTakeAllRedRings += request.RedRings
		playerHighScore := player.PlayerState.HighScore
		if request.Score > playerHighScore {
			player.PlayerState.HighScore = request.Score
		}
		_, leagueresettime, err := dbaccess.GetStartAndEndTimesForLeague(player.PlayerState.RankingLeague, player.PlayerState.RankingLeagueGroup)
		if err != nil {
			helper.InternalErr("Error getting league reset time (try restarting Outrun or sending the ResetRankingData RPC command)", err)
			return
		}
		if time.Now().UTC().Unix() < leagueresettime {
			playerLeagueHighScore := player.PlayerState.LeagueHighScore
			if request.Score > playerLeagueHighScore {
				player.PlayerState.LeagueHighScore = request.Score
			}
			player.PlayerState.TotalScore += request.Score
			if player.PlayerState.TotalScore > player.PlayerState.HighTotalScore {
				player.PlayerState.HighTotalScore = player.PlayerState.TotalScore
				player.OptionUserResult.TotalSumHighScore = player.PlayerState.HighTotalScore
			}
		}
		helper.DebugOut("request.DailyChallengeValue: %v", request.DailyChallengeValue)
		helper.DebugOut("request.DailyChallengeComplete: %v", request.DailyChallengeComplete)
		if player.PlayerState.DailyChallengeComplete == 0 && request.DailyChallengeComplete == 1 {
			if player.PlayerState.NextNumDailyChallenge <= 0 {
				player.PlayerState.NumDailyChallenge = int64(0)
				player.PlayerState.NextNumDailyChallenge = int64(1)
			}
			player.AddOperatorMessage(
				"A Daily Challenge reward.",
				obj.NewMessageItem(
					consts.DailyMissionRewards[player.PlayerState.NextNumDailyChallenge-1],
					consts.DailyMissionRewardCounts[player.PlayerState.NextNumDailyChallenge-1],
					0,
					0,
				),
				2592000,
			)
			player.PlayerState.NumDailyChallenge = player.PlayerState.NextNumDailyChallenge
		}
		if player.PlayerState.DailyChallengeComplete == 0 {
			player.PlayerState.DailyChallengeComplete = request.DailyChallengeComplete
		}
		player.PlayerState.DailyChallengeValue = request.DailyChallengeValue
		if time.Now().UTC().Unix() >= player.PlayerState.DailyMissionEndTime {
			if player.PlayerState.DailyChallengeComplete == 1 && player.PlayerState.DailyChalSetNum < 10 {
				helper.DebugOut("Advancing to next daily mission...")
				player.PlayerState.DailyChalSetNum++
			} else {
				player.PlayerState.DailyChalCatNum = int64(rand.Intn(5))
				player.PlayerState.DailyChalSetNum = int64(0)
			}
			if player.PlayerState.DailyChallengeComplete == 0 {
				player.PlayerState.NumDailyChallenge = int64(0)
				player.PlayerState.NextNumDailyChallenge = int64(1)
			} else {
				player.PlayerState.NextNumDailyChallenge++
				if int(player.PlayerState.NextNumDailyChallenge) > len(consts.DailyMissionRewards) {
					player.PlayerState.NumDailyChallenge = int64(0)
					player.PlayerState.NextNumDailyChallenge = int64(1) //restart from beginning
					player.PlayerState.DailyChalCatNum = int64(rand.Intn(5))
					player.PlayerState.DailyChalSetNum = int64(0)
				}
			}
			player.PlayerState.DailyChalPosNum = int64(1 + rand.Intn(2))
			player.PlayerState.DailyMissionID = int64((player.PlayerState.DailyChalCatNum * 33) + (player.PlayerState.DailyChalSetNum * 3) + player.PlayerState.DailyChalPosNum)
			player.PlayerState.DailyChallengeValue = int64(0)
			player.PlayerState.DailyChallengeComplete = int64(0)
			player.PlayerState.DailyMissionEndTime = now.EndOfDay().UTC().Unix() + 1
			helper.DebugOut("New daily mission ID: %v", player.PlayerState.DailyMissionID)
		}
		player.PlayerState.TotalDistance += request.Distance
		playerHighDistance := player.PlayerState.HighDistance
		if request.Distance > playerHighDistance {
			player.PlayerState.HighDistance = request.Distance
		}
		// increase character(s)'s experience
		expIncrease := request.Rings + request.FailureRings // all rings collected
		// check that increases exist
		_, ok := consts.UpgradeIncreases[mainC.ID]
		if !ok {
			helper.InternalErr("Error getting upgrade increase for main character", fmt.Errorf("no key '%s' in consts.UpgradeIncreases", mainC.ID))
			return
		}
		if hasSubCharacter {
			_, ok = consts.UpgradeIncreases[subC.ID]
			if !ok {
				helper.InternalErr("Error getting upgrade increase for sub character", fmt.Errorf("no key '%s' in consts.UpgradeIncreases", subC.ID))
				return
			}
		}
		playCharacters[0].AbilityLevelUp = []int64{}
		playCharacters[0].AbilityLevelUpExp = []int64{}
		lvupCharacters[0].AbilityLevelUp = []int64{}
		lvupCharacters[0].AbilityLevelUpExp = []int64{}

		// Recalculate level
		actualLevel := int64(0)
		index := 0
		for index < len(lvupCharacters[0].AbilityLevel) {
			actualLevel += lvupCharacters[0].AbilityLevel[index]
			index++
		}
		if lvupCharacters[0].Level != actualLevel {
			helper.DebugOut("Main character level has DESYNCED!\nStored: %v\nActual: %v", lvupCharacters[0].Level, actualLevel)
			lvupCharacters[0].Level = actualLevel
		}

		if lvupCharacters[0].Level < 100 {
			helper.DebugOut("Main character level is %v", lvupCharacters[0].Level)
			abilityIndex := 1
			for abilityIndex == 1 || lvupCharacters[0].AbilityLevel[abilityIndex] >= 10 { // unused ability is at index 1
				abilityIndex = rand.Intn(len(lvupCharacters[0].AbilityLevel))
			}
			playCharacters[0].Exp += expIncrease
			lvupCharacters[0].Exp += expIncrease
			for lvupCharacters[0].Exp >= lvupCharacters[0].Cost && lvupCharacters[0].Status != enums.CharacterStatusMaxLevel {
				// more exp than cost = level up
				if lvupCharacters[0].Level < 100 {
					lvupCharacters[0].Level++                                               // increase level
					lvupCharacters[0].AbilityLevel[abilityIndex]++                          // increase ability level
					playCharacters[0].AbilityLevel[abilityIndex]++                          // increase ability level
					lvupCharacters[0].Exp -= lvupCharacters[0].Cost                         // remove cost from exp
					lvupCharacters[0].Cost += consts.UpgradeIncreases[lvupCharacters[0].ID] // increase cost
					playCharacters[0].AbilityLevelUp = append(playCharacters[0].AbilityLevelUp, int64(abilityIndex-1))
					playCharacters[0].AbilityLevelUpExp = append(playCharacters[0].AbilityLevelUpExp, lvupCharacters[0].Cost)
					lvupCharacters[0].AbilityLevelUp = append(playCharacters[0].AbilityLevelUp, int64(abilityIndex-1))
					lvupCharacters[0].AbilityLevelUpExp = append(playCharacters[0].AbilityLevelUpExp, lvupCharacters[0].Cost)
				} else {
					lvupCharacters[0].Status = enums.CharacterStatusMaxLevel
					lvupCharacters[0].Exp = 0
				}
				abilityIndex = 1
				for abilityIndex == 1 || lvupCharacters[0].AbilityLevel[abilityIndex] >= 10 { // unused ability is at index 1
					abilityIndex = rand.Intn(len(mainC.AbilityLevel))
				}
			}
		} else {
			lvupCharacters[0].Status = enums.CharacterStatusMaxLevel
			helper.DebugOut("Main character level is maxed out!")
		}
		if hasSubCharacter {
			// Recalculate level
			actualLevel := int64(0)
			index := 0
			for index < len(lvupCharacters[1].AbilityLevel) {
				actualLevel += lvupCharacters[1].AbilityLevel[index]
				index++
			}
			if lvupCharacters[1].Level != actualLevel {
				helper.DebugOut("Sub character level has DESYNCED!\nStored: %v\nActual: %v", lvupCharacters[1].Level, actualLevel)
				lvupCharacters[1].Level = actualLevel
			}

			playCharacters[1].AbilityLevelUp = []int64{}
			playCharacters[1].AbilityLevelUpExp = []int64{}
			if lvupCharacters[1].Level < 100 {
				helper.DebugOut("Sub character level is %v", lvupCharacters[1].Level)
				abilityIndex := 1
				for abilityIndex == 1 || lvupCharacters[1].AbilityLevel[abilityIndex] >= 10 { // unused ability is at index 1
					abilityIndex = rand.Intn(len(lvupCharacters[1].AbilityLevel))
				}
				playCharacters[1].Exp += expIncrease
				lvupCharacters[1].Exp += expIncrease
				for lvupCharacters[1].Exp >= lvupCharacters[1].Cost && lvupCharacters[1].Status != enums.CharacterStatusMaxLevel {
					// more exp than cost = level up
					if lvupCharacters[1].Level < 100 {
						lvupCharacters[1].Level++                                               // increase level
						lvupCharacters[1].AbilityLevel[abilityIndex]++                          // increase ability level
						playCharacters[1].AbilityLevel[abilityIndex]++                          // increase ability level
						lvupCharacters[1].Exp -= lvupCharacters[1].Cost                         // remove cost from exp
						lvupCharacters[1].Cost += consts.UpgradeIncreases[lvupCharacters[1].ID] // increase cost
						playCharacters[1].AbilityLevelUp = append(playCharacters[1].AbilityLevelUp, int64(abilityIndex-1))
						playCharacters[1].AbilityLevelUpExp = append(playCharacters[1].AbilityLevelUpExp, lvupCharacters[1].Cost)
						lvupCharacters[1].AbilityLevelUp = append(lvupCharacters[1].AbilityLevelUp, int64(abilityIndex-1))
						lvupCharacters[1].AbilityLevelUpExp = append(lvupCharacters[1].AbilityLevelUpExp, lvupCharacters[1].Cost)
					} else {
						lvupCharacters[1].Status = enums.CharacterStatusMaxLevel
						lvupCharacters[1].Exp = 0
					}
					abilityIndex = 1
					for abilityIndex == 1 || lvupCharacters[1].AbilityLevel[abilityIndex] >= 10 { // unused ability is at index 1
						abilityIndex = rand.Intn(len(lvupCharacters[1].AbilityLevel))
					}
				}
			} else {
				lvupCharacters[1].Status = enums.CharacterStatusMaxLevel
				helper.DebugOut("Sub character level is maxed out!")
			}
		}

		// TODO: Remove these EXP debug output statements below once the level up animation works properly again
		helper.DebugOut("Old mainC Exp: %v / %v", mainC.Exp, mainC.Cost)
		helper.DebugOut("Old mainC Level: %v", mainC.Level)
		if hasSubCharacter {
			helper.DebugOut("Old subC Exp: %v / %v", subC.Exp, subC.Cost)
			helper.DebugOut("Old subC Level: %v", subC.Level)
		}
		helper.DebugOut("New mainC Exp: %v / %v", lvupCharacters[0].Exp, lvupCharacters[0].Cost)
		helper.DebugOut("New mainC Level: %v", lvupCharacters[0].Level)
		if hasSubCharacter {
			helper.DebugOut("New subC Exp: %v / %v", lvupCharacters[1].Exp, lvupCharacters[1].Cost)
			helper.DebugOut("New subC Level: %v", lvupCharacters[1].Level)
		}

		doStoryProgression := true
		helper.DebugOut("Event ID: %v", request.EventID)
		if request.EventID > 0 { // Is this an event stage?
			if strconv.Itoa(int(request.EventID))[:1] == "1" || strconv.Itoa(int(request.EventID))[:1] == "2" {
				// This is a special stage; don't do story progression since it'll screw with the current point.
				helper.DebugOut("Special Stage; disabling story progression!")
				doStoryProgression = false
			} else {
				helper.DebugOut("Event Type %s (story progression enabled)", strconv.Itoa(int(request.EventID))[1:])
			}
			helper.DebugOut("Player got %v event object(s)", request.EventValue)
			player.EventState.Param += request.EventValue
			//TODO: Send rewards to gift box
		}
		if doStoryProgression {
			player.MileageMapState.StageTotalScore += request.Score
			player.MileageMapState.StageDistance += request.Distance
			if request.Score > player.MileageMapState.StageMaxScore {
				player.MileageMapState.StageMaxScore = request.Score
			}

			goToNextChapter := request.ChapterClear == 1

			player.MileageMapState.NumBossAttack = request.NumBossAttack // TODO: This is guesswork! See if this is correct behavior!

			chaoEggs := request.GetChaoEgg
			player.PlayerState.ChaoEggs += chaoEggs
			if chaoEggs > 0 || player.PlayerState.ChaoEggs >= 10 {
				player.ChaoRouletteGroup.ChaoWheelOptions = netobj.DefaultChaoWheelOptions(player.PlayerState) // create a new wheel
				if player.PlayerState.ChaoEggs >= 10 {
					player.PlayerState.ChaoEggs = 10
				}
			}

			// TODO: Add chao eggs to player
			newPoint := request.ReachPoint

			sentIncentives = incentives

			goToNextEpisode := true
			if goToNextChapter {
				helper.DebugOut("Chapter has been cleared")
				maxChapters, episodeHasMultipleChapters := consts.EpisodeWithChapters[player.MileageMapState.Episode]
				if episodeHasMultipleChapters {
					goToNextEpisode = false
					player.MileageMapState.Chapter++
					if player.MileageMapState.Chapter > maxChapters {
						// there's no more chapters for this episode!
						goToNextEpisode = true
					}
				}
				if goToNextEpisode {
					player.MileageMapState.Episode++
					player.MileageMapState.Chapter = 1
					helper.DebugOut("goToNextEpisode -> Episode: %v", player.MileageMapState.Episode)
					if config.CFile.Debug {
						//player.MileageMapState.Episode = 11
					}
				}
				sentIncentives = incentives
				player.MileageMapState.Point = 0
				player.MileageMapState.StageMaxScore = 0
				player.MileageMapState.StageTotalScore = 0
				player.MileageMapState.StageDistance = 0
				player.MileageMapState.ChapterStartTime = time.Now().Unix()
				player.PlayerState.Rank++
				if player.PlayerState.Rank > 998 { // don't let rank go up past 999
					player.PlayerState.Rank = 998
				}
				if player.PlayerState.Energy < player.PlayerVarious.EnergyRecoveryMax {
					player.PlayerState.Energy = player.PlayerVarious.EnergyRecoveryMax //restore energy
				}
			} else {
				helper.DebugOut("Chapter has NOT been cleared")
				if newPoint >= player.MileageMapState.Point {
					player.MileageMapState.Point = newPoint
				} else {
					helper.Warn("Old point is greater than new point! Point save skipped (%v > %v)", player.MileageMapState.Point, newPoint)
				}
			}
			newRewardEpisode = player.MileageMapState.Episode
			newRewardChapter = player.MileageMapState.Chapter
			newRewardPoint = player.MileageMapState.Point
			// add rewards to PlayerState
			wonRewards := campaign.GetWonRewards(oldRewardEpisode, oldRewardChapter, oldRewardPoint, newRewardEpisode, newRewardChapter, newRewardPoint)
			helper.DebugOut("wonRewards length: %v", wonRewards)
			helper.DebugOut("Previous red rings: %v", player.PlayerState.NumRings)
			helper.DebugOut("Previous rings: %v", player.PlayerState.NumRings)
			newItems := player.PlayerState.Items
			for _, reward := range wonRewards { // TODO: This is O(n^2). Maybe alleviate this?
				helper.DebugOut("Reward: %s", reward.ItemID)
				helper.DebugOut("Reward amount: %v", reward.NumItem)
				if reward.ItemID[:2] == "12" { // ID is an item
					// check if the item is already in the player's inventory
					itemIndex := player.IndexOfItem(reward.ItemID)
					if itemIndex != -1 {
						player.PlayerState.Items[itemIndex].Amount += reward.NumItem
						if player.PlayerState.Items[itemIndex].Amount > 99 {
							helper.DebugOut("item amount is maxed out!")
							player.PlayerState.Items[itemIndex].Amount = 99
						}
					} else {
						helper.Warn("Unknown item reward '%s'", reward.ItemID)
					}
				} else if reward.ItemID == strconv.Itoa(enums.ItemIDRing) { // Rings
					player.PlayerState.NumRings += reward.NumItem
				} else if reward.ItemID == strconv.Itoa(enums.ItemIDRedRing) { // Red rings
					player.PlayerState.NumRedRings += reward.NumItem
				} else if reward.ItemID == enums.CTStrTails { // Tails node
					tailsIndex := player.IndexOfChara(enums.CTStrTails)
					player.CharacterState[tailsIndex].Status = enums.CharacterStatusUnlocked
				} else if reward.ItemID == enums.CTStrKnuckles { // Knuckles node
					knucklesIndex := player.IndexOfChara(enums.CTStrKnuckles)
					player.CharacterState[knucklesIndex].Status = enums.CharacterStatusUnlocked
				} else {
					helper.Warn("Unknown reward '%s', ignoring", reward.ItemID)
				}
				// TODO: allow for any character joining the cast, as well as possibly chao? (no story mode episodes currently do this, but maybe for new episodes in future updates?)
			}
			helper.DebugOut("Current rings: %v", player.PlayerState.NumRings)
			helper.DebugOut("Current red rings: %v", player.PlayerState.NumRedRings)
			player.PlayerState.Items = newItems

			if player.MileageMapState.Episode > 50 { // if beat game, reset to 50-1
				player.MileageMapState.Episode = 50
				player.MileageMapState.Chapter = 1
				helper.DebugOut("goToNextEpisode: Player (%s) beat the game!", player.ID)
			}
		}
	}
	helper.DebugOut("Chapter: %v", player.MileageMapState.Chapter)
	helper.DebugOut("Episode: %v", player.MileageMapState.Episode)
	helper.DebugOut("StageTotalScore: %v", player.MileageMapState.StageTotalScore)
	helper.DebugOut("Point: %v", player.MileageMapState.Point)
	helper.DebugOut("request.Score: %v", request.Score)

	mainCIndex := player.IndexOfChara(mainC.ID) // TODO: check if -1
	subCIndex := -1
	if hasSubCharacter {
		subCIndex = player.IndexOfChara(subC.ID) // TODO: check if -1
	}

	helper.DebugOut("CheatResult: %s", request.CheatResult)
	if request.CheatResult[0] != '0' {
		helper.DebugOut("(CheatResult) flag 1 set!!!")
	}
	if request.CheatResult[1] != '0' {
		helper.DebugOut("(CheatResult) flag 2 set!!!")
	}
	if request.CheatResult[2] != '0' {
		helper.DebugOut("(CheatResult) flag 3 set!!!")
	}
	if request.CheatResult[3] != '0' {
		helper.DebugOut("(CheatResult) flag 4 set!!!")
	}
	if request.CheatResult[4] != '0' {
		helper.DebugOut("(CheatResult) flag 5 set!!!")
	}
	if request.CheatResult[5] != '0' {
		helper.DebugOut("(CheatResult) flag 6 set!!!")
	}
	if request.CheatResult[6] != '0' {
		helper.DebugOut("(CheatResult) flag 7 set!!!")
	}
	if request.CheatResult[7] != '0' {
		helper.DebugOut("(CheatResult) flag 8 set!!!")
	}

	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultPostGameResults(baseInfo, player, playCharacters, sentIncentives)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	eresponse := responses.DefaultPostGameResultsEvent(baseInfo, player, playCharacters, sentIncentives, []obj.Item{}, player.EventState)
	eresponse.Seq, _ = db.BoltGetSessionIDSeq(sid)
	helper.DebugOut("Sent playCharacters[0] Exp: %v / %v", playCharacters[0].Exp, playCharacters[0].Cost)
	helper.DebugOut("Sent playCharacters[0] Level: %v", playCharacters[0].Level)
	helper.DebugOut("Sent playCharacters[0] Ability Levels: %v", playCharacters[0].AbilityLevel)
	helper.DebugOut("Sent playCharacters[0] Ability Level ups: %v", playCharacters[0].AbilityLevelUp)
	helper.DebugOut("Sent playCharacters[0] Ability Level up costs: %v", playCharacters[0].AbilityLevelUpExp)
	if hasSubCharacter {
		helper.DebugOut("Sent playCharacters[1] Exp: %v / %v", playCharacters[1].Exp, playCharacters[1].Cost)
		helper.DebugOut("Sent playCharacters[1] Level: %v", playCharacters[1].Level)
		helper.DebugOut("Sent playCharacters[1] Ability Levels: %v", playCharacters[1].AbilityLevel)
		helper.DebugOut("Sent playCharacters[1] Ability Level ups: %v", playCharacters[1].AbilityLevelUp)
		helper.DebugOut("Sent playCharacters[1] Ability Level up costs: %v", playCharacters[1].AbilityLevelUpExp)
	}
	if request.EventID > 0 {
		helper.DebugOut("Sent event response")
		err = helper.SendResponse(eresponse)
	} else {
		helper.DebugOut("Sent non-event response")
		err = helper.SendResponse(response)
	}
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
	// apply the save after the response so that we don't break the leveling
	lvupCharacters[0].AbilityLevelUp = []int64{}
	lvupCharacters[0].AbilityLevelUpExp = []int64{}
	mainC = lvupCharacters[0]
	if hasSubCharacter {
		lvupCharacters[1].AbilityLevelUp = []int64{}
		lvupCharacters[1].AbilityLevelUpExp = []int64{}
		subC = lvupCharacters[1]
	}
	player.CharacterState[mainCIndex] = mainC
	if hasSubCharacter {
		player.CharacterState[subCIndex] = subC
	}
	err = db.SavePlayer(player)
	if err != nil {
		helper.InternalErr("Error saving player", err)
		return
	}
	helper.DebugOut(fmt.Sprintf("%v", player.PlayerState.Items))

	_, err = analytics.Store(player.ID, factors.AnalyticTypeStoryEnds)
	if err != nil {
		helper.WarnErr("Error storing analytics (AnalyticTypeStoryEnds)", err)
	}
}

func GetFreeItemList(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	var response responses.FreeItemListResponse
	// TODO: maybe give the player up to 3 free uses for each item each day?
	if gameconf.CFile.AllItemsFree {
		response = responses.DefaultFreeItemList(baseInfo)
	} else {
		response = responses.FreeItemList(baseInfo, []obj.Item{}) // No free items
	}
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err := helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetMileageReward(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	recv := helper.GetGameRequest()
	var request requests.MileageRewardRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	/*
		player, err := helper.GetCallingPlayer()
		if err != nil {
			helper.InternalErr("Error getting calling player", err)
			return
		}
	*/
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultMileageReward(baseInfo, request.Chapter, request.Episode)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

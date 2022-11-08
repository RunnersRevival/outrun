package muxhandlers

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"github.com/RunnersRevival/outrun/consts"
	"github.com/RunnersRevival/outrun/db/dbaccess"
	"github.com/RunnersRevival/outrun/enums"

	"github.com/RunnersRevival/outrun/analytics"
	"github.com/RunnersRevival/outrun/analytics/factors"
	"github.com/RunnersRevival/outrun/config"
	"github.com/RunnersRevival/outrun/config/gameconf"
	"github.com/RunnersRevival/outrun/config/infoconf"
	"github.com/RunnersRevival/outrun/db"
	"github.com/RunnersRevival/outrun/emess"
	"github.com/RunnersRevival/outrun/helper"
	"github.com/RunnersRevival/outrun/localizations"
	"github.com/RunnersRevival/outrun/logic"
	"github.com/RunnersRevival/outrun/logic/conversion"
	"github.com/RunnersRevival/outrun/netobj"
	"github.com/RunnersRevival/outrun/obj"
	"github.com/RunnersRevival/outrun/obj/constobjs"
	"github.com/RunnersRevival/outrun/requests"
	"github.com/RunnersRevival/outrun/responses"
	"github.com/RunnersRevival/outrun/status"
	"github.com/jinzhu/now"
)

var (
	loginBonusDebugEnabled = false
)

func Login(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	request := requests.LoginRequest{
		RevivalVerID: 0,
	}
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}

	uid := request.LineAuth.UserID
	password := request.LineAuth.Password

	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	if uid == "0" {
		helper.Out("Entering Registration (Alpha) phase")
		// game wants to get a brand-new ID

		if config.CFile.AllowRegistrations {
			newPlayer, err := db.NewAccount()
			if err != nil {
				helper.InternalErr("Error creating account", err)
				return
			}
			err = db.SavePlayer(newPlayer)
			if err != nil {
				helper.InternalErr("Error saving player", err)
				return
			}
			baseInfo.StatusCode = status.InvalidPassword
			baseInfo.SetErrorMessage(emess.BadPassword)
			response := responses.LoginRegister(
				baseInfo,
				newPlayer.ID,
				newPlayer.Password,
				newPlayer.Key,
			)
			err = helper.SendResponse(response)
			if err != nil {
				helper.InternalErr("Error responding", err)
			}
		} else {
			// ...but registrations aren't enabled!
			helper.Out("Blocked registration attempt")

			localeStringName := "NewAccountsDisabledNotice"
			if config.CFile.IsBetaServer {
				localeStringName = "NewAccountsDisabledBetaNotice"
			}
			err = helper.SendResponse(responses.NewNextVersionResponse(baseInfo,
				0,
				0,
				"",
				localizations.GetStringByLanguage(enums.LangJapanese, localeStringName, true),
				localizations.GetStringByLanguage(enums.LangEnglish, localeStringName, true),
				"https://sonicrunners.com/",
			))
			if err != nil {
				helper.InternalErr("Error sending response", err)
				return
			}
		}

		return
	} else if uid != "0" && password == "" {
		helper.Out("Entering Pre-Login (Bravo) phase")
		// game wants to log in
		baseInfo.StatusCode = status.InvalidPassword
		baseInfo.SetErrorMessage(emess.BadPassword)
		player, err := db.GetPlayer(uid)
		if err != nil {
			if db.DoesPlayerExistInDatabase(uid) {
				helper.InternalErr("Error getting player", err)
				return
			} else {
				// likely account that wasn't found, so let's tell them that:
				response := responses.NewBaseResponse(baseInfo)
				baseInfo.StatusCode = status.MissingPlayer
				err = helper.SendResponse(response)
				if err != nil {
					helper.InternalErr("Error sending response", err)
				}
			}

			return
		}
		response := responses.LoginCheckKey(baseInfo, player.Key)
		err = helper.SendResponse(response)
		if err != nil {
			helper.InternalErr("Error sending response", err)
			return
		}
		return
	} else if uid != "0" && password != "" {
		helper.Out("Entering Login (Charlie) phase")
		// game is attempting to log in using key
		player, err := db.GetPlayer(uid)
		if err != nil {
			helper.InternalErr("Error getting player", err)
			return
		}
		if request.Password == logic.GenerateLoginPasskey(player) {
			baseInfo.StatusCode = status.OK
			baseInfo.SetErrorMessage(emess.OK)
			if time.Now().UTC().Unix() < player.SuspendedUntil {
				// player is suspended! return a NextVersion response (bring up notifications menu with specified text)
				helper.Out("Player is suspended!")
				baseInfo.StatusCode = status.ServerNextVersion
				err = helper.SendResponse(responses.NewNextVersionResponse(baseInfo,
					player.PlayerState.NumRedRings,
					player.PlayerState.NumBuyRedRings,
					player.Username,
					localizations.GetStringByLanguage(enums.LangJapanese, "SuspensionNotice_Temporary", true),
					localizations.GetStringByLanguage(player.Language, "SuspensionNotice_Temporary", true),
					"https://sonicrunners.com/",
				))
				if err != nil {
					helper.InternalErr("Error sending response", err)
					return
				}
			} else {

				if request.RevivalVerID < player.LastLoginVersionId {
					// player is suspended! return a NextVersion response (bring up notifications menu with specified text)
					helper.Out("Player is using outdated version!")
					baseInfo.StatusCode = status.ServerNextVersion
					err = helper.SendResponse(responses.NewNextVersionResponse(baseInfo,
						player.PlayerState.NumRedRings,
						player.PlayerState.NumBuyRedRings,
						player.Username,
						localizations.GetStringByLanguage(enums.LangJapanese, "VersionTooOldForAccountNotice", true),
						localizations.GetStringByLanguage(player.Language, "VersionTooOldForAccountNotice", true),
						"https://sonicrunners.com/",
					))
					if err != nil {
						helper.InternalErr("Error sending response", err)
						return
					}
				} else {
					sid, err := db.BoltAssignSessionID(uid, strconv.Itoa(int(request.Seq)))
					if err != nil {
						helper.InternalErr("Error assigning session ID", err)
						return
					}
					player.Language = request.Language
					player.LastLogin = time.Now().UTC().Unix()
					player.PlayerVarious.EnergyRecoveryMax = gameconf.CFile.EnergyRecoveryMax
					player.PlayerVarious.EnergyRecoveryTime = gameconf.CFile.EnergyRecoveryTime
					player.LastLoginDevice = request.Device
					player.LastLoginPlatform = request.Platform
					player.LastLoginVersionId = request.RevivalVerID
					helper.DebugOut("Device: %s", request.Device)
					helper.DebugOut("Platform: %v", request.Platform)
					helper.DebugOut("Store ID: %v", request.StoreID)
					err = db.SavePlayer(player)
					if err != nil {
						helper.InternalErr("Error saving player", err)
						return
					}
					response := responses.LoginSuccess(baseInfo, sid, player.Username, player.PlayerVarious.EnergyRecoveryTime, player.PlayerVarious.EnergyRecoveryMax)
					helper.DebugOut("seq = %v", request.Seq)
					response.Seq = request.Seq
					err = helper.SendResponse(response)
					if err != nil {
						helper.InternalErr("Error sending response", err)
						return
					}
					analytics.Store(player.ID, factors.AnalyticTypeLogins)
				}
			}
		} else {
			// Looks like the credentials don't match what's in the database!
			baseInfo.StatusCode = status.InvalidPassword
			baseInfo.SetErrorMessage(emess.BadPassword)
			helper.DebugOut("Incorrect passkey sent: \"%s\"", request.Password)
			err = helper.SendResponse(responses.NewBaseResponse(baseInfo))
			if err != nil {
				helper.InternalErr("Error sending response", err)
				return
			}
		}
	}
}

func LoginMaintenance(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.LoginRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}

	uid := request.LineAuth.UserID
	password := request.LineAuth.Password

	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	if uid == "0" {
		helper.Out("Entering Registration (Alpha) phase")
		baseInfo.StatusCode = status.ServerNextVersion
		baseInfo.SetErrorMessage(emess.BadPassword)
		response := responses.NewNextVersionResponse(baseInfo,
			int64(0),
			int64(0),
			"",
			localizations.GetStringByLanguage(enums.LangJapanese, "DefaultMaintenanceMessage", true),
			localizations.GetStringByLanguage(enums.LangEnglish, "DefaultMaintenanceMessage", true),
			"https://www.sonicrunners.com/",
		)
		err = helper.SendResponse(response)
		if err != nil {
			helper.InternalErr("Error responding", err)
		}
		return
	} else if uid != "0" && password == "" {
		helper.Out("Entering Pre-Login (Bravo) phase")
		// game wants to log in
		baseInfo.StatusCode = status.InvalidPassword
		baseInfo.SetErrorMessage(emess.BadPassword)
		player, err := db.GetPlayer(uid)
		if err != nil {
			if db.DoesPlayerExistInDatabase(uid) {
				helper.InternalErr("Error getting player", err)
				return
			} else {
				// likely account that wasn't found, so let's tell them that:
				response := responses.NewBaseResponse(baseInfo)
				baseInfo.StatusCode = status.MissingPlayer
				err = helper.SendResponse(response)
				if err != nil {
					helper.InternalErr("Error sending response", err)
				}
			}

			return
		}
		response := responses.LoginCheckKey(baseInfo, player.Key)
		err = helper.SendResponse(response)
		if err != nil {
			helper.InternalErr("Error sending response", err)
			return
		}
		return
	} else if uid != "0" && password != "" {
		helper.Out("Entering Login (Charlie) phase")
		// game is attempting to log in using key
		player, err := db.GetPlayer(uid)
		if err != nil {
			helper.InternalErr("Error getting player", err)
			return
		}
		if request.Password == logic.GenerateLoginPasskey(player) {
			baseInfo.StatusCode = status.ServerNextVersion
			baseInfo.SetErrorMessage(emess.OK)
			response := responses.NewNextVersionResponse(baseInfo,
				player.PlayerState.NumRedRings,
				player.PlayerState.NumBuyRedRings,
				player.Username,
				localizations.GetStringByLanguage(enums.LangJapanese, "DefaultMaintenanceMessage", true),
				localizations.GetStringByLanguage(enums.LangEnglish, "DefaultMaintenanceMessage", true),
				"https://sonicrunners.com/",
			)
			err = helper.SendResponse(response)
			if err != nil {
				helper.InternalErr("Error sending response", err)
				return
			}
		} else {
			// Looks like the credentials don't match what's in the database!
			baseInfo.StatusCode = status.InvalidPassword
			baseInfo.SetErrorMessage(emess.BadPassword)
			helper.DebugOut("Incorrect passkey sent: \"%s\"", request.Password)
			err = helper.SendResponse(responses.NewBaseResponse(baseInfo))
			if err != nil {
				helper.InternalErr("Error sending response", err)
				return
			}
		}
	}
}

func GetVariousParameter(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.VariousParameter(baseInfo, player)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
}

func GetInformation(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	uid, err := helper.GetCallingPlayerID()
	if err != nil {
		helper.InternalErr("Error getting player ID", err)
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	infos := []obj.Information{}
	helper.DebugOut("%v", infoconf.CFile.EnableInfos)
	if infoconf.CFile.EnableInfos {
		for _, ci := range infoconf.CFile.Infos {
			newInfo := conversion.ConfiguredInfoToInformation(ci)
			infos = append(infos, newInfo)
			helper.DebugOut(newInfo.Param)
		}
	}
	/*operatorInfos := []obj.OperatorInformation{
		obj.LeagueOperatorInformation(0, 123, 456, 0, now.BeginningOfWeek().UTC().Unix(), enums.RankingLeagueF_M, enums.RankingLeagueF_M, 50, 0),
		obj.EventOperatorInformation("200010000", 7, 17171, 0),
		obj.LeagueOperatorInformation(2, 123, 456, 0, now.BeginningOfWeek().UTC().Unix(), enums.RankingLeagueF_M, enums.RankingLeagueF_M, 50, 0),
	}*/
	operatorInfos, err := dbaccess.GetOperatorInfos(uid)
	if err != nil {
		helper.WarnErr("Couldn't load operator infos.", err)
	}
	numOpUnread := int64(len(operatorInfos))
	response := responses.Information(baseInfo, infos, operatorInfos, numOpUnread)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetTicker(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultTicker(baseInfo, player)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func LoginBonus(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}

	if time.Now().UTC().Unix() > player.LoginBonusState.LoginBonusEndTime {
		player.LoginBonusState = netobj.DefaultLoginBonusState(player.LoginBonusState.CurrentFirstLoginBonusDay)
	}
	doLoginBonus := false
	if time.Now().UTC().Unix() > player.LoginBonusState.NextLoginBonusTime {
		doLoginBonus = true
		player.LoginBonusState.LastLoginBonusTime = time.Now().UTC().Unix()
		player.LoginBonusState.NextLoginBonusTime = now.EndOfDay().UTC().Unix()
		player.LoginBonusState.CurrentLoginBonusDay++
		if gameconf.CFile.EnableStartDashLoginBonus {
			player.LoginBonusState.CurrentFirstLoginBonusDay++
			if player.LoginBonusState.CurrentFirstLoginBonusDay > 7 {
				player.LoginBonusState.CurrentFirstLoginBonusDay = 7
			}
		}
	}
	if !loginBonusDebugEnabled {
		err = db.SavePlayer(player)
		if err != nil {
			helper.InternalErr("Error saving player", err)
			return
		}
	} else {
		helper.DebugOut("Login bonus in debug mode; player data NOT updated!")
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultLoginBonus(baseInfo, player, doLoginBonus)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func LoginBonusSelect(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	recv := helper.GetGameRequest()
	var request requests.LoginBonusSelectRequest
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
	rewardList := []obj.Item{}
	firstRewardList := []obj.Item{}
	if request.FirstRewardDays != -1 && int(request.FirstRewardDays) < len(constobjs.DefaultFirstLoginBonusRewardList) {
		firstRewardList = constobjs.DefaultFirstLoginBonusRewardList[request.FirstRewardDays].SelectRewardList[request.FirstRewardSelect].ItemList
	}
	if request.RewardDays != -1 && int(request.RewardDays) < len(constobjs.DefaultLoginBonusRewardList) {
		rewardList = constobjs.DefaultLoginBonusRewardList[request.RewardDays].SelectRewardList[request.RewardSelect].ItemList
	}
	if !loginBonusDebugEnabled {
		for _, item := range rewardList {
			itemid, _ := strconv.Atoi(item.ID)
			player.AddOperatorMessage(
				"A Login Bonus.",
				obj.MessageItem{
					int64(itemid),
					item.Amount,
					0,
					0,
				},
				2592000,
			)
			helper.DebugOut("Sent %s x %v to gift box (Login Bonus)", item.ID, item.Amount)
		}
		for _, item := range firstRewardList {
			itemid, _ := strconv.Atoi(item.ID)
			player.AddOperatorMessage(
				"A Debut Dash Login Bonus.",
				obj.MessageItem{
					int64(itemid),
					item.Amount,
					0,
					0,
				},
				2592000,
			)
			helper.DebugOut("Sent %s x %v to gift box (Start Dash Login Bonus)", item.ID, item.Amount)
		}
		err = db.SavePlayer(player)
		if err != nil {
			helper.InternalErr("Error saving player", err)
			return
		}
	} else {
		helper.DebugOut("Login bonus in debug mode; gifts NOT sent!")
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.LoginBonusSelect(baseInfo, rewardList, firstRewardList)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetCountry(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	// TODO: Should get correct country code!
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultGetCountry(baseInfo)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err := helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetMigrationPassword(helper *helper.Helper) {
	randChar := func(charset string, length int64) string {
		runes := []rune(charset)
		final := make([]rune, 12)
		for i := range final {
			final[i] = runes[rand.Intn(len(runes))]
		}
		return string(final)
	}
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	recv := helper.GetGameRequest()
	var request requests.GetMigrationPasswordRequest
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
	if len(player.MigrationPassword) != 12 {
		player.MigrationPassword = randChar("abcdefghijklmnopqrstuvwxyz1234567890", 12)
	}
	player.UserPassword = request.UserPassword
	db.SavePlayer(player)
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.MigrationPassword(baseInfo, player)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func Migration(helper *helper.Helper) {
	randChar := func(charset string, length int64) string {
		runes := []rune(charset)
		final := make([]rune, length)
		for i := range final {
			final[i] = runes[rand.Intn(len(runes))]
		}
		return string(final)
	}
	recv := helper.GetGameRequest()
	var request requests.LoginRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	password := request.LineAuth.MigrationPassword
	migrationUserPassword := request.LineAuth.MigrationUserPassword

	baseInfo := helper.BaseInfo(emess.OK, status.OK)

	helper.DebugOut("Transfer ID: %s", password)
	playerInfo, pid, err := dbaccess.GetPlayerInfoFromMigrationPass(consts.DBMySQLTableCorePlayerInfo, password)
	if err != nil {
		// MOST LIKELY not found in database; return that
		baseInfo.StatusCode = status.InvalidPassword
		response := responses.NewBaseResponse(baseInfo)
		helper.SendResponse(response)
		helper.InternalErr("Error finding players by password", err)
		return
	}
	if migrationUserPassword == playerInfo.UserPassword {
		baseInfo.StatusCode = status.OK
		baseInfo.SetErrorMessage(emess.OK)

		// TODO: Make clearing the migration password and user password a configurable option
		playerInfo.MigrationPassword = randChar("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", 12) //generate a brand new transfer ID
		playerInfo.UserPassword = ""                                                                                  //clear user password

		if !config.CFile.AllowMultiDevice {
			playerInfo.Password = randChar("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", 10) // randomize password, preventing old device from being used
		}

		playerInfo.LastLogin = time.Now().UTC().Unix()
		err = dbaccess.SetPlayerInfo(consts.DBMySQLTableCorePlayerInfo, pid, playerInfo)
		if err != nil {
			helper.InternalErr("Error saving player info", err)
			return
		}
		sid, err := db.BoltAssignSessionID(pid, "0")
		if err != nil {
			helper.InternalErr("Error assigning session ID", err)
			return
		}
		helper.DebugOut("User ID: %s", pid)
		helper.DebugOut("Username: %s", playerInfo.Username)
		helper.DebugOut("New Transfer ID: %s", playerInfo.MigrationPassword)
		if !config.CFile.AllowMultiDevice {
			helper.DebugOut("New Internal Password: %s", playerInfo.Password)
		}
		response := responses.MigrationSuccess(baseInfo, sid, pid, playerInfo.Username, playerInfo.Password, netobj.DefaultPlayerVarious().EnergyRecoveryTime, netobj.DefaultPlayerVarious().EnergyRecoveryTime)
		helper.SendResponse(response)
	} else {
		baseInfo.StatusCode = status.InvalidPassword
		baseInfo.SetErrorMessage(emess.BadPassword)
		helper.DebugOut("Incorrect password for user ID %s", pid)
		response := responses.NewBaseResponse(baseInfo)
		helper.SendResponse(response)
	}
}

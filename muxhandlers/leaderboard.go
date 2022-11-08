package muxhandlers

import (
	"encoding/json"
	"time"

	"github.com/RunnersRevival/outrun/consts"
	"github.com/RunnersRevival/outrun/db"

	"github.com/RunnersRevival/outrun/db/dbaccess"

	"github.com/RunnersRevival/outrun/emess"
	"github.com/RunnersRevival/outrun/helper"
	"github.com/RunnersRevival/outrun/obj"
	"github.com/RunnersRevival/outrun/obj/constobjs"
	"github.com/RunnersRevival/outrun/requests"
	"github.com/RunnersRevival/outrun/responses"
	"github.com/RunnersRevival/outrun/status"
)

// Leaderboards and league endpoints

func GetWeeklyLeaderboardOptions(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	recv := helper.GetGameRequest()
	var request requests.LeaderboardRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	uid, err := helper.GetCallingPlayerID()
	if err != nil {
		helper.SendResponse(responses.NewBaseResponse(helper.BaseInfo(emess.OK, status.ExpiredSession)))
		helper.InternalErr("Error getting calling player ID", err)
		return
	}
	playerState, err := dbaccess.GetPlayerState(consts.DBMySQLTablePlayerStates, uid)
	if err != nil {
		helper.InternalErr("Error getting player state", err)
		return
	}
	mode := request.Mode
	rankingleague := playerState.RankingLeague
	rankingleaguegroup := playerState.RankingLeagueGroup
	if mode == 1 {
		rankingleague = playerState.QuickRankingLeague
		rankingleaguegroup = playerState.QuickRankingLeagueGroup
	}
	leaguestarttime, leagueendtime, err := dbaccess.GetStartAndEndTimesForLeague(rankingleague, rankingleaguegroup)
	if err != nil {
		helper.InternalErr("Error getting league start and end times (try restarting Outrun or sending the ResetRankingData RPC command)", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultWeeklyLeaderboardOptions(baseInfo, mode)
	response.StartTime = leaguestarttime
	response.ResetTime = leagueendtime
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetWeeklyLeaderboardEntries(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	recv := helper.GetGameRequest()
	var request requests.LeaderboardEntriesRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	uid, err := helper.GetCallingPlayerID()
	if err != nil {
		helper.SendResponse(responses.NewBaseResponse(helper.BaseInfo(emess.OK, status.ExpiredSession)))
		helper.InternalErr("Error getting calling player ID", err)
		return
	}
	playerState, err := dbaccess.GetPlayerState(consts.DBMySQLTablePlayerStates, uid)
	if err != nil {
		helper.InternalErr("Error getting player state", err)
		return
	}
	helper.DebugOut("Mode %v, type %v", request.Mode, request.Type)
	helper.DebugOut("Starting from %v; listing up to %v entries", request.First, request.Count)
	mode := request.Mode
	lbtype := request.Type
	rankingleague := playerState.RankingLeague
	rankingleaguegroup := playerState.RankingLeagueGroup
	if mode == 1 {
		rankingleague = playerState.QuickRankingLeague
		rankingleaguegroup = playerState.QuickRankingLeagueGroup
	}
	leaguestarttime, leagueendtime, err := dbaccess.GetStartAndEndTimesForLeague(rankingleague, rankingleaguegroup)
	if err != nil {
		helper.InternalErr("Error getting league start and end times (try restarting Outrun or sending the ResetRankingData RPC command)", err)
		return
	}
	var myEntry interface{}
	entryList := []obj.LeaderboardEntry{}
	entryCount := int64(0)
	if lbtype == 4 || lbtype == 5 {
		league := playerState.RankingLeague
		leagueGroup := playerState.RankingLeagueGroup
		if mode == 1 {
			league = playerState.QuickRankingLeague
			leagueGroup = playerState.QuickRankingLeagueGroup
		}
		entryList, myEntry, err = dbaccess.GetLeagueHighScores(mode, lbtype, league, leagueGroup, request.First-1, request.Count, uid, request.RevivalVerID)
		if err != nil {
			helper.InternalErr("Error getting high score table", err)
			return
		}
		if time.Now().UTC().Unix() < leagueendtime && time.Now().UTC().Unix() >= leaguestarttime {
			entryCount = int64(len(entryList)) // We know for sure that the game will request 71 entries for this tab, so this shouldn't be a problem, since only 50 people should be in a league/group combination at once
		} else {
			entryList = []obj.LeaderboardEntry{}
		}
	} else {
		if lbtype == 6 || lbtype == 7 || (time.Now().UTC().Unix() < leagueendtime && time.Now().UTC().Unix() >= leaguestarttime) {
			entryList, myEntry, err = dbaccess.GetHighScores(mode, lbtype, request.First-1, request.Count, uid, false)
			if err != nil {
				helper.InternalErr("Error getting high score table", err)
				return
			}
			entryCount, err = dbaccess.GetNumOfLeaderboardPlayers(mode, lbtype)
			if err != nil {
				helper.InternalErr("Error getting number of players", err)
				return
			}
		}
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.WeeklyLeaderboardEntries(
		baseInfo,
		myEntry,
		-1,
		leaguestarttime,
		leagueendtime,
		request.First,
		mode,
		entryCount,
		entryList,
	)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetLeagueData(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	recv := helper.GetGameRequest()
	var request requests.LeaderboardRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	uid, err := helper.GetCallingPlayerID()
	if err != nil {
		helper.SendResponse(responses.NewBaseResponse(helper.BaseInfo(emess.OK, status.ExpiredSession)))
		helper.InternalErr("Error getting calling player ID", err)
		return
	}
	playerState, err := dbaccess.GetPlayerState(consts.DBMySQLTablePlayerStates, uid)
	if err != nil {
		helper.InternalErr("Error getting player state", err)
		return
	}
	mode := request.Mode
	var leagueData obj.LeagueData
	if mode == 0 {
		leagueData = constobjs.LeagueDataDefinitions[playerState.RankingLeague]
		leagueData.GroupID = playerState.RankingLeagueGroup
	} else {
		leagueData = constobjs.QuickLeagueDataDefinitions[playerState.QuickRankingLeague]
		leagueData.GroupID = playerState.QuickRankingLeagueGroup
	}
	leagueData.NumGroupMember = 50
	leagueData.NumLeagueMember = 500 // TODO: Add something to dbaccess which can determine these values!
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.LeagueData(baseInfo, leagueData, mode)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetLeagueOperatorData(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	recv := helper.GetGameRequest()
	var request requests.LeaderboardRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultLeagueOperatorData(baseInfo, request.Mode)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

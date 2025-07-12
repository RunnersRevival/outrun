package responses

import (
	//	"strconv"

	//	"github.com/RunnersRevival/outrun/enums"
	"fmt"
	"strconv"

	"github.com/RunnersRevival/outrun/config/gameconf"
	"github.com/RunnersRevival/outrun/db"
	"github.com/RunnersRevival/outrun/logic"
	"github.com/RunnersRevival/outrun/netobj"
	"github.com/RunnersRevival/outrun/obj"
	"github.com/RunnersRevival/outrun/obj/constobjs"
	"github.com/RunnersRevival/outrun/responses/responseobjs"
)

type EventListResponse struct {
	BaseResponse
	EventList []obj.Event `json:"eventList"`
}

func EventList(base responseobjs.BaseInfo, eventList []obj.Event) EventListResponse {
	baseResponse := NewBaseResponse(base)
	out := EventListResponse{
		baseResponse,
		eventList,
	}
	return out
}

func DefaultEventList(base responseobjs.BaseInfo) EventListResponse {
	return EventList(
		base,
		[]obj.Event{
			/*
			   obj.NewEvent(
			       //enums.EventIDSpecialStage+10002, // game subtracts one from number?
			       //enums.EventIDAdvert+50002, // 50002 converts to ui_event_50005_Atlas_en?
			       //enums.EventIDBGM+70002, // 70002 goes to 70007
			       enums.EventIDQuick+60002, // 60002 goes to 60006
			       0,                        // event type
			       now.BeginningOfDay().Unix(),
			       now.EndOfDay().Unix(),
			       now.EndOfDay().Unix(),
			   ),
			*/
		},
	)
}

type EventRewardListResponse struct {
	BaseResponse
	EventRewardList []obj.EventReward `json:"eventRewardList"`
}

func EventRewardList(base responseobjs.BaseInfo, eventRewardList []obj.EventReward) EventRewardListResponse {
	baseResponse := NewBaseResponse(base)
	out := EventRewardListResponse{
		baseResponse,
		eventRewardList,
	}
	return out
}

func DefaultEventRewardList(base responseobjs.BaseInfo) EventRewardListResponse {
	// TODO: Maybe *don't* hardcode this mess? We'd have to recompile Outrun to change rewards.
	return EventRewardList(
		base,
		constobjs.DefaultXmasEventRewardList(),
	)
}

type EventStateResponse struct {
	BaseResponse
	netobj.EventState `json:"eventState"`
}

func EventState(base responseobjs.BaseInfo, eventState netobj.EventState) EventStateResponse {
	baseResponse := NewBaseResponse(base)
	out := EventStateResponse{
		baseResponse,
		eventState,
	}
	return out
}

// below is stuff for raid bosses

type EventUserRaidbossStateResponse struct {
	BaseResponse
	netobj.EventUserRaidbossState `json:"eventUserRaidboss"`
}

func EventUserRaidbossState(base responseobjs.BaseInfo, userRaidbossState netobj.EventUserRaidbossState) EventUserRaidbossStateResponse {
	baseResponse := NewBaseResponse(base)
	out := EventUserRaidbossStateResponse{
		baseResponse,
		userRaidbossState,
	}
	return out
}

type EventUserRaidbossListResponse struct {
	BaseResponse
	netobj.EventUserRaidbossState `json:"eventUserRaidboss"`
	EventRaidbossStates           []netobj.EventRaidbossState `json:"eventUserRaidbossList"`
}

func DefaultEventUserRaidbossList(base responseobjs.BaseInfo, userRaidbossState netobj.EventUserRaidbossState) EventUserRaidbossListResponse {
	// here we'll fetch the data from the DB in regards to the active or recently expired/defeated raid bosses

	baseResponse := NewBaseResponse(base)
	out := EventUserRaidbossListResponse{
		baseResponse,
		userRaidbossState,
		[]netobj.EventRaidbossState{
			netobj.DefaultRaidbossState(),
			netobj.DefaultRaidbossState2(),
			netobj.DefaultRaidbossState3(),
		},
	}
	return out
}

type EventActStartResponse struct {
	ActStartBaseResponse
	netobj.EventUserRaidbossState `json:"eventUserRaidboss"`
}

func EventActStart(base responseobjs.BaseInfo, playerState netobj.PlayerState, campaignList []obj.Campaign, eventUserRaidbossState netobj.EventUserRaidbossState) EventActStartResponse {
	actStartBase := ActStartBase(base, playerState, campaignList)
	return EventActStartResponse{
		actStartBase,
		eventUserRaidbossState,
	}
}

func DefaultEventActStart(base responseobjs.BaseInfo, player netobj.Player) EventActStartResponse {
	campaignList := []obj.Campaign{}
	playerState := player.PlayerState
	eventUserRaidbossState := netobj.ConvertRaidBossPlayerState(player.RaidBossPlayerState)
	return EventActStart(
		base,
		playerState,
		campaignList,
		eventUserRaidbossState,
	)
}

type EventPostGameResultsResponse struct {
	BaseResponse
	netobj.EventUserRaidbossState `json:"eventUserRaidboss"`
}

func EventPostGameResults(base responseobjs.BaseInfo, userRaidbossState netobj.EventUserRaidbossState) EventPostGameResultsResponse {
	baseResponse := NewBaseResponse(base)
	out := EventPostGameResultsResponse{
		baseResponse,
		userRaidbossState,
	}
	return out
}

type EventUpdateGameResultsResponse struct {
	BaseResponse
	PlayerState             netobj.PlayerState    `json:"playerState"`
	ChaoState               []netobj.Chao         `json:"chaoState"`
	DailyChallengeIncentive []obj.Incentive       `json:"dailyChallengeIncentive"` // should be obj.Item, but game doesn't care
	CharacterState          []netobj.Character    `json:"characterState"`
	MessageList             []obj.Message         `json:"messageList"`
	OperatorMessageList     []obj.OperatorMessage `json:"operatorMessageList"`
	TotalMessage            int64                 `json:"totalMessage"`
	TotalOperatorMessage    int64                 `json:"totalOperatorMessage"`
	PlayCharacterState      []netobj.Character    `json:"playCharacterState"`
	EventIncentiveList      []obj.Item            `json:"eventIncentiveList"`
	WheelOptions            netobj.WheelOptions   `json:"wheelOptions"`
	EventState              netobj.EventState     `json:"eventState,omitempty"`
}

func EventUpdateGameResults(base responseobjs.BaseInfo, player netobj.Player, dci []obj.Incentive, ml []obj.Message, oml []obj.OperatorMessage, pcs []netobj.Character, eil []obj.Item, wo netobj.WheelOptions, es netobj.EventState) EventUpdateGameResultsResponse {
	baseResponse := NewBaseResponse(base)
	playerState := player.PlayerState
	chaoState := player.ChaoState
	dailyChallengeIncentive := dci
	characterState := player.CharacterState
	messageList := []obj.Message{}
	operatorMessageList := []obj.OperatorMessage{}
	totalMessage := int64(len(messageList))
	totalOperatorMessage := int64(len(operatorMessageList))
	playCharacterState := pcs
	return EventUpdateGameResultsResponse{
		baseResponse,
		playerState,
		chaoState,
		dailyChallengeIncentive,
		characterState,
		messageList,
		operatorMessageList,
		totalMessage,
		totalOperatorMessage,
		playCharacterState,
		eil,
		wo,
		es,
	}
}

func DefaultEventUpdateGameResults(base responseobjs.BaseInfo, player netobj.Player, pcs []netobj.Character, es netobj.EventState) EventUpdateGameResultsResponse {
	baseResponse := NewBaseResponse(base)
	playerState := player.PlayerState
	chaoState := player.ChaoState
	dailyChallengeIncentive := []obj.Incentive{}
	characterState := player.CharacterState
	messageList := []obj.Message{}
	operatorMessageList := []obj.OperatorMessage{}
	totalMessage := int64(len(messageList))
	totalOperatorMessage := int64(len(operatorMessageList))
	eil := []obj.Item{}
	player.LastWheelOptions = logic.WheelRefreshLogic(player, player.LastWheelOptions)
	wo := player.LastWheelOptions
	return EventUpdateGameResultsResponse{
		baseResponse,
		playerState,
		chaoState,
		dailyChallengeIncentive,
		characterState,
		messageList,
		operatorMessageList,
		totalMessage,
		totalOperatorMessage,
		pcs,
		eil,
		wo,
		es,
	}
}

type EventRaidbossDesiredListResponse struct {
	BaseResponse
	DesiredList []netobj.RaidBossDesiredUser `json:"eventRaidbossDesiredList"`
}

type EventRaidBossUserListResponse struct {
	BaseResponse
	UserList      []netobj.RaidBossUserState `json:"eventRaidbossUserList"`
	RaidbossBonus netobj.RaidBossBonus       `json:"eventRaidbossBonus"`
	RaidbossState netobj.EventRaidbossState  `json:"eventRaidboss"`
}

func GetRaidBossPlayerList(state netobj.RaidBossInternalState) []netobj.RaidBossUserState {
	userStates := []netobj.RaidBossUserState{}
	for _, playerS := range state.PlayerStates {
		wrestleBeatFlg := int64(0)
		if playerS.HadFinalHit {
			wrestleBeatFlg = 1
		}
		userState := netobj.RaidBossUserState{
			WrestleID:         strconv.Itoa(int(playerS.PlayerID)),
			Name:              fmt.Sprintf("Unknown Player %X", playerS.PlayerID*2),
			Grade:             playerS.Grade,
			NumRank:           0,
			LoginTime:         1563177600,
			CharacterID:       gameconf.CFile.DefaultMainCharacter,
			CharacterLevel:    0,
			SubcharacterID:    gameconf.CFile.DefaultSubCharacter,
			SubcharacterLevel: 0,
			MainChaoID:        gameconf.CFile.DefaultMainChao,
			MainChaoLevel:     0,
			SubChaoID:         gameconf.CFile.DefaultSubChao,
			SubChaoLevel:      0,
			Language:          0,
			League:            0,
			WrestleCount:      playerS.NumBattles,
			WrestleDamage:     playerS.DamageTotal,
			WrestleBeatFlg:    wrestleBeatFlg,
		}
		player, err := db.GetPlayer(strconv.Itoa(int(playerS.PlayerID)))
		if err == nil {
			mainCharaInfo, _ := player.GetChara(player.PlayerState.MainCharaID)
			mainCharaLevel := mainCharaInfo.Level
			subCharaLevel := int64(0)
			mainChaoLevel := int64(0)
			subChaoLevel := int64(0)
			if player.PlayerState.SubCharaID != "-1" {
				subCharaInfo, _ := player.GetChara(player.PlayerState.SubCharaID)
				subCharaLevel = subCharaInfo.Level
			}
			if player.PlayerState.MainChaoID != "-1" {
				mainChaoInfo, _ := player.GetChao(player.PlayerState.MainChaoID)
				mainChaoLevel = mainChaoInfo.Level
			}
			if player.PlayerState.SubChaoID != "-1" {
				subChaoInfo, _ := player.GetChao(player.PlayerState.SubChaoID)
				subChaoLevel = subChaoInfo.Level
			}
			userState = netobj.RaidBossUserState{
				WrestleID:         strconv.Itoa(int(playerS.PlayerID)),
				Name:              player.Username,
				Grade:             playerS.Grade,
				NumRank:           player.PlayerState.Rank,
				LoginTime:         player.LastLogin,
				CharacterID:       player.PlayerState.MainCharaID,
				CharacterLevel:    mainCharaLevel,
				SubcharacterID:    player.PlayerState.SubCharaID,
				SubcharacterLevel: subCharaLevel,
				MainChaoID:        player.PlayerState.MainChaoID,
				MainChaoLevel:     mainChaoLevel,
				SubChaoID:         player.PlayerState.SubChaoID,
				SubChaoLevel:      subChaoLevel,
				Language:          0,
				League:            player.PlayerState.RankingLeague,
				WrestleCount:      playerS.NumBattles,
				WrestleDamage:     playerS.DamageTotal,
				WrestleBeatFlg:    wrestleBeatFlg,
			}
		}
		userStates = append(userStates, userState)
	}
	return userStates
}

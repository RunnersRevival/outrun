package dbaccess

import (
	"log"

	"github.com/Mtbcooler/outrun/config/eventconf"
	"github.com/Mtbcooler/outrun/consts"
	"github.com/Mtbcooler/outrun/enums"
	"github.com/Mtbcooler/outrun/netobj"
	"github.com/Mtbcooler/outrun/netobj/constnetobjs"
	"github.com/Mtbcooler/outrun/obj"
)

func GetPlayerFromDB(id string) (netobj.Player, error) {
	// TODO: This is slow! Phase out the usage of this function ASAP!
	playerinfo, err := GetPlayerInfo(consts.DBMySQLTableCorePlayerInfo, id)
	if err != nil {
		return constnetobjs.BlankPlayer, err
	}
	playerstate, err := GetPlayerState(consts.DBMySQLTablePlayerStates, id)
	if err != nil {
		return constnetobjs.BlankPlayer, err
	}
	mileagemapstate, err := GetMileageMapState(consts.DBMySQLTableMileageMapStates, id)
	if err != nil {
		return constnetobjs.BlankPlayer, err
	}
	optionuserresult, err := GetOptionUserResult(consts.DBMySQLTableOptionUserResults, id)
	if err != nil {
		return constnetobjs.BlankPlayer, err
	}
	wheeloptions := netobj.DefaultWheelOptions(playerstate.NumRouletteTicket, 0, enums.WheelRankNormal, 5, 0) // TODO: REPLACE ME! FOR TESTING ONLY!
	rouletteinfo, err := GetRouletteInfo(consts.DBMySQLTableRouletteInfos, id)
	if err != nil {
		return constnetobjs.BlankPlayer, err
	}
	allowedCharacters := []string{}
	allowedChao := []string{}
	for _, chao := range playerinfo.ChaoState {
		if chao.Level < 10 { // not max level
			allowedChao = append(allowedChao, chao.ID)
		}
	}
	for _, character := range playerinfo.CharacterState {
		if character.Star < 10 { // not max star
			allowedCharacters = append(allowedCharacters, character.ID)
		}
	}
	chaoroulettedata := netobj.DefaultChaoRouletteGroup(playerstate, allowedCharacters, allowedChao) // TODO: REPLACE ME! FOR TESTING ONLY!
	loginbonusstate, err := GetLoginBonusState(consts.DBMySQLTableLoginBonusStates, id)
	if err != nil {
		return constnetobjs.BlankPlayer, err
	}
	player := netobj.NewPlayer(
		id,
		playerinfo.Username,
		playerinfo.Password,
		playerinfo.MigrationPassword,
		playerinfo.UserPassword,
		playerinfo.Key,
		playerinfo.Language,
		playerstate,
		playerinfo.CharacterState,
		playerinfo.ChaoState,
		mileagemapstate,
		[]netobj.MileageFriend{},
		netobj.DefaultPlayerVarious(),
		optionuserresult,
		wheeloptions,
		rouletteinfo,
		chaoroulettedata,
		[]eventconf.ConfiguredEvent{},
		[]obj.Message{},
		[]obj.OperatorMessage{},
		loginbonusstate,
	)
	return player, nil
}

func InitializeTablesIfNecessary() error {
	log.Println("[INFO] Initializing analytics table... (1/12)")
	_, err := db.Exec(consts.SQLAnalyticsSchema)
	if err != nil {
		return err
	}
	log.Println("[INFO] Initializing player info table... (2/12)")
	_, err = db.Exec(consts.SQLCorePlayerInfoSchema)
	if err != nil {
		return err
	}
	log.Println("[INFO] Initializing player states table... (3/12)")
	_, err = db.Exec(consts.SQLPlayerStatesSchema)
	if err != nil {
		return err
	}
	log.Println("[INFO] Initializing mileage map states table... (4/12)")
	_, err = db.Exec(consts.SQLMileageMapStatesSchema)
	if err != nil {
		return err
	}
	log.Println("[INFO] Initializing user stats table... (5/12)")
	_, err = db.Exec(consts.SQLOptionUserResultsSchema)
	if err != nil {
		return err
	}
	log.Println("[INFO] Initializing roulette info states table... (6/12)")
	_, err = db.Exec(consts.SQLRouletteInfosSchema)
	if err != nil {
		return err
	}
	log.Println("[INFO] Initializing login bonus states table... (7/12)")
	_, err = db.Exec(consts.SQLLoginBonusStatesSchema)
	if err != nil {
		return err
	}
	log.Println("[INFO] Initializing operator messages table... (8/12)")
	_, err = db.Exec(consts.SQLOperatorMessagesSchema)
	if err != nil {
		return err
	}
	log.Println("[INFO] Initializing ranking league data table... (9/12)")
	_, err = db.Exec(consts.SQLRankingLeagueDataSchema)
	if err != nil {
		return err
	}
	log.Println("[INFO] Initializing session IDs table... (10/12)")
	_, err = db.Exec(consts.SQLSessionIDsSchema)
	if err != nil {
		return err
	}
	log.Println("[INFO] Initializing operator infos table... (11/12)")
	_, err = db.Exec(consts.SQLOperatorInfosSchema)
	if err != nil {
		return err
	}
	log.Println("[INFO] Initializing event states table... (12/12)")
	_, err = db.Exec(consts.SQLEventStatesSchema)
	return err
}
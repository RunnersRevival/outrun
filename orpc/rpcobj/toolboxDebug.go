package rpcobj

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/RunnersRevival/outrun/config/gameconf"
	"github.com/RunnersRevival/outrun/consts"
	"github.com/RunnersRevival/outrun/db"
	"github.com/RunnersRevival/outrun/db/dbaccess"
	"github.com/RunnersRevival/outrun/enums"
	"github.com/RunnersRevival/outrun/logic"
	"github.com/RunnersRevival/outrun/netobj"
	"github.com/RunnersRevival/outrun/netobj/constnetobjs"
	"github.com/RunnersRevival/outrun/obj"
	"github.com/RunnersRevival/outrun/obj/constobjs"
)

func (t *Toolbox) Debug_GetCampaignStatus(uid string, reply *ToolboxReply) error {
	player, err := db.GetPlayer(uid)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = "unable to get player: " + err.Error()
		return err
	}
	reply.Status = StatusOK
	reply.Info = strconv.Itoa(int(player.MileageMapState.Chapter)) + "," + strconv.Itoa(int(player.MileageMapState.Episode)) + "," + strconv.Itoa(int(player.MileageMapState.Point))
	return nil
}

func (t *Toolbox) Debug_GetAllPlayerIDs(nothing bool, reply *ToolboxReply) error {
	playerIDs := []string{}
	dbaccess.ForEachKey(consts.DBBucketPlayers, func(k, v []byte) error {
		playerIDs = append(playerIDs, string(k))
		return nil
	})
	final := strings.Join(playerIDs, ",")
	reply.Status = StatusOK
	reply.Info = final
	return nil
}

func (t *Toolbox) Debug_GetAllBattlePlayerIDs(nothing bool, reply *ToolboxReply) error {
	playerIDs := []string{}
	dbaccess.BattleDBForEachKey(consts.BattleDBBucketWaiting, func(k, v []byte) error {
		playerIDs = append(playerIDs, string(k))
		return nil
	})
	dbaccess.BattleDBForEachKey(consts.BattleDBBucketMatched, func(k, v []byte) error {
		playerIDs = append(playerIDs, string(k))
		return nil
	})
	final := strings.Join(playerIDs, ",")
	reply.Status = StatusOK
	reply.Info = final
	return nil
}

func (t *Toolbox) Debug_ResetPlayer(uid string, reply *ToolboxReply) error {
	player, err := db.GetPlayer(uid)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = "unable to get player: " + err.Error()
		return err
	}
	_ = db.NewAccountWithID(uid, player.ResetCount+1)
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_DeletePlayer(uid string, reply *ToolboxReply) error {
	err := db.DeletePlayer(uid)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = "unable to delete player: " + err.Error()
		return err
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_GetRouletteInfo(uid string, reply *ToolboxReply) error {
	player, err := db.GetPlayer(uid)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = "unable to get player: " + err.Error()
		return err
	}
	rouletteInfo := player.RouletteInfo
	jri, err := json.Marshal(rouletteInfo)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = "unable to marshal RouletteInfo: " + err.Error()
		return err
	}
	reply.Status = StatusOK
	reply.Info = string(jri)
	return nil
}

func (t *Toolbox) Debug_ResetChaoRouletteGroup(uid string, reply *ToolboxReply) error {
	player, err := db.GetPlayer(uid)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = "unable to get player: " + err.Error()
		return err
	}
	chaoRouletteGroup := netobj.DefaultChaoRouletteGroup(player.PlayerState, player.GetAllNonMaxedCharacters(), player.GetAllNonMaxedChao())
	player.ChaoRouletteGroup = chaoRouletteGroup
	err = db.SavePlayer(player)
	if err != nil {
		reply.Status = StatusOK
		reply.Info = "OK"
		return err
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_ResetCharactersAndCompensate(uid string, reply *ToolboxReply) error {
	player, err := db.GetPlayer(uid)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = "unable to get player: " + err.Error()
		return err
	}
	toAdd := int64(0)
	for _, char := range player.CharacterState {
		toAdd += char.Level * 7
	}
	player.PlayerState.NumRedRings += toAdd
	player.CharacterState = netobj.DefaultCharacterState()
	err = db.SavePlayer(player)
	if err != nil {
		reply.Status = StatusOK
		reply.Info = "OK"
		return err
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_ResetChao(uid string, reply *ToolboxReply) error {
	player, err := db.GetPlayer(uid)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = "unable to get player: " + err.Error()
		return err
	}
	player.ChaoState = constnetobjs.NetChaoList
	err = db.SavePlayer(player)
	if err != nil {
		reply.Status = StatusOK
		reply.Info = "OK"
		return err
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_UnlockAllInChaoState(uid string, reply *ToolboxReply) error {
	player, err := db.GetPlayer(uid)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = "unable to get player: " + err.Error()
		return err
	}
	player.ChaoState = constnetobjs.UnlockedChaoState()
	err = db.SavePlayer(player)
	if err != nil {
		reply.Status = StatusOK
		reply.Info = "OK"
		return err
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_MigrateUser(uidToUID string, reply *ToolboxReply) error {
	uidSrc := strings.Split(uidToUID, "->")
	if len(uidSrc) != 2 {
		reply.Status = StatusOtherError
		reply.Info = "improperly formatted string (Example: 1234567890->1987654321)"
	}

	fromUID := uidSrc[0]
	toUID := uidSrc[1]
	oldPlayer, err := db.GetPlayer(fromUID)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = err.Error()
		return err
	}
	currentPlayer, err := db.GetPlayer(toUID)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = err.Error()
		return err
	}
	currentPlayer.PlayerState = oldPlayer.PlayerState
	currentPlayer.Username = oldPlayer.Username
	currentPlayer.LastLogin = oldPlayer.LastLogin
	currentPlayer.CharacterState = oldPlayer.CharacterState
	currentPlayer.ChaoState = oldPlayer.ChaoState
	currentPlayer.MileageMapState = oldPlayer.MileageMapState
	currentPlayer.MileageFriends = oldPlayer.MileageFriends
	currentPlayer.PlayerVarious = oldPlayer.PlayerVarious
	currentPlayer.LastWheelOptions = oldPlayer.LastWheelOptions
	currentPlayer.ChaoRouletteGroup = oldPlayer.ChaoRouletteGroup
	currentPlayer.RouletteInfo = oldPlayer.RouletteInfo

	err = db.SavePlayer(currentPlayer)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = err.Error()
		return err
	}

	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_UsernameSearch(username string, reply *ToolboxReply) error {
	playerIDs := []string{}
	dbaccess.ForEachKey(consts.DBBucketPlayers, func(k, v []byte) error {
		playerIDs = append(playerIDs, string(k))
		return nil
	})
	sameUsernames := []string{}
	for _, uid := range playerIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = "error getting ID " + uid + ": " + err.Error()
			return err
		}
		if player.Username == username {
			sameUsernames = append(sameUsernames, player.ID)
		}
	}
	if len(sameUsernames) == 0 {
		reply.Status = StatusOtherError
		reply.Info = "unable to find ID for username " + username
		return nil
	}
	reply.Status = StatusOK
	reply.Info = strings.Join(sameUsernames, ",")
	return nil
}

func (t *Toolbox) Debug_RawPlayer(uid string, reply *ToolboxReply) error {
	playerSrc, err := dbaccess.Get(consts.DBBucketPlayers, uid)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = err.Error()
		return err
	}
	reply.Status = StatusOK
	reply.Info = string(playerSrc)
	return nil
}

func (t *Toolbox) Debug_ResetCharacterState(uid string, reply *ToolboxReply) error {
	player, err := db.GetPlayer(uid)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = "unable to get player: " + err.Error()
		return err
	}
	player.CharacterState = netobj.DefaultCharacterState()
	err = db.SavePlayer(player)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = err.Error()
		return err
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_UnlockAllInCharacterState(uid string, reply *ToolboxReply) error {
	player, err := db.GetPlayer(uid)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = "unable to get player: " + err.Error()
		return err
	}
	player.CharacterState = netobj.UnlockedCharacterState()
	err = db.SavePlayer(player)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = err.Error()
		return err
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_MatchPlayersToGameConf(uids string, reply *ToolboxReply) error {
	allUIDs := strings.Split(uids, ",")
	for _, uid := range allUIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}
		player.CharacterState = netobj.DefaultCharacterState() // already uses AllCharactersUnlocked
		player.ChaoState = constnetobjs.DefaultChaoState()     // already uses AllChaoUnlocked
		player.PlayerState.MainCharaID = gameconf.CFile.DefaultMainCharacter
		player.PlayerState.SubChaoID = gameconf.CFile.DefaultSubChao
		player.PlayerState.MainChaoID = gameconf.CFile.DefaultMainChao
		player.PlayerState.SubCharaID = gameconf.CFile.DefaultSubCharacter
		player.PlayerState.NumRings = gameconf.CFile.StartingRings
		player.PlayerState.NumRedRings = gameconf.CFile.StartingRedRings
		player.PlayerState.Energy = gameconf.CFile.StartingEnergy
		err = db.SavePlayer(player)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("error saving player %s: ", uid) + err.Error()
			return err
		}
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_PrepTag1p0(uids string, reply *ToolboxReply) error {
	allUIDs := strings.Split(uids, ",")
	sqrt := func(n int64) int64 {
		fn := float64(n)
		result := math.Sqrt(fn)
		return int64(result)
	}

	_ = sqrt

	for _, uid := range allUIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}
		// compensate first
		amountPerLevel := int64(7) // Red Rings offered per level
		newRedRingAmount := int64(5)
		for _, char := range player.CharacterState {
			newRedRingAmount += char.Level * amountPerLevel
		}
		player.CharacterState = netobj.DefaultCharacterState() // already uses AllCharactersUnlocked
		player.ChaoState = constnetobjs.DefaultChaoState()     // already uses AllChaoUnlocked
		player.PlayerState.MainCharaID = gameconf.CFile.DefaultMainCharacter
		player.PlayerState.SubChaoID = gameconf.CFile.DefaultSubChao
		player.PlayerState.MainChaoID = gameconf.CFile.DefaultMainChao
		player.PlayerState.SubCharaID = gameconf.CFile.DefaultSubCharacter
		player.PlayerState.NumRings = int64(10000)
		player.PlayerState.NumRedRings = newRedRingAmount
		player.PlayerState.Energy = gameconf.CFile.StartingEnergy
		player.PlayerState.Items = constobjs.DefaultPlayerStateItems

		err = db.SavePlayer(player)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("error saving player %s: ", uid) + err.Error()
			return err
		}
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_PlayersByPassword(password string, reply *ToolboxReply) error {
	foundPlayers, err := logic.FindPlayersByPassword(password, false)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = "error finding players by password: " + err.Error()
		return err
	}
	playerIDs := []string{}
	for _, player := range foundPlayers {
		playerIDs = append(playerIDs, player.ID)
	}
	final := strings.Join(playerIDs, ",")
	reply.Status = StatusOK
	reply.Info = final
	return nil
}

func (t *Toolbox) Debug_ResetPlayersRank(uids string, reply *ToolboxReply) error {
	allUIDs := strings.Split(uids, ",")
	for _, uid := range allUIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}
		player.PlayerState.Rank = 0 // for some reason, this gets incremented 1 by the game
		err = db.SavePlayer(player)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("error saving player %s: ", uid) + err.Error()
			return err
		}
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_FixMarineRedRings(uids string, reply *ToolboxReply) error {
	wh := constobjs.CharacterMarine
	whid := wh.ID
	whrr := wh.PriceRedRings
	whrrpr := wh.Price
	allUIDs := strings.Split(uids, ",")
	for _, uid := range allUIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}
		i := player.IndexOfChara(whid)
		if i == -1 {
			reply.Status = StatusOK
			reply.Info = "index not found!"
			return fmt.Errorf("index not found!")
		}
		player.CharacterState[i].Character.PriceRedRings = whrr // TODO: check if needed
		player.CharacterState[i].PriceRedRings = whrr
		player.CharacterState[i].Price = whrrpr
		err = db.SavePlayer(player)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("error saving player %s: ", uid) + err.Error()
			return err
		}
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_FixWerehogRedRings(uids string, reply *ToolboxReply) error {
	wh := constobjs.CharacterWerehog
	whid := wh.ID
	whrr := wh.PriceRedRings
	allUIDs := strings.Split(uids, ",")
	for _, uid := range allUIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}
		i := player.IndexOfChara(whid)
		if i == -1 {
			reply.Status = StatusOK
			reply.Info = "index not found!"
			return fmt.Errorf("index not found!")
		}
		player.CharacterState[i].Character.PriceRedRings = whrr // TODO: check if needed
		player.CharacterState[i].PriceRedRings = whrr
		err = db.SavePlayer(player)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("error saving player %s: ", uid) + err.Error()
			return err
		}
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_RemoveCharacter(character ChangeCharacter, reply *ToolboxReply) error {
	allUIDs := strings.Split(character.UIDs, ",")

	for _, uid := range allUIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}
		charaState := player.CharacterState
		index := player.IndexOfChara(character.ID)
		if index != -1 {
			player.CharacterState = append(charaState[:index], charaState[index+1:]...)
		}
		println("Removed Character ", character.ID, " to ", uid)
		err = db.SavePlayer(player)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("error saving player %s: ", uid) + err.Error()
			return err
		}
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_RemoveCharacterIfLocked(character ChangeCharacter, reply *ToolboxReply) error {
	allUIDs := strings.Split(character.UIDs, ",")

	for _, uid := range allUIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}
		charaState := player.CharacterState
		index := player.IndexOfChara(character.ID)
		if index != -1 {
			if player.CharacterState[index].Status == enums.CharacterStatusLocked {
				player.CharacterState = append(charaState[:index], charaState[index+1:]...)
				println("Removed Character ", character.ID, " to ", uid)
			} else {
				println("Player ", uid, " has unlocked character ID ", character.ID, ", ignoring. ")
			}
		}
		err = db.SavePlayer(player)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("error saving player %s: ", uid) + err.Error()
			return err
		}
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_AddCharacter(character ChangeCharacter, reply *ToolboxReply) error {
	allUIDs := strings.Split(character.UIDs, ",")
	var newCharacter = obj.Character{
		character.ID,
		character.Cost,
		character.NumRedRings,
		character.Price,         // used for limit breaking
		character.PriceRedRings, // red rings used for limit breaking
	}
	for _, uid := range allUIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}
		index := player.IndexOfChara(character.ID)
		if index == -1 {
			switch character.LockCondition {
			case 0:
				player.CharacterState = append(player.CharacterState, netobj.DefaultCharacter(newCharacter))
			case 1:
				player.CharacterState = append(player.CharacterState, netobj.DefaultStageLockedCharacter(newCharacter))
			case 2:
				player.CharacterState = append(player.CharacterState, netobj.DefaultLockedCharacter(newCharacter))
			case 3:
				player.CharacterState = append(player.CharacterState, netobj.DefaultRouletteOnlyLockedCharacter(newCharacter))
			}
			println("Added Character ", character.ID, " to ", uid)
			err = db.SavePlayer(player)
			if err != nil {
				reply.Status = StatusOtherError
				reply.Info = fmt.Sprintf("error saving player %s: ", uid) + err.Error()
				return err
			}
		} else {
			println("Character is already found in player's CharacterState!")
		}
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_FixGothicAmyPrices(uids string, reply *ToolboxReply) error {
	gc := constobjs.CharacterGothicAmy
	gcid := gc.ID
	gcrr := int64(0)
	gcc := gc.Cost
	allUIDs := strings.Split(uids, ",")
	for _, uid := range allUIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}
		i := player.IndexOfChara(gcid)
		if i != -1 {
			player.CharacterState[i].LockCondition = int64(0)
			player.CharacterState[i].Character.PriceRedRings = gcrr // TODO: check if needed
			player.CharacterState[i].PriceRedRings = gcrr
			player.CharacterState[i].Price = int64(0)
			player.CharacterState[i].Cost = gcc
			reply.Status = StatusOK
			println("Changed Cost values for ", uid)
			//return fmt.Errorf("index not found!")
			err = db.SavePlayer(player)
			if err != nil {
				reply.Status = StatusOtherError
				reply.Info = fmt.Sprintf("error saving player %s: ", uid) + err.Error()
				return err
			}
		} else {
			println("Skipped")
		}
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_UpdateCharacterLockCondition(character ChangeCharacter, reply *ToolboxReply) error {
	allUIDs := strings.Split(character.UIDs, ",")
	for _, uid := range allUIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}
		i := player.IndexOfChara(character.ID)
		if i != -1 {
			player.CharacterState[i].LockCondition = character.LockCondition
			// player.CharacterState[i].Character.PriceRedRings = gcrr // TODO: check if needed
			player.CharacterState[i].PriceRedRings = character.PriceRedRings
			player.CharacterState[i].Price = character.Price
			// player.CharacterState[i].Cost = gcc
			reply.Status = StatusOK
			println("Changed Lock Condition for ", uid)
			//return fmt.Errorf("index not found!")
			err = db.SavePlayer(player)
			if err != nil {
				reply.Status = StatusOtherError
				reply.Info = fmt.Sprintf("error saving player %s: ", uid) + err.Error()
				return err
			}
		} else {
			println("Skipped")
		}
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_FixVariantPricing(uids string, reply *ToolboxReply) error {
	cmap := map[string]obj.Character{
		"301000": constobjs.CharacterAmitieAmy,
		"301001": constobjs.CharacterGothicAmy,
		"301002": constobjs.CharacterHalloweenShadow,
		"301003": constobjs.CharacterHalloweenRouge,
		"301004": constobjs.CharacterHalloweenOmega,
		"301005": constobjs.CharacterXMasSonic,
		"301006": constobjs.CharacterXMasTails,
		"301007": constobjs.CharacterXMasKnuckles,
	}
	allUIDs := strings.Split(uids, ",")
	for _, uid := range allUIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}
		for i, netchar := range player.CharacterState {
			cid := netchar.ID
			char, ok := cmap[cid]
			if ok {
				defaultCost := char.Cost
				defaultNumRedRings := char.NumRedRings
				defaultPrice := char.Price
				defaultPriceRedRings := char.PriceRedRings
				player.CharacterState[i].Character.Cost = defaultCost // TODO: check if needed
				player.CharacterState[i].Cost = defaultCost
				player.CharacterState[i].Character.NumRedRings = defaultNumRedRings // TODO: check if needed
				player.CharacterState[i].NumRedRings = defaultNumRedRings
				player.CharacterState[i].Character.Price = defaultPrice // TODO: check if needed
				player.CharacterState[i].Price = defaultPrice
				player.CharacterState[i].Character.PriceRedRings = defaultPriceRedRings // TODO: check if needed
				player.CharacterState[i].PriceRedRings = defaultPriceRedRings
				reply.Status = StatusOtherError
				println("character with ID '" + cid + "' was updated successfully for player ID '" + uid + "'")
				//return fmt.Errorf(reply.Info)
			}
			err = db.SavePlayer(player)
			if err != nil {
				reply.Status = StatusOtherError
				reply.Info = fmt.Sprintf("error saving player %s: ", uid) + err.Error()
				return err
			}
		}
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

// This code ain't pretty, and takes a long time to execute depending on how many players are in the database!
func (t *Toolbox) Debug_SendOperatorMessageToAll(args SendOperatorMessageToAllArgs, reply *ToolboxReply) error {
	playerIDs := []string{}
	dbaccess.ForEachKey(consts.DBBucketPlayers, func(k, v []byte) error {
		playerIDs = append(playerIDs, string(k))
		return nil
	})
	for _, uid := range playerIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}
		if player.Messages == nil {
			player.Messages = []obj.Message{}
		}
		if player.OperatorMessages == nil {
			player.OperatorMessages = []obj.OperatorMessage{}
		}
		player.AddOperatorMessage(args.MessageContents, args.Item, args.ExpiresAfter)
		err = db.SavePlayer(player)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("error saving player %s: ", uid) + err.Error()
			return err
		}
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_SendOperatorMessage(args SendOperatorMessageArgs, reply *ToolboxReply) error {
	player, err := db.GetPlayer(args.UID)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = fmt.Sprintf("unable to get player %s: ", args.UID) + err.Error()
		return err
	}
	if player.Messages == nil {
		player.Messages = []obj.Message{}
	}
	if player.OperatorMessages == nil {
		player.OperatorMessages = []obj.OperatorMessage{}
	}
	player.AddOperatorMessage(args.MessageContents, args.Item, args.ExpiresAfter)
	err = db.SavePlayer(player)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = fmt.Sprintf("error saving player %s: ", args.UID) + err.Error()
		return err
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_SendOperatorMessageToSome(args SendOperatorMessageToSomeArgs, reply *ToolboxReply) error {
	uidList := strings.Split(args.UIDs, ",")

	for _, uid := range uidList {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}

		if player.Messages == nil {
			player.Messages = []obj.Message{}
		}
		if player.OperatorMessages == nil {
			player.OperatorMessages = []obj.OperatorMessage{}
		}
		player.AddOperatorMessage(args.MessageContents, args.Item, args.ExpiresAfter)
		err = db.SavePlayer(player)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("error saving player %s: ", uid) + err.Error()
			return err
		}
	}

	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_FixCharacterPrices(uids string, reply *ToolboxReply) error {
	// TODO: This function possibly needs to be adjusted if event characters are ever added.
	cmap := map[string]obj.Character{
		"300000": constobjs.CharacterSonic,
		"300001": constobjs.CharacterTails,
		"300002": constobjs.CharacterKnuckles,
		"300003": constobjs.CharacterAmy,
		"300004": constobjs.CharacterShadow,
		"300005": constobjs.CharacterBlaze,
		"300006": constobjs.CharacterRouge,
		"300007": constobjs.CharacterOmega,
		"300008": constobjs.CharacterBig,
		"300009": constobjs.CharacterCream,
		"300010": constobjs.CharacterEspio,
		"300011": constobjs.CharacterCharmy,
		"300012": constobjs.CharacterVector,
		"300013": constobjs.CharacterSilver,
		"300014": constobjs.CharacterMetalSonic,
		"300015": constobjs.CharacterClassicSonic,
		"300016": constobjs.CharacterWerehog,
		"300017": constobjs.CharacterSticks,
		"300018": constobjs.CharacterTikal,
		"300019": constobjs.CharacterMephiles,
		"300020": constobjs.CharacterPSISilver,
		"300021": constobjs.CharacterMarine,
		"301000": constobjs.CharacterAmitieAmy,
		"301001": constobjs.CharacterGothicAmy,
		"301002": constobjs.CharacterHalloweenShadow,
		"301003": constobjs.CharacterHalloweenRouge,
		"301004": constobjs.CharacterHalloweenOmega,
		"301005": constobjs.CharacterXMasSonic,
		"301006": constobjs.CharacterXMasTails,
		"301007": constobjs.CharacterXMasKnuckles,
	}
	allUIDs := strings.Split(uids, ",")
	for _, uid := range allUIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}
		log.Printf("[RPC-DEBUG] Resetting character prices for player %s\n", uid)
		for i, netchar := range player.CharacterState {
			cid := netchar.ID
			char, ok := cmap[cid]
			if !ok {
				reply.Status = StatusOtherError
				reply.Info = "character with ID '" + cid + "' was not found in CharacterState for player ID '" + uid + "'"
				return fmt.Errorf(reply.Info)
			}
			defaultCost := char.Cost
			defaultNumRedRings := char.NumRedRings
			defaultPrice := char.Price
			defaultPriceRedRings := char.PriceRedRings
			player.CharacterState[i].Character.Cost = defaultCost // TODO: check if needed
			player.CharacterState[i].Cost = defaultCost
			player.CharacterState[i].Character.NumRedRings = defaultNumRedRings // TODO: check if needed
			player.CharacterState[i].NumRedRings = defaultNumRedRings
			player.CharacterState[i].Character.Price = defaultPrice // TODO: check if needed
			player.CharacterState[i].Price = defaultPrice
			player.CharacterState[i].Character.PriceRedRings = defaultPriceRedRings // TODO: check if needed
			player.CharacterState[i].PriceRedRings = defaultPriceRedRings
		}
		err = db.SavePlayer(player)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("error saving player %s: ", uid) + err.Error()
			return err
		}
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

// Purges players who haven't logged in for 11 months
// If testmode is true, it only counts how many players could be purged
func (t *Toolbox) Debug_PurgeInactivePlayers(testmode bool, reply *ToolboxReply) error {
	playerIDs := []string{}
	dbaccess.ForEachKey(consts.DBBucketPlayers, func(k, v []byte) error {
		playerIDs = append(playerIDs, string(k))
		return nil
	})
	numberOfPurgedPlayers := 0
	for _, uid := range playerIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}
		if player.LastLogin < time.Now().AddDate(0, -11, 0).UTC().Unix() {
			log.Printf("[RPC-DEBUG] Player %s hasn't logged in for six or more months! (Last Login: %v) Purge disallowed = %t\n", uid, player.LastLogin, player.DisallowInactivePurge)
			if !player.DisallowInactivePurge {
				numberOfPurgedPlayers++
				if !testmode {
					err := db.DeletePlayer(uid)
					if err != nil {
						reply.Status = StatusOtherError
						reply.Info = "unable to delete player " + uid + ": " + err.Error() + " - purged " + strconv.Itoa(numberOfPurgedPlayers) + " inactive players"
						return err
					}
				}
			}
		}
	}
	reply.Status = StatusOK
	if testmode {
		reply.Info = "OK - found " + strconv.Itoa(numberOfPurgedPlayers) + " inactive players that are purgable"
	} else {
		reply.Info = "OK - purged " + strconv.Itoa(numberOfPurgedPlayers) + " inactive players"
	}
	return nil
}

// Returns how many players are on the server and how many of them are active
func (t *Toolbox) Debug_CountPlayers(nothing bool, reply *ToolboxReply) error {
	playerIDs := []string{}
	dbaccess.ForEachKey(consts.DBBucketPlayers, func(k, v []byte) error {
		playerIDs = append(playerIDs, string(k))
		return nil
	})
	numberOfActivePlayers := 0
	for _, uid := range playerIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}
		if player.LastLogin > time.Now().AddDate(0, -2, 0).UTC().Unix() {
			numberOfActivePlayers++
		}
	}
	reply.Status = StatusOK
	reply.Info = "OK - there are " + strconv.Itoa(len(playerIDs)) + " players on this Outrun for Revival instance, and of those, " + strconv.Itoa(numberOfActivePlayers) + " are active (have logged in during the past 2 months)"
	return nil
}

func (t *Toolbox) Debug_FixAllChaoStateNumAcquiredValues(nothing bool, reply *ToolboxReply) error {
	playerIDs := []string{}
	dbaccess.ForEachKey(consts.DBBucketPlayers, func(k, v []byte) error {
		playerIDs = append(playerIDs, string(k))
		return nil
	})
	for _, uid := range playerIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}
		chaoState := player.ChaoState

		for i, chao := range chaoState {
			if chao.Status != enums.ChaoStatusNotOwned {
				player.ChaoState[i].Acquired = chao.Level + 1
			} else {
				player.ChaoState[i].Acquired = 0
			}
		}
		err = db.SavePlayer(player)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("error saving player %s: ", uid) + err.Error()
			return err
		}
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

// Returns how many players are in the Battle database
func (t *Toolbox) Debug_CountBattlePlayers(nothing bool, reply *ToolboxReply) error {
	waiting := 0
	matched := 0
	waitingPlayerIDs := []string{}
	dbaccess.BattleDBForEachKey(consts.BattleDBBucketWaiting, func(k, v []byte) error {
		waitingPlayerIDs = append(waitingPlayerIDs, string(k))
		waiting++
		return nil
	})
	matchedPlayerIDs := []string{}
	dbaccess.BattleDBForEachKey(consts.BattleDBBucketMatched, func(k, v []byte) error {
		matchedPlayerIDs = append(matchedPlayerIDs, string(k))
		matched++
		return nil
	})
	reply.Status = StatusOK
	reply.Info = "OK - there are " + strconv.Itoa(waiting) + " players waiting for a battle and there are " + strconv.Itoa(matched) + " players that are already in a match."
	return nil
}

func (t *Toolbox) Debug_PlayersByMigrationPassword(mpassword string, reply *ToolboxReply) error {
	foundPlayers, err := logic.FindPlayersByMigrationPassword(mpassword, false)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = "error finding players by migration password: " + err.Error()
		return err
	}
	playerIDs := []string{}
	for _, player := range foundPlayers {
		playerIDs = append(playerIDs, player.ID)
	}
	final := strings.Join(playerIDs, ",")
	if len(playerIDs) == 0 {
		final = "-nothing found-"
	}
	reply.Status = StatusOK
	reply.Info = final
	return nil
}

func (t *Toolbox) Debug_ResetEventProgressForAll(nothing bool, reply *ToolboxReply) error {
	playerIDs := []string{}
	dbaccess.ForEachKey(consts.DBBucketPlayers, func(k, v []byte) error {
		playerIDs = append(playerIDs, string(k))
		return nil
	})
	for _, uid := range playerIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}
		player.EventState = netobj.DefaultEventState()
		err = db.SavePlayer(player)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("error saving player %s: ", uid) + err.Error()
			return err
		}
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

// Credit: FairPlay137 for starting this code
func (t *Toolbox) Debug_RemoveDuplicateCharaEntries(uids string, reply *ToolboxReply) error {
	allUIDs := strings.Split(uids, ",")
	playerIDsWithDupe := make([]string, 0)

	for _, uid := range allUIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}

		charaState := player.CharacterState
		seenIDs := make([]string, 0)
		dupedIndexArr := make([]int, 0)

		for index, charaState := range charaState {
			if doesElementExist(seenIDs, charaState.Character.ID) {
				// we found a dupe, add it to the array
				dupedIndexArr = append(dupedIndexArr, index)
				if !doesElementExist(playerIDsWithDupe, player.ID) {
					playerIDsWithDupe = append(playerIDsWithDupe, player.ID)
				}
			} else {
				seenIDs = append(seenIDs, charaState.Character.ID)
			}
		}

		// reverse dupe index array so we don't shift the array during removal
		for i, j := 0, len(dupedIndexArr)-1; i < j; i, j = i+1, j-1 {
			dupedIndexArr[i], dupedIndexArr[j] = dupedIndexArr[j], dupedIndexArr[i]
		}

		// make sure we're not rewriting character states that aren't broken, just to be safe
		if len(dupedIndexArr) != 0 {
			for index2 := 0; index2 < len(dupedIndexArr); index2++ {
				if dupedIndexArr[index2] < 0 || dupedIndexArr[index2] >= len(charaState) {
					log.Println("[RPC-DEBUG] The given index is out of bounds.")
				} else {
					// write the modified charaState to the player object
					player.CharacterState = append(charaState[:dupedIndexArr[index2]], charaState[dupedIndexArr[index2]+1:]...)
					charaState = player.CharacterState
				}
			}
		}

		err = db.SavePlayer(player)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("error saving player %s: ", uid) + err.Error()
			return err
		}
	}

	// convert our string array of IDs to byte array
	playerIDsDupeJoined := strings.Join(playerIDsWithDupe, ",")
	playerIDsDupeByte := []byte(playerIDsDupeJoined)

	os.MkdirAll("logging/ids_with_dupes/", 0644)
	path := "logging/ids_with_dupes/" + "ids.txt"
	err := ioutil.WriteFile(path, playerIDsDupeByte, 0644)

	if err != nil {
		log.Printf("[RPC-DEBUG] Unable to log IDs", err)
	}

	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func doesElementExist(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func (t *Toolbox) Debug_SetDisallowInactivePurgeFlag(uids string, reply *ToolboxReply) error {
	allUIDs := strings.Split(uids, ",")
	for _, uid := range allUIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}
		player.DisallowInactivePurge = true
		err = db.SavePlayer(player)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("error saving player %s: ", uid) + err.Error()
			return err
		}
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_ResetDisallowInactivePurgeFlag(uids string, reply *ToolboxReply) error {
	allUIDs := strings.Split(uids, ",")
	for _, uid := range allUIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}
		player.DisallowInactivePurge = false
		err = db.SavePlayer(player)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("error saving player %s: ", uid) + err.Error()
			return err
		}
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_UpdateAllContinueCounts(nothing bool, reply *ToolboxReply) error {
	playerIDs := []string{}
	dbaccess.ForEachKey(consts.DBBucketPlayers, func(k, v []byte) error {
		playerIDs = append(playerIDs, string(k))
		return nil
	})
	for _, uid := range playerIDs {
		player, err := db.GetPlayer(uid)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("unable to get player %s: ", uid) + err.Error()
			return err
		}
		player.PlayerVarious.OnePlayContinueCount = int64(2)
		err = db.SavePlayer(player)
		if err != nil {
			reply.Status = StatusOtherError
			reply.Info = fmt.Sprintf("error saving player %s: ", uid) + err.Error()
			return err
		}
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

func (t *Toolbox) Debug_AddTransferCredentialsForPlayer(uid string, reply *ToolboxReply) error {
	player, err := db.GetPlayer(uid)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = "unable to get player: " + err.Error()
		return err
	}
	if len(player.UserPassword) == 0 {
		reply.Status = StatusOtherError
		reply.Info = "specified player does not have a user password set"
		return nil
	}
	transferCreds := netobj.PlayerToTransferCredentials(player)
	err = db.SaveTransferCredentials(transferCreds)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = "unable to save transfer credentials: " + err.Error()
		return err
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

// exports a player's full JSON data based on UID and username. to protect player data, a full player data export for a given player requires both the user ID and the username.
func (t *Toolbox) Debug_ExportPlayerData(uidplususername string, reply *ToolboxReply) error {
	if len(uidplususername) <= 10 {
		reply.Status = StatusOtherError
		reply.Info = "need uid and username to proceed. ex: \"1234567890Runners\""
		return nil
	}
	uid := uidplususername[:10]
	username := uidplususername[10:]
	player, err := db.GetPlayer(uid)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = "unable to get player: " + err.Error()
		return err
	}
	if player.Username != username {
		reply.Status = StatusOtherError
		reply.Info = "verification failed: \"" + username + "\" is not the correct username for uid " + uid
		return nil
	}
	result, err := json.Marshal(player)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = "unable to marshal player data: " + err.Error()
		return err
	}
	reply.Status = StatusOK
	reply.Info = string(result)
	return nil
}

// imports a player into the database, just as long as it doesn't conflict with another player
func (t *Toolbox) Debug_ImportPlayerData(jsondata string, reply *ToolboxReply) error {
	var importedPlayer netobj.Player
	err := json.Unmarshal([]byte(jsondata), &importedPlayer)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = "unable to unmarshal JSON data: " + err.Error()
		return err
	}
	if db.CheckPlayerExists(importedPlayer.ID) {
		reply.Status = StatusOtherError
		reply.Info = "a player with that UID already exists in the database"
		return nil
	}
	err = db.SavePlayer(importedPlayer)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = fmt.Sprintf("error importing player %s: ", importedPlayer.ID) + err.Error()
		return err
	}
	reply.Status = StatusOK
	reply.Info = "OK"
	return nil
}

// forcibly imports a player into the database, overwriting anything that might have been there before
func (t *Toolbox) Debug_ForceImportPlayerData(jsondata string, reply *ToolboxReply) error {
	var importedPlayer netobj.Player
	err := json.Unmarshal([]byte(jsondata), &importedPlayer)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = "unable to unmarshal JSON data: " + err.Error()
		return err
	}
	overwroteData := db.CheckPlayerExists(importedPlayer.ID)
	err = db.SavePlayer(importedPlayer)
	if err != nil {
		reply.Status = StatusOtherError
		reply.Info = fmt.Sprintf("error importing player %s: ", importedPlayer.ID) + err.Error()
		return err
	}
	reply.Status = StatusOK
	if overwroteData {
		reply.Info = "OK - overwrote another player!"
	} else {
		reply.Info = "OK"
	}
	return nil
}

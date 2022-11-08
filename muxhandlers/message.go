package muxhandlers

import (
	"encoding/json"
	"strconv"

	"github.com/RunnersRevival/outrun/db/dbaccess"

	"github.com/RunnersRevival/outrun/db"
	"github.com/RunnersRevival/outrun/emess"
	"github.com/RunnersRevival/outrun/enums"
	"github.com/RunnersRevival/outrun/helper"
	"github.com/RunnersRevival/outrun/logic"
	"github.com/RunnersRevival/outrun/netobj"
	"github.com/RunnersRevival/outrun/obj"
	"github.com/RunnersRevival/outrun/requests"
	"github.com/RunnersRevival/outrun/responses"
	"github.com/RunnersRevival/outrun/status"
)

func GetMessageList(helper *helper.Helper) {
	sid, _ := helper.GetSessionID()
	if !helper.CheckSession(true) {
		return
	}
	uid, err := helper.GetCallingPlayerID()
	if err != nil {
		helper.InternalErr("Error getting player ID", err)
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	messages := []obj.Message{}
	opmessages, _ := dbaccess.GetOperatorMessages(uid)

	//response := responses.DefaultMessageList(baseInfo)
	response := responses.MessageList(baseInfo, messages, opmessages)
	response.Seq, _ = db.BoltGetSessionIDSeq(sid)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetMessage(helper *helper.Helper) {
	if !helper.CheckSession(true) {
		return
	}
	data := helper.GetGameRequest()
	var request requests.GetMessageRequest
	err := json.Unmarshal(data, &request)
	if err != nil {
		helper.InternalErr("Error unmarshalling", err)
		return
	}
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	messages := []obj.Message{}
	opmessages, _ := dbaccess.GetOperatorMessages(player.ID)

	presentList := []obj.Present{}
	acceptingMessages := false
	acceptingOperatorMessages := false

	switch messageIds := request.MessageIDs.(type) {
	case []interface{}:
		acceptingMessages = true
		helper.DebugOut("%v", messageIds)
		for _, msgid := range messageIds {
			helper.DebugOut("Accepting message ID %v", msgid)
			// TODO: Implement actual logic here!
		}
	case string:
		if request.MessageIDs.(string) == "0" {
			helper.DebugOut("No messages to accept")
		} else {
			helper.Warn("Unexpected string value \"%s\" for request.MessageIDs", request.MessageIDs.(string))
			helper.InvalidRequest()
			return
		}
	default:
		helper.Warn("Unexpected type of request.MessageIDs")
		helper.InvalidRequest()
		return
	}

	switch operatorMessageIds := request.OperatorMessageIDs.(type) {
	case []interface{}:
		acceptingOperatorMessages = true
		helper.DebugOut("%v", operatorMessageIds)
		for _, omsgid := range operatorMessageIds {
			helper.DebugOut("Accepting operator message ID %v", omsgid)
			// TODO: Implement actual logic here!
		}
	case string:
		if request.OperatorMessageIDs.(string) == "0" {
			helper.DebugOut("No operator messages to accept")
		} else {
			helper.Warn("Unexpected string value \"%s\" for request.OperatorMessageIDs", request.OperatorMessageIDs.(string))
			helper.InvalidRequest()
			return
		}
	default:
		helper.Warn("Unexpected type of request.OperatorMessageIDs")
		helper.InvalidRequest()
		return
	}

	if !acceptingMessages && !acceptingOperatorMessages {
		helper.DebugOut("Assuming this is an 'Accept All Gifts' command...")
		for _, msg := range messages {
			helper.DebugOut("Accepting message ID %s", msg.ID)
			// TODO: Put actual code here!
		}
		for _, omsg := range opmessages {
			helper.DebugOut("Accepting operator message ID %s", omsg.ID)
			// TODO: Put actual code here!
		}
	}
	// TODO: Combine stackable items in present list, as that's how the original server went about creating the present list

	helper.DebugOut("%v", presentList)
	for _, currentPresent := range presentList {
		itemid := strconv.Itoa(int(currentPresent.ItemID))
		helper.DebugOut("Present: %s", itemid)
		helper.DebugOut("Present amount: %v", currentPresent.NumItem)
		if itemid[:2] == "12" { // ID is an item
			// check if the item is already in the player's inventory
			for _, item := range player.PlayerState.Items {
				if item.ID == itemid { // item found, increment amount
					item.Amount += currentPresent.NumItem
					break
				}
			}
		} else if itemid == strconv.Itoa(enums.ItemIDRing) { // Rings
			player.PlayerState.NumRings += currentPresent.NumItem
		} else if itemid == strconv.Itoa(enums.ItemIDRedRing) { // Red rings
			player.PlayerState.NumRedRings += currentPresent.NumItem
		} else if itemid == strconv.Itoa(enums.ItemIDEnergy) { // Revive tokens
			player.PlayerState.Energy += currentPresent.NumItem
		} else if itemid == strconv.Itoa(enums.IDSpecialEgg) {
			player.PlayerState.ChaoEggs += currentPresent.NumItem
			player.ChaoRouletteGroup.ChaoWheelOptions = netobj.DefaultChaoWheelOptions(player.PlayerState) //refresh chao wheel
		} else if itemid == strconv.Itoa(enums.IDRouletteTicketPremium) {
			player.PlayerState.NumChaoRouletteTicket += currentPresent.NumItem
			player.ChaoRouletteGroup.ChaoWheelOptions = netobj.DefaultChaoWheelOptions(player.PlayerState) //refresh chao wheel
		} else if itemid == strconv.Itoa(enums.IDRouletteTicketItem) {
			player.PlayerState.NumRouletteTicket += currentPresent.NumItem
			player.LastWheelOptions = logic.WheelRefreshLogic(player, player.LastWheelOptions) //refresh wheel
		} else if itemid[:2] == "40" { // ID is a Chao
			chaoIndex := player.IndexOfChao(itemid)
			if chaoIndex == -1 { // chao index not found, should never happen
				helper.InternalErr("cannot get index of chao '"+strconv.Itoa(chaoIndex)+"'", err)
				return
			}
			if player.ChaoState[chaoIndex].Status == enums.ChaoStatusNotOwned {
				// earn the Chao
				player.ChaoState[chaoIndex].Status = enums.ChaoStatusOwned
				player.ChaoState[chaoIndex].Acquired = currentPresent.NumItem
				player.ChaoState[chaoIndex].Level = currentPresent.NumItem - 1
			} else {
				player.ChaoState[chaoIndex].Acquired += currentPresent.NumItem
				player.ChaoState[chaoIndex].Level += currentPresent.NumItem
				if player.ChaoState[chaoIndex].Level > 10 { // if max chao level
					player.ChaoState[chaoIndex].Level = 10                        // reset to maximum
					player.ChaoState[chaoIndex].Status = enums.ChaoStatusMaxLevel // set status to MaxLevel
				}
			}

		} else if itemid[:2] == "30" { // ID is a character
			charIndex := player.IndexOfChara(itemid)
			if charIndex == -1 { // character index not found - append to character state
				helper.DebugOut("Character ID %s is not in CharacterState; adding...", itemid)
				player.CharacterState = append(player.CharacterState, netobj.DefaultCharacter(netobj.GenerateCharacterFromCharacterID(itemid)))
			} else {
				if player.CharacterState[charIndex].Status == enums.CharacterStatusLocked {
					// unlock the character
					player.CharacterState[charIndex].Status = enums.CharacterStatusUnlocked
				} else {
					starUpCount := currentPresent.NumItem
					for starUpCount > 0 && player.CharacterState[charIndex].Star < 10 { // 10 is max amount of stars a character can have before game breaks
						starUpCount--
						player.CharacterState[charIndex].Star++
					}
				}
			}
		} else {
			helper.Out("Unknown present ID %s", itemid)
		}
	}
	var response interface{}
	if baseInfo.StatusCode == status.OK {
		response = responses.GetMessage(baseInfo, player, presentList, player.GetAllOperatorMessageIDs(), player.GetAllOperatorMessageIDs())
	} else {
		response = responses.NewBaseResponse(baseInfo)
	}
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
	db.SavePlayer(player)
}

func SendEnergy(helper *helper.Helper) {
	// agnostic; sending energy is not possible without having facebook integration working
	baseInfo := helper.BaseInfo(emess.OK, status.AlreadySentEnergy)
	response := responses.NewBaseResponse(baseInfo)
	err := helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

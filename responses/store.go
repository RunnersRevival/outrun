package responses

import (
	"github.com/RunnersRevival/outrun/netobj"
	"github.com/RunnersRevival/outrun/obj"
	"github.com/RunnersRevival/outrun/responses/responseobjs"
)

type RedStarExchangeListResponse struct {
	BaseResponse
	ItemList      []obj.RedStarItem `json:"itemList"`
	TotalItems    int64             `json:"totalItems"`
	MonthPurchase int64             `json:"monthPurchase"`
	Birthday      string            `json:"birthday"`
}

func RedStarExchangeList(base responseobjs.BaseInfo, itemList []obj.RedStarItem, monthPurchase int64, birthday string) RedStarExchangeListResponse {
	baseResponse := NewBaseResponse(base)
	totalItems := int64(len(itemList))
	return RedStarExchangeListResponse{
		baseResponse,
		itemList,
		totalItems,
		monthPurchase,
		birthday,
	}
}

func DefaultRedStarExchangeList(base responseobjs.BaseInfo) RedStarExchangeListResponse {
	itemList := []obj.RedStarItem{}
	monthPurchase := int64(0)
	birthday := "1900-1-1"
	return RedStarExchangeList(base, itemList, monthPurchase, birthday)
}

type RedStarExchangeResponse struct {
	BaseResponse
	PlayerState netobj.PlayerState `json:"playerState"`
}

func RedStarExchange(base responseobjs.BaseInfo, playerState netobj.PlayerState) RedStarExchangeResponse {
	baseResponse := NewBaseResponse(base)
	return RedStarExchangeResponse{
		baseResponse,
		playerState,
	}
}

func DefaultRedStarExchange(base responseobjs.BaseInfo, player netobj.Player) RedStarExchangeResponse {
	playerState := player.PlayerState
	return RedStarExchange(
		base,
		playerState,
	)
}

type SetBirthdayResponse struct {
	BaseResponse
	Birthday string `json:"birthday"`
}

func SetBirthday(base responseobjs.BaseInfo, birthday string) SetBirthdayResponse {
	baseResponse := NewBaseResponse(base)
	return SetBirthdayResponse{
		baseResponse,
		birthday,
	}
}

func DefaultSetBirthday(base responseobjs.BaseInfo) SetBirthdayResponse {
	birthday := "1900-1-1"
	return SetBirthday(
		base,
		birthday,
	)
}

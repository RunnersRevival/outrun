package muxhandlers

import (
    "github.com/RunnersRevival/outrun/emess"
    "github.com/RunnersRevival/outrun/helper"
    "github.com/RunnersRevival/outrun/responses"
    "github.com/RunnersRevival/outrun/status"
)

func GetItemStockNum(helper *helper.Helper) {
    // TODO: Flesh out properly! The game responds with
    // [IDRouletteTicketPremium, IDRouletteTicketItem, IDSpecialEgg]
    // for item IDs, along with an event ID, likely for event characters.
    baseInfo := helper.BaseInfo(emess.OK, status.OK)
    response := responses.DefaultItemStockNum(baseInfo)
    err := helper.SendResponse(response)
    if err != nil {
        helper.InternalErr("Error sending response", err)
    }
}

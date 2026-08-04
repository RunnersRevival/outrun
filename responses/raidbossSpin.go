package responses

import (
    "github.com/RunnersRevival/outrun/obj"
    "github.com/RunnersRevival/outrun/obj/constobjs"
    "github.com/RunnersRevival/outrun/responses/responseobjs"
)

type ItemStockNumResponse struct {
    BaseResponse
    ItemStockList []obj.Item `json:"itemStockList"`
}

func ItemStockNum(base responseobjs.BaseInfo, itemStockList []obj.Item) ItemStockNumResponse {
    baseResponse := NewBaseResponse(base)
    return ItemStockNumResponse{
        baseResponse,
        itemStockList,
    }
}

func DefaultItemStockNum(base responseobjs.BaseInfo) ItemStockNumResponse {
    return ItemStockNum(
        base,
        constobjs.DefaultSpinItems,
    )
}

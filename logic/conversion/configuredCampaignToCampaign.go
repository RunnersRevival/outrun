package conversion

import (
	"time"

	"github.com/RunnersRevival/outrun/config/campaignconf"
	"github.com/RunnersRevival/outrun/obj"
	"github.com/jinzhu/now"
)

func ConfiguredCampaignToCampaign(cc campaignconf.ConfiguredCampaign) obj.Campaign {
	// Should be used by the game as soon as possible
	startTime := cc.StartTime
	switch startTime {
	case -2:
		startTime = now.BeginningOfDay().Unix()
	case -3:
		startTime = now.EndOfDay().Unix()
	case -4:
		startTime = time.Now().Unix() - 1
	}
	endTime := cc.EndTime
	switch endTime {
	case -2:
		endTime = now.BeginningOfDay().Unix()
	case -3:
		endTime = now.EndOfDay().Unix()
	case -4:
		endTime = time.Now().Add(24 * time.Hour).Unix()
	}
	newEvent := obj.Campaign{
		cc.RealType(),
		cc.Content,
		cc.SubContent,
		startTime,
		endTime,
	}
	return newEvent
}

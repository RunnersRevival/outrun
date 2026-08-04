package constobjs

import (
	"github.com/RunnersRevival/outrun/enums"
	"github.com/RunnersRevival/outrun/obj"
)

var DefaultLoginBonusRewardList = func() []obj.LoginBonusReward {
	return []obj.LoginBonusReward{
		obj.NewLoginBonusReward(
			[]obj.SelectReward{
				obj.NewSelectReward(
					[]obj.Item{
						obj.NewItem(enums.ItemIDStrRing, 3000),
						obj.NewItem(enums.ItemIDStrItemRouletteTicket, 1),
					},
				),
			},
		),
		obj.NewLoginBonusReward(
			[]obj.SelectReward{
				obj.NewSelectReward(
					[]obj.Item{
						obj.NewItem(enums.ItemIDStrRing, 3000),
						obj.NewItem(enums.ItemIDStrItemRouletteTicket, 1),
					},
				),
			},
		),
		obj.NewLoginBonusReward(
			[]obj.SelectReward{
				obj.NewSelectReward(
					[]obj.Item{
						obj.NewItem(enums.ItemIDStrRing, 5000),
						obj.NewItem(enums.ItemIDStrItemRouletteTicket, 1),
					},
				),
			},
		),
		obj.NewLoginBonusReward(
			[]obj.SelectReward{
				obj.NewSelectReward(
					[]obj.Item{
						obj.NewItem(enums.ItemIDStrRing, 5000),
						obj.NewItem(enums.ItemIDStrItemRouletteTicket, 1),
					},
				),
			},
		),
		obj.NewLoginBonusReward(
			[]obj.SelectReward{
				obj.NewSelectReward(
					[]obj.Item{
						obj.NewItem(enums.ItemIDStrRedRing, 10),
						obj.NewItem(enums.ItemIDStrRing, 5000),
						obj.NewItem(enums.ItemIDStrItemRouletteTicket, 2),
					},
				),
			},
		),
		obj.NewLoginBonusReward(
			[]obj.SelectReward{
				obj.NewSelectReward(
					[]obj.Item{
						obj.NewItem(enums.ItemIDStrRing, 10000),
						obj.NewItem(enums.ItemIDStrItemRouletteTicket, 2),
					},
				),
			},
		),
		obj.NewLoginBonusReward(
			[]obj.SelectReward{
				obj.NewSelectReward(
					[]obj.Item{
						obj.NewItem(enums.ItemIDStrRedRing, 20),
						obj.NewItem(enums.ItemIDStrRing, 15000),
						obj.NewItem(enums.ItemIDStrItemRouletteTicket, 2),
					},
				),
			},
		),
	}
}()

var EventLoginBonusRewardList = func() []obj.LoginBonusReward { // It's beginning to look a lot like Christmas...
	return []obj.LoginBonusReward{
		obj.NewLoginBonusReward( // Day 1
			[]obj.SelectReward{
				obj.NewSelectReward(
					[]obj.Item{
						obj.NewItem(enums.ItemIDStrRing, 3500),
						obj.NewItem(enums.ItemIDStrItemRouletteTicket, 1),
					},
				),
			},
		),
		obj.NewLoginBonusReward( // Day 2
			[]obj.SelectReward{
				obj.NewSelectReward(
					[]obj.Item{
						obj.NewItem(enums.ItemIDStrRedRing, 5),
						obj.NewItem(enums.ItemIDStrRing, 4000),
						obj.NewItem(enums.ItemIDStrItemRouletteTicket, 1),
						obj.NewItem(enums.ItemIDStrPremiumRouletteTicket, 1),
					},
				),
			},
		),
		obj.NewLoginBonusReward( // Day 3
			[]obj.SelectReward{
				obj.NewSelectReward(
					[]obj.Item{
						obj.NewItem(enums.ItemIDStrRing, 5000),
						obj.NewItem(enums.ItemIDStrRedRing, 5),
						obj.NewItem(enums.ItemIDStrPremiumRouletteTicket, 1),
					},
				),
			},
		),
		obj.NewLoginBonusReward( // Day 4
			[]obj.SelectReward{
				obj.NewSelectReward(
					[]obj.Item{
						obj.NewItem(enums.ItemIDStrRedRing, 10),
						obj.NewItem(enums.ItemIDStrRing, 6500),
						obj.NewItem(enums.ItemIDStrPremiumRouletteTicket, 1),
					},
				),
			},
		),
		obj.NewLoginBonusReward( // Day 5
			[]obj.SelectReward{
				obj.NewSelectReward(
					[]obj.Item{
						obj.NewItem(enums.ItemIDStrRing, 8000),
						obj.NewItem(enums.ItemIDStrRedRing, 10),
						obj.NewItem(enums.ItemIDStrPremiumRouletteTicket, 2),
						obj.NewItem(enums.ItemIDStrItemRouletteTicket, 2),
					},
				),
			},
		),
		obj.NewLoginBonusReward( // Day 6
			[]obj.SelectReward{
				obj.NewSelectReward(
					[]obj.Item{
						obj.NewItem(enums.ItemIDStrRedRing, 20),
						obj.NewItem(enums.ItemIDStrRing, 12000),
						obj.NewItem(enums.ItemIDStrPremiumRouletteTicket, 2),
						obj.NewItem(enums.ItemIDStrItemRouletteTicket, 5),
					},
				),
			},
		),
		obj.NewLoginBonusReward( // Day 7
			[]obj.SelectReward{
				obj.NewSelectReward(
					[]obj.Item{
						obj.NewItem(enums.CTStrXMasSonic, 1),
						obj.NewItem(enums.ItemIDStrRedRing, 25),
						obj.NewItem(enums.ItemIDStrRing, 20000),
						obj.NewItem(enums.ItemIDStrPremiumRouletteTicket, 2),
						obj.NewItem(enums.ItemIDStrItemRouletteTicket, 5),
					},
				),
			},
		),
	}
}()

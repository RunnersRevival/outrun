package consts

import "github.com/RunnersRevival/outrun/enums"

type PrizeInfo struct {
	AppearanceChance float64 // % chance for it to be chosen to be in wheel by the server
	Type             int64   // 0 for Chao, 1 for Character
}

// A 'load' as depicted below is the chance for the server to pick
// the associated item, where chosen is if randFloat(0, 100) < load.
// IMPORTANT: This load is exclusive to the rarity of the Chao that
// is being chosen by the server.

var RandomChaoWheelCharacterPrizes = map[string]float64{
	// characterID: load
	// Hopefully this should sum up to 100 just for
	// simplicity, but it shouldn't be a requirement.
	enums.CTStrSonic:    0.3, // Initial character
	enums.CTStrTails:    0.3, // Obtained in story mode
	enums.CTStrKnuckles: 0.3, // Obtained in story mode
	enums.CTStrAmy:      0.3,
	enums.CTStrBig:      0.3,
	//enums.CTStrBlaze:           1.0, // Revival Event (Sonic Rush)
	enums.CTStrCharmy: 0.3,
	enums.CTStrCream:  0.3,
	enums.CTStrEspio:  0.3,
	//enums.CTStrMephiles:        0.0, // Revival Event
	enums.CTStrOmega: 0.3,
	//enums.CTStrPSISilver:       0.0, // Revival Event
	enums.CTStrRouge:  0.3,
	enums.CTStrShadow: 0.3,
	enums.CTStrMarine:          1.0, // Revival Event
	//enums.CTStrTangle: 1.0, // Revival Event
	//enums.CTStrWhisper: 1.0, // Revival Event
	//enums.CTStrSticks:          0.3, // Revival Event
	enums.CTStrTikal:           1.5, // Event (Sonic Adventure)
	enums.CTStrVector: 0.3,
	enums.CTStrWerehog:         0.3, // Revival Event
	//enums.CTStrClassicSonic: 1.0, // Event (Birthday)
	//enums.CTStrMetalSonic:      0.0, // Revival Event

	// enums.CTStrAmitieAmy:       1.0, // Event (Puyo Puyo Quest)
	// enums.CTStrGothicAmy:       0.0, // Revival Event
	// enums.CTStrHalloweenShadow: 1.0, // Event (Halloween)
	// enums.CTStrHalloweenRouge:  1.0, // Event (Halloween)
	// enums.CTStrHalloweenOmega:  1.0, // Event (Halloween)
	// enums.CTStrXMasSonic:       1.0, // Event (Christmas)
	// enums.CTStrXMasTails:       1.0, // Event (Christmas)
	// enums.CTStrXMasKnuckles:    1.0, // Event (Christmas)
	// enums.CTStrXT:              1.2, // Revival Event (Christmas)
}

var RandomChaoWheelChaoPrizes = map[string]float64{
	// TODO: Balance these
	enums.ChaoIDStrHeroChao:             2.2, // Event (Animal Rescue event 1.0)
	enums.ChaoIDStrGoldChao:             2.2, // Event (Animal Rescue event 1.0)
	enums.ChaoIDStrDarkChao:             2.2, // Event (Animal Rescue event 1.0)
	enums.ChaoIDStrJewelChao:            2.2, // Event (Animal Rescue event 1.0)
	enums.ChaoIDStrNormalChao:           3.0, // Event (Animal Rescue event 1.0)
	enums.ChaoIDStrOmochao:              3.0, // Event (Animal Rescue event 1.0)
	enums.ChaoIDStrRCMonkey:             1.0, // Event (Animal Rescue event 1.0)
	enums.ChaoIDStrRCSpring:             2.0,
	enums.ChaoIDStrRCElectromagnet:      2.0,
	enums.ChaoIDStrBabyCyanWisp:         3.0,
	enums.ChaoIDStrBabyIndigoWisp:       3.0,
	enums.ChaoIDStrBabyYellowWisp:       3.0,
	enums.ChaoIDStrRCPinwheel:           2.0,
	enums.ChaoIDStrRCPiggyBank:          1.2,
	enums.ChaoIDStrRCBalloon:            2.0,
	//enums.ChaoIDStrEasterChao:           1.0, // Event (Easter; Increase Odds During Event)
	//enums.ChaoIDStrEasterBunny:          1.0, // Event (Easter; Increase Odds During Event)
	//enums.ChaoIDStrMerlina:              1.5, // Event (Easter: Premium Roulette for Timed Mode Event Only; Obtainable through Rewards List Only for Story Event)
	//enums.ChaoIDStrPurplePapurisu:       2.5, // Event (Puyo Puyo Quest)
	//enums.ChaoIDStrSuketoudara:          2.5, // Event (Puyo Puyo Quest)
	//enums.ChaoIDStrCarbuncle:            2.5, // Event (Puyo Puyo Quest: Premium Roulette for Timed Mode Event Only; Obtainable through Rewards List Only for Story Event)
	enums.ChaoIDStrEggChao:           1.0,
	enums.ChaoIDStrPumpkinChao:       1.0,
	enums.ChaoIDStrSkullChao:         1.1,
	enums.ChaoIDStrYacker:            1.0,
	enums.ChaoIDStrRCGoldenPiggyBank: 1.0,
	enums.ChaoIDStrWizardChao:        1.0,
	enums.ChaoIDStrRCTurtle:          1.1,
	enums.ChaoIDStrRCUFO:             1.1,
	enums.ChaoIDStrRCBomber:          1.0,
	//enums.ChaoIDStrStarShapedMissile:    0.0, // Event (Zazz Raid Boss; Increase Odds During Event)
	//enums.ChaoIDStrRCSatellite:          0.0, // Event (Zazz Raid Boss; Increase Odds During Event)
	//enums.ChaoIDStrRCMoonMech:           0.0, // Event (Zazz Raid Boss; Only Obtainable through the Raid Boss Roulette, which is currently unavailable)
	//enums.ChaoIDStrRappy:                1.0, // Event (Phantasy Star Online 2)
	//enums.ChaoIDStrKuna:                 1.0, // Event (Phantasy Star Online 2)
	//enums.ChaoIDStrMagLv1:               1.0, // Event (Phantasy Star Online 2)
	enums.ChaoIDStrBlowfishTransporter:  1.0, // Event (Tropical Coast; Increase Odds During Event)
	//enums.ChaoIDStrMotherWisp:           1.2, // Event (Tropical Coast: Premium Roulette for Timed Mode Event Only; Obtainable through Rewards List Only for Story Event)
	enums.ChaoIDStrMarineChao:           1.0, // Event (Tropical Coast; Increase Odds During Event)
	//enums.ChaoIDStrGenesis:              1.5, // Event (Birthday: Premium Roulette for Timed Mode Event Only; Obtainable through Rewards List Only for Story Event)
	//enums.ChaoIDStrCartridge:            1.5, // Event (Birthday; Increase Odds During Event)
	//enums.ChaoIDStrDeathEgg:             1.0, // Event (Birthday; Increase Odds During Event)
	enums.ChaoIDStrRCFighter:            1.0,
	enums.ChaoIDStrRCHovercraft:         1.0,
	enums.ChaoIDStrRCHelicopter:         1.0,
	enums.ChaoIDStrGreenCrystalMonsterS: 1.0,
	enums.ChaoIDStrGreenCrystalMonsterL: 1.0,
	enums.ChaoIDStrRCAirship:            1.0,
	enums.ChaoIDStrMagicLamp:            2.0, // Event (Desert Ruins and Animal Rescue 2.0; Increase Odds During Event)
	//enums.ChaoIDStrDesertChao:           3.0, // Event (Desert Ruins; Increase Odds During Event)
	//enums.ChaoIDStrErazorDjinn:          2.5, // Event (Desert Ruins: Premium Roulette for Timed Mode Event Only; Obtainable through Rewards List Only for Story Event)
	//enums.ChaoIDStrNightopian:           2.0, // Event (NiGHTS)
	//enums.ChaoIDStrNiGHTS:               2.0, // Event (NiGHTS)
	//enums.ChaoIDStrReala:                2.0, // Event (NiGHTS)
	enums.ChaoIDStrSonicOmochao:         1.0, // Event (Team Sonic Omochao)
	enums.ChaoIDStrTailsOmochao:         1.0, // Event (Team Sonic Omochao)
	enums.ChaoIDStrKnucklesOmochao:      1.0, // Event (Team Sonic Omochao)
	//enums.ChaoIDStrKingBoomBoo:          1.5, // Event (Halloween: Premium Roulette for Timed Mode Event Only; Obtainable through Rewards List Only for Story Event)
	enums.ChaoIDStrBoo:                  1.0, // Event (Halloween; Increase Odds During Event)
	enums.ChaoIDStrHalloweenChao:        1.0, // Event (Halloween; Increase Odds During Event)
	//enums.ChaoIDStrHeavyBomb:            1.2, // Event (Fantasy Zone)
	//enums.ChaoIDStrOPapa:                1.5, // Event (Fantasy Zone)
	//enums.ChaoIDStrOpaOpa:               1.5, // Event (Fantasy Zone)
	enums.ChaoIDStrBlockBomb:  1.0,
	enums.ChaoIDStrHunkofMeat: 1.3,
	// enums.ChaoIDStrYeti:                 1.5, // Event (Christmas)
	// enums.ChaoIDStrSnowChao:             2.0, // Event (Christmas)
	// enums.ChaoIDStrChristmasYeti:        1.5, // Event (Christmas)
	// enums.ChaoIDStrChristmasNiGHTS:      1.5, // Event (Christmas NiGHTS)
	// enums.ChaoIDStrIdeya:                1.5, // Event (Christmas NiGHTS)
	// enums.ChaoIDStrChristmasNightopian:  1.5, // Event (Christmas NiGHTS)
	enums.ChaoIDStrOrbot:      1.0,
	enums.ChaoIDStrCubot:      1.0,
	enums.ChaoIDStrLightChaos: 1.5,
	enums.ChaoIDStrHeroChaos:  1.5,
	enums.ChaoIDStrDarkChaos:  1.5,
	enums.ChaoIDStrChip:       1.5,
	//enums.ChaoIDStrShahra:               0.0, // Runners' League Story Mode
	enums.ChaoIDStrCaliburn:         1.2,
	enums.ChaoIDStrKingArthursGhost: 1.0,
	enums.ChaoIDStrRCTornado:        1.0,
	enums.ChaoIDStrRCBattleCruiser:  3.0,
	enums.ChaoIDStrRedCrystalMonsterS: 1.0,
	enums.ChaoIDStrRedCrystalMonsterL: 1.0,
	enums.ChaoIDStrGoldenGoose:        1.0,
	enums.ChaoIDStrRCPirateSpaceship: 3.0,
	enums.ChaoIDStrGoldenAngel:       1.0,
	enums.ChaoIDStrRCTornado2:           3.5, // Event (Sonic Adventure; Increase Odds During Event)
	enums.ChaoIDStrChaos:                3.5, // Event (Sonic Adventure; Increase Odds During Event)
	enums.ChaoIDStrOrca:                 3.5, // Event (Sonic Adventure; Increase Odds During Event)
	//enums.ChaoIDStrChaoWalker:           0.0, // Daily Battle
	//enums.ChaoIDStrDarkQueen:            0.0, // Runners' League Timed Mode
	enums.ChaoIDStrRCBlockFace: 1.0,
	//enums.ChaoIDStrDFekt:                0.0, // Revival Event (assets TBD)
	//enums.ChaoIDStrDarkChaoWalker:       0.0, // Daily Battle?
	//enums.ChaoIDStrPrideChaoL:           2.2, // Revival Event (Pride Month Celebration)
	//enums.ChaoIDStrPrideChaoG:           2.2, // Revival Event (Pride Month Celebration)
	//enums.ChaoIDStrPrideChaoB:           2.2, // Revival Event (Pride Month Celebration)
	//enums.ChaoIDStrPrideChaoT:           2.2, // Revival Event (Pride Month Celebration)
	//enums.ChaoIDStrPrideChaoP:           2.2, // Revival Event (Pride Month Celebration)
	//enums.ChaoIDStrPrideChaoA:           2.2, // Revival Event (Pride Month Celebration)
	//enums.ChaoIDStrPrideChaoNB:          2.2, // Revival Event (Pride Month Celebration)
}

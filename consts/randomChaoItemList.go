package consts

import "github.com/RunnersRevival/outrun/enums"
h
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
	enums.CTStrSonic:    0.65, // Initial character - Permanent Pool
	enums.CTStrTails:    0.65, // Obtained in story mode - Permanent Pool
	enums.CTStrKnuckles: 0.65, // Obtained in story mode - Permanent Pool
	// enums.CTStrAmy:      0.65, // Permanent Pool; Purchasable Only
	// enums.CTStrBig:      0.65, // Permanent Pool; Purchasable Only
	enums.CTStrBlaze:           0.65, // Permanent Pool
	// enums.CTStrCharmy: 0.65, // Permanent Pool; Purchasable Only
	// enums.CTStrCream:  0.65, // Permanent Pool; Purchasable Only
	// enums.CTStrEspio:  0.65, // Permanent Pool; Purchasable Only
	enums.CTStrMephiles:        0.65, // Permanent Pool
	// enums.CTStrOmega: 0.65, // Permanent Pool; Purchasable Only
	enums.CTStrPSISilver:       0.65, // Permanent Pool
	// enums.CTStrRouge:  0.65, // Permanent Pool; Purchasable Only
	// enums.CTStrShadow: 0.65, // Permanent Pool; Purchasable Only
	enums.CTStrMarine:          0.65, // Permanent Pool
	enums.CTStrTangle: 0.65, // Permanent Pool
	enums.CTStrWhisper: 0.65, // Permanent Pool
	enums.CTStrSticks:          0.65, // Permanent Pool
	enums.CTStrTikal:           0.65, // Permanent Pool
	// enums.CTStrVector: 0.65, // Permanent Pool; Purchasable Only
	enums.CTStrWerehog:         0.65, // Permanent Pool
	enums.CTStrClassicSonic: 1.5, // Permanent Pool
	enums.CTStrMetalSonic:      0.65, // Permanent Pool
	// enums.CTStrSilver:      0.65, // Permanent Pool; Purchasable Only

	// enums.CTStrAmitieAmy:       1.81, // Event (Puyo Puyo Quest)
	// enums.CTStrGothicAmy:       1.81, // Revival Event
	// enums.CTStrHalloweenShadow: 1.81, // Event (Halloween)
	// enums.CTStrHalloweenRouge:  1.81, // Event (Halloween)
	// enums.CTStrHalloweenOmega:  1.81, // Event (Halloween)
	// enums.CTStrXMasSonic:       1.43, // Event (Christmas)
	// enums.CTStrXMasTails:       1.43, // Event (Christmas)
	// enums.CTStrXMasKnuckles:    1.43, // Event (Christmas)
	// enums.CTStrXT:              1.43, // Revival Event (Christmas)
}

var RandomChaoWheelChaoPrizes = map[string]float64{
	// TODO: Balance these
	enums.ChaoIDStrHeroChao:             2.13, // Event (Animal Rescue event 1.0) - Permanent Pool
	enums.ChaoIDStrGoldChao:             2.13, // Event (Animal Rescue event 1.0) - Permanent Pool
	enums.ChaoIDStrDarkChao:             2.13, // Event (Animal Rescue event 1.0) - Permanent Pool
	enums.ChaoIDStrJewelChao:            2.13, // Event (Animal Rescue event 1.0) - Permanent Pool
	enums.ChaoIDStrNormalChao:           2.13, // Event (Animal Rescue event 1.0) - Permanent Pool
	enums.ChaoIDStrOmochao:              2.13, // Event (Animal Rescue event 1.0) - Permanent Pool
	enums.ChaoIDStrRCMonkey:             2.13, // Event (Animal Rescue event 1.0) - Permanent Pool
	enums.ChaoIDStrRCSpring:             2.13, // Permanent Pool
	enums.ChaoIDStrRCElectromagnet:      2.13, // Permanent Pool
	enums.ChaoIDStrBabyCyanWisp:         2.13, // Permanent Pool
	enums.ChaoIDStrBabyIndigoWisp:       2.13, // Permanent Pool
	enums.ChaoIDStrBabyYellowWisp:       2.13, // Permanent Pool
	enums.ChaoIDStrRCPinwheel:           2.13, // Permanent Pool
	enums.ChaoIDStrRCPiggyBank:          2.13, // Permanent Pool
	enums.ChaoIDStrRCBalloon:            2.13, // Permanent Pool
	enums.ChaoIDStrEasterChao:           2.13, // Event (Easter; Increase Odds During Event) - Permanent Pool
	enums.ChaoIDStrEasterBunny:          0.21, // Event (Easter; Increase Odds During Event) - Permanent Pool
	//enums.ChaoIDStrMerlina:              1.5, // Event (Easter: Premium Roulette for Timed Mode Event Only; Obtainable through Rewards List Only for Story Event)
	// enums.ChaoIDStrPurplePapurisu:       6.8, // Event (Puyo Puyo Quest)
	//enums.ChaoIDStrSuketoudara:          2.5, // Event (Puyo Puyo Quest)
	//enums.ChaoIDStrCarbuncle:            3.67, // Event (Puyo Puyo Quest: Premium Roulette for Timed Mode Event Only; Obtainable through Rewards List Only for Story Event)
	enums.ChaoIDStrEggChao:           0.21, // Permanent Pool
	enums.ChaoIDStrPumpkinChao:       0.21, // Permanent Pool
	enums.ChaoIDStrSkullChao:         0.21, // Permanent Pool
	enums.ChaoIDStrYacker:            0.21, // Permanent Pool
	enums.ChaoIDStrRCGoldenPiggyBank: 0.21, // Permanent Pool
	enums.ChaoIDStrWizardChao:        0.21, // Permanent Pool
	enums.ChaoIDStrRCTurtle:          0.21, // Permanent Pool
	enums.ChaoIDStrRCUFO:             0.21, // Permanent Pool
	enums.ChaoIDStrRCBomber:          0.21, // Permanent Pool
	enums.ChaoIDStrStarShapedMissile:    0.21, // Event (Zazz Raid Boss; Increase Odds During Event) - Permanent Pool
	enums.ChaoIDStrRCSatellite:          0.21, // Event (Zazz Raid Boss; Increase Odds During Event) - Permanent Pool
	// enums.ChaoIDStrRCMoonMech:           3.67, // Event (Zazz Raid Boss; Only Obtainable through the Raid Boss Roulette, which is currently unavailable)
	// enums.ChaoIDStrRappy:                4.25, // Event (Phantasy Star Online 2)
	// enums.ChaoIDStrKuna:                 3.67, // Event (Phantasy Star Online 2)
	// enums.ChaoIDStrMagLv1:               6.18, // Event (Phantasy Star Online 2)
	enums.ChaoIDStrBlowfishTransporter:  0.21, // Event (Tropical Coast; Increase Odds During Event) - Permanent Pool
	//enums.ChaoIDStrMotherWisp:           1.2, // Event (Tropical Coast: Premium Roulette for Timed Mode Event Only; Obtainable through Rewards List Only for Story Event)
	enums.ChaoIDStrMarineChao:           0.21, // Event (Tropical Coast; Increase Odds During Event) - Permanent Pool
	//enums.ChaoIDStrGenesis:              1.5, // Event (Birthday: Premium Roulette for Timed Mode Event Only; Obtainable through Rewards List Only for Story Event)
	enums.ChaoIDStrCartridge:            3.4, // Event (Birthday; Increase Odds During Event) - Permanent Pool
 	enums.ChaoIDStrDeathEgg:             2.2, // Event (Birthday; Increase Odds During Event) - Permanent Pool
	enums.ChaoIDStrRCFighter:            0.21, // Permanent Pool
	enums.ChaoIDStrRCHovercraft:         0.21, // Permanent Pool
	enums.ChaoIDStrRCHelicopter:         0.21, // Permanent Pool
	enums.ChaoIDStrGreenCrystalMonsterS: 0.21, // Permanent Pool
	enums.ChaoIDStrGreenCrystalMonsterL: 0.21, // Permanent Pool
	enums.ChaoIDStrRCAirship:            0.21, // Permanent Pool
	enums.ChaoIDStrMagicLamp:            0.21, // Event (Desert Ruins and Animal Rescue 2.0; Increase Odds During Event) - Permanent Pool
	enums.ChaoIDStrDesertChao:           0.21, // Event (Desert Ruins; Increase Odds During Event) - Permanent Pool
	// enums.ChaoIDStrErazorDjinn:          4.0, // Event (Desert Ruins: Premium Roulette for Timed Mode Event Only; Obtainable through Rewards List Only for Story Event)
	// enums.ChaoIDStrNightopian:           4.25, // Event (NiGHTS)
	// enums.ChaoIDStrNiGHTS:               3.67, // Event (NiGHTS)
	// enums.ChaoIDStrReala:                3.67, // Event (NiGHTS)
	enums.ChaoIDStrSonicOmochao:         0.21, // Event (Team Sonic Omochao) - Permanent Pool
	enums.ChaoIDStrTailsOmochao:         0.21, // Event (Team Sonic Omochao) - Permanent Pool
	enums.ChaoIDStrKnucklesOmochao:      0.21, // Event (Team Sonic Omochao) - Permanent Pool
	//enums.ChaoIDStrKingBoomBoo:          1.5, // Event (Halloween: Premium Roulette for Timed Mode Event Only; Obtainable through Rewards List Only for Story Event)
	enums.ChaoIDStrBoo:                  0.21, // Event (Halloween; Increase Odds During Event) - Permanent Pool
	enums.ChaoIDStrHalloweenChao:        0.21, // Event (Halloween; Increase Odds During Event) - Permanent Pool
	// enums.ChaoIDStrHeavyBomb:            4.25, // Event (Fantasy Zone)
	// enums.ChaoIDStrOPapa:                3.4, // Event (Fantasy Zone)
	// enums.ChaoIDStrOpaOpa:               3.67, // Event (Fantasy Zone)
	enums.ChaoIDStrBlockBomb:  0.21, // Permanent Pool
	enums.ChaoIDStrHunkofMeat: 0.21, // Permanent Pool
	// enums.ChaoIDStrYeti:                 3.58, // Event (Christmas)
	// enums.ChaoIDStrSnowChao:             3.58, // Event (Christmas)
	// enums.ChaoIDStrChristmasYeti:        1.5, // Event (Christmas)
	// enums.ChaoIDStrChristmasNiGHTS:      4.0, // Event (Christmas NiGHTS)
	// enums.ChaoIDStrIdeya:                3.58, // Event (Christmas NiGHTS)
	// enums.ChaoIDStrChristmasNightopian:  3.58, // Event (Christmas NiGHTS)
	enums.ChaoIDStrOrbot:      0.21, // Permanent Pool
	enums.ChaoIDStrCubot:      0.21, // Permanent Pool
	enums.ChaoIDStrLightChaos: 1.16, // Permanent Pool
	enums.ChaoIDStrHeroChaos:  1.16, // Permanent Pool
	enums.ChaoIDStrDarkChaos:  1.16, // Permanent Pool
	enums.ChaoIDStrChip:       1.16, // Permanent Pool
	// enums.ChaoIDStrShahra:               4.0, // Runners' League Story Mode
	enums.ChaoIDStrCaliburn:         1.16, // Permanent Pool
	enums.ChaoIDStrKingArthursGhost: 1.16, // Permanent Pool
	enums.ChaoIDStrRCTornado:        1.16, // Permanent Pool
	enums.ChaoIDStrRCBattleCruiser:  1.16, // Permanent Pool
	enums.ChaoIDStrRedCrystalMonsterS: 1.16, // Permanent Pool
	enums.ChaoIDStrRedCrystalMonsterL: 1.16, // Permanent Pool
	enums.ChaoIDStrGoldenGoose:        1.16, // Permanent Pool
	enums.ChaoIDStrRCPirateSpaceship: 1.16, // Permanent Pool
	enums.ChaoIDStrGoldenAngel:       1.16, // Permanent Pool
	enums.ChaoIDStrRCTornado2:           1.16, // Event (Sonic Adventure; Increase Odds During Event) - Permanent Pool
	enums.ChaoIDStrChaos:                1.16, // Event (Sonic Adventure; Increase Odds During Event) - Permanent Pool
	enums.ChaoIDStrOrca:                 0.21, // Event (Sonic Adventure; Increase Odds During Event) - Permanent Pool
	//enums.ChaoIDStrChaoWalker:           0.0, // Daily Battle
	// enums.ChaoIDStrDarkQueen:            3.67, // Runners' League Timed Mode
	enums.ChaoIDStrRCBlockFace: 1.16, // Permanent Pool
	//enums.ChaoIDStrDFekt:                0.0, // Revival Event (assets TBD)
	enums.ChaoIDStrDarkChaoWalker:       1.16, // Permanent Pool
	enums.ChaoIDStrPrideChaoL:           3.4, // Revival Event (Pride Month Celebration)
	enums.ChaoIDStrPrideChaoG:           3.4, // Revival Event (Pride Month Celebration)
	enums.ChaoIDStrPrideChaoB:           3.4, // Revival Event (Pride Month Celebration)
	enums.ChaoIDStrPrideChaoT:           3.4, // Revival Event (Pride Month Celebration)
	enums.ChaoIDStrPrideChaoP:           3.4, // Revival Event (Pride Month Celebration)
	enums.ChaoIDStrPrideChaoA:           3.4, // Revival Event (Pride Month Celebration)
	enums.ChaoIDStrPrideChaoNB:          3.4, // Revival Event (Pride Month Celebration)
}

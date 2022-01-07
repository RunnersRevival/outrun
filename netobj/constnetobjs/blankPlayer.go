package constnetobjs

import (
	"math/rand"
	"strconv"

	"github.com/Mtbcooler/outrun/obj"

	"github.com/Mtbcooler/outrun/config/eventconf"
	"github.com/Mtbcooler/outrun/consts"
	"github.com/Mtbcooler/outrun/enums"
	"github.com/Mtbcooler/outrun/netobj"
)

var BlankPlayer = func() netobj.Player {
	randChar := func(charset string, length int64) string {
		runes := []rune(charset)
		final := make([]rune, length)
		for i := range final {
			final[i] = runes[rand.Intn(len(runes))]
		}
		return string(final)
	}
	// create ID
	uid := ""
	for i := range make([]byte, 10) {
		if i == 0 { // if first character
			uid += strconv.Itoa(rand.Intn(9) + 1)
		} else {
			uid += strconv.Itoa(rand.Intn(10))
		}
	}
	username := ""
	password := randChar("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", 10)
	migrationPassword := randChar("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", 12)
	userPassword := ""
	key := randChar("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", 12)
	playerState := netobj.DefaultPlayerState()
	characterState := netobj.DefaultCharacterState()
	chaoState := GetAllNetChaoList() // needed for chaoRouletteAllowed
	mileageMapState := netobj.DefaultMileageMapState()
	mileageFriends := []netobj.MileageFriend{}
	playerVarious := netobj.DefaultPlayerVarious()
	optionUserResult := netobj.DefaultOptionUserResult()
	rouletteInfo := netobj.DefaultRouletteInfo()
	wheelOptions := netobj.DefaultWheelOptions(playerState.NumRouletteTicket, rouletteInfo.RouletteCountInPeriod, enums.WheelRankNormal, consts.RouletteFreeSpins, 0, consts.RouletteStartingJackpotRings)
	// TODO: get rid of logic here?
	allowedCharacters := []string{}
	allowedChao := []string{}
	for _, chao := range chaoState {
		if chao.Level < 10 { // not max level
			allowedChao = append(allowedChao, chao.ID)
		}
	}
	for _, character := range characterState {
		if character.Star < 10 { // not max star
			allowedCharacters = append(allowedCharacters, character.ID)
		}
	}
	chaoRouletteGroup := netobj.DefaultChaoRouletteGroup(playerState, allowedCharacters, allowedChao)
	personalEvents := []eventconf.ConfiguredEvent{}
	messages := []obj.Message{}
	operatorMessages := []obj.OperatorMessage{}
	loginBonusState := netobj.DefaultLoginBonusState(0)
	language := int64(enums.LangEnglish)
	eventState := netobj.DefaultEventState()
	eventUserRaidbossState := netobj.DefaultUserRaidbossState()
	suspendedUntil := int64(0)
	suspendReason := int64(0)
	lastLoginDevice := ""
	lastLoginPlatform := int64(0)
	lastLoginVersionId := int64(0)
	return netobj.NewPlayer(
		uid,
		username,
		password,
		migrationPassword,
		userPassword,
		key,
		language,
		playerState,
		characterState,
		chaoState,
		mileageMapState,
		mileageFriends,
		playerVarious,
		optionUserResult,
		wheelOptions,
		rouletteInfo,
		chaoRouletteGroup,
		personalEvents,
		messages,
		operatorMessages,
		loginBonusState,
		eventState,
		eventUserRaidbossState,
		suspendedUntil,
		suspendReason,
		lastLoginDevice,
		lastLoginPlatform,
		lastLoginVersionId,
	)
}() // TODO: Solve duplication requirement with db/assistants.go

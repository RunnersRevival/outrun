package constobjs

import (
	"strconv"

	"github.com/RunnersRevival/outrun/enums"
	"github.com/RunnersRevival/outrun/obj"
)

/*
All values are placeholders unless otherwise marked (Ex.: Sonic).
This should be changed when real values are found, or if we decide
that having custom values would be better for the balance of the game.

Multiple fields also have no currently known purposes, so these fields
are replaced with numbers that should be very easy to spot as 'abnormal'
in gameplay, thus giving credence to the idea that these values are
being actively used in gameplay. They may also have underlying issues,
which can be detected through a logcat reading.
*/

const NumRedRings = 1337
const PriceRedRings = 9001

// TODO: replace strconv.Itoa conversions to their string equivalents in enums. This should be done after #10 is solved and closed!

var CharacterSonic = obj.Character{
	strconv.Itoa(enums.CharaTypeSonic),
	0,           // unlocked from the start, no cost
	NumRedRings, // ?
	100000,       // used for limit breaking
	50,          // red rings used for limit breaking
}

var CharacterTails = obj.Character{
	strconv.Itoa(enums.CharaTypeTails),
	350,
	NumRedRings,
	100000, // used for limit breaking
	50,    // red rings used for limit breaking
}

var CharacterKnuckles = obj.Character{
	strconv.Itoa(enums.CharaTypeKnuckles),
	350,
	NumRedRings,
	100000, // used for limit breaking
	50,    // red rings used for limit breaking
}

var CharacterAmy = obj.Character{
	strconv.Itoa(enums.CharaTypeAmy),
	400,
	NumRedRings,
	100000, // used for limit breaking
	50,    // red rings used for limit breaking
}

var CharacterShadow = obj.Character{
	strconv.Itoa(enums.CharaTypeShadow),
	500,
	NumRedRings,
	500000, // used for limit breaking
	200,    // red rings used for limit breaking
}

var CharacterBlaze = obj.Character{
	strconv.Itoa(enums.CharaTypeBlaze),
	550,
	NumRedRings,
	500000, // used for limit breaking
	200,    // red rings used for limit breaking
}

var CharacterRouge = obj.Character{
	strconv.Itoa(enums.CharaTypeRouge),
	550,
	NumRedRings,
	500000, // used for limit breaking
	200,    // red rings used for limit breaking
}

var CharacterOmega = obj.Character{
	strconv.Itoa(enums.CharaTypeOmega),
	650,
	NumRedRings,
	500000, // used for limit breaking
	200,    // red rings used for limit breaking
}

var CharacterBig = obj.Character{
	strconv.Itoa(enums.CharaTypeBig),
	700,
	NumRedRings,
	10000, // used for limit breaking
	5,    // red rings used for limit breaking
}

var CharacterCream = obj.Character{
	strconv.Itoa(enums.CharaTypeCream),
	750,
	NumRedRings,
	10000, // used for limit breaking
	5,    // red rings used for limit breaking
}
var CharacterEspio = obj.Character{
	strconv.Itoa(enums.CharaTypeEspio),
	650,
	NumRedRings,
	500000, // used for limit breaking
	200,    // red rings used for limit breaking
}

var CharacterCharmy = obj.Character{
	strconv.Itoa(enums.CharaTypeCharmy),
	650,
	NumRedRings,
	500000, // used for limit breaking
	200,    // red rings used for limit breaking
}

var CharacterVector = obj.Character{
	strconv.Itoa(enums.CharaTypeVector),
	700,
	NumRedRings,
	500000, // used for limit breaking
	200,    // red rings used for limit breaking
}

var CharacterSilver = obj.Character{
	strconv.Itoa(enums.CharaTypeSilver),
	800,
	NumRedRings,
	10000, // used for limit breaking
	5,    // red rings used for limit breaking
}

var CharacterMetalSonic = obj.Character{
	strconv.Itoa(enums.CharaTypeMetalSonic),
	900,
	NumRedRings,
	500000, // used for limit breaking
	200,    // red rings used for limit breaking
}

var CharacterAmitieAmy = obj.Character{
	strconv.Itoa(enums.CharaTypeAmitieAmy),
	32000,
	NumRedRings,
	1200000, // used for limit breaking
	500,    // red rings used for limit breaking
}

var CharacterClassicSonic = obj.Character{
	strconv.Itoa(enums.CharaTypeClassicSonic),
	3000,
	NumRedRings,
	500000, // used for limit breaking
	200,    // red rings used for limit breaking
}

var CharacterTikal = obj.Character{
	strconv.Itoa(enums.CharaTypeTikal),
	1300,
	NumRedRings,
	500000, // used for limit breaking
	200,    // red rings used for limit breaking
}

var CharacterGothicAmy = obj.Character{
	strconv.Itoa(enums.CharaTypeGothicAmy),
	50000,
	NumRedRings,
	1200000, // used for limit breaking
	500,    // red rings used for limit breaking
}

var CharacterHalloweenShadow = obj.Character{
	strconv.Itoa(enums.CharaTypeHalloweenShadow),
	35000,
	NumRedRings,
	1200000, // used for limit breaking
	500,    // red rings used for limit breaking
}

var CharacterHalloweenRouge = obj.Character{
	strconv.Itoa(enums.CharaTypeHalloweenRouge),
	35000,
	NumRedRings,
	1200000, // used for limit breaking
	500,    // red rings used for limit breaking
}

var CharacterHalloweenOmega = obj.Character{
	strconv.Itoa(enums.CharaTypeHalloweenOmega),
	35000,
	NumRedRings,
	1200000, // used for limit breaking
	500,    // red rings used for limit breaking
}

var CharacterMephiles = obj.Character{
	strconv.Itoa(enums.CharaTypeMephiles),
	1550,
	NumRedRings,
	500000, // used for limit breaking
	200,    // red rings used for limit breaking
}

var CharacterPSISilver = obj.Character{
	strconv.Itoa(enums.CharaTypePSISilver),
	2300,
	NumRedRings,
	500000, // used for limit breaking
	200,    // red rings used for limit breaking
}

var CharacterXMasSonic = obj.Character{
	strconv.Itoa(enums.CharaTypeXMasSonic),
	30000,
	NumRedRings,
	1200000, // used for limit breaking
	500,    // red rings used for limit breaking
}

var CharacterXMasTails = obj.Character{
	strconv.Itoa(enums.CharaTypeXMasTails),
	30000,
	NumRedRings,
	1200000, // used for limit breaking
	500,    // red rings used for limit breaking
}

var CharacterXMasKnuckles = obj.Character{
	strconv.Itoa(enums.CharaTypeXMasKnuckles),
	30000,
	NumRedRings,
	1200000, // used for limit breaking
	500,    // red rings used for limit breaking
}

var CharacterWerehog = obj.Character{
	strconv.Itoa(enums.CharaTypeWerehog),
	800,
	NumRedRings,
	500000, // used for limit breaking
	200,    // red rings used for limit breaking
}

var CharacterSticks = obj.Character{
	strconv.Itoa(enums.CharaTypeSticks),
	750,
	NumRedRings,
	500000, // used for limit breaking
	200,    // red rings used for limit breaking
}

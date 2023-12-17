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
	NumRedRings, // unused? characters can only be unlocked and leveled up thru rings
	100000,      // used for limit smashing
	50,          // red rings used for limit smashing
}

var CharacterTails = obj.Character{
	strconv.Itoa(enums.CharaTypeTails),
	250,
	NumRedRings,
	100000, // used for limit smashing
	50,     // red rings used for limit smashing
}

var CharacterKnuckles = obj.Character{
	strconv.Itoa(enums.CharaTypeKnuckles),
	250,
	NumRedRings,
	100000, // used for limit smashing
	50,     // red rings used for limit smashing
}

var CharacterAmy = obj.Character{
	strconv.Itoa(enums.CharaTypeAmy),
	300,
	NumRedRings,
	100000, // used for limit smashing
	50,     // red rings used for limit smashing
}

var CharacterShadow = obj.Character{
	strconv.Itoa(enums.CharaTypeShadow),
	450,
	NumRedRings,
	500000, // used for limit smashing
	200,    // red rings used for limit smashing
}

var CharacterBlaze = obj.Character{
	strconv.Itoa(enums.CharaTypeBlaze),
	450,
	NumRedRings,
	500000, // used for limit smashing
	200,    // red rings used for limit smashing
}

var CharacterRouge = obj.Character{
	strconv.Itoa(enums.CharaTypeRouge),
	450,
	NumRedRings,
	500000, // used for limit smashing
	200,    // red rings used for limit smashing
}

var CharacterOmega = obj.Character{
	strconv.Itoa(enums.CharaTypeOmega),
	600,
	NumRedRings,
	500000, // used for limit smashing
	200,    // red rings used for limit smashing
}

var CharacterBig = obj.Character{
	strconv.Itoa(enums.CharaTypeBig),
	450,
	NumRedRings,
	500000, // used for limit smashing
	200,    // red rings used for limit smashing
}

var CharacterCream = obj.Character{
	strconv.Itoa(enums.CharaTypeCream),
	450,
	NumRedRings,
	500000, // used for limit smashing
	200,    // red rings used for limit smashing
}

var CharacterEspio = obj.Character{
	strconv.Itoa(enums.CharaTypeEspio),
	600,
	NumRedRings,
	500000, // used for limit smashing
	200,    // red rings used for limit smashing
}

var CharacterCharmy = obj.Character{
	strconv.Itoa(enums.CharaTypeCharmy),
	600,
	NumRedRings,
	500000, // used for limit smashing
	200,    // red rings used for limit smashing
}

var CharacterVector = obj.Character{
	strconv.Itoa(enums.CharaTypeVector),
	600,
	NumRedRings,
	500000, // used for limit smashing
	200,    // red rings used for limit smashing
}

var CharacterSilver = obj.Character{
	strconv.Itoa(enums.CharaTypeSilver),
	750,
	NumRedRings,
	500000, // used for limit smashing
	200,    // red rings used for limit smashing
}

var CharacterMetalSonic = obj.Character{
	strconv.Itoa(enums.CharaTypeMetalSonic),
	600,
	NumRedRings,
	500000, // used for limit smashing
	200,    // red rings used for limit smashing
}

var CharacterAmitieAmy = obj.Character{
	strconv.Itoa(enums.CharaTypeAmitieAmy),
	900,
	NumRedRings,
	0, // used for limit smashing
	0, // red rings used for limit smashing
}

var CharacterClassicSonic = obj.Character{
	strconv.Itoa(enums.CharaTypeClassicSonic),
	750,
	NumRedRings,
	500000, // used for limit smashing
	200,    // red rings used for limit smashing
}

var CharacterTikal = obj.Character{
	strconv.Itoa(enums.CharaTypeTikal),
	750,
	NumRedRings,
	500000, // used for limit smashing
	200,    // red rings used for limit smashing
}

var CharacterGothicAmy = obj.Character{
	strconv.Itoa(enums.CharaTypeGothicAmy),
	900,
	NumRedRings,
	0, // used for limit smashing
	0, // red rings used for limit smashing
}

var CharacterHalloweenShadow = obj.Character{
	strconv.Itoa(enums.CharaTypeHalloweenShadow),
	900,
	NumRedRings,
	0, // used for limit smashing
	0, // red rings used for limit smashing
}

var CharacterHalloweenRouge = obj.Character{
	strconv.Itoa(enums.CharaTypeHalloweenRouge),
	900,
	NumRedRings,
	0, // used for limit smashing
	0, // red rings used for limit smashing
}

var CharacterHalloweenOmega = obj.Character{
	strconv.Itoa(enums.CharaTypeHalloweenOmega),
	900,
	NumRedRings,
	0, // used for limit smashing
	0, // red rings used for limit smashing
}

var CharacterMephiles = obj.Character{
	strconv.Itoa(enums.CharaTypeMephiles),
	750,
	NumRedRings,
	500000, // used for limit smashing
	200,    // red rings used for limit smashing
}

var CharacterPSISilver = obj.Character{
	strconv.Itoa(enums.CharaTypePSISilver),
	750,
	NumRedRings,
	500000, // used for limit smashing
	200,    // red rings used for limit smashing
}

var CharacterXMasSonic = obj.Character{
	enums.CTStrXMasSonic,
	900,
	NumRedRings,
	0, // used for limit smashing
	0, // red rings used for limit smashing
}

var CharacterXMasTails = obj.Character{
	enums.CTStrXMasTails,
	900,
	NumRedRings,
	0, // used for limit smashing
	0, // red rings used for limit smashing
}

var CharacterXMasKnuckles = obj.Character{
	enums.CTStrXMasKnuckles,
	900,
	NumRedRings,
	0, // used for limit smashing
	0, // red rings used for limit smashing
}

var CharacterWerehog = obj.Character{
	strconv.Itoa(enums.CharaTypeWerehog),
	750,
	NumRedRings,
	500000, // used for limit smashing
	200,    // red rings used for limit smashing
}

var CharacterSticks = obj.Character{
	strconv.Itoa(enums.CharaTypeSticks),
	750,
	NumRedRings,
	500000, // used for limit smashing
	200,    // red rings used for limit smashing
}

var CharacterMarine = obj.Character{
	strconv.Itoa(enums.CharaTypeMarine),
	1200,
	NumRedRings,
	1200000, // used for limit smashing
	200,     // red rings used for limit smashing
}

var CharacterWhisper = obj.Character{
	strconv.Itoa(enums.CharaTypeWhisper),
	1200,
	NumRedRings,
	0, // used for limit smashing
	0, // red rings used for limit smashing
}

var CharacterTangle = obj.Character{
	strconv.Itoa(enums.CharaTypeTangle),
	1200,
	NumRedRings,
	0, // used for limit smashing
	0, // red rings used for limit smashing
}

package bacot

var (
	advanceLeetSpeaks = map[rune]struct{}{
		'-': {},
		'_': {},
		'|': {},
		'.': {},
		',': {},
		'(': {},
		')': {},
		'>': {},
		'"': {},
		'`': {},
		'~': {},
		'*': {},
		'&': {},
		'%': {},
		'?': {},
	}

	simpleLeetSpeaks = map[rune]rune{
		'4': 'a',
		'@': 'a',
		'8': 'b',
		'<': 'c',
		'(': 'c',
		'[': 'c',
		'©': 'c',
		'3': 'e',
		'€': 'e',
		'6': 'g',
		'#': 'h',
		'!': 'i',
		'|': 'i',
		'1': 'i',
		'0': 'o',
		'2': 'r',
		'®': 'r',
		'5': 's',
		'$': 's',
		'§': 's',
		'7': 't',
		'†': 't',
		'µ': 'u',
		'×': 'x',
		'¥': 'y',
	}

	whiteSpace = map[rune]struct{}{
		9:  {}, // Horizontal tab
		10: {}, // Line feed
		11: {}, // Vertical Tabulation
		12: {}, // Form feed
		13: {}, // Carriage return
		32: {}, // Space
	}

	replaceToSpace = map[rune]struct{}{
		'-': {},
		'_': {},
		'~': {},
		'|': {},
		'=': {},
		'+': {},
		'.': {},
		',': {},
	}

	replaceToClear = map[rune]struct{}{
		'\'': {},
		'"':  {},
		'`':  {},
		'^':  {},
		'*':  {},
	}
)

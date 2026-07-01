package bacot

var whiteSpace = []rune{
	'-',
	'_',
	'|',
	'.',
	',',
	'(',
	')',
	'>',
	'"',
	'`',
	'~',
	'*',
	'&',
	'%',
	'?',
}

type charReplacement struct {
	WhiteSpace map[rune]struct{}
}

func charReplacemnetCreate() *charReplacement {

	ws := make(map[rune]struct{}, 100)
	for _, c := range whiteSpace {
		ws[c] = struct{}{}
	}

	return &charReplacement{
		WhiteSpace: ws,
	}
}

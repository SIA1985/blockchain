package script

type token uint8
type trimmedToToken map[string]token

const (
	/*данные*/
	DATA token = iota

	/*операции*/
	HASH256
	IS_EQUAL
	CHECK_SIG
)

func IsOperationalToken(tk token) bool {
	return HASH256 <= tk && tk <= CHECK_SIG
}

type Tokenizer struct {
	stringToToken trimmedToToken
}

func NewTokenizer() *Tokenizer {
	stringToToken := make(trimmedToToken)

	stringToToken["hash256"] = HASH256
	stringToToken["is_equal"] = IS_EQUAL
	stringToToken["check_sig"] = CHECK_SIG

	return &Tokenizer{stringToToken}
}

func (t Tokenizer) Tokenize(trimmed string) token {
	if tk, ok := t.stringToToken[trimmed]; ok {
		return tk
	} else {
		return DATA
	}
}

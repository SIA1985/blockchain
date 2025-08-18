package script

import "regexp"

var exec map[token]operation = map[token]operation{
	HASH256:   Hash256,
	IS_EQUAL:  IsEqual,
	CHECK_SIG: CheckSig,
}

func Run(text string) (ret bool, err error) {
	ret = false
	err = nil

	/*препроцесс*/
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	trimmeds := reg.Split(text, -1)

	/*токенизация и исполнение*/
	t := NewTokenizer()
	stack := NewStack()
	for _, trimmed := range trimmeds {
		tk := t.Tokenize(trimmed)
		if !IsOperationalToken(tk) {
			stack.Push(trimmed)
			continue
		}

		err = exec[tk](stack)
		if err != nil {
			return
		}
	}

	/*проверка результата*/
	if stack.Size() != 1 {
		return
	}

	/*todo: определения true или false*/
	// ret = stack.Pop()
	return
}

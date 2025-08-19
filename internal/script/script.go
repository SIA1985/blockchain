package script

import (
	"fmt"
	"regexp"
	"strconv"
)

var exec map[token]operation = map[token]operation{
	HASH256:   Hash256,
	IS_EQUAL:  IsEqual,
	CHECK_SIG: CheckSig,
}

func Run(text string, witness []string) (ret bool, err error) {
	ret = false
	err = nil

	/*препроцесс*/
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	trimmeds := reg.Split(text, -1)

	/*подготовка стэка*/
	stack := NewStack()
	for _, data := range witness {
		stack.Push(data)
	}

	/*интерпретация*/
	t := NewTokenizer()
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
		err = fmt.Errorf("stack size != 1 at the end")
		return
	}

	var result int
	result, err = strconv.Atoi(stack.Pop())
	if err != nil {
		return
	}

	ret = (result != 0)
	return
}

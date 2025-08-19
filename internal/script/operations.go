package script

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type operation func(*Stack) error

var Hash256 operation = func(s *Stack) (err error) {
	if s.Size() < 1 {
		err = fmt.Errorf("not enought args for Hash256")
		return
	}

	var toHash []byte
	toHash, err = hex.DecodeString(s.Pop())
	if err != nil {
		return
	}

	hash := sha256.Sum256(toHash)

	s.Push(hex.EncodeToString(hash[:]))

	return
}

var IsEqual operation = func(s *Stack) (err error) {
	if s.Size() < 2 {
		err = fmt.Errorf("not enought args for IsEqual")
		return
	}

	if s.Pop() == s.Pop() {
		s.Push("1")
	} else {
		s.Push("0")
	}

	return
}

var CheckSig operation = func(s *Stack) (err error) {
	if s.Size() < 3 {
		err = fmt.Errorf("not enought args for CheckSig")
		return
	}

	return
}

package script

import (
	"blockchain/internal/algorythms"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type operation func(*Stack) error

type DataType uint8

const (
	Unknown = iota
	Hex
	Dec
	Bin
	Str

	HexPref = "0x"
	BinPref = "0b"
	PrefLen = 2
)

func dataToType(data string) (dataType DataType) {
	dataType = Unknown

	if len(data) > PrefLen {
		switch data[:PrefLen] {
		case HexPref:
			dataType = Hex
		case BinPref:
			dataType = Bin
		}

		data = data[PrefLen+1:]
	}

	if dataType == Unknown {
		var reg = regexp.MustCompile(`^[0-9]+$`)

		if reg.MatchString(data) {
			dataType = Dec
		} else {
			dataType = Str
		}
	}

	return
}

var Hash256 operation = func(s *Stack) (err error) {
	if s.Size() < 1 {
		err = fmt.Errorf("not enought args for Hash256")
		return
	}

	var hash [32]byte
	toHash := make([]byte, 0)

	data := s.Pop()
	switch dataToType(data) {
	case Hex:
		//todo: bool 2
		data, _ = strings.CutPrefix(data, HexPref)

		var hex int64
		hex, err = strconv.ParseInt(data, 16, 64)
		if err != nil {
			return
		}

		toHash = algorythms.Int64ToByteArr(hex)
	case Dec:
		var dec int64
		dec, err = strconv.ParseInt(data, 10, 64)
		if err != nil {
			return
		}

		toHash = algorythms.Int64ToByteArr(dec)
	case Bin:
		//todo: bool 2
		data, _ = strings.CutPrefix(data, BinPref)

		var bin int64
		bin, err = strconv.ParseInt(data, 2, 64)
		if err != nil {
			return
		}

		toHash = algorythms.Int64ToByteArr(bin)
	case Str:
		toHash, err = hex.DecodeString(data)
		if err != nil {
			return
		}
	}

	hash = sha256.Sum256(toHash)

	s.Push(HexPref + hex.EncodeToString(hash[:]))

	return
}

var IsEqual operation = func(s *Stack) (err error) {
	if s.Size() < 2 {
		err = fmt.Errorf("not enought args for IsEqual")
		return
	}

	a := s.Pop()
	b := s.Pop()
	if a == b {
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

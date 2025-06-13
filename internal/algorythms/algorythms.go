package algorythms

import (
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

func Int64ToByteArr(v int64) []byte {
	return []byte{
		byte(0xff & v),
		byte(0xff & (v >> 8)),
		byte(0xff & (v >> 16)),
		byte(0xff & (v >> 24)),
		byte(0xff & (v >> 32)),
		byte(0xff & (v >> 40)),
		byte(0xff & (v >> 48)),
		byte(0xff & (v >> 56)),
	}
}

func UInt64ToByteArr(v uint64) []byte {
	return []byte{
		byte(0xff & v),
		byte(0xff & (v >> 8)),
		byte(0xff & (v >> 16)),
		byte(0xff & (v >> 24)),
		byte(0xff & (v >> 32)),
		byte(0xff & (v >> 40)),
		byte(0xff & (v >> 48)),
		byte(0xff & (v >> 56)),
	}
}

func ProofOfWork(blockByteArr []byte) (blockHash []byte, nonce int64, targetBits uint64) {
	nonce = 0
	targetBits = 24 /*todo*/

	var hash [32]byte
	var hashInt big.Int
	var toCompare *big.Int = big.NewInt(1) /*todo*/
	toCompare.Lsh(toCompare, uint(256-targetBits))

	/*todo: распараллелить с каналом или вэйтгруппой*/
	fmt.Println("Майнинг блока...")
	for nonce < math.MaxInt64 {
		data := append(blockByteArr, Int64ToByteArr(nonce)...)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(toCompare) == -1 {
			break
		} else {
			nonce++
		}
	}
	/*todo: nonce == math.MaxInt64 - что делать?*/

	blockHash = hash[:]

	fmt.Printf("%d, %x\n", nonce, hash)
	fmt.Println()

	return
}

/*todo: hashInt.Cmp(toCompare) == -1 - вынести в функцию*/

func Validate(blockByteArrWithNonce []byte, targetBits uint64) bool {
	var hashInt big.Int
	var toCompare *big.Int = big.NewInt(1)
	toCompare.Lsh(toCompare, uint(256-targetBits))

	hash := sha256.Sum256(blockByteArrWithNonce)
	hashInt.SetBytes(hash[:])

	return (hashInt.Cmp(toCompare) == -1)
}

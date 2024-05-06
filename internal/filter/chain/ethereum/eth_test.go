package ethereum

import (
	"math/big"
	"testing"
)

func Test_HexStr(t *testing.T) {
	i, ok := big.NewInt(0).SetString("00000000000000000000000000000000000000000000000000000000000058f8", 16)
	t.Log("num -------------------- ", i)
	t.Log("ok -------------------- ", ok)
}

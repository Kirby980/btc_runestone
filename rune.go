package btc_runestone

import (
	"math/big"
)

var char = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func ToRune(x *big.Int) string {
	var s string
	flag := false
	zero := big.NewInt(0)
	one := big.NewInt(1)
	t := new(big.Int)

	for x.Cmp(zero) != 0 {
		if flag {
			t.Sub(x, one)
			t.Mod(t, big.NewInt(26))
			x.Sub(x, one)
			x.Div(x, big.NewInt(26))
		} else {
			t.Mod(x, big.NewInt(26))
			x.Div(x, big.NewInt(26))
		}
		flag = true
		index := t.Int64()
		s = char[index] + s
	}
	return s
}
func ToBigint(s string) {
	t := big.NewInt(0)
	p := big.NewInt(1)

	base := big.NewInt(26)
	for i := len(s) - 1; i >= 0; i-- {
		temp := big.NewInt(int64(s[i] - 'A'))
		if i != len(s)-1 {
			temp.Add(temp, big.NewInt(1))
		}
		temp.Mul(temp, p)
		t.Add(t, temp)
		p.Mul(p, base)
	}
}

package btc_runestone

import (
	"github.com/btcsuite/btcd/txscript"
	"math/big"
)

type Rune big.Int

type Flag big.Int

const (
	OP_RETURN    = txscript.OP_RETURN
	MAGIC_NUMBER = txscript.OP_13
)

type Message struct {
	cenotaph bool
	edicts   []Edict
	fields   map[*big.Int]*big.Int
}

type Etching struct {
	divisibility *big.Int
	terms        *Term
	runes        *Rune
	spacers      *big.Int
	symbol       string
	premine      *big.Int
}
type Term struct {
	amount *big.Int
	cap    *big.Int
	Height StartAndEnd
	Offset StartAndEnd
}

type StartAndEnd struct {
	Start *big.Int
	End   *big.Int
}
type RuneStone struct {
	mint    *RuneId
	pointer uint32
	edicts  []Edict
	etching *Etching
}
type RuneId struct {
	block *big.Int
	tx    *big.Int
}
type Edict struct {
	id     RuneId
	amount *big.Int
	output *big.Int
}

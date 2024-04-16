package btc_runestone

import (
	"fmt"
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
	fields   map[*big.Int][]*big.Int
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

func (id *RuneId) String() string {
	return fmt.Sprintf("block:%v,tx:%v", id.block, id.tx)
}
func (e Edict) String() string {
	return fmt.Sprintf("EdictBigInt{id: %v, amount: %v, output: %v}", e.id.String(), e.amount, e.output)
}
func (m Message) String() string {
	return fmt.Sprintf("MessageBigInt{cenotaph: %t, edicts: %v, fields: %v}", m.cenotaph, m.edicts, m.fields)
}
func (e Etching) String() string {
	return fmt.Sprintf("EtchingBigInt{divisibility: %v, terms: %v, runes: %v, spacers: %v, symbol: %s, premine: %v}", e.divisibility, e.terms, e.runes, e.spacers, e.symbol, e.premine)
}
func (r RuneStone) String() string {
	return fmt.Sprintf("RuneStoneBigInt{mint: %v, \npointer: %v, \nedicts: %v ,\netching: %v}", r.mint.String(), r.pointer, r.edicts, r.etching.String())
}

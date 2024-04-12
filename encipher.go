package btc_runestone

import (
	"fmt"
	"github.com/btcsuite/btcd/txscript"
	"math/big"
	"sort"
)

func Encipher(runeStone RuneStone) *txscript.ScriptBuilder {
	var payload []byte
	if runeStone.etching != nil {
		flags := big.NewInt(0)
		Etch := new(Flag)
		Etch = (*Flag)(big.NewInt(0))
		Etch.Set(flags)
		if runeStone.etching.terms != nil {
			Terms := new(Flag)
			Terms = (*Flag)(big.NewInt(1))
			Terms.Set(flags)
			fmt.Println(flags)
		}
		payload = append(payload, Encode(big.NewInt(2))...)
		payload = append(payload, Encode(flags)...)
		if (*big.Int)(runeStone.etching.runes) != nil {
			payload = append(payload, Encode(big.NewInt(4))...)
			payload = append(payload, Encode((*big.Int)(runeStone.etching.runes))...)
		}
		if runeStone.etching.divisibility != nil {
			payload = append(payload, Encode(big.NewInt(1))...)
			payload = append(payload, Encode(runeStone.etching.divisibility)...)
		}
		if runeStone.etching.spacers != nil {
			payload = append(payload, Encode(big.NewInt(3))...)
			payload = append(payload, Encode(runeStone.etching.spacers)...)
		}
		if runeStone.etching.symbol != "" {
			payload = append(payload, Encode(big.NewInt(5))...)
			payload = append(payload, []byte(runeStone.etching.symbol)...)
		}
		if runeStone.etching.premine != nil {
			payload = append(payload, Encode(big.NewInt(6))...)
			payload = append(payload, Encode(runeStone.etching.premine)...)
		}
		if runeStone.etching.terms != nil {
			payload = append(payload, Encode(big.NewInt(10))...)
			payload = append(payload, Encode(runeStone.etching.terms.amount)...)
			payload = append(payload, Encode(big.NewInt(8))...)
			payload = append(payload, Encode(runeStone.etching.terms.cap)...)
			if runeStone.etching.terms.Height.Start != nil {
				payload = append(payload, Encode(big.NewInt(12))...)
				payload = append(payload, Encode(runeStone.etching.terms.Height.Start)...)
			}
			if runeStone.etching.terms.Height.End != nil {
				payload = append(payload, Encode(big.NewInt(14))...)
				payload = append(payload, Encode(runeStone.etching.terms.Height.End)...)
			}
			if runeStone.etching.terms.Offset.Start != nil {
				payload = append(payload, Encode(big.NewInt(16))...)
				payload = append(payload, Encode(runeStone.etching.terms.Offset.Start)...)
			}
			if runeStone.etching.terms.Offset.End != nil {
				payload = append(payload, Encode(big.NewInt(18))...)
				payload = append(payload, Encode(runeStone.etching.terms.Offset.End)...)
			}
		}
	}

	if runeStone.mint != nil {
		payload = append(payload, Encode(big.NewInt(20))...)
		payload = append(payload, Encode(runeStone.mint.block)...)
		payload = append(payload, Encode(big.NewInt(20))...)
		payload = append(payload, Encode(runeStone.mint.tx)...)
	}
	if runeStone.pointer != 0 {
		payload = append(payload, Encode(big.NewInt(22))...)
		payload = append(payload, Encode(new(big.Int).SetUint64(uint64(runeStone.pointer)))...)
	}
	if runeStone.edicts != nil {
		payload = append(payload, Encode(big.NewInt(0))...)
		edicts := runeStone.edicts
		sort.Slice(edicts, func(i, j int) bool {
			if edicts[i].id.block.Cmp(edicts[j].id.block) < 0 {
				return true
			}
			if edicts[i].id.block == edicts[j].id.block && edicts[i].id.block.Cmp(edicts[j].id.block) < 0 {
				return true
			}
			return false
		})

		var previous = RuneId{big.NewInt(0), big.NewInt(0)}
		for _, edict := range edicts {
			block, tx := previous.Delta(edict.id)
			payload = append(payload, Encode(block)...)
			payload = append(payload, Encode(tx)...)
			payload = append(payload, Encode(edict.amount)...)
			payload = append(payload, Encode(edict.output)...)
			previous = edict.id
		}
	}
	var builder = txscript.NewScriptBuilder().AddOp(txscript.OP_RETURN).AddOp(MAGIC_NUMBER)

	for len(payload) > 0 {
		chunkSize := txscript.MaxScriptElementSize
		if len(payload) < chunkSize {
			chunkSize = len(payload)
		}
		chunk := payload[:chunkSize]
		builder.AddData(chunk)
		payload = payload[chunkSize:]
	}
	return builder
}

func (r *RuneId) Delta(next RuneId) (block *big.Int, tx *big.Int) {
	block = next.block.Sub(next.block, r.block)

	if block.Cmp(big.NewInt(0)) == 0 {
		tx = next.tx.Sub(next.tx, r.tx)
	} else {
		tx = next.tx
	}
	return block, tx
}

func (f *Flag) Set(flags *big.Int) {
	flags.Or(flags, big.NewInt(1).Lsh(big.NewInt(1), uint((*big.Int)(f).Uint64())))
}
func Encode(n *big.Int) []byte {
	var result []byte
	for n.Cmp(big.NewInt(128)) > 0 {
		temp := new(big.Int).Set(n)
		last := temp.And(n, new(big.Int).SetUint64(0b0111_1111))
		result = append(result, last.Or(last, new(big.Int).SetUint64(0b1000_0000)).Bytes()...)
		n.Rsh(n, 7)
	}
	result = append(result, n.Bytes()...)
	return result
}

package btc_runestone

import (
	"fmt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/spf13/cast"
	"math/big"
)

func Decode(buffer []byte) (*big.Int, int) {
	n := big.NewInt(0)
	for i, tick := range buffer {
		if i > 18 {
			return new(big.Int), 0
		}
		value := uint64(tick) & 0b0111_1111
		if i == 18 && value&0b0111_1100 != 0 {
			return new(big.Int), 0
		}
		temp := new(big.Int).SetUint64(value)
		n.Or(n, temp.Lsh(temp, uint(7*i)))
		if tick&0b1000_0000 == 0 {
			return n, i + 1
		}
	}
	return new(big.Int), 0
}
func Integers(payload []byte) []*big.Int {
	integers := make([]*big.Int, 0)
	i := 0
	for i < len(payload) {
		uint64s, intLen := Decode(payload[i:])
		if intLen != 0 {
			integers = append(integers, uint64s)
			i += intLen
		}
	}
	return integers
}
func (e Edict) FromIntegers(tx wire.MsgTx, id *RuneId, amount, output *big.Int) *Edict {
	if id.block.Cmp(big.NewInt(0)) == 0 && id.tx.Cmp(big.NewInt(0)) > 0 {
		return nil
	}
	if output.Cmp(big.NewInt(int64(len(tx.TxOut)))) > 0 {
		return nil
	}
	fmt.Println(id)
	return &Edict{
		id:     *id,
		amount: amount,
		output: output,
	}
}

func Decipher(transaction wire.MsgTx) *RuneStone {
	payload := Payload(transaction)
	if payload == nil {
		return nil
	}
	integers := Integers(payload)
	message := FromIntegers(transaction, integers)
	//claim := takeTag(&message.fields, 14)
	_ = TakeTag(&message.fields, big.NewInt(TagAmount))
	_ = TakeTag(&message.fields, big.NewInt(TagPremine))

	_ = TakeTag(&message.fields, big.NewInt(TagCap))
	mint := TakeTag(&message.fields, big.NewInt(TagMint))
	mintID := &RuneId{}
	if len(mint) > 0 {
		mintID = &RuneId{
			block: mint[0],
			tx:    mint[1],
		}
	}
	flags := TakeTag(&message.fields, big.NewInt(TagFlags))
	spacers := TakeTag(&message.fields, big.NewInt(TagSpacers))
	runes := TakeTag(&message.fields, big.NewInt(TagRune))
	symbol := TakeTag(&message.fields, big.NewInt(TagSymbol))
	divisibility := TakeTag(&message.fields, big.NewInt(TagDivisibility))
	var etching Etching
	if len(flags) != 0 {
		EtchBigInt := new(Flag)
		EtchBigInt = (*Flag)(big.NewInt(FlagEtching))
		etch := EtchBigInt.TakeFlag(flags[0])
		TermsBigint := new(Flag)
		TermsBigint = (*Flag)(big.NewInt(FlagTerms))
		_ = TermsBigint.TakeFlag(flags[0])
		if !etch {
			etching = Etching{
				symbol: cast.ToString(symbol),
			}
		}
		if len(runes) != 0 {
			etching.runes = (*Rune)(runes[0])
		}
		if len(divisibility) != 0 {
			etching.divisibility = divisibility[0]
		}
		if len(spacers) != 0 {
			etching.spacers = spacers[0]
		}
	}
	for _, value := range message.fields {
		if value[0].Mod(value[0], big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
			message.cenotaph = true
			break
		}
	}
	/*c := RuneId{}
	if claim < uint64(len(message.edicts)) {
		c = message.edicts[claim].id
	}*/
	return &RuneStone{
		mint:    mintID,
		edicts:  message.edicts,
		etching: &etching,
	}
}

func Payload(transaction wire.MsgTx) []byte {
	for _, out := range transaction.TxOut {
		instructions := out.PkScript
		if OP_RETURN != instructions[0] {
			continue
		}
		if txscript.OP_13 != instructions[1] {
			continue
		}
		if len(instructions) < 3 {
			continue
		}
		var payload []byte
		for _, instruction := range instructions[3:] {
			payload = append(payload, instruction)
		}
		return payload
	}
	return nil
}
func FromIntegers(tx wire.MsgTx, payload []*big.Int) Message {
	var edicts []Edict
	fields := make(map[*big.Int][]*big.Int)
	cenotaph := false
	for i := 0; i < len(payload); i += 2 {
		tag := payload[i]
		if tag.Cmp(big.NewInt(TagBody)) == 0 {
			id := RuneId{
				block: big.NewInt(0),
				tx:    big.NewInt(0),
			}
			for j := i + 1; j < len(payload); j += 4 {
				if j+3 > len(payload) {
					cenotaph = true
					break
				}
				next, err := id.nextBigInt(payload[j], payload[j+1])
				if err != nil {
					cenotaph = true
					break
				}
				edict := &Edict{}
				if e := edict.FromIntegers(tx, next, payload[j+2], payload[j+3]); e != nil {
					edicts = append(edicts, *e)
				} else {
					cenotaph = true
				}
			}
			break
		}

		if i+1 >= len(payload) {
			break
		}
		fields[tag] = append(fields[tag], payload[i+1])
	}
	return Message{
		cenotaph: cenotaph,
		edicts:   edicts,
		fields:   fields,
	}
}

func (id *RuneId) nextBigInt(block *big.Int, tx *big.Int) (*RuneId, error) {
	newBlock := id.block.Add(id.block, block)

	var newTx *big.Int
	if block.Cmp(big.NewInt(0)) == 0 {
		newTx = id.tx.Add(id.tx, tx)
	} else {
		newTx = tx
	}

	return &RuneId{block: newBlock, tx: newTx}, nil
}

func TakeTag(fields *map[*big.Int][]*big.Int, tag *big.Int) []*big.Int {
	var result []*big.Int
	for key, value := range *fields {
		if key.Cmp(tag) == 0 {
			result = append(result, value...)
			delete(*fields, key)
		}
	}
	return result
}
func (f *Flag) TakeFlag(flags *big.Int) bool {
	set := flags.And(flags, (*big.Int)(f)).Cmp(big.NewInt(0)) != 0

	not := new(big.Int).Not((*big.Int)(f))
	flags.And(flags, not)
	return set
}

func (id *RuneId) Next(block *big.Int, tx *big.Int) (*RuneId, error) {
	newBlock := id.block.Add(id.block, block)

	var newTx *big.Int
	if block.Cmp(big.NewInt(0)) == 0 {
		newTx = id.tx.Add(id.tx, tx)
	} else {
		newTx = tx
	}

	return &RuneId{block: newBlock, tx: newTx}, nil
}

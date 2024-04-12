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
	fmt.Println(message)
	//claim := takeTag(&message.fields, 14)
	_ = TakeTag(&message.fields, big.NewInt(10))
	divisibility := TakeTag(&message.fields, big.NewInt(1))
	_ = TakeTag(&message.fields, big.NewInt(6))
	runes := TakeTag(&message.fields, big.NewInt(4))
	spacers := TakeTag(&message.fields, big.NewInt(3))
	symbol := TakeTag(&message.fields, big.NewInt(5))
	_ = TakeTag(&message.fields, big.NewInt(8))
	flags := TakeTag(&message.fields, big.NewInt(2))
	Etch := new(Flag)
	Etch = (*Flag)(big.NewInt(0))
	etch := Etch.TakeFlag(flags)
	Terms := new(Flag)
	Terms = (*Flag)(big.NewInt(1))
	_ = Terms.TakeFlag(flags)
	var etching Etching
	if !etch {
		etching = Etching{
			divisibility: divisibility,
			runes:        (*Rune)(runes),
			spacers:      spacers,
			symbol:       cast.ToString(symbol),
		}
	}

	for _, value := range message.fields {
		if value.Mod(value, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
			message.cenotaph = true
			break
		}
	}
	/*c := RuneId{}
	if claim < uint64(len(message.edicts)) {
		c = message.edicts[claim].id
	}*/
	return &RuneStone{
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
	fields := make(map[*big.Int]*big.Int)
	cenotaph := false
	for i := 0; i < len(payload); i += 2 {
		tag := payload[i]

		if tag.Cmp(big.NewInt(0)) == 0 {
			id := RuneId{
				block: big.NewInt(0),
				tx:    big.NewInt(0),
			}
			for j := i + 1; j < len(payload); j += 4 {
				if j+3 >= len(payload) {
					cenotaph = true
					break
				}
				next, err := id.Next(payload[j], payload[j+1])
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

		fields[tag] = payload[i+1]
	}
	fmt.Println(cenotaph)
	fmt.Println(edicts)
	fmt.Println(fields)

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

func TakeTag(fields *map[*big.Int]*big.Int, tag *big.Int) *big.Int {
	value, ok := (*fields)[tag]
	if !ok {
		return big.NewInt(0)
	}
	delete(*fields, tag)
	return value
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

package btc_runestone

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"math/big"
	"testing"
)

func TestTest(t *testing.T) {

	runeStone := RuneStone{
		mint: &RuneId{
			block: new(big.Int).SetUint64(2586500),
			tx:    new(big.Int).SetUint64(2279),
		},
		edicts: []Edict{
			{
				id: RuneId{
					block: new(big.Int).SetUint64(2586500),
					tx:    new(big.Int).SetUint64(2279),
				},
				amount: new(big.Int).SetUint64(100),
				output: new(big.Int).SetUint64(0),
			},
			{
				id: RuneId{
					block: new(big.Int).SetUint64(2586500),
					tx:    new(big.Int).SetUint64(2279),
				},
				amount: new(big.Int).SetUint64(100),
				output: new(big.Int).SetUint64(1),
			},
		},
	}
	builder := encipherBigInt(runeStone)
	bytes, _ := builder.Script()
	fmt.Println(bytes)
	decodeString, _ := hex.DecodeString("00eae29d0111e80702")
	fmt.Println(decodeString)
	script, _ := txscript.NewScriptBuilder().AddOp(OP_RETURN).AddOp(txscript.OP_13).AddData(bytes[3:]).Script()
	//script, _ := txscript.NewScriptBuilder().AddOp(OP_RETURN).AddOp(txscript.OP_13).AddData(decodeString).Script()
	tx := wire.MsgTx{
		TxOut: []*wire.TxOut{
			{
				PkScript: script,
				Value:    0,
			},
			{
				PkScript: []byte{},
				Value:    1000,
			},
		},
		LockTime: 0,
		Version:  2,
	}
	stoneBigint := Decipher(tx)
	fmt.Printf("%v", stoneBigint)
}

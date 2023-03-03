package ethgo

import (
	"fmt"
	"math/big"

	"github.com/umbracle/fastrlp"
)

// GetHash returns the Hash of the transaction
func (t *Transaction) GetHash() (hash Hash, err error) {
	var rlpEncode []byte
	if rlpEncode, err = t.MarshalRLPTo(nil); err != nil {
		return Hash{}, err
	}
	return BytesToHash(Keccak256(rlpEncode)), nil
}

// MarshalRLPTo marshals the transaction to a []byte destination
func (t *Transaction) MarshalRLPTo(dst []byte) ([]byte, error) {
	raw, err := fastrlp.MarshalRLP(t)
	if err != nil {
		return nil, err
	}
	if t.Type == TransactionLegacy {
		return raw, nil
	}
	// append type byte
	return append([]byte{byte(t.Type)}, raw...), nil
}

// MarshalRLPWith marshals the transaction to RLP with a specific fastrlp.Arena
func (t *Transaction) MarshalRLPWith(arena *fastrlp.Arena) (*fastrlp.Value, error) {
	vv := arena.NewArray()

	if t.Type != 0 {
		// either dynamic and access type
		vv.Set(arena.NewBigInt(t.ChainID))
	}

	vv.Set(arena.NewUint(t.Nonce))

	if t.Type == TransactionDynamicFee {
		// dynamic fee uses
		vv.Set(arena.NewBigInt(t.MaxPriorityFeePerGas))
		vv.Set(arena.NewBigInt(t.MaxFeePerGas))
	} else {
		// legacy and access type use gas price
		vv.Set(arena.NewUint(t.GasPrice))
	}

	vv.Set(arena.NewUint(t.Gas))

	// Address may be empty
	if t.To != nil {
		vv.Set(arena.NewBytes((*t.To)[:]))
	} else {
		vv.Set(arena.NewNull())
	}

	vv.Set(arena.NewBigInt(t.Value))
	vv.Set(arena.NewCopyBytes(t.Input))

	if t.Type != 0 {
		// either dynamic and access type
		accessList, err := t.AccessList.MarshalRLPWith(arena)
		if err != nil {
			return nil, err
		}
		vv.Set(accessList)
	}

	// signature values
	vv.Set(arena.NewCopyBytes(t.V))
	vv.Set(arena.NewCopyBytes(t.R))
	vv.Set(arena.NewCopyBytes(t.S))

	if t.Type == TransactionLegacy {
		return vv, nil
	}
	return vv, nil
}

func (t *Transaction) UnmarshalRLP(buf []byte) error {
	t.Hash = BytesToHash(Keccak256(buf))

	if len(buf) < 1 {
		return fmt.Errorf("expecting 1 byte but 0 byte provided")
	}
	if buf[0] <= 0x7f {
		// it includes a type byte
		switch typ := buf[0]; typ {
		case 1:
			t.Type = TransactionAccessList
		case 2:
			t.Type = TransactionDynamicFee
		default:
			return fmt.Errorf("type byte %d not found", typ)
		}
		buf = buf[1:]
	}
	if err := fastrlp.UnmarshalRLP(buf, t); err != nil {
		return err
	}
	return nil
}

func (t *Transaction) UnmarshalRLPWith(v *fastrlp.Value) error {
	elems, err := v.GetElems()
	if err != nil {
		return err
	}

	getElem := func() *fastrlp.Value {
		v := elems[0]
		elems = elems[1:]
		return v
	}

	var num int
	switch t.Type {
	case TransactionLegacy:
		num = 9
	case TransactionAccessList:
		// legacy + chain id + access list
		num = 11
	case TransactionDynamicFee:
		// access list txn + gas fee 1 + gas fee 2 - gas price
		num = 12
	default:
		return fmt.Errorf("transaction type %d not found", t.Type)
	}
	if numElems := len(elems); numElems != num {
		return fmt.Errorf("not enough elements to decode transaction, expected %d but found %d", num, numElems)
	}

	if t.Type != 0 {
		t.ChainID = new(big.Int)
		if err := getElem().GetBigInt(t.ChainID); err != nil {
			return err
		}
	}

	// nonce
	if t.Nonce, err = getElem().GetUint64(); err != nil {
		return err
	}

	if t.Type == TransactionDynamicFee {
		// dynamic fee uses
		t.MaxPriorityFeePerGas = new(big.Int)
		if err := getElem().GetBigInt(t.MaxPriorityFeePerGas); err != nil {
			return err
		}
		t.MaxFeePerGas = new(big.Int)
		if err := getElem().GetBigInt(t.MaxFeePerGas); err != nil {
			return err
		}
	} else {
		// legacy and access type use gas price
		if t.GasPrice, err = getElem().GetUint64(); err != nil {
			return err
		}
	}

	// gas
	if t.Gas, err = getElem().GetUint64(); err != nil {
		return err
	}
	// to
	vv, _ := getElem().Bytes()
	if len(vv) == 20 {
		// address
		addr := BytesToAddress(vv)
		t.To = &addr
	} else {
		// reset To
		t.To = nil
	}
	// value
	t.Value = new(big.Int)
	if err := getElem().GetBigInt(t.Value); err != nil {
		return err
	}
	// input
	if t.Input, err = getElem().GetBytes(t.Input[:0]); err != nil {
		return err
	}

	if t.Type != 0 {
		if err := t.AccessList.UnmarshalRLPWith(getElem()); err != nil {
			return err
		}
	}

	// V
	if t.V, err = getElem().GetBytes(t.V); err != nil {
		return err
	}
	// R
	if t.R, err = getElem().GetBytes(t.R); err != nil {
		return err
	}
	// S
	if t.S, err = getElem().GetBytes(t.S); err != nil {
		return err
	}

	return nil
}

func (a *AccessList) MarshalRLPTo(dst []byte) ([]byte, error) {
	return fastrlp.MarshalRLP(a)
}

func (a *AccessList) MarshalRLPWith(arena *fastrlp.Arena) (*fastrlp.Value, error) {
	if len(*a) == 0 {
		return arena.NewNullArray(), nil
	}
	v := arena.NewArray()
	for _, i := range *a {
		acct := arena.NewArray()
		acct.Set(arena.NewCopyBytes(i.Address[:]))
		if len(i.Storage) == 0 {
			acct.Set(arena.NewNullArray())
		} else {
			strV := arena.NewArray()
			for _, v := range i.Storage {
				strV.Set(arena.NewCopyBytes(v[:]))
			}
			acct.Set(strV)
		}
		v.Set(acct)
	}
	return v, nil
}

func (a *AccessList) UnmarshalRLP(buf []byte) error {
	return fastrlp.UnmarshalRLP(buf, a)
}

func (a *AccessList) UnmarshalRLPWith(v *fastrlp.Value) error {
	if v.Type() == fastrlp.TypeArrayNull {
		// empty
		return nil
	}

	elems, err := v.GetElems()
	if err != nil {
		return err
	}
	for _, elem := range elems {
		entry := AccessEntry{}

		acctElems, err := elem.GetElems()
		if err != nil {
			return err
		}
		if len(acctElems) != 2 {
			return fmt.Errorf("two elems expected but %d found", len(acctElems))
		}

		// decode 'address'
		if err = acctElems[0].GetAddr(entry.Address[:]); err != nil {
			return err
		}

		// decode 'storage'
		if acctElems[1].Type() != fastrlp.TypeArrayNull {
			storageElems, err := acctElems[1].GetElems()
			if err != nil {
				return err
			}

			entry.Storage = make([]Hash, len(storageElems))
			for indx, storage := range storageElems {
				// decode storage
				if err = storage.GetHash(entry.Storage[indx][:]); err != nil {
					return err
				}
			}
		}
		(*a) = append((*a), entry)
	}
	return nil
}

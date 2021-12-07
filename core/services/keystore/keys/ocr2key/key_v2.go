package ocr2key

type Raw []byte

func (raw Raw) Key() KeyBundle {
	// offchain private key 64 bytes || offchain encryption key 32 bytes || onchain 32 bytes private key
	var kb KeyBundle
	err := kb.Unmarshal(raw)
	if err != nil {
		panic(err)
	}
	return kb
}

func (raw Raw) String() string {
	return "<OCR2 Raw Private Key>"
}

func (raw Raw) GoString() string {
	return raw.String()
}

func NewV2() (KeyBundle, error) {
	return KeyBundle{}, nil
}

func (key KeyBundle) Raw() Raw {
	b, err := key.Marshal()
	if err != nil {
		panic(err)
	}
	return b
}

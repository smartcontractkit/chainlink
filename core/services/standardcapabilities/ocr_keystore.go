package standardcapabilities

import (
	"context"
	"fmt"
	"slices"
	"strings"

	capkeystore "github.com/smartcontractkit/capabilities/libs/capability/keystore"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

// ocrKeystore serves the node's OCR2 key bundles to a v2 process.
type ocrKeystore struct {
	ks keystore.OCR2
}

var _ capkeystore.OCRKeystore = ocrKeystore{}

func (s ocrKeystore) List(_ context.Context, chainTypes ...string) ([]capkeystore.OCRKeyBundle, error) {
	all, err := s.ks.GetAll()
	if err != nil {
		return nil, err
	}

	var kbs []capkeystore.OCRKeyBundle
	for _, kb := range all {
		if len(chainTypes) == 0 || slices.Contains(chainTypes, string(kb.ChainType())) {
			kbs = append(kbs, ocrKeyBundle{kb})
		}
	}

	slices.SortFunc(kbs, func(a, b capkeystore.OCRKeyBundle) int { return strings.Compare(a.ID(), b.ID()) })
	return kbs, nil
}

// Get searches every bundle rather than calling keystore.OCR2.Get, whose error does not say whether the bundle is
// missing or the keystore failed.
func (s ocrKeystore) Get(_ context.Context, id string) (capkeystore.OCRKeyBundle, error) {
	all, err := s.ks.GetAll()
	if err != nil {
		return nil, err
	}

	i := slices.IndexFunc(all, func(kb ocr2key.KeyBundle) bool { return kb.ID() == id })
	if i < 0 {
		return nil, fmt.Errorf("%s: %w", id, capkeystore.ErrKeyBundleNotFound)
	}

	return ocrKeyBundle{all[i]}, nil
}

// ocrKeyBundle is an ocr2key.KeyBundle with its chain type as a string; every other method is promoted.
type ocrKeyBundle struct {
	ocr2key.KeyBundle
}

func (b ocrKeyBundle) ChainType() string { return string(b.KeyBundle.ChainType()) }

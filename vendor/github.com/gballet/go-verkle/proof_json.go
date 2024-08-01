// This is free and unencumbered software released into the public domain.
//
// Anyone is free to copy, modify, publish, use, compile, sell, or
// distribute this software, either in source code form or as a compiled
// binary, for any purpose, commercial or non-commercial, and by any
// means.
//
// In jurisdictions that recognize copyright laws, the author or authors
// of this software dedicate any and all copyright interest in the
// software to the public domain. We make this dedication for the benefit
// of the public at large and to the detriment of our heirs and
// successors. We intend this dedication to be an overt act of
// relinquishment in perpetuity of all present and future rights to this
// software under copyright law.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
//
// For more information, please refer to <https://unlicense.org>

package verkle

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// HexToPrefixedString turns a byte slice into its hex representation
// and prefixes it with `0x`.
func HexToPrefixedString(data []byte) string {
	return "0x" + hex.EncodeToString(data)
}

// PrefixedHexStringToBytes does the opposite of HexToPrefixedString.
func PrefixedHexStringToBytes(input string) ([]byte, error) {
	if input[0:2] == "0x" {
		input = input[2:]
	}
	return hex.DecodeString(input)
}

type ipaproofMarshaller struct {
	CL              [IPA_PROOF_DEPTH]string `json:"cl"`
	CR              [IPA_PROOF_DEPTH]string `json:"cr"`
	FinalEvaluation string                  `json:"finalEvaluation"`
}

func (ipp *IPAProof) MarshalJSON() ([]byte, error) {
	return json.Marshal(&ipaproofMarshaller{
		CL: [IPA_PROOF_DEPTH]string{
			HexToPrefixedString(ipp.CL[0][:]),
			HexToPrefixedString(ipp.CL[1][:]),
			HexToPrefixedString(ipp.CL[2][:]),
			HexToPrefixedString(ipp.CL[3][:]),
			HexToPrefixedString(ipp.CL[4][:]),
			HexToPrefixedString(ipp.CL[5][:]),
			HexToPrefixedString(ipp.CL[6][:]),
			HexToPrefixedString(ipp.CL[7][:]),
		},
		CR: [IPA_PROOF_DEPTH]string{
			HexToPrefixedString(ipp.CR[0][:]),
			HexToPrefixedString(ipp.CR[1][:]),
			HexToPrefixedString(ipp.CR[2][:]),
			HexToPrefixedString(ipp.CR[3][:]),
			HexToPrefixedString(ipp.CR[4][:]),
			HexToPrefixedString(ipp.CR[5][:]),
			HexToPrefixedString(ipp.CR[6][:]),
			HexToPrefixedString(ipp.CR[7][:]),
		},
		FinalEvaluation: HexToPrefixedString(ipp.FinalEvaluation[:]),
	})
}

func (ipp *IPAProof) UnmarshalJSON(data []byte) error {
	aux := &ipaproofMarshaller{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if len(aux.FinalEvaluation) != 64 && len(aux.FinalEvaluation) != 66 {
		return fmt.Errorf("invalid hex string for final evaluation: %s", aux.FinalEvaluation)
	}

	currentValueBytes, err := PrefixedHexStringToBytes(aux.FinalEvaluation)
	if err != nil {
		return fmt.Errorf("error decoding hex string for current value: %v", err)
	}
	copy(ipp.FinalEvaluation[:], currentValueBytes)

	for i := range ipp.CL {
		if len(aux.CL[i]) != 64 && len(aux.CL[i]) != 66 {
			return fmt.Errorf("invalid hex string for CL[%d]: %s", i, aux.CL[i])
		}
		val, err := PrefixedHexStringToBytes(aux.CL[i])
		if err != nil {
			return fmt.Errorf("error decoding hex string for CL[%d]: %s", i, aux.CL[i])
		}
		copy(ipp.CL[i][:], val)
		if len(aux.CR[i]) != 64 && len(aux.CR[i]) != 66 {
			return fmt.Errorf("invalid hex string for CR[%d]: %s", i, aux.CR[i])
		}
		val, err = PrefixedHexStringToBytes(aux.CR[i])
		if err != nil {
			return fmt.Errorf("error decoding hex string for CR[%d]: %s", i, aux.CR[i])
		}
		copy(ipp.CR[i][:], val)
	}
	copy(ipp.FinalEvaluation[:], currentValueBytes)

	return nil
}

type verkleProofMarshaller struct {
	OtherStems            []string  `json:"otherStems"`
	DepthExtensionPresent string    `json:"depthExtensionPresent"`
	CommitmentsByPath     []string  `json:"commitmentsByPath"`
	D                     string    `json:"d"`
	IPAProof              *IPAProof `json:"ipaProof"`
}

func (vp *VerkleProof) MarshalJSON() ([]byte, error) {
	aux := &verkleProofMarshaller{
		OtherStems:            make([]string, len(vp.OtherStems)),
		DepthExtensionPresent: HexToPrefixedString(vp.DepthExtensionPresent),
		CommitmentsByPath:     make([]string, len(vp.CommitmentsByPath)),
		D:                     HexToPrefixedString(vp.D[:]),
		IPAProof:              vp.IPAProof,
	}

	for i, s := range vp.OtherStems {
		aux.OtherStems[i] = HexToPrefixedString(s[:])
	}
	for i, c := range vp.CommitmentsByPath {
		aux.CommitmentsByPath[i] = HexToPrefixedString(c[:])
	}
	return json.Marshal(aux)
}

func (vp *VerkleProof) UnmarshalJSON(data []byte) error {
	var aux verkleProofMarshaller
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return fmt.Errorf("verkle proof unmarshal error: %w", err)
	}

	vp.DepthExtensionPresent, err = PrefixedHexStringToBytes(aux.DepthExtensionPresent)
	if err != nil {
		return fmt.Errorf("error decoding hex string for depth and extension present: %v", err)
	}

	vp.CommitmentsByPath = make([][32]byte, len(aux.CommitmentsByPath))
	for i, c := range aux.CommitmentsByPath {
		val, err := PrefixedHexStringToBytes(c)
		if err != nil {
			return fmt.Errorf("error decoding hex string for commitment #%d: %w", i, err)
		}
		copy(vp.CommitmentsByPath[i][:], val)
	}

	currentValueBytes, err := PrefixedHexStringToBytes(aux.D)
	if err != nil {
		return fmt.Errorf("error decoding hex string for D: %w", err)
	}
	copy(vp.D[:], currentValueBytes)

	vp.OtherStems = make([][31]byte, len(aux.OtherStems))
	for i, c := range aux.OtherStems {
		val, err := PrefixedHexStringToBytes(c)
		if err != nil {
			return fmt.Errorf("error decoding hex string for other stem #%d: %w", i, err)
		}
		copy(vp.OtherStems[i][:], val)
	}

	vp.IPAProof = aux.IPAProof
	return nil
}

type stemStateDiffMarshaller struct {
	Stem        string           `json:"stem"`
	SuffixDiffs SuffixStateDiffs `json:"suffixDiffs"`
}

func (ssd StemStateDiff) MarshalJSON() ([]byte, error) {
	return json.Marshal(&stemStateDiffMarshaller{
		Stem:        HexToPrefixedString(ssd.Stem[:]),
		SuffixDiffs: ssd.SuffixDiffs,
	})
}

func (ssd *StemStateDiff) UnmarshalJSON(data []byte) error {
	var aux stemStateDiffMarshaller
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("stemdiff unmarshal error: %w", err)
	}

	stem, err := PrefixedHexStringToBytes(aux.Stem)
	if err != nil {
		return fmt.Errorf("invalid hex string for stem: %w", err)
	}
	*ssd = StemStateDiff{
		SuffixDiffs: aux.SuffixDiffs,
	}
	copy(ssd.Stem[:], stem)
	return nil
}

type suffixStateDiffMarshaller struct {
	Suffix       byte    `json:"suffix"`
	CurrentValue *string `json:"currentValue"`
	NewValue     *string `json:"newValue"`
}

func (ssd SuffixStateDiff) MarshalJSON() ([]byte, error) {
	var cvstr, nvstr *string
	if ssd.CurrentValue != nil {
		tempstr := HexToPrefixedString(ssd.CurrentValue[:])
		cvstr = &tempstr
	}
	if ssd.NewValue != nil {
		tempstr := HexToPrefixedString(ssd.NewValue[:])
		nvstr = &tempstr
	}
	return json.Marshal(&suffixStateDiffMarshaller{
		Suffix:       ssd.Suffix,
		CurrentValue: cvstr,
		NewValue:     nvstr,
	})
}

func (ssd *SuffixStateDiff) UnmarshalJSON(data []byte) error {
	aux := &suffixStateDiffMarshaller{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("suffix diff unmarshal error: %w", err)
	}

	if aux.CurrentValue != nil && len(*aux.CurrentValue) != 64 && len(*aux.CurrentValue) != 0 && len(*aux.CurrentValue) != 66 {
		return fmt.Errorf("invalid hex string for current value: %s", *aux.CurrentValue)
	}

	*ssd = SuffixStateDiff{
		Suffix: aux.Suffix,
	}

	if aux.CurrentValue != nil && len(*aux.CurrentValue) != 0 {
		currentValueBytes, err := PrefixedHexStringToBytes(*aux.CurrentValue)
		if err != nil {
			return fmt.Errorf("error decoding hex string for current value: %v", err)
		}

		ssd.CurrentValue = &[32]byte{}
		copy(ssd.CurrentValue[:], currentValueBytes)
	}

	if aux.NewValue != nil && len(*aux.NewValue) != 0 {
		newValueBytes, err := PrefixedHexStringToBytes(*aux.NewValue)
		if err != nil {
			return fmt.Errorf("error decoding hex string for current value: %v", err)
		}

		ssd.NewValue = &[32]byte{}
		copy(ssd.NewValue[:], newValueBytes)
	}

	return nil
}

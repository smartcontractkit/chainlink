package protocol

import (
	"bytes"
)

// The functions in this file are only used in tests, hence
// the name "TestEqual" to make that more clear

func (rep AttestedReportOne) TestEqual(rep2 AttestedReportOne) bool {
	return rep.Skip == rep2.Skip && bytes.Equal(rep.Report, rep2.Report) && bytes.Equal(rep.Signature, rep2.Signature)
}

func (rep AttestedReportMany) TestEqual(c2 AttestedReportMany) bool {
	if !bytes.Equal(rep.Report, c2.Report) {
		return false
	}

	if len(rep.AttributedSignatures) != len(c2.AttributedSignatures) {
		return false
	}

	for i := range rep.AttributedSignatures {
		if !rep.AttributedSignatures[i].Equal(c2.AttributedSignatures[i]) {
			return false
		}
	}

	return true
}

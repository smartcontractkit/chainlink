package web

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"

	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/vault"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// VaultController exposes utilities for the Vault
type VaultController struct {
	App chainlink.Application
}

type VerifyDKGResultRequest struct {
	InstanceID      string `json:"instanceId"`
	MasterPublicKey string `json:"masterPublicKey"`
}

// VerifyDKGResult verifies that the DKGResult corresponds to the given public key
// and contains a share that can be decrypted by the node's DKG recipient key.
// Example:
// "POST <application>/vault/dkg_results/verify"
func (vc *VaultController) VerifyDKGResult(c *gin.Context) {
	var req VerifyDKGResultRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, errors.New("could not parse request body"))
		return
	}

	if req.InstanceID == "" || req.MasterPublicKey == "" {
		jsonAPIError(c, http.StatusBadRequest, errors.New("instanceId and masterPublicKey are required"))
		return
	}

	keys, err := vc.App.GetKeyStore().DKGRecipient().GetAll()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	if len(keys) == 0 {
		jsonAPIError(c, http.StatusBadRequest, errors.New("no DKG recipient key found"))
		return
	}
	if len(keys) > 1 {
		jsonAPIError(c, http.StatusBadRequest, errors.New("multiple DKG recipient keys found"))
		return
	}

	orm := vault.NewVaultORM(vc.App.GetDB())
	v, err := orm.ReadResultPackage(c.Request.Context(), dkgocrtypes.InstanceID(req.InstanceID))
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	}

	if v == nil {
		jsonAPIError(c, http.StatusNotFound, errors.New("DKG result not found"))
		return
	}

	err = vaultcap.VerifyDKGResult(v.ReportWithResultPackage, req.MasterPublicKey, keys[0])
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, fmt.Errorf("DKG result verification failed: %w", err))
		return
	}

	sha := sha256.Sum256(v.ReportWithResultPackage)
	shaStr := hex.EncodeToString(sha[:])
	jsonAPIResponse(c, presenters.NewVerifyDKGResultResource(shaStr), "verifyDKGResult")
}

type ExportDKGResultRequest struct {
	InstanceID string `json:"instanceId"`
}

// ExportDKGResult returns the DKGResult corresponding to the given instance ID
// "POST <application>/vault/dkg_results/export"
func (vc *VaultController) ExportDKGResult(c *gin.Context) {
	var req ExportDKGResultRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, errors.New("could not parse request body"))
		return
	}

	if req.InstanceID == "" {
		jsonAPIError(c, http.StatusBadRequest, errors.New("instanceId is required"))
		return
	}

	orm := vault.NewVaultORM(vc.App.GetDB())
	v, err := orm.ReadResultPackage(c.Request.Context(), dkgocrtypes.InstanceID(req.InstanceID))
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	}

	if v == nil {
		jsonAPIError(c, http.StatusNotFound, errors.New("DKG result not found"))
		return
	}

	hexPackage := hex.EncodeToString(v.ReportWithResultPackage)
	sha := sha256.Sum256(v.ReportWithResultPackage)
	shaStr := hex.EncodeToString(sha[:])
	jsonAPIResponse(c, presenters.NewExportDKGResultResource(hexPackage, shaStr), "exportDKGResult")
}

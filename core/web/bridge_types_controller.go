package web

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgconn"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// ValidateBridgeTypeNotExist checks that a bridge has not already been created
func ValidateBridgeTypeNotExist(bt *bridges.BridgeTypeRequest, orm bridges.ORM) error {
	fe := models.NewJSONAPIErrors()
	_, err := orm.FindBridge(bt.Name)
	if err == nil {
		fe.Add(fmt.Sprintf("Bridge Type %v already exists", bt.Name))
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		fe.Add(fmt.Sprintf("Error determining if bridge type %v already exists", bt.Name))
	}
	return fe.CoerceEmptyToNil()
}

// ValidateBridgeType checks that the bridge type has the required field with valid values.
func ValidateBridgeType(bt *bridges.BridgeTypeRequest) error {
	fe := models.NewJSONAPIErrors()
	if len(bt.Name.String()) < 1 {
		fe.Add("No name specified")
	}
	if _, err := bridges.ParseBridgeName(bt.Name.String()); err != nil {
		fe.Merge(err)
	}
	u := bt.URL.String()
	if len(strings.TrimSpace(u)) == 0 {
		fe.Add("URL must be present")
	}
	if bt.MinimumContractPayment != nil &&
		bt.MinimumContractPayment.Cmp(assets.NewLinkFromJuels(0)) < 0 {
		fe.Add("MinimumContractPayment must be positive")
	}
	return fe.CoerceEmptyToNil()
}

// BridgeTypesController manages BridgeType requests in the node.
type BridgeTypesController struct {
	App chainlink.Application
}

// Create adds the BridgeType to the given context.
func (btc *BridgeTypesController) Create(c *gin.Context) {
	btr := &bridges.BridgeTypeRequest{}

	if err := c.ShouldBindJSON(btr); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	bta, bt, err := bridges.NewBridgeType(btr)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	if e := ValidateBridgeType(btr); e != nil {
		jsonAPIError(c, http.StatusBadRequest, e)
		return
	}
	orm := btc.App.BridgeORM()
	if e := ValidateBridgeTypeNotExist(btr, orm); e != nil {
		jsonAPIError(c, http.StatusBadRequest, e)
		return
	}
	if e := orm.CreateBridgeType(bt); e != nil {
		jsonAPIError(c, http.StatusInternalServerError, e)
		return
	}
	switch e := err.(type) {
	case *pgconn.PgError:
		var apiErr error
		if e.ConstraintName == "external_initiators_name_key" {
			apiErr = fmt.Errorf("bridge Type %v conflict", bt.Name)
		} else {
			apiErr = err
		}
		jsonAPIError(c, http.StatusConflict, apiErr)
		return
	default:
		resource := presenters.NewBridgeResource(*bt)
		resource.IncomingToken = bta.IncomingToken

		jsonAPIResponse(c, resource, "bridge")
	}
}

// Index lists Bridges, one page at a time.
func (btc *BridgeTypesController) Index(c *gin.Context, size, page, offset int) {
	bridges, count, err := btc.App.BridgeORM().BridgeTypes(offset, size)

	var resources []presenters.BridgeResource
	for _, bridge := range bridges {
		resources = append(resources, *presenters.NewBridgeResource(bridge))
	}

	paginatedResponse(c, "Bridges", size, page, resources, count, err)
}

// Show returns the details of a specific Bridge.
func (btc *BridgeTypesController) Show(c *gin.Context) {
	name := c.Param("BridgeName")

	taskType, err := bridges.ParseBridgeName(name)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	bt, err := btc.App.BridgeORM().FindBridge(taskType)
	if errors.Is(err, sql.ErrNoRows) {
		jsonAPIError(c, http.StatusNotFound, errors.New("bridge not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewBridgeResource(bt), "bridge")
}

// Update can change the restricted attributes for a bridge
func (btc *BridgeTypesController) Update(c *gin.Context) {
	name := c.Param("BridgeName")
	btr := &bridges.BridgeTypeRequest{}

	taskType, err := bridges.ParseBridgeName(name)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	orm := btc.App.BridgeORM()
	bt, err := orm.FindBridge(taskType)
	if errors.Is(err, sql.ErrNoRows) {
		jsonAPIError(c, http.StatusNotFound, errors.New("bridge not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	if err := c.ShouldBindJSON(btr); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	if err := ValidateBridgeType(btr); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	if err := orm.UpdateBridgeType(&bt, btr); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewBridgeResource(bt), "bridge")
}

// Destroy removes a specific Bridge.
func (btc *BridgeTypesController) Destroy(c *gin.Context) {
	name := c.Param("BridgeName")

	taskType, err := bridges.ParseBridgeName(name)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	orm := btc.App.BridgeORM()
	bt, err := orm.FindBridge(taskType)
	if errors.Is(err, sql.ErrNoRows) {
		jsonAPIError(c, http.StatusNotFound, errors.New("bridge not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("error searching for bridge: %+v", err))
		return
	}
	jobsUsingBridge, err := btc.App.JobORM().FindJobIDsWithBridge(name)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("error searching for associated v2 jobs: %+v", err))
		return
	}
	if len(jobsUsingBridge) > 0 {
		jsonAPIError(c, http.StatusConflict, fmt.Errorf("can't remove the bridge because jobs %v are associated with it", jobsUsingBridge))
		return
	}
	if err = orm.DeleteBridgeType(&bt); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("failed to delete bridge: %+v", err))
		return
	}

	jsonAPIResponse(c, presenters.NewBridgeResource(bt), "bridge")
}

package web

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// JobProposalsController manages the job proposals
type JobProposalsController struct {
	App chainlink.Application
}

// Approve approves a job proposal.
// Example:
// "POST <application>/job_proposals/<id>/reject"
func (jpc *JobProposalsController) Approve(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	feedsSvc := jpc.App.GetFeedsService()

	err = feedsSvc.ApproveJobProposal(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			jsonAPIError(c, http.StatusNotFound, errors.New("job proposal not found"))
			return
		}

		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jp, err := feedsSvc.GetJobProposal(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			jsonAPIError(c, http.StatusNotFound, errors.New("job proposal not found"))
			return
		}

		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c,
		presenters.NewJobProposalResource(*jp),
		"job_proposals",
		http.StatusOK,
	)
}

// Reject rejects a job proposal.
// Example:
// "POST <application>/job_proposals/<id>/reject"
func (jpc *JobProposalsController) Reject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	feedsSvc := jpc.App.GetFeedsService()

	err = feedsSvc.RejectJobProposal(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			jsonAPIError(c, http.StatusNotFound, errors.New("job proposal not found"))
			return
		}

		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jp, err := feedsSvc.GetJobProposal(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			jsonAPIError(c, http.StatusNotFound, errors.New("job proposal not found"))
			return
		}

		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c,
		presenters.NewJobProposalResource(*jp),
		"job_proposals",
		http.StatusOK,
	)
}

type UpdateSpecRequest struct {
	Spec string `json:"spec"`
}

// UpdateSpec updates the spec of a job proposal
// Example:
// "POST <application>/job_proposals/<id>/reject"
func (jpc *JobProposalsController) UpdateSpec(c *gin.Context) {
	request := UpdateSpecRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	feedsSvc := jpc.App.GetFeedsService()

	err = feedsSvc.UpdateJobProposalSpec(c.Request.Context(), id, request.Spec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			jsonAPIError(c, http.StatusNotFound, errors.New("job proposal not found"))
			return
		}

		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jp, err := feedsSvc.GetJobProposal(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			jsonAPIError(c, http.StatusNotFound, errors.New("job proposal not found"))
			return
		}

		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c,
		presenters.NewJobProposalResource(*jp),
		"job_proposals",
		http.StatusOK,
	)
}

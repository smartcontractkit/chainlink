import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'

export const ENDPOINT = '/v2/job_proposals'
export const SHOW_ENDPOINT = `${ENDPOINT}/:id`
export const REJECT_ENDPOINT = `${ENDPOINT}/:id/reject`
export const APPROVE_ENDPOINT = `${ENDPOINT}/:id/approve`
export const CANCEL_ENDPOINT = `${ENDPOINT}/:id/cancel`
export const UPDATE_SPEC_ENDPOINT = `${ENDPOINT}/:id/spec`

// Job Proposals represents the job proposals
export class JobProposals {
  constructor(private api: jsonapi.Api) {}

  @boundMethod
  public getJobProposal(
    id: string,
  ): Promise<jsonapi.ApiResponse<models.JobProposal>> {
    return this.show({}, { id })
  }

  @boundMethod
  public approveJobProposal(
    id: string,
  ): Promise<jsonapi.ApiResponse<models.JobProposal>> {
    return this.approve({}, { id })
  }

  @boundMethod
  public rejectJobProposal(
    id: string,
  ): Promise<jsonapi.ApiResponse<models.JobProposal>> {
    return this.reject({}, { id })
  }

  @boundMethod
  public cancelJobProposal(
    id: string,
  ): Promise<jsonapi.ApiResponse<models.JobProposal>> {
    return this.cancel({}, { id })
  }

  @boundMethod
  public updateJobProposalSpec(
    id: string,
    req: models.UpdateJobProposalSpecRequest,
  ): Promise<jsonapi.ApiResponse<models.JobProposal>> {
    return this.updateSpec(req, { id })
  }

  private show = this.api.fetchResource<
    {},
    models.JobProposal,
    {
      id: string
    }
  >(SHOW_ENDPOINT)

  private reject = this.api.createResource<
    {},
    models.JobProposal,
    {
      id: string
    }
  >(REJECT_ENDPOINT)

  private approve = this.api.createResource<
    {},
    models.JobProposal,
    {
      id: string
    }
  >(APPROVE_ENDPOINT)

  private cancel = this.api.createResource<
    {},
    models.JobProposal,
    {
      id: string
    }
  >(CANCEL_ENDPOINT)

  private updateSpec = this.api.updateResource<
    models.UpdateJobProposalSpecRequest,
    models.JobProposal,
    {
      id: string
    }
  >(UPDATE_SPEC_ENDPOINT)
}

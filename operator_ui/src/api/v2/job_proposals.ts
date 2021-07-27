import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'

export const ENDPOINT = '/v2/job_proposals'
export const SHOW_ENDPOINT = `${ENDPOINT}/:id`

// Job Proposals represents the job proposals
export class JobProposals {
  constructor(private api: jsonapi.Api) {}

  @boundMethod
  public getJobProposals(): Promise<jsonapi.ApiResponse<models.JobProposal[]>> {
    return this.index()
  }

  @boundMethod
  public getJobProposal(
    id: string,
  ): Promise<jsonapi.ApiResponse<models.JobProposal>> {
    return this.show({}, { id })
  }

  private index = this.api.fetchResource<{}, models.JobProposal[]>(ENDPOINT)
  private show = this.api.fetchResource<
    {},
    models.JobProposal,
    {
      id: string
    }
  >(SHOW_ENDPOINT)
}

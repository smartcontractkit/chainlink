import * as jsonapi from '@chainlink/json-api-client'
import { boundMethod } from 'autobind-decorator'

/**
 * Destroy deletes a JobSpecError, effectively dismissing the notification
 *
 * @example "<application>/specs/:SpecID"
 */
interface DestroyPathParams {
  jobSpecErrorID: number
}

const DESTROY_ENDPOINT = '/v2/job_spec_errors/:jobSpecErrorID'

export class JobSpecErrors {
  constructor(private api: jsonapi.Api) {}

  @boundMethod
  public destroyJobSpecError(id: number): Promise<jsonapi.ApiResponse<null>> {
    return this.destroy(undefined, { jobSpecErrorID: id })
  }

  private destroy = this.api.deleteResource<undefined, null, DestroyPathParams>(
    DESTROY_ENDPOINT,
  )
}

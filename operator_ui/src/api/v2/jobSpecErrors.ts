import * as jsonapi from '@chainlink/json-api-client'
import { boundMethod } from 'autobind-decorator'

/**
 * Destroy deletes a JobSpecError, effectively dismissing the notification
 *
 * @example "<application>/specs/:SpecID"
 */
interface DestroyPathParams {
  id: string
}

const DESTROY_ENDPOINT = '/v2/job_spec_errors/:id'

export class JobSpecErrors {
  constructor(private api: jsonapi.Api) {}

  @boundMethod
  public destroyJobSpecError(id: string): Promise<jsonapi.ApiResponse<null>> {
    return this.destroy(undefined, { id })
  }

  private destroy = this.api.deleteResource<undefined, null, DestroyPathParams>(
    DESTROY_ENDPOINT,
  )
}

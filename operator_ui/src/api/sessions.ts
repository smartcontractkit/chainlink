import * as jsonapi from 'utils/json-api-client'
import { Api } from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'
import * as sessionsController from 'core/web/sessions_controller'

/**
 * Create creates a session ID for the given user credentials
 * and returns it in a cookie.
 */
const CREATE_ENDPOINT = '/sessions'

/**
 * Destroy erases the session ID for the sole API user.
 */
const DESTROY_ENDPOINT = '/sessions'

export class Sessions {
  constructor(private api: Api) {}

  @boundMethod
  public createSession(
    sessionRequest: models.SessionRequest,
  ): Promise<jsonapi.ApiResponse<sessionsController.Session>> {
    return this.create(sessionRequest)
  }

  @boundMethod
  public destroySession(): Promise<
    jsonapi.ApiResponse<sessionsController.Session>
  > {
    return this.destroy()
  }

  private create = this.api.createResource<
    models.SessionRequest,
    sessionsController.Session
  >(CREATE_ENDPOINT)

  private destroy = this.api.deleteResource<
    undefined,
    sessionsController.Session
  >(DESTROY_ENDPOINT)
}

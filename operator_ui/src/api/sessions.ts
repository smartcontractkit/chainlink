import * as jsonapi from 'api/transport/json'
import * as models from 'core/store/models'
import * as sessionsController from 'core/web/sessions_controller'

/**
 * Create creates a session ID for the given user credentials
 * and returns it in a cookie.
 */
const CREATE_ENDPOINT = '/sessions'
const create = jsonapi.createResource<
  models.SessionRequest,
  sessionsController.Session
>(CREATE_ENDPOINT)

/**
 * Destroy erases the session ID for the sole API user.
 */
const DESTROY_ENDPOINT = '/sessions'
const destroy = jsonapi.deleteResource<undefined, sessionsController.Session>(
  DESTROY_ENDPOINT,
)

export function createSession(
  sessionRequest: models.SessionRequest,
): Promise<jsonapi.ApiResponse<sessionsController.Session>> {
  return create(sessionRequest)
}

export function destroySession(): Promise<
  jsonapi.ApiResponse<sessionsController.Session>
> {
  return destroy()
}

import * as jsonapi from 'api/transport/json'
import * as models from 'core/store/models'
import * as sessions_controller from 'core/web/sessions_controller'

/**
 * Create creates a session ID for the given user credentials
 * and returns it in a cookie.
 */
const CREATE_ENDPOINT = '/sessions'
const create = jsonapi.createResource<
  models.SessionRequest,
  sessions_controller.Session
>(CREATE_ENDPOINT)

/**
 * Destroy erases the session ID for the sole API user.
 */
const DESTROY_ENDPOINT = '/sessions'
const destroy = jsonapi.deleteResource<undefined, sessions_controller.Session>(
  DESTROY_ENDPOINT
)

export const createSession = (sessionRequest: models.SessionRequest) =>
  create(sessionRequest)

export const destroySession = () => destroy()

import { Middleware } from 'redux'
import { UpdateQueryAction } from '../reducers/actions'
import { Query } from '../reducers/query'

/**
 * Parse the current query string out of the browser location
 *
 * @param location The location value to use, not hardcoded so we can inject
 * mock values for testing
 */
function parseQuery(location: Location): Query {
  const searchParams = new URL(location.toString()).searchParams
  const search = searchParams.get('search')

  if (search) {
    return search
  }

  return
}

/**
 * Create a redux middleware responsible for updating the current query for every action
 *
 * @param location The location value to use, can be injected for testing
 */
export function createQueryMiddleware(
  location: Location = document.location,
): Middleware {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const queryMiddleware: Middleware = _store => next => action => {
    // dispatch original action right away
    next(action)

    // parse query and dispatch an update
    // TODO: throttle this?
    const query = parseQuery(location)
    const updateQueryAction: UpdateQueryAction = {
      type: 'QUERY_UPDATED',
      data: query,
    }
    next(updateQueryAction)
  }

  return queryMiddleware
}

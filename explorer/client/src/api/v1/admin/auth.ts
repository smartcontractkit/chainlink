import * as jsonapi from '@chainlink/json-api-client'

const SIGN_IN_ENDPOINT = '/api/v1/admin/login'

interface SignInRequestParams {
  username: string
  password: string
}

/**
 * SignIn authenticates an admin user.
 * @param username The admin username
 * @param password The plain text password
 */
export function signIn(
  username: string,
  password: string,
): Promise<jsonapi.PaginatedApiResponse<number>> {
  return jsonapi.createResource<SignInRequestParams, number>(SIGN_IN_ENDPOINT)({
    username,
    password,
  })
}

const SIGN_OUT_ENDPOINT = '/api/v1/admin/logout'

/**
 * SignOut authenticates an admin user.
 */
export function signOut(): Promise<{}> {
  return jsonapi.deleteResource<{}, {}>(SIGN_OUT_ENDPOINT)()
}

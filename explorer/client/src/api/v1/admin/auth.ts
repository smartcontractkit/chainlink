import * as jsonapi from '@chainlink/json-api-client'
import { Api } from '@chainlink/json-api-client'
import { boundMethod } from 'autobind-decorator'
const SIGN_IN_ENDPOINT = '/api/v1/admin/login'

interface SignInRequestParams {
  username: string
  password: string
}

const SIGN_OUT_ENDPOINT = '/api/v1/admin/logout'

export class Auth {
  constructor(private api: Api) {}

  /**
   * SignIn authenticates an admin user.
   * @param username The admin username
   * @param password The plain text password
   */
  @boundMethod
  public signIn(
    username: string,
    password: string,
  ): Promise<jsonapi.PaginatedApiResponse<number>> {
    return this.api.POST<SignInRequestParams, number>(SIGN_IN_ENDPOINT)({
      username,
      password,
    })
  }

  /**
   * SignOut signs out an admin user.
   */
  @boundMethod
  public signOut(): Promise<{}> {
    return this.api.DELETE<{}, {}>(SIGN_OUT_ENDPOINT)()
  }
}

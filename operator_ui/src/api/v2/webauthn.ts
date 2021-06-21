import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'

const REGISTRATION_ENDPOINT = '/v2/enroll_webauthn'

export class WebAuthn {
  constructor(private api: jsonapi.Api) {}

  @boundMethod
  public beginKeyRegistration(
    request: models.BeginWebAuthnRegistrationV2Request,
  ): Promise<jsonapi.ApiResponse<models.BeginWebAuthnRegistrationV2>> {
    return this.create(request)
  }

  @boundMethod
  public finishKeyRegistration(
    request: models.FinishWebAuthnRegistrationV2Request,
  ): Promise<jsonapi.ApiResponse<models.FinishWebAuthnRegistrationV2>> {
    return this.put(request)
  }

  private create = this.api.fetchResource<
    models.BeginWebAuthnRegistrationV2Request,
    models.BeginWebAuthnRegistrationV2
  >(REGISTRATION_ENDPOINT)

  private put = this.api.createResource<
    models.FinishWebAuthnRegistrationV2Request,
    models.FinishWebAuthnRegistrationV2
  >(REGISTRATION_ENDPOINT)
}

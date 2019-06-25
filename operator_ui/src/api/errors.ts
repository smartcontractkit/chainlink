interface IJsonApiError {
  errors: any
}

export class AuthenticationError extends Error {
  errors: any[]

  constructor(response: Response) {
    // super(`Lorem "${message}" ipsum dolor.`)
    super('TODO: AuthenticationError')
    this.errors = [
      {
        status: response.status,
        detail: response.statusText
      }
    ]
  }
}

export class BadRequestError extends Error {
  errors: any[]

  constructor({ errors }: IJsonApiError) {
    // super(`Lorem "${message}" ipsum dolor.`)
    super('TODO: BadRequestError')
    this.errors = errors
  }
}

export class ServerError extends Error {
  errors: any[]

  constructor(response: Response) {
    // super(`Lorem "${message}" ipsum dolor.`)
    super('TODO: ServerError')
    this.errors = [
      {
        status: response.status,
        detail: response.statusText
      }
    ]
  }
}

export class UnknownResponseError extends Error {
  errors: any[]

  constructor(response: Response) {
    // super(`Lorem "${message}" ipsum dolor.`)
    super('TODO: UnknownResponseError')
    this.errors = [
      {
        status: response.status,
        detail: response.statusText
      }
    ]
  }
}

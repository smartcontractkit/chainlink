import { JsonApiResponse } from 'json-api-normalizer'

interface Error {
  status: number
  detail: any
}

export interface DocumentWithErrors {
  errors: any
}

export class AuthenticationError extends Error {
  errors: Error[]

  constructor(response: Response) {
    super(`AuthenticationError(${response.statusText})`)
    this.errors = [
      {
        status: response.status,
        detail: response.statusText,
      },
    ]
  }
}

export class BadRequestError extends Error {
  errors: JsonApiResponse['errors']

  constructor({ errors }: Pick<JsonApiResponse, 'errors'>) {
    super('BadRequestError')
    this.errors = errors
  }
}

export class ServerError extends Error {
  errors: Error[]

  constructor(response: Response) {
    super(`ServerError(${response.statusText})`)
    this.errors = [
      {
        status: response.status,
        detail: response.statusText,
      },
    ]
  }
}

export class UnknownResponseError extends Error {
  errors: Error[]

  constructor(response: Response) {
    super(`UnknownResponseError(${response.statusText})`)
    this.errors = [
      {
        status: response.status,
        detail: response.statusText,
      },
    ]
  }
}

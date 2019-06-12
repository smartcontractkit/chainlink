export function AuthenticationError(response: Response) {
  this.errors = [
    {
      status: response.status,
      detail: response.statusText
    }
  ]
}

export function BadRequestError({ errors }) {
  this.errors = errors
}

export function ServerError(response: Response) {
  this.errors = [
    {
      status: response.status,
      detail: response.statusText
    }
  ]
}

export function UnknownResponseError(response: Response) {
  this.errors = [
    {
      status: response.status,
      detail: response.statusText
    }
  ]
}

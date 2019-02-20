export function AuthenticationError(response) {
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

export function ServerError(response) {
  this.errors = [
    {
      status: response.status,
      detail: response.statusText
    }
  ]
}

export function UnknownResponseError(response) {
  this.errors = [
    {
      status: response.status,
      detail: response.statusText
    }
  ]
}

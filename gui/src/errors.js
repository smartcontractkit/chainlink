export function AuthenticationError (response) {
  this.errors = [{
    status: response.status,
    detail: response.statusText
  }]
}

export function CreateError ({errors}) {
  this.errors = errors
}

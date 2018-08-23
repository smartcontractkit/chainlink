export function AuthenticationError (message) {
  this.message = message
}

export function CreateError ({errors}) {
  this.errors = errors
}

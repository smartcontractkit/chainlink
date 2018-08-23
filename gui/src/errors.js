export function AuthenticationError (message) {
  this.message = message
}

export function CreateError (error) {
  this.errors = error.errors
}

import React from 'react'
import { useHistory } from 'react-router-dom'
import { AuthenticationError } from '@chainlink/json-api-client'
import { setAuthentication } from 'utils/storage'

export const useErrorHandler = (): {
  error: unknown
  ErrorComponent: React.FC
  setError: React.Dispatch<unknown>
} => {
  const [error, setError] = React.useState<unknown>()
  const history = useHistory()

  React.useEffect(() => {
    if (error instanceof AuthenticationError) {
      setAuthentication({ allowed: false })
      history.push('/signin')
    }
  }, [error, history])

  const ErrorComponent: React.FC = error
    ? () => <div>Error: {JSON.stringify(error)}</div>
    : () => null

  return { error, ErrorComponent, setError }
}

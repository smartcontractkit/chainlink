import React from 'react'
import { useHistory } from 'react-router-dom'
import { AuthenticationError } from 'utils/json-api-client'
import { useDispatch } from 'react-redux'
import { receiveSignoutSuccess } from 'actionCreators'

export const useErrorHandler = (): {
  error: unknown
  ErrorComponent: React.FC
  setError: React.Dispatch<unknown>
} => {
  const [error, setError] = React.useState<unknown>()
  const history = useHistory()
  const dispatch = useDispatch()

  React.useEffect(() => {
    if (error instanceof AuthenticationError) {
      /**
       * Because sign in page is using redux to figure out whether the
       * user is logged in we need to dispatch a redux action. The reducer
       * updates the store and syncs with local storage (which is a bad
       * practice as it's a side-effect, but let's focus on solving one
       * problem at a time 😅).
       */
      dispatch(receiveSignoutSuccess())
      history.push('/signin')
    }
  }, [dispatch, error, history])

  const ErrorComponent: React.FC = error
    ? () => (
        <div>
          Error:{' '}
          {error instanceof Error ? error.message : JSON.stringify(error)}
        </div>
      )
    : () => null

  return { error, ErrorComponent, setError }
}

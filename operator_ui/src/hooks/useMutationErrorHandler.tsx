import React from 'react'
import { useHistory } from 'react-router-dom'
import { useDispatch } from 'react-redux'

import { notifyErrorMsg } from 'actionCreators'
import { ApolloError } from '@apollo/client'
import { receiveSignoutSuccess } from 'actionCreators'

// useMutationErrorHandler handles an unknown error which is caught from a
// mutation operation. If the error returned is an authentication error, it
// signs the user out and redirects them to the sign in page, otherwise it
// displays an alert with the error message.
export const useMutationErrorHandler = () => {
  const [error, setError] = React.useState<unknown>()
  const history = useHistory()
  const dispatch = useDispatch()

  React.useEffect(() => {
    if (error === null || error === undefined) {
      return
    } else if (error instanceof ApolloError) {
      // Check for an authentication error and logout
      error.graphQLErrors.forEach((err) => {
        if (err.extensions?.code == 'UNAUTHORIZED') {
          dispatch(receiveSignoutSuccess())

          history.push('/signin')

          return
        }
      })

      dispatch(notifyErrorMsg(error.message))
    } else if (error instanceof Error) {
      dispatch(notifyErrorMsg(error.message))
    } else {
      dispatch(notifyErrorMsg('An error occured'))
    }
  }, [dispatch, error, history])

  return { handleMutationError: setError }
}

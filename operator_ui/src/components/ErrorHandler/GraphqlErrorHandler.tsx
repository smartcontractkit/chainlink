import React from 'react'

import { ApolloError } from '@apollo/client'
import { useDispatch } from 'react-redux'
import { useHistory } from 'react-router-dom'

import { receiveSignoutSuccess } from 'actionCreators'

interface Props {
  error?: ApolloError
}

// GraphqlErrorHandler takes an Apollo error and renders the error message. If
// it is an authentication error, it will clear the redux cache of the auth
// credentials and redirect back to the sign in page.
//
// This is a temporary solution until we can move authentication away from using
// redux.
export const GraphqlErrorHandler: React.FC<Props> = ({ error }) => {
  const history = useHistory()
  const dispatch = useDispatch()

  if (error) {
    // Check for an authentication error
    error.graphQLErrors.forEach((err) => {
      if (err.extensions?.code == 'UNAUTHORIZED') {
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
    })

    return <div>Error: {error.message}</div>
  }

  return null
}

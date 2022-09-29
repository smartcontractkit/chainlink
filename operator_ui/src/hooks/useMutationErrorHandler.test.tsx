import React from 'react'

import { ApolloError } from '@apollo/client'
import { GraphQLError } from 'graphql'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import { getAuthentication } from 'utils/storage'

import Notifications from 'pages/Notifications'
import { useMutationErrorHandler } from './useMutationErrorHandler'

const { getByText } = screen

const StubComponent = ({ mockError }: { mockError?: unknown }) => {
  const { handleMutationError } = useMutationErrorHandler()

  React.useEffect(() => {
    handleMutationError(mockError)
  }, [mockError, handleMutationError])

  return null
}

function renderComponent(mockError?: unknown) {
  renderWithRouter(
    <>
      <Notifications />
      <Route exact path="/">
        <StubComponent mockError={mockError} />
      </Route>

      <Route exact path="/signin">
        Redirect Success
      </Route>
    </>,
  )
}

describe('useMutationErrorHandler', () => {
  it('renders an empty component if error undefined', () => {
    renderComponent()

    expect(document.documentElement).toHaveTextContent('')
  })

  it('renders the apollo error message', () => {
    const graphQLErrors = [new GraphQLError('GraphQL error')]
    const errorMessage = 'Something went wrong'
    const apolloError = new ApolloError({
      graphQLErrors,
      errorMessage,
    })

    renderComponent(apolloError)

    expect(getByText('Something went wrong')).toBeInTheDocument()
  })

  it('redirects an authenticated error', () => {
    const graphQLErrors = [
      new GraphQLError(
        'Unauthorized',
        undefined,
        undefined,
        undefined,
        undefined,
        undefined,
        { code: 'UNAUTHORIZED' },
      ),
    ]
    const errorMessage = 'Something went wrong'
    const apolloError = new ApolloError({
      graphQLErrors,
      errorMessage,
    })

    renderComponent(apolloError)

    expect(getByText('Redirect Success')).toBeInTheDocument()
    expect(getAuthentication()).toEqual({ allowed: false })
  })

  it('renders the message in an alert when it is a simple error', () => {
    renderComponent(new Error('Something went wrong'))

    expect(getByText('Something went wrong')).toBeInTheDocument()
  })

  it('renders a generic message in an alert as a default', () => {
    renderComponent('generic message') // A string type is not handled and falls to the default

    expect(getByText('An error occured')).toBeInTheDocument()
  })
})

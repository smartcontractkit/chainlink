import * as React from 'react'

import { ApolloError } from '@apollo/client'
import { GraphQLError } from 'graphql'
import { renderWithRouter, screen } from 'support/test-utils'

import { GraphqlErrorHandler } from './GraphqlErrorHandler'
import { Route } from 'react-router'

const { findByText, queryByText } = screen

function renderComponent(error?: ApolloError) {
  renderWithRouter(
    <>
      <Route exact path="/">
        <GraphqlErrorHandler error={error} />
      </Route>

      <Route exact path="/signin">
        Redirect Success
      </Route>
    </>,
  )
}

describe('GraphQLErrorHandler', () => {
  it('renders nothing when error is nil', async () => {
    renderComponent()

    expect(expect(document.documentElement).toHaveTextContent(''))
  })

  it('renders the error', async () => {
    const graphQLErrors = [
      new GraphQLError('Something went wrong with GraphQL'),
    ]
    const errorMessage = 'this is an error message'
    const apolloError = new ApolloError({
      graphQLErrors,
      errorMessage,
    })

    renderComponent(apolloError)

    expect(queryByText('Error: this is an error message')).toBeInTheDocument()
  })

  it('redirects when the error is unauthorized', async () => {
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
    const errorMessage = 'Unauthorized'
    const apolloError = new ApolloError({
      graphQLErrors,
      errorMessage,
    })

    renderComponent(apolloError)

    expect(await findByText('Redirect Success')).toBeInTheDocument()
  })
})

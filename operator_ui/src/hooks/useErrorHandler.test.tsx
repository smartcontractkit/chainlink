import React from 'react'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import { AuthenticationError } from 'utils/json-api-client'
import { getAuthentication } from 'utils/storage'

import { useErrorHandler } from './useErrorHandler'

const { getByText } = screen

const StubComponent = ({ mockError }: { mockError?: unknown }) => {
  const { ErrorComponent, setError } = useErrorHandler()

  React.useEffect(() => {
    setError(mockError)
  }, [mockError, setError])

  return <ErrorComponent />
}

describe('useErrorHandler', () => {
  it('renders an empty component if everything is fine', () => {
    renderWithRouter(<StubComponent />)

    expect(document.documentElement.textContent).toEqual('')
  })

  it('renders the error message if something goes wrong', () => {
    renderWithRouter(<StubComponent mockError="Something went wrong" />)

    expect(getByText('Error: "Something went wrong"')).toBeInTheDocument()
  })

  it('logs the user out and redirects to the signin page on authentication error', () => {
    renderWithRouter(
      <>
        <Route
          path="/some-path"
          render={() => (
            <StubComponent
              mockError={new AuthenticationError({} as Response)}
            />
          )}
        />
        <Route path="/signin">Redirect Success</Route>
      </>,
      {
        initialEntries: [`/some-path`],
      },
    )

    expect(getByText('Redirect Success')).toBeInTheDocument()
    expect(getAuthentication()).toEqual({ allowed: false })
  })

  it('Shows "Not found" message if the resource is not found', () => {
    renderWithRouter(
      <StubComponent
        mockError={{
          errors: [
            {
              code: 404,
              message: 'Not found',
            },
          ],
        }}
      />,
    )

    expect(
      getByText('Error: {"errors":[{"code":404,"message":"Not found"}]}'),
    ).toBeInTheDocument()
  })
})

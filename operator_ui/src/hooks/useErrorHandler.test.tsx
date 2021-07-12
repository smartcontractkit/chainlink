import React from 'react'
import { Route } from 'react-router-dom'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { AuthenticationError } from 'utils/json-api-client'
import { getAuthentication } from 'utils/storage'

import { useErrorHandler } from './useErrorHandler'

const StubComponent = ({ mockError }: { mockError?: unknown }) => {
  const { ErrorComponent, setError } = useErrorHandler()

  React.useEffect(() => {
    setError(mockError)
  }, [mockError, setError])

  return <ErrorComponent />
}

describe('useErrorHandler', () => {
  it('renders an empty component if everything is fine', () => {
    const wrapper = mountWithProviders(<StubComponent />)

    expect(wrapper.text()).toBe('')
  })

  it('renders the error message if something goes wrong', () => {
    const wrapper = mountWithProviders(
      <StubComponent mockError="Something went wrong" />,
    )

    expect(wrapper.text()).toContain('Error: "Something went wrong"')
  })

  it('logs the user out and redirects to the signin page on authentication error', () => {
    const wrapper = mountWithProviders(
      <Route
        path="/some-path"
        render={() => (
          <StubComponent mockError={new AuthenticationError({} as Response)} />
        )}
      />,
      {
        initialEntries: [`/some-path`],
      },
    )
    const routerCopmonentProps: any = wrapper.find('Router').props()

    expect(routerCopmonentProps?.history?.location?.pathname).toEqual('/signin')
    expect(getAuthentication()).toEqual({ allowed: false })
  })

  it('Shows "Not found" message if the resource is not found', () => {
    const wrapper = mountWithProviders(
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

    expect(wrapper.text()).toContain(
      'Error: {"errors":[{"code":404,"message":"Not found"}]}',
    )
  })
})

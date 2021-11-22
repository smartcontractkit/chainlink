import React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'

import { useQueryParams } from './useQueryParams'

const { getByText } = screen

const StubComponent = () => {
  const qp = useQueryParams()

  return <div>{qp.toString()}</div>
}

describe('useQueryParams', () => {
  it('extracts the query params', () => {
    renderWithRouter(
      <Route exact path="/">
        <StubComponent />
      </Route>,
      { initialEntries: ['/?foo=bar&baz=qux'] },
    )

    expect(getByText('foo=bar&baz=qux')).toBeInTheDocument()
  })
})

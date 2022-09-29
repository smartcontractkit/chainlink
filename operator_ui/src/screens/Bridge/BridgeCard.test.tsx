import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { BridgeCard } from './BridgeCard'
import { buildBridgePayloadFields } from 'support/factories/gql/fetchBridge'

const { getByRole, queryByText } = screen

describe('BridgeCard', () => {
  let handleDelete: jest.Mock

  function renderComponent(bridge: BridgePayload_Fields) {
    renderWithRouter(
      <>
        <Route exact path="/">
          <BridgeCard bridge={bridge} onDelete={handleDelete} />
        </Route>
        <Route path="/bridges/:id/edit">Edit Link Success</Route>
      </>,
    )
  }

  beforeEach(() => {
    handleDelete = jest.fn()
  })

  it('renders a bridge', () => {
    const bridge = buildBridgePayloadFields()

    renderComponent(bridge)

    expect(queryByText(bridge.name)).toBeInTheDocument()
    expect(queryByText(bridge.url)).toBeInTheDocument()
    expect(queryByText(bridge.outgoingToken)).toBeInTheDocument()
    expect(queryByText(bridge.confirmations)).toBeInTheDocument()
    expect(queryByText(bridge.minimumContractPayment)).toBeInTheDocument()
  })

  it('calls delete', () => {
    const bridge = buildBridgePayloadFields()

    renderComponent(bridge)

    userEvent.click(getByRole('button', { name: /open-menu/i }))
    userEvent.click(getByRole('menuitem', { name: /delete/i }))

    expect(handleDelete).toHaveBeenCalled()
  })

  it('navigates to the edit bridge page', () => {
    const bridge = buildBridgePayloadFields()

    renderComponent(bridge)

    userEvent.click(getByRole('button', { name: /open-menu/i }))
    userEvent.click(getByRole('menuitem', { name: /edit/i }))

    expect(queryByText(/edit link success/i)).toBeInTheDocument()
  })
})

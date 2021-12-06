import React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'test-utils'
import userEvent from '@testing-library/user-event'

import { BridgeView } from './BridgeView'
import { buildBridgePayloadFields } from 'support/factories/gql/fetchBridge'

const { getAllByText, getByRole, getByText } = screen

describe('BridgeView', () => {
  let handleDelete: jest.Mock

  beforeEach(() => {
    handleDelete = jest.fn()
  })

  function renderComponent(bridge: BridgePayload_Fields) {
    renderWithRouter(
      <>
        <Route exact path="/bridges/:id">
          <BridgeView bridge={bridge} onDelete={handleDelete} />
        </Route>

        <Route path="/bridges/:id/edit">Edit Page</Route>
      </>,
      { initialEntries: [`/bridges/bridge-api`] },
    )
  }

  it('renders the details of the bridge spec', async () => {
    renderComponent(buildBridgePayloadFields())

    expect(getByText('Name')).toBeInTheDocument()
    expect(getByText('URL')).toBeInTheDocument()
    expect(getByText('Confirmations')).toBeInTheDocument()
    expect(getByText('Min. Contract Payment')).toBeInTheDocument()
    expect(getByText('Outgoing Token')).toBeInTheDocument()

    expect(getAllByText('bridge-api')).toHaveLength(2) // In the heading and details
    expect(getByText('http://bridge.com')).toBeInTheDocument()
    expect(getByText(1)).toBeInTheDocument()
    expect(getByText(0)).toBeInTheDocument()
    expect(getByText('outgoing1')).toBeInTheDocument()
  })

  it('links to the edit page', async () => {
    renderComponent(buildBridgePayloadFields())

    userEvent.click(getByRole('button', { name: /open-menu/i }))
    userEvent.click(getByRole('menuitem', { name: /edit/i }))

    expect(getByText('Edit Page')).toBeInTheDocument()
  })

  it('handles a bridge delete', async () => {
    renderComponent(buildBridgePayloadFields())

    userEvent.click(getByRole('button', { name: /open-menu/i }))
    userEvent.click(getByRole('menuitem', { name: /delete/i }))
    userEvent.click(getByRole('button', { name: /confirm/i }))

    expect(handleDelete).toHaveBeenCalled()
  })
})

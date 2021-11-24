import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { EditBridgeView } from './EditBridgeView'
import { buildBridgePayloadFields } from 'support/factories/gql/fetchBridge'

const { findByText, getByRole, getByTestId, getByText } = screen

function renderComponent(bridge: BridgePayload_Fields) {
  const handleSubmit = jest.fn()

  renderWithRouter(
    <>
      <Route exact path="/bridges/:id/edit">
        <EditBridgeView bridge={bridge} onSubmit={handleSubmit} />
      </Route>

      <Route exact path="/bridges/:id">
        Link Success
      </Route>
    </>,
    { initialEntries: ['/bridges/1/edit'] },
  )
}

describe('EditBridgeView', () => {
  it('renders with initial values', async () => {
    const bridge = buildBridgePayloadFields()

    renderComponent(bridge)

    expect(getByText('Edit Bridge')).toBeInTheDocument()
    expect(getByTestId('bridge-form')).toHaveFormValues({
      name: 'bridge-api',
      url: 'http://bridge.com',
      minimumContractPayment: '0',
      confirmations: 1,
    })
  })

  it('return to the show page', async () => {
    const bridge = buildBridgePayloadFields()

    renderComponent(bridge)

    userEvent.click(getByRole('link', { name: /cancel/i }))

    expect(await findByText('Link Success')).toBeInTheDocument()
  })
})

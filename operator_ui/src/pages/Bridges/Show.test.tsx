import React from 'react'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'test-utils'
import userEvent from '@testing-library/user-event'
import globPath from 'test-helpers/globPath'
import Notifications from 'pages/Notifications'

import { Show } from 'pages/Bridges/Show'

const { findByText, getByRole, getByText } = screen

describe('pages/Bridges/Show', () => {
  it('renders the details of the bridge spec', async () => {
    const response = {
      data: {
        id: 'tallbridge',
        type: 'bridges',
        attributes: {
          name: 'Tall Bridge',
          url: 'https://localhost.com:712/endpoint',
          confirmations: 9,
          outgoingToken: 'outgoingToken',
        },
      },
    }

    global.fetch.getOnce(globPath(`/v2/bridge_types/tallbridge`), response)

    renderWithRouter(
      <Route path="/bridges/:bridgeId">
        <Show />
      </Route>,
      { initialEntries: [`/bridges/tallbridge`] },
    )

    await findByText('Bridge Info')

    expect(getByText('Tall Bridge')).toBeInTheDocument()
    expect(getByText('9')).toBeInTheDocument()
    expect(getByText('https://localhost.com:712/endpoint')).toBeInTheDocument()
    expect(getByText('outgoingToken')).toBeInTheDocument()
    expect(getByText('Tall Bridge')).toBeInTheDocument()
  })

  it('deletes a bridge', async () => {
    const response = {
      data: {
        id: 'tallbridge',
        type: 'bridges',
        attributes: {
          name: 'Tall Bridge',
          url: 'https://localhost.com:712/endpoint',
          confirmations: 9,
          outgoingToken: 'outgoingToken',
        },
      },
    }

    global.fetch.getOnce(globPath(`/v2/bridge_types/tallbridge`), response)

    renderWithRouter(
      <>
        <Notifications />
        <Route path="/bridges/:bridgeId">
          <Show />
        </Route>

        <Route exact path="/bridges">
          Redirect Success
        </Route>
      </>,
      { initialEntries: [`/bridges/tallbridge`] },
    )

    await findByText('Bridge Info')

    userEvent.click(getByRole('button', { name: /Delete/i }))

    await findByText('Confirm')

    global.fetch.deleteOnce(globPath(`/v2/bridge_types/tallbridge`), {})

    userEvent.click(getByRole('button', { name: /Confirm/i }))

    expect(await findByText('Redirect Success')).toBeInTheDocument()
  })

  it('fails to delete a bridge', async () => {
    const response = {
      data: {
        id: 'tallbridge',
        type: 'bridges',
        attributes: {
          name: 'Tall Bridge',
          url: 'https://localhost.com:712/endpoint',
          confirmations: 9,
          outgoingToken: 'outgoingToken',
        },
      },
    }

    global.fetch.getOnce(globPath(`/v2/bridge_types/tallbridge`), response)

    renderWithRouter(
      <>
        <Notifications />
        <Route path="/bridges/:bridgeId">
          <Show />
        </Route>

        <Route exact path="/bridges">
          Redirect Success
        </Route>
      </>,
      { initialEntries: [`/bridges/tallbridge`] },
    )

    await findByText('Bridge Info')

    userEvent.click(getByRole('button', { name: /Delete/i }))

    await findByText('Confirm')

    global.fetch.deleteOnce(globPath(`/v2/bridge_types/tallbridge`), {
      body: {
        errors: [
          {
            detail: "can't remove the bridge",
          },
        ],
      },
      status: 409,
    })

    userEvent.click(getByRole('button', { name: /Confirm/i }))

    expect(await findByText("can't remove the bridge")).toBeInTheDocument()
  })
})

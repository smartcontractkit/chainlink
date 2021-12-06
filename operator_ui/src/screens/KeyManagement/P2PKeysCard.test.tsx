import * as React from 'react'

import { render, screen, waitForElementToBeRemoved } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { buildP2PKey, buildP2PKeys } from 'support/factories/gql/fetchP2PKeys'
import { P2PKeysCard, Props as P2PKeysProps } from './P2PKeysCard'

const { getAllByRole, getByRole, queryByRole, queryByText } = screen

function renderComponent(cardProps: P2PKeysProps) {
  render(<P2PKeysCard {...cardProps} />)
}

describe('CSAKeysCard', () => {
  let promise: Promise<any>
  let handleCreate: jest.Mock
  let handleDelete: jest.Mock

  beforeEach(() => {
    promise = Promise.resolve()
    handleCreate = jest.fn()
    handleDelete = jest.fn(() => promise)
  })

  it('renders the p2p keys', () => {
    const p2pKeys = buildP2PKeys()

    renderComponent({
      loading: false,
      data: {
        p2pKeys: {
          results: p2pKeys,
        },
      },
      onCreate: handleCreate,
      onDelete: handleDelete,
    })

    expect(getAllByRole('row')).toHaveLength(3)

    expect(queryByText(`Peer ID: ${p2pKeys[0].peerID}`)).toBeInTheDocument()
    expect(queryByText(`Peer ID: ${p2pKeys[1].peerID}`)).toBeInTheDocument()
  })

  it('renders no content', () => {
    renderComponent({
      loading: false,
      data: {
        p2pKeys: {
          results: [],
        },
      },
      onCreate: handleCreate,
      onDelete: handleDelete,
    })

    expect(queryByText('No entries to show')).toBeInTheDocument()
  })

  it('renders a loading spinner', () => {
    renderComponent({
      loading: true,
      onCreate: handleCreate,
      onDelete: handleDelete,
    })

    expect(queryByRole('progressbar')).toBeInTheDocument()
  })

  it('renders an error message', () => {
    renderComponent({
      loading: false,
      errorMsg: 'error message',
      onCreate: handleCreate,
      onDelete: handleDelete,
    })

    expect(queryByText('error message')).toBeInTheDocument()
  })

  it('calls onCreate', () => {
    renderComponent({
      loading: false,
      data: {
        p2pKeys: {
          results: [],
        },
      },
      onCreate: handleCreate,
      onDelete: handleDelete,
    })

    userEvent.click(getByRole('button', { name: /new p2p key/i }))

    expect(handleCreate).toHaveBeenCalled()
  })

  it('calls onDelete', async () => {
    const bundle = buildP2PKey()
    renderComponent({
      loading: false,
      data: {
        p2pKeys: {
          results: [bundle],
        },
      },
      onCreate: handleCreate,
      onDelete: handleDelete,
    })

    userEvent.click(getByRole('button', { name: /delete/i }))
    expect(queryByText(bundle.id)).toBeInTheDocument()

    userEvent.click(getByRole('button', { name: /confirm/i }))

    await waitForElementToBeRemoved(getByRole('dialog'))

    expect(handleDelete).toHaveBeenCalled()
  })
})

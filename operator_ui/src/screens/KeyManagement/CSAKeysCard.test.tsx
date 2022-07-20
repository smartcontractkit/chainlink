import * as React from 'react'

import { render, screen } from 'support/test-utils'

import { buildCSAKeys } from 'support/factories/gql/fetchCSAKeys'
import { CSAKeysCard, Props as CSAKeysCardProps } from './CSAKeysCard'
import userEvent from '@testing-library/user-event'

const { getByRole, queryByRole, queryByText } = screen

function renderComponent(cardProps: CSAKeysCardProps) {
  render(<CSAKeysCard {...cardProps} />)
}

describe('CSAKeysCard', () => {
  let handleCreate: jest.Mock

  beforeEach(() => {
    handleCreate = jest.fn()
  })

  it('renders the keys', () => {
    const csaKeys = buildCSAKeys()

    renderComponent({
      loading: false,
      data: {
        csaKeys: {
          results: csaKeys,
        },
      },
      onCreate: handleCreate,
    })

    expect(queryByText(csaKeys[0].publicKey)).toBeInTheDocument()
    expect(queryByText(csaKeys[1].publicKey)).toBeInTheDocument()

    // Button should not appear when there are keys present
    expect(queryByRole('button', { name: /new csa key/i })).toBeNull()
  })

  it('renders no content', () => {
    renderComponent({
      loading: false,
      data: {
        csaKeys: {
          results: [],
        },
      },
      onCreate: handleCreate,
    })

    expect(queryByText('No entries to show')).toBeInTheDocument()

    // Button should appear when there are no keys
    expect(queryByRole('button', { name: /new csa key/i })).toBeInTheDocument()
  })

  it('renders a loading spinner', () => {
    renderComponent({
      loading: true,
      onCreate: handleCreate,
    })

    expect(queryByRole('progressbar')).toBeInTheDocument()
  })

  it('renders an error message', () => {
    renderComponent({
      loading: false,
      errorMsg: 'error message',
      onCreate: handleCreate,
    })

    expect(queryByText('error message')).toBeInTheDocument()
  })

  it('calls onCreate', () => {
    renderComponent({
      loading: false,
      data: {
        csaKeys: {
          results: [],
        },
      },
      onCreate: handleCreate,
    })

    userEvent.click(getByRole('button', { name: /new csa key/i }))

    expect(handleCreate).toHaveBeenCalled()
  })
})

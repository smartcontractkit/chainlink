import React from 'react'

import { render, screen } from '@testing-library/react'

import {
  KeyValueListCard,
  Props as KeyValueListCardProps,
} from './KeyValueListCard'

const { queryByText } = screen

const renderComponent = (props: KeyValueListCardProps) =>
  render(<KeyValueListCard {...props} />)

describe('KeyValueListCard', () => {
  it('renders a key value', () => {
    renderComponent({
      loading: false,
      entries: [['FOO', 'bar']],
    })

    expect(queryByText('FOO')).toBeInTheDocument()
    expect(queryByText('bar')).toBeInTheDocument()

    // No table header
    expect(queryByText('Key')).toBeNull()
    expect(queryByText('Value')).toBeNull()
  })

  it('renders loading', () => {
    renderComponent({
      loading: true,
      entries: [],
    })

    expect(queryByText('...')).toBeInTheDocument()
  })

  it('renders error message', () => {
    renderComponent({
      loading: false,
      entries: [],
      error: 'Error!',
    })

    expect(queryByText('Error!')).toBeInTheDocument()
  })

  it('displays the table header', () => {
    renderComponent({
      loading: false,
      entries: [['FOO', 'bar']],
      showHead: true,
    })

    expect(queryByText('Key')).toBeInTheDocument()
    expect(queryByText('Value')).toBeInTheDocument()
  })
})

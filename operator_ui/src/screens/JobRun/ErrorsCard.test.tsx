import * as React from 'react'

import { render, screen } from 'support/test-utils'

import { ErrorsCard } from './ErrorsCard'

const { queryByText } = screen

describe('ErrorsCard', () => {
  it('renders the errors', async () => {
    render(<ErrorsCard errors={['error 1', 'error 2']} />)

    expect(queryByText('Errors')).toBeInTheDocument()
    expect(queryByText('error 1')).toBeInTheDocument()
    expect(queryByText('error 2')).toBeInTheDocument()
  })
})

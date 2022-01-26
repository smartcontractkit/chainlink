import React from 'react'
import { render, screen } from 'support/test-utils'
import TokenBalanceCard from '../../../src/components/Cards/TokenBalance'

const { queryByText } = screen

describe('components/Cards/TokenBalance', () => {
  it('renders the title and a loading indicator when it is fetching', () => {
    render(<TokenBalanceCard title="Ether Balance" value={undefined} />)

    expect(queryByText('Ether Balance')).toBeInTheDocument()
    expect(queryByText('...')).toBeInTheDocument()
  })

  it('renders the title and the error message', () => {
    render(<TokenBalanceCard title="Ether Balance" error="An Error" />)

    expect(queryByText('Ether Balance')).toBeInTheDocument()
    expect(queryByText('An Error')).toBeInTheDocument()
  })

  it('renders the title and the formatted balance', () => {
    render(
      <TokenBalanceCard title="Ether Balance" value="7779070000000000000000" />,
    )

    expect(queryByText('Ether Balance')).toBeInTheDocument()
    expect(queryByText('7.779070k')).toBeInTheDocument()
  })
})

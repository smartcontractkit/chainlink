import * as React from 'react'

import userEvent from '@testing-library/user-event'
import { renderWithRouter, screen } from 'support/test-utils'

import { buildJobProposals } from 'support/factories/gql/fetchFeedsManagersWithProposals'
import { JobProposalsCard } from './JobProposalsCard'

const { findAllByRole, getByRole } = screen

describe('JobProposalsCard', () => {
  it('renders the pending job proposals ', async () => {
    const proposals = buildJobProposals()

    renderWithRouter(<JobProposalsCard proposals={proposals} />)

    const rows = await findAllByRole('row')
    expect(rows).toHaveLength(2)

    expect(rows[1]).toHaveTextContent('1')
    expect(rows[1]).toHaveTextContent('N/A')
    expect(rows[1]).toHaveTextContent('1 minute ago')
  })

  it('renders the approved job proposals', async () => {
    const proposals = buildJobProposals()

    renderWithRouter(<JobProposalsCard proposals={proposals} />)

    userEvent.click(getByRole('tab', { name: /approved/i }))

    const rows = await findAllByRole('row')
    expect(rows).toHaveLength(2)

    expect(rows[1]).toHaveTextContent('2')
    expect(rows[1]).toHaveTextContent('00000000-0000-0000-0000-000000000002')
    expect(rows[1]).toHaveTextContent('1 minute ago')
  })

  it('renders the rejected job proposals', async () => {
    const proposals = buildJobProposals()

    renderWithRouter(<JobProposalsCard proposals={proposals} />)

    userEvent.click(getByRole('tab', { name: /rejected/i }))

    const rows = await findAllByRole('row')
    expect(rows).toHaveLength(2)

    expect(rows[1]).toHaveTextContent('3')
    expect(rows[1]).toHaveTextContent('N/A')
    expect(rows[1]).toHaveTextContent('1 minute ago')
  })

  it('renders the rejected job proposals', async () => {
    const proposals = buildJobProposals()

    renderWithRouter(<JobProposalsCard proposals={proposals} />)

    userEvent.click(getByRole('tab', { name: /cancelled/i }))

    const rows = await findAllByRole('row')
    expect(rows).toHaveLength(2)

    expect(rows[1]).toHaveTextContent('4')
    expect(rows[1]).toHaveTextContent('N/A')
    expect(rows[1]).toHaveTextContent('1 minute ago')
  })
})

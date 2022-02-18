import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { PendingTable } from './PendingTable'
import { buildPendingJobProposal } from 'support/factories/gql/fetchFeedsManagersWithProposals'

const { getByRole, queryByText } = screen

function renderComponent(proposals: FeedsManager_JobProposalsFields[]) {
  renderWithRouter(
    <>
      <Route path="/">
        <PendingTable proposals={proposals} />
      </Route>
      <Route path="/job_proposals/:id">Redirect Success</Route>
    </>,
  )
}

describe('PendingTable', () => {
  it('renders the table', () => {
    const proposal = buildPendingJobProposal()

    renderComponent([proposal])

    expect(queryByText('ID')).toBeInTheDocument()
    expect(queryByText('Last Proposed')).toBeInTheDocument()

    expect(queryByText(proposal.id)).toBeInTheDocument()
    expect(queryByText('1 minute ago')).toBeInTheDocument()
  })

  it('navigates to edit', () => {
    const proposal = buildPendingJobProposal()

    renderComponent([proposal])

    userEvent.click(getByRole('link', { name: proposal.id }))

    expect(queryByText('Redirect Success')).toBeInTheDocument()
  })
})

import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { ApprovedTable } from './ApprovedTable'
import { buildApprovedJobProposal } from 'support/factories/gql/fetchFeedsManagersWithProposals'

const { getByRole, queryByText } = screen

function renderComponent(proposals: FeedsManager_JobProposalsFields[]) {
  renderWithRouter(
    <>
      <Route path="/">
        <ApprovedTable proposals={proposals} />
      </Route>
      <Route path="/job_proposals/:id">Redirect Success</Route>
    </>,
  )
}

describe('ApprovedTable', () => {
  it('renders the table', () => {
    const proposal = buildApprovedJobProposal()

    renderComponent([proposal])

    expect(queryByText('ID')).toBeInTheDocument()
    expect(queryByText('External Job ID')).toBeInTheDocument()
    expect(queryByText('Latest Version')).toBeInTheDocument()
    expect(queryByText('Last Proposed')).toBeInTheDocument()

    expect(queryByText(proposal.id)).toBeInTheDocument()
    expect(queryByText(proposal.externalJobID as string)).toBeInTheDocument()
    expect(queryByText(proposal.latestSpec.version)).toBeInTheDocument()
    expect(queryByText('1 minute ago')).toBeInTheDocument()

    expect(queryByText('Update available')).toBeNull()
  })

  it('displays an update available', () => {
    const proposal = buildApprovedJobProposal({ pendingUpdate: true })

    renderComponent([proposal])

    expect(queryByText('Update available')).toBeNull()
  })

  it('navigates to edit', () => {
    const proposal = buildApprovedJobProposal()

    renderComponent([proposal])

    userEvent.click(getByRole('link', { name: proposal.id }))

    expect(queryByText('Redirect Success')).toBeInTheDocument()
  })
})

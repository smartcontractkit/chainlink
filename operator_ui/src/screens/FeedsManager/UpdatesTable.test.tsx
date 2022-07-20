import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { UpdatesTable } from './UpdatesTable'
import {
  buildApprovedJobProposal,
  buildRejectedJobProposal,
} from 'support/factories/gql/fetchFeedsManagersWithProposals'

const { getAllByRole, getByRole, queryByText } = screen

function renderComponent(proposals: FeedsManager_JobProposalsFields[]) {
  renderWithRouter(
    <>
      <Route path="/">
        <UpdatesTable proposals={proposals} />
      </Route>
      <Route path="/job_proposals/:id">Redirect Success</Route>
    </>,
  )
}

describe('UpdatesTable', () => {
  it('renders the table', () => {
    const approvedProposal = buildApprovedJobProposal({ pendingUpdate: true })
    const rejectedProposal = buildRejectedJobProposal({ pendingUpdate: true })

    renderComponent([approvedProposal, rejectedProposal])

    expect(queryByText('ID')).toBeInTheDocument()
    expect(queryByText('External Job ID')).toBeInTheDocument()
    expect(queryByText('Latest Version')).toBeInTheDocument()
    expect(queryByText('Last Proposed')).toBeInTheDocument()

    const rows = getAllByRole('row')
    expect(rows).toHaveLength(3)

    expect(rows[1]).toHaveTextContent(approvedProposal.id)
    expect(rows[1]).toHaveTextContent(approvedProposal.externalJobID as string)
    expect(rows[1]).toHaveTextContent(
      approvedProposal.latestSpec.version.toString(),
    )
    expect(rows[1]).toHaveTextContent('1 minute ago')

    expect(rows[2]).toHaveTextContent(rejectedProposal.id)
    expect(rows[2]).toHaveTextContent('--')
    expect(rows[2]).toHaveTextContent(
      rejectedProposal.latestSpec.version.toString(),
    )
    expect(rows[2]).toHaveTextContent('1 minute ago')
  })

  it('navigates to edit', () => {
    const proposal = buildApprovedJobProposal({ pendingUpdate: true })

    renderComponent([proposal])

    userEvent.click(getByRole('link', { name: proposal.id }))

    expect(queryByText('Redirect Success')).toBeInTheDocument()
  })
})

import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'

import { JobProposalCard } from './JobProposalCard'
import {
  buildJobProposal,
  buildJobProposalSpec,
} from 'support/factories/gql/fetchJobProposal'

const { queryAllByText, queryByText } = screen

describe('JobProposalCard', () => {
  function renderComponent(proposal: JobProposalPayloadFields) {
    renderWithRouter(
      <>
        <Route exact path="/">
          <table>
            <tbody>
              <JobProposalCard proposal={proposal} />
            </tbody>
          </table>
        </Route>

        <Route exact path="/jobs/:id">
          Run Page
        </Route>
      </>,
    )
  }

  it('renders a pending job proposal card', () => {
    const proposal = buildJobProposal()

    renderComponent(proposal)

    expect(queryByText('Status')).toBeInTheDocument()
    expect(queryByText('FMS UUID')).toBeInTheDocument()
    expect(queryByText('External Job ID')).toBeInTheDocument()
    expect(queryByText('Approved Version')).toBeInTheDocument()

    expect(queryByText('Pending')).toBeInTheDocument()
    expect(queryByText(proposal.remoteUUID)).toBeInTheDocument()
    expect(queryAllByText('--')).toHaveLength(2)
  })

  it('renders a pending job proposal card with an approved spec', () => {
    const proposal = buildJobProposal({
      externalJobID: '10000000-0000-0000-0000-000000000001',
      jobID: '1',
      specs: [
        buildJobProposalSpec({
          status: 'APPROVED',
          version: 2,
        }),
      ],
    })

    renderComponent(proposal)

    expect(queryByText('Status')).toBeInTheDocument()
    expect(queryByText('FMS UUID')).toBeInTheDocument()
    expect(queryByText('External Job ID')).toBeInTheDocument()
    expect(queryByText('Approved Version')).toBeInTheDocument()

    expect(queryByText('Pending')).toBeInTheDocument()
    expect(queryByText(proposal.externalJobID as string)).toBeInTheDocument()
    expect(queryByText(proposal.remoteUUID)).toBeInTheDocument()
    expect(queryByText('2')).toBeInTheDocument()
  })
})

import * as React from 'react'
import { render, screen, waitFor } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import {
  buildJobProposal,
  buildJobProposalSpec,
} from 'support/factories/gql/fetchJobProposal'

import { JobProposalView } from './JobProposalView'

const { findByRole, getByRole, getByTestId, queryByText } = screen

describe('JobProposalView', () => {
  let handleUpdateSpec: jest.Mock
  let handleApprove: jest.Mock
  let handleCancel: jest.Mock
  let handleReject: jest.Mock

  function renderComponent(proposal: JobProposalPayloadFields) {
    render(
      <JobProposalView
        proposal={proposal}
        onUpdateSpec={handleUpdateSpec}
        onApprove={handleApprove}
        onCancel={handleCancel}
        onReject={handleReject}
      />,
    )
  }

  beforeEach(() => {
    handleUpdateSpec = jest.fn()
    handleApprove = jest.fn()
    handleCancel = jest.fn()
    handleReject = jest.fn()
  })

  describe('pending proposal', () => {
    let proposal: JobProposalPayloadFields

    beforeEach(() => {
      proposal = buildJobProposal({ status: 'PENDING' })
    })

    it('renders a pending job proposal', async () => {
      renderComponent(proposal)

      expect(queryByText('Job Proposal #1')).toBeInTheDocument()
      expect(queryByText('Pending')).toBeInTheDocument()

      expect(getByTestId('codeblock')).toHaveTextContent(
        proposal.specs[0].definition,
      )
      expect(queryByText(/edit/i)).toBeInTheDocument()
    })

    it('approves the proposal', async () => {
      renderComponent(proposal)

      userEvent.click(getByRole('button', { name: /approve/i }))
      userEvent.click(getByRole('button', { name: /confirm/i }))
      await waitFor(() => expect(handleApprove).toHaveBeenCalled())
    })

    it('rejects the propoosal', async () => {
      renderComponent(proposal)

      userEvent.click(getByRole('button', { name: /reject/i }))
      userEvent.click(getByRole('button', { name: /confirm/i }))
      await waitFor(() => expect(handleReject).toHaveBeenCalled())
    })

    it('updates the spec', async () => {
      renderComponent(proposal)

      userEvent.click(getByRole('button', { name: /edit/i }))

      expect(await findByRole('heading', { name: /edit job spec/i }))

      userEvent.click(getByRole('button', { name: /submit/i }))

      await waitFor(() => expect(handleUpdateSpec).toHaveBeenCalled())
    })
  })

  describe('approved proposal', () => {
    let proposal: JobProposalPayloadFields

    beforeEach(() => {
      proposal = buildJobProposal({
        status: 'APPROVED',
        specs: [buildJobProposalSpec({ status: 'APPROVED' })],
      })
    })

    it('renders an approved job proposal', async () => {
      renderComponent(proposal)

      expect(queryByText('Approved')).toBeInTheDocument()

      expect(getByTestId('codeblock')).toHaveTextContent(
        proposal.specs[0].definition,
      )
      expect(queryByText(/edit/i)).toBeNull()
    })

    it('cancels the proposal', async () => {
      renderComponent(proposal)

      userEvent.click(getByRole('button', { name: /cancel/i }))
      userEvent.click(getByRole('button', { name: /confirm/i }))
      await waitFor(() => expect(handleCancel).toHaveBeenCalled())
    })
  })

  describe('cancelled proposal', () => {
    let proposal: JobProposalPayloadFields

    beforeEach(() => {
      proposal = buildJobProposal({
        status: 'CANCELLED',
        specs: [buildJobProposalSpec({ status: 'CANCELLED' })],
      })
    })

    it('renders a cancelled job proposal', async () => {
      renderComponent(proposal)

      expect(queryByText('Cancelled')).toBeInTheDocument()

      expect(getByTestId('codeblock')).toHaveTextContent(
        proposal.specs[0].definition,
      )
      expect(queryByText(/edit/i)).toBeInTheDocument()
    })

    it('approves the proposal', async () => {
      renderComponent(proposal)

      userEvent.click(getByRole('button', { name: /approve/i }))
      userEvent.click(getByRole('button', { name: /confirm/i }))
      await waitFor(() => expect(handleApprove).toHaveBeenCalled())
    })

    it('updates the spec', async () => {
      renderComponent(proposal)

      userEvent.click(getByRole('button', { name: /edit/i }))

      expect(await findByRole('heading', { name: /edit job spec/i }))

      userEvent.click(getByRole('button', { name: /submit/i }))

      await waitFor(() => expect(handleUpdateSpec).toHaveBeenCalled())
    })
  })

  describe('rejected proposal', () => {
    let proposal: JobProposalPayloadFields

    beforeEach(() => {
      proposal = buildJobProposal({
        status: 'REJECTED',
        specs: [buildJobProposalSpec({ status: 'REJECTED' })],
      })
    })

    it('renders a rejected job proposal', async () => {
      renderComponent(proposal)

      expect(queryByText('Rejected')).toBeInTheDocument()

      expect(getByTestId('codeblock')).toHaveTextContent(
        proposal.specs[0].definition,
      )
      expect(queryByText(/edit/i)).toBeNull()
    })
  })
})

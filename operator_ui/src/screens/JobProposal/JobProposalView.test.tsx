import userEvent from '@testing-library/user-event'
import * as React from 'react'
import { buildJobProposal } from 'support/factories/gql/fetchJobProposal'
import { render, screen, waitFor } from 'support/test-utils'

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

      expect(queryByText('Job proposal #1')).toBeInTheDocument()
      expect(queryByText('Status: Pending')).toBeInTheDocument()

      const codeblock = getByTestId('codeblock')
      expect(codeblock).toHaveTextContent(proposal.spec)
      expect(queryByText(/edit job spec/i)).toBeInTheDocument()
    })

    it('approves the propoosal', async () => {
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

    it('opens the edit modal', async () => {
      renderComponent(proposal)

      userEvent.click(getByRole('button', { name: /edit job spec/i }))

      expect(await findByRole('heading', { name: /edit job spec/i }))
    })
  })

  describe('approved proposal', () => {
    let proposal: JobProposalPayloadFields

    beforeEach(() => {
      proposal = buildJobProposal({ status: 'APPROVED' })
    })

    it('renders an approved job proposal', async () => {
      renderComponent(proposal)

      expect(queryByText('Job proposal #1')).toBeInTheDocument()
      expect(queryByText('Status: Approved')).toBeInTheDocument()

      const codeblock = getByTestId('codeblock')
      expect(codeblock).toHaveTextContent(proposal.spec)
      expect(queryByText(/edit job spec/i)).toBeNull()
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
      proposal = buildJobProposal({ status: 'CANCELLED' })
    })

    it('renders a cancelled job proposal', async () => {
      renderComponent(proposal)

      expect(queryByText('Job proposal #1')).toBeInTheDocument()
      expect(queryByText('Status: Cancelled')).toBeInTheDocument()

      const codeblock = getByTestId('codeblock')
      expect(codeblock).toHaveTextContent(proposal.spec)
      expect(queryByText(/edit job spec/i)).toBeInTheDocument()
    })

    it('approves the proposal', async () => {
      renderComponent(proposal)

      userEvent.click(getByRole('button', { name: /approve/i }))
      userEvent.click(getByRole('button', { name: /confirm/i }))
      await waitFor(() => expect(handleApprove).toHaveBeenCalled())
    })

    it('opens the edit modal', async () => {
      renderComponent(proposal)

      userEvent.click(getByRole('button', { name: /edit job spec/i }))

      expect(await findByRole('heading', { name: /edit job spec/i }))
    })
  })

  describe('rejected proposal', () => {
    let proposal: JobProposalPayloadFields

    beforeEach(() => {
      proposal = buildJobProposal({ status: 'REJECTED' })
    })

    it('renders a rejected job proposal', async () => {
      renderComponent(proposal)

      expect(queryByText('Job proposal #1')).toBeInTheDocument()
      expect(queryByText('Status: Rejected')).toBeInTheDocument()

      const codeblock = getByTestId('codeblock')
      expect(codeblock).toHaveTextContent(proposal.spec)
      expect(queryByText(/edit job spec/i)).toBeNull()
    })
  })
})

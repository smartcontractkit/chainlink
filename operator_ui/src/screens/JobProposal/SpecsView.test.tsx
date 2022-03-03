import * as React from 'react'
import { render, screen, waitFor } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { buildJobProposalSpec } from 'support/factories/gql/fetchJobProposal'

import { SpecsView } from './SpecsView'

const { findByRole, getByRole, getByTestId, queryByText } = screen

describe('SpecsView', () => {
  let handleUpdateSpec: jest.Mock
  let handleApprove: jest.Mock
  let handleCancel: jest.Mock
  let handleReject: jest.Mock

  function renderComponent(specs: ReadonlyArray<JobProposal_SpecsFields>) {
    render(
      <SpecsView
        specs={specs}
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
    let specs: ReadonlyArray<JobProposal_SpecsFields>

    beforeEach(() => {
      specs = [buildJobProposalSpec({ status: 'PENDING' })]
    })

    it('renders a pending job proposal', async () => {
      renderComponent(specs)

      expect(getByTestId('codeblock')).toHaveTextContent(specs[0].definition)
      expect(queryByText(/edit/i)).toBeInTheDocument()
    })

    it('approves the proposal', async () => {
      renderComponent(specs)

      userEvent.click(getByRole('button', { name: /approve/i }))
      userEvent.click(getByRole('button', { name: /confirm/i }))
      await waitFor(() => expect(handleApprove).toHaveBeenCalled())
    })

    it('rejects the propoosal', async () => {
      renderComponent(specs)

      userEvent.click(getByRole('button', { name: /reject/i }))
      userEvent.click(getByRole('button', { name: /confirm/i }))
      await waitFor(() => expect(handleReject).toHaveBeenCalled())
    })

    it('updates the spec', async () => {
      renderComponent(specs)

      userEvent.click(getByRole('button', { name: /edit/i }))

      expect(await findByRole('heading', { name: /edit job spec/i }))

      userEvent.click(getByRole('button', { name: /submit/i }))

      await waitFor(() => expect(handleUpdateSpec).toHaveBeenCalled())
    })
  })

  describe('approved proposal', () => {
    let specs: ReadonlyArray<JobProposal_SpecsFields>

    beforeEach(() => {
      specs = [buildJobProposalSpec({ status: 'APPROVED' })]
    })

    it('renders an approved job proposal', async () => {
      renderComponent(specs)

      expect(getByTestId('codeblock')).toHaveTextContent(specs[0].definition)
      expect(queryByText(/edit/i)).toBeNull()
    })

    it('cancels the proposal', async () => {
      renderComponent(specs)

      userEvent.click(getByRole('button', { name: /cancel/i }))
      userEvent.click(getByRole('button', { name: /confirm/i }))
      await waitFor(() => expect(handleCancel).toHaveBeenCalled())
    })
  })

  describe('cancelled proposal', () => {
    let specs: ReadonlyArray<JobProposal_SpecsFields>

    beforeEach(() => {
      specs = [buildJobProposalSpec({ status: 'CANCELLED' })]
    })

    it('renders a cancelled job proposal', async () => {
      renderComponent(specs)

      expect(getByTestId('codeblock')).toHaveTextContent(specs[0].definition)
      expect(queryByText(/edit/i)).toBeInTheDocument()
    })

    it('approves the proposal', async () => {
      renderComponent(specs)

      userEvent.click(getByRole('button', { name: /approve/i }))
      userEvent.click(getByRole('button', { name: /confirm/i }))
      await waitFor(() => expect(handleApprove).toHaveBeenCalled())
    })

    it('updates the spec', async () => {
      renderComponent(specs)

      userEvent.click(getByRole('button', { name: /edit/i }))

      expect(await findByRole('heading', { name: /edit job spec/i }))

      userEvent.click(getByRole('button', { name: /submit/i }))

      await waitFor(() => expect(handleUpdateSpec).toHaveBeenCalled())
    })
  })

  describe('rejected proposal', () => {
    let specs: ReadonlyArray<JobProposal_SpecsFields>

    beforeEach(() => {
      specs = [buildJobProposalSpec({ status: 'REJECTED' })]
    })

    it('renders a rejected job proposal', async () => {
      renderComponent(specs)

      expect(getByTestId('codeblock')).toHaveTextContent(specs[0].definition)
      expect(queryByText(/edit/i)).toBeNull()
    })
  })
})

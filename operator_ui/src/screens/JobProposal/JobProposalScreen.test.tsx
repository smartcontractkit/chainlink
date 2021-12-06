import * as React from 'react'

import { MockedProvider, MockedResponse } from '@apollo/client/testing'
import { GraphQLError } from 'graphql'
import {
  renderWithRouter,
  screen,
  waitForElementToBeRemoved,
} from 'support/test-utils'
import userEvent from '@testing-library/user-event'
import { Route } from 'react-router-dom'

import { buildJobProposal } from 'support/factories/gql/fetchJobProposal'
import {
  APPROVE_JOB_PROPOSAL_MUTATION,
  CANCEL_JOB_PROPOSAL_MUTATION,
  JOB_PROPOSAL_QUERY,
  REJECT_JOB_PROPOSAL_MUTATION,
  UPDATE_JOB_PROPOSAL_SPEC_MUTATION,
  JobProposalScreen,
} from './JobProposalScreen'
import Notifications from 'pages/Notifications'

const { findByTestId, findByText, getByRole, queryByRole } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <>
      <Notifications />
      <Route exact path="/job_proposals/:id">
        <MockedProvider mocks={mocks} addTypename={false}>
          <JobProposalScreen />
        </MockedProvider>
      </Route>
    </>,
    { initialEntries: ['/job_proposals/1'] },
  )
}

describe('JobProposalScreen', () => {
  it('renders the screen', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: JOB_PROPOSAL_QUERY,
          variables: { id: '1' },
        },
        result: {
          data: {
            jobProposal: buildJobProposal(),
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForElementToBeRemoved(() => queryByRole('progressbar'))

    expect(await findByText('Job proposal #1')).toBeInTheDocument()
  })

  it('updates the spec', async () => {
    const proposal = buildJobProposal()

    const mocks: MockedResponse[] = [
      {
        request: {
          query: JOB_PROPOSAL_QUERY,
          variables: { id: '1' },
        },
        result: {
          data: {
            jobProposal: proposal,
          },
        },
      },
      {
        request: {
          query: UPDATE_JOB_PROPOSAL_SPEC_MUTATION,
          variables: {
            id: proposal.id,
            input: {
              spec: 'name="updated spec"',
            },
          },
        },
        result: {
          data: {
            updateJobProposalSpec: {
              __typename: 'UpdateJobProposalSpecSuccess',
              jobProposal: {
                __typename: 'JobProposal',
                id: proposal.id,
              },
            },
          },
        },
      },
      {
        request: {
          query: JOB_PROPOSAL_QUERY,
          variables: { id: '1' },
        },
        result: {
          data: {
            jobProposal: buildJobProposal({ spec: 'name="updated spec"' }),
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForElementToBeRemoved(() => queryByRole('progressbar'))

    userEvent.click(getByRole('button', { name: /edit job spec/i }))

    const specInput = screen.getByRole('textbox', {
      name: /job spec \*/i,
    })
    userEvent.clear(specInput)
    userEvent.type(specInput, 'name="updated spec"')

    userEvent.click(getByRole('button', { name: /submit/i }))

    expect(await findByText('Spec updated')).toBeInTheDocument()
  })

  it('proposal not found for update', async () => {
    const proposal = buildJobProposal()

    const mocks: MockedResponse[] = [
      {
        request: {
          query: JOB_PROPOSAL_QUERY,
          variables: { id: '1' },
        },
        result: {
          data: {
            jobProposal: proposal,
          },
        },
      },
      {
        request: {
          query: UPDATE_JOB_PROPOSAL_SPEC_MUTATION,
          variables: {
            id: proposal.id,
            input: {
              spec: 'name="updated spec"',
            },
          },
        },
        result: {
          data: {
            updateJobProposalSpec: {
              __typename: 'NotFoundError',
              message: 'job proposal not found',
            },
          },
        },
      },
      {
        request: {
          query: JOB_PROPOSAL_QUERY,
          variables: { id: '1' },
        },
        result: {
          data: {
            jobProposal: proposal,
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForElementToBeRemoved(() => queryByRole('progressbar'))

    userEvent.click(getByRole('button', { name: /edit job spec/i }))

    const specInput = screen.getByRole('textbox', {
      name: /job spec \*/i,
    })
    userEvent.clear(specInput)
    userEvent.type(specInput, 'name="updated spec"')

    userEvent.click(getByRole('button', { name: /submit/i }))

    expect(await findByText('job proposal not found')).toBeInTheDocument()
  })

  it('approves a pending proposal', async () => {
    const proposal = buildJobProposal()

    const mocks: MockedResponse[] = [
      {
        request: {
          query: JOB_PROPOSAL_QUERY,
          variables: { id: '1' },
        },
        result: {
          data: {
            jobProposal: proposal,
          },
        },
      },
      {
        request: {
          query: APPROVE_JOB_PROPOSAL_MUTATION,
          variables: { id: proposal.id },
        },
        result: {
          data: {
            approveJobProposal: {
              __typename: 'ApproveJobProposalSuccess',
              jobProposal: {
                id: proposal.id,
              },
            },
          },
        },
      },
      {
        request: {
          query: JOB_PROPOSAL_QUERY,
          variables: { id: '1' },
        },
        result: {
          data: {
            jobProposal: buildJobProposal({
              status: 'APPROVED',
              externalJobID: '00000000-0000-0000-0000-000000000001',
            }),
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForElementToBeRemoved(() => queryByRole('progressbar'))

    userEvent.click(getByRole('button', { name: /approve/i }))
    userEvent.click(getByRole('button', { name: /confirm/i }))

    expect(await findByText('Job Proposal approved')).toBeInTheDocument()
    expect(await findByText('Status: Approved')).toBeInTheDocument()
  })

  it('rejects a pending proposal', async () => {
    const proposal = buildJobProposal()

    const mocks: MockedResponse[] = [
      {
        request: {
          query: JOB_PROPOSAL_QUERY,
          variables: { id: '1' },
        },
        result: {
          data: {
            jobProposal: proposal,
          },
        },
      },
      {
        request: {
          query: REJECT_JOB_PROPOSAL_MUTATION,
          variables: { id: proposal.id },
        },
        result: {
          data: {
            rejectJobProposal: {
              __typename: 'RejectJobProposalSuccess',
              jobProposal: {
                id: proposal.id,
              },
            },
          },
        },
      },
      {
        request: {
          query: JOB_PROPOSAL_QUERY,
          variables: { id: '1' },
        },
        result: {
          data: {
            jobProposal: buildJobProposal({
              status: 'REJECTED',
            }),
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForElementToBeRemoved(() => queryByRole('progressbar'))

    userEvent.click(getByRole('button', { name: /reject/i }))
    userEvent.click(getByRole('button', { name: /confirm/i }))

    expect(await findByText('Job Proposal rejected')).toBeInTheDocument()
    expect(await findByText('Status: Rejected')).toBeInTheDocument()
  })

  it('cancels an approved proposal', async () => {
    const proposal = buildJobProposal({ status: 'APPROVED' })

    const mocks: MockedResponse[] = [
      {
        request: {
          query: JOB_PROPOSAL_QUERY,
          variables: { id: '1' },
        },
        result: {
          data: {
            jobProposal: proposal,
          },
        },
      },
      {
        request: {
          query: CANCEL_JOB_PROPOSAL_MUTATION,
          variables: { id: proposal.id },
        },
        result: {
          data: {
            cancelJobProposal: {
              __typename: 'CancelJobProposalSuccess',
              jobProposal: {
                id: proposal.id,
              },
            },
          },
        },
      },
      {
        request: {
          query: JOB_PROPOSAL_QUERY,
          variables: { id: '1' },
        },
        result: {
          data: {
            jobProposal: buildJobProposal({
              status: 'CANCELLED',
            }),
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForElementToBeRemoved(() => queryByRole('progressbar'))

    userEvent.click(getByRole('button', { name: /cancel/i }))
    userEvent.click(getByRole('button', { name: /confirm/i }))

    expect(await findByText('Job Proposal cancelled')).toBeInTheDocument()
    expect(await findByText('Status: Cancelled')).toBeInTheDocument()
  })

  it('approves a cancelled proposal', async () => {
    const proposal = buildJobProposal({ status: 'CANCELLED' })

    const mocks: MockedResponse[] = [
      {
        request: {
          query: JOB_PROPOSAL_QUERY,
          variables: { id: '1' },
        },
        result: {
          data: {
            jobProposal: proposal,
          },
        },
      },
      {
        request: {
          query: APPROVE_JOB_PROPOSAL_MUTATION,
          variables: { id: proposal.id },
        },
        result: {
          data: {
            approveJobProposal: {
              __typename: 'ApproveJobProposalSuccess',
              jobProposal: {
                id: proposal.id,
              },
            },
          },
        },
      },
      {
        request: {
          query: JOB_PROPOSAL_QUERY,
          variables: { id: '1' },
        },
        result: {
          data: {
            jobProposal: buildJobProposal({
              status: 'APPROVED',
            }),
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForElementToBeRemoved(() => queryByRole('progressbar'))

    userEvent.click(getByRole('button', { name: /approve/i }))
    userEvent.click(getByRole('button', { name: /confirm/i }))

    expect(await findByText('Job Proposal approved')).toBeInTheDocument()
    expect(await findByText('Status: Approved')).toBeInTheDocument()
  })

  it('renders a not found page', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: JOB_PROPOSAL_QUERY,
          variables: { id: '1' },
        },
        result: {
          data: {
            jobProposal: {
              __typename: 'NotFoundError',
              message: 'job proposal not found',
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForElementToBeRemoved(() => queryByRole('progressbar'))

    expect(await findByTestId('not-found-page')).toBeInTheDocument()
  })

  it('renders GQL errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: JOB_PROPOSAL_QUERY,
          variables: { id: '1' },
        },
        result: {
          errors: [new GraphQLError('Error!')],
        },
      },
    ]

    renderComponent(mocks)

    await waitForElementToBeRemoved(() => queryByRole('progressbar'))

    expect(await findByText('Error: Error!')).toBeInTheDocument()
  })
})

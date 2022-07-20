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

import {
  buildJobProposal,
  buildJobProposalSpec,
} from 'support/factories/gql/fetchJobProposal'
import {
  APPROVE_JOB_PROPOSAL_SPEC_MUTATION,
  CANCEL_JOB_PROPOSAL_SPEC_MUTATION,
  JOB_PROPOSAL_QUERY,
  REJECT_JOB_PROPOSAL_SPEC_MUTATION,
  UPDATE_JOB_PROPOSAL_SPEC_DEFINITION_MUTATION,
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

      <Route exact path="/feeds_manager">
        Root Redirect Success
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

    expect(await findByText('Job Proposal #1')).toBeInTheDocument()
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
          query: UPDATE_JOB_PROPOSAL_SPEC_DEFINITION_MUTATION,
          variables: {
            id: proposal.id,
            input: {
              definition: 'name="updated spec"',
            },
          },
        },
        result: {
          data: {
            updateJobProposalSpecDefinition: {
              __typename: 'UpdateJobProposalSpecDefinitionSuccess',
              jobProposalSpec: {
                __typename: 'JobProposalSpec',
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
            jobProposal: buildJobProposal(),
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForElementToBeRemoved(() => queryByRole('progressbar'))

    userEvent.click(getByRole('button', { name: /edit/i }))

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
          query: UPDATE_JOB_PROPOSAL_SPEC_DEFINITION_MUTATION,
          variables: {
            id: proposal.id,
            input: {
              definition: 'name="updated spec"',
            },
          },
        },
        result: {
          data: {
            updateJobProposalSpecDefinition: {
              __typename: 'NotFoundError',
              message: 'job proposal spec not found',
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

    userEvent.click(getByRole('button', { name: /edit/i }))

    const specInput = screen.getByRole('textbox', {
      name: /job spec \*/i,
    })
    userEvent.clear(specInput)
    userEvent.type(specInput, 'name="updated spec"')

    userEvent.click(getByRole('button', { name: /submit/i }))

    expect(await findByText('job proposal spec not found')).toBeInTheDocument()
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
          query: APPROVE_JOB_PROPOSAL_SPEC_MUTATION,
          variables: { id: proposal.specs[0].id },
        },
        result: {
          data: {
            approveJobProposalSpec: {
              __typename: 'ApproveJobProposalSpecSuccess',
              spec: {
                id: proposal.specs[0].id,
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

    expect(await findByText('Root Redirect Success')).toBeInTheDocument()
    expect(await findByText('Spec approved')).toBeInTheDocument()
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
          query: REJECT_JOB_PROPOSAL_SPEC_MUTATION,
          variables: { id: proposal.specs[0].id },
        },
        result: {
          data: {
            rejectJobProposalSpec: {
              __typename: 'RejectJobProposalSpecSuccess',
              spec: {
                id: proposal.specs[0].id,
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

    expect(await findByText('Spec rejected')).toBeInTheDocument()
    expect(await findByText('Rejected')).toBeInTheDocument()
  })

  it('cancels an approved proposal', async () => {
    const proposal = buildJobProposal({
      status: 'APPROVED',
      specs: [
        buildJobProposalSpec({
          status: 'APPROVED',
        }),
      ],
    })

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
          query: CANCEL_JOB_PROPOSAL_SPEC_MUTATION,
          variables: { id: proposal.specs[0].id },
        },
        result: {
          data: {
            cancelJobProposalSpec: {
              __typename: 'CancelJobProposalSpecSuccess',
              spec: {
                id: proposal.specs[0].id,
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

    expect(await findByText('Spec cancelled')).toBeInTheDocument()
    expect(await findByText('Cancelled')).toBeInTheDocument()
  })

  it('approves a cancelled proposal', async () => {
    const proposal = buildJobProposal({
      status: 'CANCELLED',
      specs: [
        buildJobProposalSpec({
          status: 'CANCELLED',
        }),
      ],
    })

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
          query: APPROVE_JOB_PROPOSAL_SPEC_MUTATION,
          variables: { id: proposal.specs[0].id },
        },
        result: {
          data: {
            approveJobProposalSpec: {
              __typename: 'ApproveJobProposalSpecSuccess',
              spec: {
                id: proposal.specs[0].id,
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

    expect(await findByText('Root Redirect Success')).toBeInTheDocument()
    expect(await findByText('Spec approved')).toBeInTheDocument()
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

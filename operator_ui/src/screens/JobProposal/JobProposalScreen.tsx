import React from 'react'

import { gql, useMutation, useQuery } from '@apollo/client'
import { useParams } from 'react-router'
import { useDispatch } from 'react-redux'

import { notifySuccessMsg, notifyErrorMsg } from 'actionCreators'
import { FormValues } from './EditJobSpecDialog'
import { GraphqlErrorHandler } from 'src/components/ErrorHandler/GraphqlErrorHandler'
import { JobProposalView, JOB_PROPOSAL_PAYLOAD_FIELDS } from './JobProposalView'
import { Loading } from 'src/components/Feedback/Loading'
import NotFound from 'src/pages/NotFound'
import { useMutationErrorHandler } from 'src/hooks/useMutationErrorHandler'

export const JOB_PROPOSAL_QUERY = gql`
  ${JOB_PROPOSAL_PAYLOAD_FIELDS}
  query FetchJobProposal($id: ID!) {
    jobProposal(id: $id) {
      __typename
      ... on JobProposal {
        ...JobProposalPayloadFields
      }
      ... on NotFoundError {
        message
      }
    }
  }
`

export const UPDATE_JOB_PROPOSAL_SPEC_MUTATION = gql`
  mutation UpdateJobProposalSpec(
    $id: ID!
    $input: UpdateJobProposalSpecInput!
  ) {
    updateJobProposalSpec(id: $id, input: $input) {
      ... on UpdateJobProposalSpecSuccess {
        __typename
      }
      ... on NotFoundError {
        message
      }
    }
  }
`

export const REJECT_JOB_PROPOSAL_MUTATION = gql`
  mutation RejectJobProposal($id: ID!) {
    rejectJobProposal(id: $id) {
      ... on RejectJobProposalSuccess {
        jobProposal {
          id
        }
      }
      ... on NotFoundError {
        message
      }
    }
  }
`

export const APPROVE_JOB_PROPOSAL_MUTATION = gql`
  mutation ApproveJobProposal($id: ID!) {
    approveJobProposal(id: $id) {
      ... on ApproveJobProposalSuccess {
        jobProposal {
          id
        }
      }
      ... on NotFoundError {
        message
      }
    }
  }
`

export const CANCEL_JOB_PROPOSAL_MUTATION = gql`
  mutation CancelJobProposal($id: ID!) {
    cancelJobProposal(id: $id) {
      ... on CancelJobProposalSuccess {
        jobProposal {
          id
        }
      }
      ... on NotFoundError {
        message
      }
    }
  }
`

interface RouteParams {
  id: string
}

export const JobProposalScreen: React.FC = () => {
  const { id } = useParams<RouteParams>()
  const dispatch = useDispatch()
  const { handleMutationError } = useMutationErrorHandler()
  const { data, loading, error } = useQuery<
    FetchJobProposal,
    FetchJobProposalVariables
  >(JOB_PROPOSAL_QUERY, { variables: { id } })

  const [updateJobProposalSpec] = useMutation<
    UpdateJobProposalSpec,
    UpdateJobProposalSpecVariables
  >(UPDATE_JOB_PROPOSAL_SPEC_MUTATION, {
    refetchQueries: [JOB_PROPOSAL_QUERY],
  })

  const [rejectJobProposal] = useMutation<
    RejectJobProposal,
    RejectJobProposalVariables
  >(REJECT_JOB_PROPOSAL_MUTATION, {
    refetchQueries: [JOB_PROPOSAL_QUERY],
  })

  const [approveJobProposal] = useMutation<
    ApproveJobProposal,
    ApproveJobProposalVariables
  >(APPROVE_JOB_PROPOSAL_MUTATION, {
    refetchQueries: [JOB_PROPOSAL_QUERY],
  })

  const [cancelJobProposal] = useMutation<
    CancelJobProposal,
    CancelJobProposalVariables
  >(CANCEL_JOB_PROPOSAL_MUTATION, {
    refetchQueries: [JOB_PROPOSAL_QUERY],
  })

  if (loading) {
    return <Loading />
  }

  if (error) {
    return <GraphqlErrorHandler error={error} />
  }

  const handleUpdateSpec = async (values: FormValues) => {
    try {
      const result = await updateJobProposalSpec({
        variables: { id, input: { ...values } },
      })

      const payload = result.data?.updateJobProposalSpec
      switch (payload?.__typename) {
        case 'UpdateJobProposalSpecSuccess':
          dispatch(notifySuccessMsg('Spec updated'))

          break
        case 'NotFoundError':
          dispatch(notifyErrorMsg(payload.message))

          break
      }
    } catch (e) {
      handleMutationError(e)
    }
  }

  const handleRejectJobProposal = async () => {
    try {
      const result = await rejectJobProposal({
        variables: { id },
      })
      const payload = result.data?.rejectJobProposal
      switch (payload?.__typename) {
        case 'RejectJobProposalSuccess':
          dispatch(notifySuccessMsg('Job Proposal rejected'))

          break
        case 'NotFoundError':
          dispatch(notifyErrorMsg(payload.message))

          break
      }
    } catch (e) {
      handleMutationError(e)
    }
  }

  const handleCancelJobProposal = async () => {
    try {
      const result = await cancelJobProposal({
        variables: { id },
      })
      const payload = result.data?.cancelJobProposal
      switch (payload?.__typename) {
        case 'CancelJobProposalSuccess':
          dispatch(notifySuccessMsg('Job Proposal cancelled'))

          break
        case 'NotFoundError':
          dispatch(notifyErrorMsg(payload.message))

          break
      }
    } catch (e) {
      handleMutationError(e)
    }
  }

  const handleApproveJobProposal = async () => {
    try {
      const result = await approveJobProposal({
        variables: { id },
      })
      const payload = result.data?.approveJobProposal
      switch (payload?.__typename) {
        case 'ApproveJobProposalSuccess':
          dispatch(notifySuccessMsg('Job Proposal approved'))

          break
        case 'NotFoundError':
          dispatch(notifyErrorMsg(payload.message))

          break
      }
    } catch (e) {
      handleMutationError(e)
    }
  }

  const payload = data?.jobProposal
  switch (payload?.__typename) {
    case 'JobProposal':
      return (
        <JobProposalView
          proposal={payload}
          onApprove={handleApproveJobProposal}
          onCancel={handleCancelJobProposal}
          onReject={handleRejectJobProposal}
          onUpdateSpec={handleUpdateSpec}
        />
      )
    case 'NotFoundError':
      return <NotFound />
    default:
      return null
  }
}

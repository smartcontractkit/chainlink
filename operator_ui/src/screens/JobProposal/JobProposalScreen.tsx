import React from 'react'

import { gql, useMutation, useQuery } from '@apollo/client'
import { useHistory, useParams } from 'react-router'
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

export const UPDATE_JOB_PROPOSAL_SPEC_DEFINITION_MUTATION = gql`
  mutation UpdateJobProposalSpecDefinition(
    $id: ID!
    $input: UpdateJobProposalSpecDefinitionInput!
  ) {
    updateJobProposalSpecDefinition(id: $id, input: $input) {
      ... on UpdateJobProposalSpecDefinitionSuccess {
        __typename
      }
      ... on NotFoundError {
        message
      }
    }
  }
`

export const REJECT_JOB_PROPOSAL_SPEC_MUTATION = gql`
  mutation RejectJobProposalSpec($id: ID!) {
    rejectJobProposalSpec(id: $id) {
      ... on RejectJobProposalSpecSuccess {
        spec {
          id
        }
      }
      ... on NotFoundError {
        message
      }
    }
  }
`

export const APPROVE_JOB_PROPOSAL_SPEC_MUTATION = gql`
  mutation ApproveJobProposalSpec($id: ID!) {
    approveJobProposalSpec(id: $id) {
      ... on ApproveJobProposalSpecSuccess {
        spec {
          id
        }
      }
      ... on NotFoundError {
        message
      }
    }
  }
`

export const CANCEL_JOB_PROPOSAL_SPEC_MUTATION = gql`
  mutation CancelJobProposalSpec($id: ID!) {
    cancelJobProposalSpec(id: $id) {
      ... on CancelJobProposalSpecSuccess {
        spec {
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
  const history = useHistory()
  const { handleMutationError } = useMutationErrorHandler()
  const { data, loading, error } = useQuery<
    FetchJobProposal,
    FetchJobProposalVariables
  >(JOB_PROPOSAL_QUERY, { variables: { id } })

  const [updateJobProposalSpecDefinition] = useMutation<
    UpdateJobProposalSpecDefinition,
    UpdateJobProposalSpecDefinitionVariables
  >(UPDATE_JOB_PROPOSAL_SPEC_DEFINITION_MUTATION, {
    refetchQueries: [JOB_PROPOSAL_QUERY],
  })

  const [rejectJobProposalSpec] = useMutation<
    RejectJobProposalSpec,
    RejectJobProposalSpecVariables
  >(REJECT_JOB_PROPOSAL_SPEC_MUTATION, {
    refetchQueries: [JOB_PROPOSAL_QUERY],
  })

  const [approveJobProposalSpec] = useMutation<
    ApproveJobProposalSpec,
    ApproveJobProposalSpecVariables
  >(APPROVE_JOB_PROPOSAL_SPEC_MUTATION, {
    refetchQueries: [JOB_PROPOSAL_QUERY],
  })

  const [cancelJobProposalSpec] = useMutation<
    CancelJobProposalSpec,
    CancelJobProposalSpecVariables
  >(CANCEL_JOB_PROPOSAL_SPEC_MUTATION, {
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
      const result = await updateJobProposalSpecDefinition({
        variables: { id: values.id, input: { definition: values.definition } },
      })

      const payload = result.data?.updateJobProposalSpecDefinition
      switch (payload?.__typename) {
        case 'UpdateJobProposalSpecDefinitionSuccess':
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

  const handleRejectJobProposal = async (specID: string) => {
    try {
      const result = await rejectJobProposalSpec({
        variables: { id: specID },
      })
      const payload = result.data?.rejectJobProposalSpec
      switch (payload?.__typename) {
        case 'RejectJobProposalSpecSuccess':
          dispatch(notifySuccessMsg('Spec rejected'))

          break
        case 'NotFoundError':
          dispatch(notifyErrorMsg(payload.message))

          break
      }
    } catch (e) {
      handleMutationError(e)
    }
  }

  const handleCancelJobProposal = async (specID: string) => {
    try {
      const result = await cancelJobProposalSpec({
        variables: { id: specID },
      })
      const payload = result.data?.cancelJobProposalSpec
      switch (payload?.__typename) {
        case 'CancelJobProposalSpecSuccess':
          dispatch(notifySuccessMsg('Spec cancelled'))

          break
        case 'NotFoundError':
          dispatch(notifyErrorMsg(payload.message))

          break
      }
    } catch (e) {
      handleMutationError(e)
    }
  }

  const handleApproveJobProposal = async (specID: string) => {
    try {
      const result = await approveJobProposalSpec({
        variables: { id: specID },
      })
      const payload = result.data?.approveJobProposalSpec
      switch (payload?.__typename) {
        case 'ApproveJobProposalSpecSuccess':
          history.push('/feeds_manager')

          setTimeout(() => dispatch(notifySuccessMsg('Spec approved')), 200)

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

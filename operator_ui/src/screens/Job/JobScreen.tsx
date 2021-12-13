import React from 'react'

import { gql, useMutation, useQuery } from '@apollo/client'
import { useHistory, useLocation, useParams } from 'react-router-dom'
import { useDispatch } from 'react-redux'

import {
  createJobRunV2,
  notifyErrorMsg,
  notifySuccessMsg,
} from 'src/actionCreators'
import BaseLink from 'src/components/BaseLink'
import ErrorMessage from 'components/Notifications/DefaultError'
import { GraphqlErrorHandler } from 'src/components/ErrorHandler/GraphqlErrorHandler'
import { JobView, JOB_PAYLOAD_FIELDS } from './JobView'
import { Loading } from 'src/components/Feedback/Loading'
import NotFound from 'src/pages/NotFound'
import { useMutationErrorHandler } from 'src/hooks/useMutationErrorHandler'

// Defines the number of records to show on the Overview Tab
const RECENT_RUNS_PAGE_SIZE = 5

export const JOB_QUERY = gql`
  ${JOB_PAYLOAD_FIELDS}
  query FetchJob($id: ID!, $offset: Int, $limit: Int) {
    job(id: $id) {
      __typename
      ... on Job {
        ...JobPayload_Fields
      }
      ... on NotFoundError {
        message
      }
    }
  }
`

export const DELETE_JOB_MUTATION = gql`
  mutation DeleteJob($id: ID!) {
    deleteJob(id: $id) {
      ... on DeleteJobSuccess {
        job {
          id
        }
      }
      ... on NotFoundError {
        message
      }
    }
  }
`

const CreateRunSuccessNotification = ({ data }: any) => (
  <React.Fragment>
    Successfully created job run{' '}
    <BaseLink href={`/runs/${data.id}`}>{data.id}</BaseLink>
  </React.Fragment>
)

interface RouteParams {
  id: string
}

export const JobScreen = () => {
  const { id } = useParams<RouteParams>()
  const dispatch = useDispatch()
  const history = useHistory()
  const { search } = useLocation()
  const { handleMutationError } = useMutationErrorHandler()

  // This doesn't use useQueryParms because we only want to use the initial page
  // load query param to fetch the data.
  //eslint-disable-next-line react-hooks/exhaustive-deps
  const qp = React.useMemo(() => new URLSearchParams(search), [])
  const page = parseInt(qp.get('page') || '1', 10)
  const pageSize = parseInt(
    qp.get('per') || RECENT_RUNS_PAGE_SIZE.toString(),
    10,
  )

  const { data, loading, error, refetch } = useQuery<
    FetchJob,
    FetchJobVariables
  >(JOB_QUERY, {
    variables: { id, offset: (page - 1) * pageSize, limit: pageSize },
    fetchPolicy: 'network-only',
  })

  const [deleteJob] = useMutation<DeleteJob, DeleteJobVariables>(
    DELETE_JOB_MUTATION,
  )

  if (loading) {
    return <Loading />
  }

  if (error) {
    return <GraphqlErrorHandler error={error} />
  }

  const handleDelete = async () => {
    try {
      const result = await deleteJob({
        variables: { id },
      })

      const payload = result.data?.deleteJob
      switch (payload?.__typename) {
        case 'DeleteJobSuccess':
          history.push('/jobs')

          setTimeout(
            () => dispatch(notifySuccessMsg(`Successfully deleted job ${id}`)),
            200,
          )

          break
        case 'NotFoundError':
          dispatch(notifyErrorMsg(payload.message))

          break
      }
    } catch (e) {
      handleMutationError(e)
    }
  }

  const handleRun = async (pipelineInput: string) => {
    await dispatch(
      createJobRunV2(
        id,
        pipelineInput,
        CreateRunSuccessNotification,
        ErrorMessage,
      ),
    )

    refetch()
  }

  const handleRefetch = (page: number, per: number) => {
    refetch({ offset: (page - 1) * pageSize, limit: per })
  }

  const handleRefetchRecentRuns = () => {
    refetch({ offset: 0, limit: RECENT_RUNS_PAGE_SIZE })
  }

  const payload = data?.job
  switch (payload?.__typename) {
    case 'Job':
      return (
        <JobView
          job={payload}
          runsCount={payload.runs.metadata.total}
          onDelete={handleDelete}
          onRun={handleRun}
          refetch={handleRefetch}
          refetchRecentRuns={handleRefetchRecentRuns}
        />
      )
    case 'NotFoundError':
      return <NotFound />
    default:
      return null
  }
}

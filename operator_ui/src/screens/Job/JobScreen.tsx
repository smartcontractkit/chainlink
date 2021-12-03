import React from 'react'

import { gql, useQuery } from '@apollo/client'
import { useHistory, useLocation, useParams } from 'react-router-dom'
import { useDispatch } from 'react-redux'

import {
  createJobRunV2,
  notifyError,
  notifySuccessMsg,
} from 'src/actionCreators'
import ErrorMessage from 'components/Notifications/DefaultError'
import { GraphqlErrorHandler } from 'src/components/ErrorHandler/GraphqlErrorHandler'
import { JobView, JOB_PAYLOAD_FIELDS } from './JobView'
import { Loading } from 'src/components/Feedback/Loading'
import NotFound from 'src/pages/NotFound'

import { v2 } from 'api'
import BaseLink from 'src/components/BaseLink'

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

const CreateRunSuccessNotification = ({ data }: any) => (
  <React.Fragment>
    Successfully created job run{' '}
    <BaseLink href={`/jobs/${data.attributes.jobId}/runs/${data.id}`}>
      {data.id}
    </BaseLink>
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

  // This doesn't use useQueryParms because we only want to use the initial page
  // load query param to fetch the data.
  const qp = React.useMemo(() => new URLSearchParams(search), []) // eslint-disable-next-line react-hooks/exhaustive-deps
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

  if (loading) {
    return <Loading />
  }

  if (error) {
    return <GraphqlErrorHandler error={error} />
  }

  // TODO - Convert this to GQL
  const handleDelete = async () => {
    try {
      await v2.jobs.destroyJobSpec(id)

      history.push('/jobs')

      // This setTimeout is needed because notifications are cleared on a route
      // change. Fix this once notifications are refactored.
      setTimeout(
        () => dispatch(notifySuccessMsg(`Successfully deleted job ${id}`)),
        500,
      )
    } catch (e: any) {
      dispatch(notifyError(ErrorMessage, e))
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

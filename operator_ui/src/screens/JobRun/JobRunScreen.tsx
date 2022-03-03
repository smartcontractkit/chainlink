import React from 'react'

import { gql, useQuery } from '@apollo/client'
import { useParams } from 'react-router-dom'

import { GraphqlErrorHandler } from 'src/components/ErrorHandler/GraphqlErrorHandler'
import { Loading } from 'src/components/Feedback/Loading'
import { JobRunView, JOB_RUN_PAYLOAD_FIELDS } from './JobRunView'
import NotFound from 'src/pages/NotFound'

export const JOB_RUN_QUERY = gql`
  ${JOB_RUN_PAYLOAD_FIELDS}
  query FetchJobRun($id: ID!) {
    jobRun(id: $id) {
      __typename
      ... on JobRun {
        ...JobRunPayload_Fields
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

export const JobRunScreen = () => {
  const { id } = useParams<RouteParams>()

  const { data, loading, error } = useQuery<FetchJobRun, FetchJobRunVariables>(
    JOB_RUN_QUERY,
    { variables: { id } },
  )

  if (loading) {
    return <Loading />
  }

  if (error) {
    return <GraphqlErrorHandler error={error} />
  }

  const payload = data?.jobRun
  switch (payload?.__typename) {
    case 'JobRun':
      return <JobRunView run={payload} />
    case 'NotFoundError':
      return <NotFound />
    default:
      return null
  }
}

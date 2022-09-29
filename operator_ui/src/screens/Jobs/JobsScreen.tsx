import React from 'react'

import { gql, useQuery } from '@apollo/client'

import { GraphqlErrorHandler } from 'src/components/ErrorHandler/GraphqlErrorHandler'
import { JobsView, JOBS_PAYLOAD__RESULTS_FIELDS } from './JobsView'
import { Loading } from 'src/components/Feedback/Loading'
import { useQueryParams } from 'src/hooks/useQueryParams'

export const JOBS_QUERY = gql`
  ${JOBS_PAYLOAD__RESULTS_FIELDS}
  query FetchJobs($offset: Int, $limit: Int) {
    jobs(offset: $offset, limit: $limit) {
      results {
        ...JobsPayload_ResultsFields
      }
      metadata {
        total
      }
    }
  }
`

export const JobsScreen = () => {
  const qp = useQueryParams()
  const page = parseInt(qp.get('page') || '1', 10)
  // Default set to 1000 until we can implement a server side search
  const pageSize = parseInt(qp.get('per') || '1000', 10)

  const { data, loading, error } = useQuery<FetchJobs, FetchJobsVariables>(
    JOBS_QUERY,
    {
      variables: { offset: (page - 1) * pageSize, limit: pageSize },
      fetchPolicy: 'network-only',
    },
  )

  if (loading) {
    return <Loading />
  }

  if (error) {
    return <GraphqlErrorHandler error={error} />
  }

  if (data) {
    return (
      <JobsView
        jobs={data.jobs.results}
        page={page}
        pageSize={pageSize}
        total={data.jobs.metadata.total}
      />
    )
  }

  return null
}

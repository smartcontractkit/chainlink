import React from 'react'

import { gql, useQuery } from '@apollo/client'

import { GraphqlErrorHandler } from 'src/components/ErrorHandler/GraphqlErrorHandler'
import { JobRunsView, JOB_RUNS_PAYLOAD__RESULTS_FIELDS } from './JobRunsView'
import { useQueryParams } from 'src/hooks/useQueryParams'

export const JOB_RUNS_QUERY = gql`
  ${JOB_RUNS_PAYLOAD__RESULTS_FIELDS}
  query FetchJobRuns($offset: Int, $limit: Int) {
    jobRuns(offset: $offset, limit: $limit) {
      results {
        ...JobRunsPayload_ResultsFields
      }
      metadata {
        total
      }
    }
  }
`

export const JobRunsScreen = () => {
  const qp = useQueryParams()
  const page = parseInt(qp.get('page') || '1', 10)
  const pageSize = parseInt(qp.get('per') || '25', 10)

  const { data, loading, error } = useQuery<
    FetchJobRuns,
    FetchJobRunsVariables
  >(JOB_RUNS_QUERY, {
    variables: { offset: (page - 1) * pageSize, limit: pageSize },
    fetchPolicy: 'cache-and-network',
  })

  if (error) {
    return <GraphqlErrorHandler error={error} />
  }

  return (
    <JobRunsView
      loading={loading}
      data={data}
      page={page}
      pageSize={pageSize}
    />
  )
}

import React from 'react'

import { gql, useQuery } from '@apollo/client'

import {
  RecentJobsCard,
  RECENT_JOBS_PAYLOAD__RESULTS_FIELDS,
} from './RecentJobsCard'

const RECENT_JOBS_SIZE = 5

export const RECENT_JOBS_QUERY = gql`
  ${RECENT_JOBS_PAYLOAD__RESULTS_FIELDS}
  query FetchRecentJobs($offset: Int, $limit: Int) {
    jobs(offset: $offset, limit: $limit) {
      results {
        ...RecentJobsPayload_ResultsFields
      }
    }
  }
`

export const RecentJobs = () => {
  const { data, loading, error } = useQuery<
    FetchRecentJobs,
    FetchRecentJobsVariables
  >(RECENT_JOBS_QUERY, {
    variables: { offset: 0, limit: RECENT_JOBS_SIZE },
    fetchPolicy: 'cache-and-network',
  })

  return (
    <RecentJobsCard data={data} errorMsg={error?.message} loading={loading} />
  )
}

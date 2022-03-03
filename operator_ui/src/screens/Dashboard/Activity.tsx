import React from 'react'

import { gql, useQuery } from '@apollo/client'

import {
  ActivityCard,
  RECENT_JOB_RUNS_PAYLOAD__RESULTS_FIELDS,
} from './ActivityCard'

const RECENT_JOB_RUNS_SIZE = 5

export const RECENT_JOB_RUNS_QUERY = gql`
  ${RECENT_JOB_RUNS_PAYLOAD__RESULTS_FIELDS}
  query FetchRecentJobRuns($offset: Int, $limit: Int) {
    jobRuns(offset: $offset, limit: $limit) {
      results {
        ...RecentJobRunsPayload_ResultsFields
      }
      metadata {
        total
      }
    }
  }
`

export const Activity = () => {
  const { data, loading, error } = useQuery<
    FetchRecentJobRuns,
    FetchRecentJobRunsVariables
  >(RECENT_JOB_RUNS_QUERY, {
    variables: { offset: 0, limit: RECENT_JOB_RUNS_SIZE },
    fetchPolicy: 'cache-and-network',
  })

  return (
    <ActivityCard
      data={data}
      errorMsg={error?.message}
      loading={loading}
      maxRunsSize={RECENT_JOB_RUNS_SIZE}
    />
  )
}

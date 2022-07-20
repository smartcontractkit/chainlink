import React from 'react'

import { ChainsView } from './ChainsView'
import { GraphqlErrorHandler } from 'src/components/ErrorHandler/GraphqlErrorHandler'
import { Loading } from 'src/components/Feedback/Loading'
import { useQueryParams } from 'src/hooks/useQueryParams'
import { useChainsQuery } from 'src/hooks/queries/useChainsQuery'

export const ChainsScreen = () => {
  const qp = useQueryParams()
  const page = parseInt(qp.get('page') || '1', 10)
  const pageSize = parseInt(qp.get('per') || '50', 10)

  const { data, loading, error } = useChainsQuery({
    variables: { offset: (page - 1) * pageSize, limit: pageSize },
    fetchPolicy: 'network-only',
  })

  if (loading) {
    return <Loading />
  }

  if (error) {
    return <GraphqlErrorHandler error={error} />
  }

  if (data) {
    return (
      <ChainsView
        chains={data.chains.results}
        page={page}
        pageSize={pageSize}
        total={data.chains.metadata.total}
      />
    )
  }

  return null
}

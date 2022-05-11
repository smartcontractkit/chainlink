import React from 'react'

import { EVMAccountsCard } from './EVMAccountsCard'
import { useEVMAccountsQuery } from 'src/hooks/queries/useEVMAccountsQuery'

export const EVMAccounts = () => {
  const { data, loading, error } = useEVMAccountsQuery({
    fetchPolicy: 'cache-and-network',
  })

  return (
    <EVMAccountsCard loading={loading} data={data} errorMsg={error?.message} />
  )
}

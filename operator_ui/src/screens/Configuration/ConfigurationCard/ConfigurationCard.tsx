import React from 'react'

import { gql, useQuery } from '@apollo/client'

import { KeyValueListCard } from 'src/components/Cards/KeyValueListCard'

export const CONFIG__ITEMS_FIELDS = gql`
  fragment Config_ItemsFields on ConfigItem {
    key
    value
  }
`

export const CONFIG_QUERY = gql`
  ${CONFIG__ITEMS_FIELDS}
  query FetchConfig {
    config {
      items {
        ...Config_ItemsFields
      }
    }
  }
`

export const ConfigurationCard = () => {
  const { data, loading, error } = useQuery<FetchConfig, FetchConfigVariables>(
    CONFIG_QUERY,
    {
      fetchPolicy: 'cache-and-network',
    },
  )

  const entries = React.useMemo(() => {
    if (!data) {
      return []
    }

    return data.config.items
      .map((k) => {
        return [k.key, k.value]
      })
      .sort()
  }, [data])

  return (
    <KeyValueListCard
      title="Configuration"
      error={error?.message}
      loading={loading}
      entries={entries}
      showHead
    />
  )
}

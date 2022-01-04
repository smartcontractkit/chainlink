import React from 'react'

import { gql, useMutation, useQuery } from '@apollo/client'
import { useDispatch } from 'react-redux'

import { notifySuccessMsg } from 'actionCreators'
import { CSAKeysCard, CSA_KEYS_PAYLOAD__RESULTS_FIELDS } from './CSAKeysCard'
import { useMutationErrorHandler } from 'src/hooks/useMutationErrorHandler'

export const CSA_KEYS_QUERY = gql`
  ${CSA_KEYS_PAYLOAD__RESULTS_FIELDS}
  query FetchCSAKeys {
    csaKeys {
      results {
        ...CSAKeysPayload_ResultsFields
      }
    }
  }
`

export const CREATE_CSA_KEY_MUTATION = gql`
  mutation CreateCSAKey {
    createCSAKey {
      ... on CreateCSAKeySuccess {
        csaKey {
          id
        }
      }
      ... on CSAKeyExistsError {
        message
      }
    }
  }
`

export const CSAKeys = () => {
  const dispatch = useDispatch()
  const { handleMutationError } = useMutationErrorHandler()
  const { data, loading, error, refetch } = useQuery<
    FetchCsaKeys,
    FetchCsaKeysVariables
  >(CSA_KEYS_QUERY, {
    fetchPolicy: 'network-only',
  })
  const [createCSAKey] = useMutation<CreateCsaKey, CreateCsaKeyVariables>(
    CREATE_CSA_KEY_MUTATION,
  )

  const handleCreate = async () => {
    try {
      const result = await createCSAKey()

      const payload = result.data?.createCSAKey
      switch (payload?.__typename) {
        case 'CreateCSAKeySuccess':
          dispatch(notifySuccessMsg('CSA Key created'))

          refetch()

          break
      }
    } catch (e) {
      handleMutationError(e)
    }
  }

  return (
    <CSAKeysCard
      loading={loading}
      data={data}
      errorMsg={error?.message}
      onCreate={handleCreate}
    />
  )
}

import React from 'react'

import { gql, useMutation, useQuery } from '@apollo/client'
import { useDispatch } from 'react-redux'

import { useMutationErrorHandler } from 'src/hooks/useMutationErrorHandler'
import {
  createSuccessNotification,
  deleteSuccessNotification,
} from './notifications'
import { P2PKeysCard, P2P_KEYS_PAYLOAD__RESULTS_FIELDS } from './P2PKeysCard'

export const P2P_KEYS_QUERY = gql`
  ${P2P_KEYS_PAYLOAD__RESULTS_FIELDS}
  query FetchP2PKeys {
    p2pKeys {
      results {
        ...P2PKeysPayload_ResultsFields
      }
    }
  }
`

export const CREATE_P2P_KEY_MUTATION = gql`
  ${P2P_KEYS_PAYLOAD__RESULTS_FIELDS}
  mutation CreateP2PKey {
    createP2PKey {
      ... on CreateP2PKeySuccess {
        p2pKey {
          ...P2PKeysPayload_ResultsFields
        }
      }
    }
  }
`

export const DELETE_P2P_KEY_MUTATION = gql`
  ${P2P_KEYS_PAYLOAD__RESULTS_FIELDS}
  mutation DeleteP2PKey($id: ID!) {
    deleteP2PKey(id: $id) {
      ... on DeleteP2PKeySuccess {
        p2pKey {
          ...P2PKeysPayload_ResultsFields
        }
      }
    }
  }
`

export const P2PKeys = () => {
  const dispatch = useDispatch()
  const { handleMutationError } = useMutationErrorHandler()
  const { data, loading, error, refetch } = useQuery<
    FetchP2PKeys,
    FetchP2PKeysVariables
  >(P2P_KEYS_QUERY)
  const [createP2PKey] = useMutation<CreateP2PKey, CreateP2PKeyVariables>(
    CREATE_P2P_KEY_MUTATION,
  )

  const [deleteP2PKey] = useMutation<DeleteP2PKey, DeleteP2PKeyVariables>(
    DELETE_P2P_KEY_MUTATION,
  )

  const handleCreate = async () => {
    try {
      const result = await createP2PKey()
      const payload = result.data?.createP2PKey
      switch (payload?.__typename) {
        case 'CreateP2PKeySuccess':
          dispatch(
            createSuccessNotification({
              keyType: 'P2P Key',
              keyValue: payload.p2pKey.id,
            }),
          )

          refetch()

          break
      }
    } catch (e) {
      handleMutationError(e)
    }
  }

  const handleDelete = async (id: string) => {
    try {
      const result = await deleteP2PKey({ variables: { id } })
      const payload = result.data?.deleteP2PKey
      switch (payload?.__typename) {
        case 'DeleteP2PKeySuccess':
          dispatch(
            deleteSuccessNotification({
              keyType: 'P2P',
            }),
          )

          refetch()

          break
      }
    } catch (e) {
      handleMutationError(e)
    }
  }

  return (
    <P2PKeysCard
      loading={loading}
      data={data}
      errorMsg={error?.message}
      onCreate={handleCreate}
      onDelete={handleDelete}
    />
  )
}

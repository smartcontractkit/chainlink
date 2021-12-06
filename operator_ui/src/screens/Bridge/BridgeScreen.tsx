import React from 'react'

import { gql, useMutation, useQuery } from '@apollo/client'
import { useDispatch } from 'react-redux'
import { useHistory, useParams } from 'react-router-dom'

import { notifySuccessMsg, notifyErrorMsg } from 'actionCreators'
import { BridgeView, BRIDGE_PAYLOAD_FIELDS } from './BridgeView'
import { GraphqlErrorHandler } from 'src/components/ErrorHandler/GraphqlErrorHandler'
import { Loading } from 'src/components/Feedback/Loading'
import NotFound from 'src/pages/NotFound'
import { useMutationErrorHandler } from 'src/hooks/useMutationErrorHandler'

export const BRIDGE_QUERY = gql`
  ${BRIDGE_PAYLOAD_FIELDS}
  query FetchBridge($id: ID!) {
    bridge(id: $id) {
      __typename
      ... on Bridge {
        ...BridgePayload_Fields
      }
      ... on NotFoundError {
        message
      }
    }
  }
`

export const DELETE_BRIDGE_MUTATION = gql`
  mutation DeleteBridge($id: ID!) {
    deleteBridge(id: $id) {
      ... on DeleteBridgeSuccess {
        bridge {
          id
        }
      }
      ... on NotFoundError {
        message
      }
      ... on DeleteBridgeInvalidNameError {
        message
      }
      ... on DeleteBridgeConflictError {
        message
      }
    }
  }
`

interface RouteParams {
  id: string
}

export const BridgeScreen = () => {
  const { id } = useParams<RouteParams>()
  const history = useHistory()
  const dispatch = useDispatch()
  const { handleMutationError } = useMutationErrorHandler()

  const { data, loading, error } = useQuery<FetchBridge, FetchBridgeVariables>(
    BRIDGE_QUERY,
    { variables: { id } },
  )
  const [deleteBridge] = useMutation<DeleteBridge, DeleteBridgeVariables>(
    DELETE_BRIDGE_MUTATION,
  )

  if (loading) {
    return <Loading />
  }

  if (error) {
    return <GraphqlErrorHandler error={error} />
  }

  const handleDelete = async () => {
    try {
      const result = await deleteBridge({
        variables: { id },
      })

      const payload = result.data?.deleteBridge
      switch (payload?.__typename) {
        case 'DeleteBridgeSuccess':
          history.push('/bridges')

          dispatch(notifySuccessMsg('Bridge Deleted'))

          break
        case 'NotFoundError':
        case 'DeleteBridgeInvalidNameError':
        case 'DeleteBridgeConflictError':
          dispatch(notifyErrorMsg(payload.message))

          break
      }
    } catch (e) {
      handleMutationError(e)
    }
  }

  const payload = data?.bridge
  switch (payload?.__typename) {
    case 'Bridge':
      return <BridgeView bridge={payload} onDelete={handleDelete} />
    case 'NotFoundError':
      return <NotFound />
    default:
      return null
  }
}

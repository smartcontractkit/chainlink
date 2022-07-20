import React from 'react'

import { gql, useMutation, useQuery } from '@apollo/client'
import { useDispatch } from 'react-redux'
import { useHistory, useParams } from 'react-router-dom'

import { notifySuccessMsg, notifyErrorMsg } from 'actionCreators'
import { GraphqlErrorHandler } from 'src/components/ErrorHandler/GraphqlErrorHandler'
import { Loading } from 'src/components/Feedback/Loading'
import { NodeView, NODE_PAYLOAD_FIELDS } from './NodeView'
import NotFound from 'src/pages/NotFound'
import { useMutationErrorHandler } from 'src/hooks/useMutationErrorHandler'

export const NODE_QUERY = gql`
  ${NODE_PAYLOAD_FIELDS}
  query FetchNode($id: ID!) {
    node(id: $id) {
      __typename
      ... on Node {
        ...NodePayload_Fields
      }
      ... on NotFoundError {
        message
      }
    }
  }
`

export const DELETE_NODE_MUTATION = gql`
  mutation DeleteNode($id: ID!) {
    deleteNode(id: $id) {
      ... on DeleteNodeSuccess {
        node {
          id
          chain {
            id
          }
        }
      }
      ... on NotFoundError {
        message
      }
    }
  }
`

interface RouteParams {
  id: string
}

export const NodeScreen = () => {
  const { id } = useParams<RouteParams>()
  const dispatch = useDispatch()
  const history = useHistory()
  const { handleMutationError } = useMutationErrorHandler()

  const { data, loading, error } = useQuery<FetchNode, FetchNodeVariables>(
    NODE_QUERY,
    { variables: { id } },
  )

  const [deleteNode] = useMutation<DeleteNode, DeleteNodeVariables>(
    DELETE_NODE_MUTATION,
  )

  if (loading) {
    return <Loading />
  }

  if (error) {
    return <GraphqlErrorHandler error={error} />
  }

  const handleDeleteNode = async () => {
    try {
      const result = await deleteNode({
        variables: { id },
      })
      const payload = result.data?.deleteNode
      switch (payload?.__typename) {
        case 'DeleteNodeSuccess':
          history.push(`/chains/${payload.node.chain.id}`)

          dispatch(notifySuccessMsg('Node deleted'))

          break
        case 'NotFoundError':
          dispatch(notifyErrorMsg(payload.message))

          break
      }
    } catch (e) {
      handleMutationError(e)
    }
  }

  const payload = data?.node
  switch (payload?.__typename) {
    case 'Node':
      return <NodeView node={payload} onDelete={handleDeleteNode} />
    case 'NotFoundError':
      return <NotFound />
    default:
      return null
  }
}

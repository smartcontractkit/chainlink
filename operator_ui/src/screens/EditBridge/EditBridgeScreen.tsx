import React from 'react'

import { gql, useMutation, useQuery } from '@apollo/client'
import { useDispatch } from 'react-redux'
import { useParams } from 'react-router-dom'

import { notifyErrorMsg, notifySuccess } from 'actionCreators'
import BaseLink from 'src/components/BaseLink'
import { BRIDGE_QUERY } from '../Bridge/BridgeScreen'
import { FormValues } from 'components/Form/BridgeForm'
import { EditBridgeView } from './EditBridgeView'
import { useMutationErrorHandler } from 'src/hooks/useMutationErrorHandler'
import { Loading } from 'src/components/Feedback/Loading'
import { GraphqlErrorHandler } from 'src/components/ErrorHandler/GraphqlErrorHandler'
import NotFound from 'src/pages/NotFound'

export const UPDATE_BRIDGE_MUTATION = gql`
  mutation UpdateBridge($id: ID!, $input: UpdateBridgeInput!) {
    updateBridge(id: $id, input: $input) {
      ... on UpdateBridgeSuccess {
        bridge {
          id
        }
      }
      ... on NotFoundError {
        message
      }
    }
  }
`

const SuccessNotification = ({ id }: { id: string }) => {
  return (
    <>
      <span>Successfully updated</span>
      <BaseLink href={`/bridges/${id}`}>{id}</BaseLink>
    </>
  )
}

interface RouteParams {
  id: string
}

export const EditBridgeScreen = () => {
  const { id } = useParams<RouteParams>()
  const dispatch = useDispatch()
  const { handleMutationError } = useMutationErrorHandler()
  const { data, loading, error } = useQuery<FetchBridge, FetchBridgeVariables>(
    BRIDGE_QUERY,
    { variables: { id } },
  )

  const [updateBridge] = useMutation<UpdateBridge, UpdateBridgeVariables>(
    UPDATE_BRIDGE_MUTATION,
    {
      refetchQueries: [{ query: BRIDGE_QUERY, variables: { id } }],
    },
  )

  if (loading) {
    return <Loading />
  }

  if (error) {
    return <GraphqlErrorHandler error={error} />
  }

  const handleSubmit = async (values: FormValues) => {
    try {
      const result = await updateBridge({
        variables: { id, input: { ...values } },
      })

      const payload = result.data?.updateBridge
      switch (payload?.__typename) {
        case 'UpdateBridgeSuccess':
          dispatch(
            notifySuccess(
              () => <SuccessNotification id={payload.bridge.id} />,
              {},
            ),
          )

          break
        case 'NotFoundError':
          dispatch(notifyErrorMsg(payload.message))
      }
    } catch (e) {
      handleMutationError(e)
    }
  }

  const payload = data?.bridge
  switch (payload?.__typename) {
    case 'Bridge':
      return <EditBridgeView onSubmit={handleSubmit} bridge={payload} />
    case 'NotFoundError':
      return <NotFound />
    default:
      return null
  }
}

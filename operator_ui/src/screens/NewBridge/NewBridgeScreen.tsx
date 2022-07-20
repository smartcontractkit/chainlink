import React from 'react'

import { gql, useMutation } from '@apollo/client'
import { useDispatch } from 'react-redux'

import { notifySuccess } from 'actionCreators'
import BaseLink from 'src/components/BaseLink'
import { FormValues } from 'components/Form/BridgeForm'
import { NewBridgeView } from './NewBridgeView'
import { useMutationErrorHandler } from 'src/hooks/useMutationErrorHandler'

export const CREATE_BRIDGE_MUTATION = gql`
  mutation CreateBridge($input: CreateBridgeInput!) {
    createBridge(input: $input) {
      ... on CreateBridgeSuccess {
        bridge {
          id
        }
        incomingToken
      }
    }
  }
`

const SuccessNotification = ({
  id,
  incomingToken,
}: {
  id: string
  incomingToken: string
}) => {
  return (
    <>
      <span>Successfully created bridge&nbsp;</span>
      <BaseLink href={`/bridges/${id}`}>{id}</BaseLink>
      <span>&nbsp;with incoming access token: {incomingToken}</span>
    </>
  )
}

export const NewBridgeScreen = () => {
  const dispatch = useDispatch()
  const { handleMutationError } = useMutationErrorHandler()
  const [createBridge] = useMutation<CreateBridge, CreateBridgeVariables>(
    CREATE_BRIDGE_MUTATION,
  )

  const handleSubmit = async (values: FormValues) => {
    try {
      const result = await createBridge({
        variables: { input: { ...values } },
      })

      const payload = result.data?.createBridge
      switch (payload?.__typename) {
        case 'CreateBridgeSuccess':
          dispatch(
            notifySuccess(
              () => (
                <SuccessNotification
                  id={payload.bridge.id}
                  incomingToken={payload.incomingToken}
                />
              ),
              {},
            ),
          )

          break
      }
    } catch (e) {
      handleMutationError(e)
    }
  }

  return <NewBridgeView onSubmit={handleSubmit} />
}

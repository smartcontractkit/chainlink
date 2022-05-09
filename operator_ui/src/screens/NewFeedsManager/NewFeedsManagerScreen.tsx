import React from 'react'

import { useMutation, gql } from '@apollo/client'
import { FormikHelpers } from 'formik'
import { useDispatch } from 'react-redux'
import { Redirect, useHistory, useLocation } from 'react-router-dom'

import { notifySuccessMsg, notifyErrorMsg } from 'actionCreators'
import { GraphqlErrorHandler } from 'src/components/ErrorHandler/GraphqlErrorHandler'
import { FormValues } from 'components/Form/FeedsManagerForm'
import { parseInputErrors } from 'src/utils/inputErrors'
import { Loading } from 'src/components/Feedback/Loading'
import { NewFeedsManagerView } from './NewFeedsManagerView'
import {
  FEEDS_MANAGERS_QUERY,
  useFeedsManagersQuery,
} from 'src/hooks/queries/useFeedsManagersQuery'
import { useMutationErrorHandler } from 'src/hooks/useMutationErrorHandler'

export const CREATE_FEEDS_MANAGER_MUTATION = gql`
  mutation CreateFeedsManager($input: CreateFeedsManagerInput!) {
    createFeedsManager(input: $input) {
      ... on CreateFeedsManagerSuccess {
        feedsManager {
          id
          name
          uri
          publicKey
          isConnectionActive
          createdAt
        }
      }
      ... on SingleFeedsManagerError {
        message
        code
      }
      ... on NotFoundError {
        message
        code
      }
      ... on InputErrors {
        errors {
          path
          message
          code
        }
      }
    }
  }
`

export const NewFeedsManagerScreen: React.FC = () => {
  const history = useHistory()
  const location = useLocation()
  const dispatch = useDispatch()
  const { handleMutationError } = useMutationErrorHandler()
  const { data, loading, error } = useFeedsManagersQuery()
  const [createFeedsManager] = useMutation<
    CreateFeedsManager,
    CreateFeedsManagerVariables
  >(CREATE_FEEDS_MANAGER_MUTATION, {
    refetchQueries: [FEEDS_MANAGERS_QUERY],
  })

  if (loading) {
    return <Loading />
  }

  if (error) {
    return <GraphqlErrorHandler error={error} />
  }

  // We currently only support a single feeds manager, but plan to support more
  // in the future.
  const manager =
    data != undefined && data.feedsManagers.results[0]
      ? data.feedsManagers.results[0]
      : undefined

  const handleSubmit = async (
    values: FormValues,
    { setErrors }: FormikHelpers<FormValues>,
  ) => {
    try {
      const result = await createFeedsManager({
        variables: { input: { ...values } },
      })

      const payload = result.data?.createFeedsManager
      switch (payload?.__typename) {
        case 'CreateFeedsManagerSuccess':
          history.push('/feeds_manager')

          dispatch(notifySuccessMsg('Feeds Manager Created'))

          break
        case 'SingleFeedsManagerError':
        case 'NotFoundError':
          dispatch(notifyErrorMsg(payload.message))

          break
        case 'InputErrors':
          dispatch(notifyErrorMsg('Invalid Input'))

          setErrors(parseInputErrors(payload))

          break
      }
    } catch (e) {
      handleMutationError(e)
    }
  }

  if (manager) {
    return (
      <Redirect
        to={{
          pathname: '/feeds_manager',
          state: { from: location },
        }}
      />
    )
  }

  return <NewFeedsManagerView onSubmit={handleSubmit} />
}

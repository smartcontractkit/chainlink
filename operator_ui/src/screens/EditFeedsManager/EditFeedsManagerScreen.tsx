import React from 'react'

import { useMutation, gql } from '@apollo/client'
import { FormikHelpers } from 'formik'
import { useDispatch } from 'react-redux'
import { Redirect, useLocation, useHistory } from 'react-router-dom'

import { notifySuccessMsg, notifyErrorMsg } from 'actionCreators'
import { EditFeedsManagerView } from './EditFeedsManagerView'
import { GraphqlErrorHandler } from 'src/components/ErrorHandler/GraphqlErrorHandler'
import { FormValues } from 'components/Form/FeedsManagerForm'
import { parseInputErrors } from 'src/utils/inputErrors'
import { Loading } from 'components/Feedback/Loading'
import {
  useFeedsManagersQuery,
  FEEDS_MANAGERS_QUERY,
} from 'src/hooks/queries/useFeedsManagersQuery'
import { useMutationErrorHandler } from 'src/hooks/useMutationErrorHandler'

export const UPDATE_FEEDS_MANAGER_MUTATION = gql`
  mutation UpdateFeedsManager($id: ID!, $input: UpdateFeedsManagerInput!) {
    updateFeedsManager(id: $id, input: $input) {
      ... on UpdateFeedsManagerSuccess {
        feedsManager {
          id
          name
          uri
          publicKey
          isConnectionActive
          createdAt
        }
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

export const EditFeedsManagerScreen: React.FC = () => {
  const history = useHistory()
  const location = useLocation()
  const dispatch = useDispatch()
  const { handleMutationError } = useMutationErrorHandler()
  const { data, loading, error } = useFeedsManagersQuery()
  const [updateFeedsManager] = useMutation<
    UpdateFeedsManager,
    UpdateFeedsManagerVariables
  >(UPDATE_FEEDS_MANAGER_MUTATION, {
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

  if (!manager) {
    return (
      <Redirect
        to={{
          pathname: '/feeds_manager/new',
          state: { from: location },
        }}
      />
    )
  }

  const handleSubmit = async (
    values: FormValues,
    { setErrors }: FormikHelpers<FormValues>,
  ) => {
    try {
      const result = await updateFeedsManager({
        variables: { id: manager.id, input: { ...values } },
      })

      const payload = result.data?.updateFeedsManager
      switch (payload?.__typename) {
        case 'UpdateFeedsManagerSuccess':
          history.push('/feeds_manager')

          dispatch(notifySuccessMsg('Feeds Manager Updated'))

          break
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

  return <EditFeedsManagerView data={manager} onSubmit={handleSubmit} />
}

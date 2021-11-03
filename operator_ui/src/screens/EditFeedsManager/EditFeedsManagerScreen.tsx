import React from 'react'
import { useMutation, gql } from '@apollo/client'
import { FormikHelpers } from 'formik'
import { Redirect, useLocation, useHistory } from 'react-router-dom'

import { EditFeedsManagerView } from './EditFeedsManagerView'
import { FormValues } from 'components/Forms/FeedsManagerForm'
import { Loading } from 'components/Feedback/Loading'
import {
  useFetchFeedsManagers,
  FETCH_FEEDS_MANAGERS,
} from 'src/hooks/useFetchFeedsManager'

// NOTE: To be refactored to not use redux
import { useDispatch } from 'react-redux'
import { notifySuccessMsg, notifyErrorMsg } from 'actionCreators'

export const UPDATE_FEEDS_MANAGER = gql`
  mutation UpdateFeedsManager($id: ID!, $input: UpdateFeedsManagerInput!) {
    updateFeedsManager(id: $id, input: $input) {
      ... on UpdateFeedsManagerSuccess {
        feedsManager {
          id
          name
          uri
          publicKey
          jobTypes
          isBootstrapPeer
          isConnectionActive
          bootstrapPeerMultiaddr
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
  const { data, loading, error } = useFetchFeedsManagers()
  const [updateFeedsManager] = useMutation<
    UpdateFeedsManager,
    UpdateFeedsManagerVariables
  >(UPDATE_FEEDS_MANAGER, {
    refetchQueries: [FETCH_FEEDS_MANAGERS],
  })

  if (loading) {
    return <Loading />
  }

  if (error) {
    return <div>error</div>
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
          const errs = payload.errors.reduce((obj, item) => {
            const key = item['path'].replace(/^input\//, '')

            return {
              ...obj,
              [key]: item.message,
            }
          }, {})

          setErrors(errs)

          break
      }
    } catch (e) {
      // TODO - Handle errors
      console.log(e)
    }
  }

  return <EditFeedsManagerView data={manager} onSubmit={handleSubmit} />
}

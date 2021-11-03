import React from 'react'
import { FormikHelpers } from 'formik'
import { useMutation, gql } from '@apollo/client'
import { Redirect, useHistory, useLocation } from 'react-router-dom'

import { FormValues } from 'components/Forms/FeedsManagerForm'
import { Loading } from 'src/components/Feedback/Loading'
import { NewFeedsManagerView } from './NewFeedsManagerView'
import {
  FETCH_FEEDS_MANAGERS,
  useFetchFeedsManagers,
} from 'src/hooks/useFetchFeedsManager'

// NOTE: To be refactored to not use redux
import { useDispatch } from 'react-redux'
import { notifySuccessMsg, notifyErrorMsg } from 'actionCreators'

export const CREATE_FEEDS_MANAGER = gql`
  mutation CreateFeedsManager($input: CreateFeedsManagerInput!) {
    createFeedsManager(input: $input) {
      ... on CreateFeedsManagerSuccess {
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
  const { data, loading, error } = useFetchFeedsManagers()
  const [createFeedsManager] = useMutation<
    CreateFeedsManager,
    CreateFeedsManagerVariables
  >(CREATE_FEEDS_MANAGER, {
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

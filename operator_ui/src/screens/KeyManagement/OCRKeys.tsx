import React from 'react'

import { gql, useMutation } from '@apollo/client'
import { useDispatch } from 'react-redux'

import { useMutationErrorHandler } from 'src/hooks/useMutationErrorHandler'
import {
  createSuccessNotification,
  deleteSuccessNotification,
} from './notifications'
import { OCRKeysCard } from './OCRKeysCard'
import { useOCRKeysQuery } from 'src/hooks/queries/useOCRKeysQuery'

export const CREATE_OCR_KEY_BUNDLE_MUTATION = gql`
  mutation CreateOCRKeyBundle {
    createOCRKeyBundle {
      ... on CreateOCRKeyBundleSuccess {
        bundle {
          id
        }
      }
    }
  }
`

export const DELETE_OCR_KEY_BUNDLE_MUTATION = gql`
  mutation DeleteOCRKeyBundle($id: ID!) {
    deleteOCRKeyBundle(id: $id) {
      ... on DeleteOCRKeyBundleSuccess {
        bundle {
          id
        }
      }
    }
  }
`

export const OCRKeys = () => {
  const dispatch = useDispatch()
  const { handleMutationError } = useMutationErrorHandler()
  const { data, loading, error, refetch } = useOCRKeysQuery({
    fetchPolicy: 'network-only',
  })
  const [createOCRKeyBundle] = useMutation<
    CreateOcrKeyBundle,
    CreateOcrKeyBundleVariables
  >(CREATE_OCR_KEY_BUNDLE_MUTATION)

  const [deleteOCRKeyBundle] = useMutation<
    DeleteOcrKeyBundle,
    DeleteOcrKeyBundleVariables
  >(DELETE_OCR_KEY_BUNDLE_MUTATION)

  const handleCreate = async () => {
    try {
      const result = await createOCRKeyBundle()

      const payload = result.data?.createOCRKeyBundle
      switch (payload?.__typename) {
        case 'CreateOCRKeyBundleSuccess':
          dispatch(
            createSuccessNotification({
              keyType: 'Off-ChainReporting Key Bundle',
              keyValue: payload.bundle.id,
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
      const result = await deleteOCRKeyBundle({ variables: { id } })

      const payload = result.data?.deleteOCRKeyBundle
      switch (payload?.__typename) {
        case 'DeleteOCRKeyBundleSuccess':
          dispatch(
            deleteSuccessNotification({
              keyType: 'Off-ChainReporting Key Bundle',
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
    <OCRKeysCard
      loading={loading}
      data={data}
      errorMsg={error?.message}
      onCreate={handleCreate}
      onDelete={handleDelete}
    />
  )
}

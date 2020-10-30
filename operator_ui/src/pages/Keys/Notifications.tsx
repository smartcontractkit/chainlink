import React from 'react'
import { notifySuccess } from 'actionCreators'

export const deleteNotification = ({ keyType }: { keyType: string }) =>
  notifySuccess(() => <>Successfully deleted {keyType} Key</>, {})

export const createNotification = ({
  keyType,
  keyValue,
}: {
  keyType: string
  keyValue: string
}) =>
  notifySuccess(
    () => (
      <>
        Successfully created {keyType} Key: {keyValue}
      </>
    ),
    {},
  )

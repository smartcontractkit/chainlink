import { notifySuccessMsg } from 'actionCreators'

// createSuccessNotification generates a action creator which provides a create
// success message.
export const createSuccessNotification = ({
  keyType,
  keyValue,
}: {
  keyType: string
  keyValue: string
}) => notifySuccessMsg(`Successfully created ${keyType}: ${keyValue}`)

// deleteSuccessNotification generates a action creator which provides a delete
// success message.
export const deleteSuccessNotification = ({ keyType }: { keyType: string }) =>
  notifySuccessMsg(`Successfully deleted ${keyType} Key`)

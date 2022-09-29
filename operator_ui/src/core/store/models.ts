/**
 * RunStatus is a string that represents the run status
 */
export enum RunStatus {
  IN_PROGRESS = 'in_progress',
  PENDING_INCOMING_CONFIRMATIONS = 'pending_incoming_confirmations',
  PENDING_CONNECTION = 'pending_connection',
  PENDING_BRIDGE = 'pending_bridge',
  PENDING_SLEEP = 'pending_sleep',
  ERRORED = 'errored',
  COMPLETED = 'completed',
}

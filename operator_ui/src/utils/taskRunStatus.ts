// A task run status is inferred from the state of a task run as we are not
// provided a status from the API.
export enum TaskRunStatus {
  UNKNOWN = 'UNKNOWN',
  PENDING = 'PENDING',
  ERROR = 'ERROR',
  COMPLETE = 'COMPLETE',
}

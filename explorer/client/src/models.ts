export interface IJobRun {
  chainlinkNode: any
  finishedAt: string
  createdAt: string
  error?: string
  id: string
  jobId: string
  requestId?: string
  requester?: string
  runId: string
  status: string
  txHash?: string
  type: string
}

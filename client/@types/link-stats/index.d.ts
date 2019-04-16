interface IJobRun {
  id: number
  runId: string
  jobId: string
  status: string
  error?: string
  createdAt: string
  completedAt?: string
  initiator?: IInitiator
  taskRuns: ITaskRun[]
}

interface ITaskRun {
  id: number
  type: string
  status: string
  error?: string
}

interface IInitiator {
  id: number
  type: string
  requester: string
  requestId: string
}

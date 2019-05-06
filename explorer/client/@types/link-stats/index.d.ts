interface IChainlinkNode {
  name: string
}

interface IJobRun {
  id: string
  runId: string
  jobId: string
  status: string
  type: string
  requester: string
  requestId: string
  txHash: string
  error?: string
  createdAt: string
  completedAt?: string
  chainlinkNode: IChainlinkNode
  taskRuns: ITaskRun[]
}

interface ITaskRun {
  id: number
  type: string
  status: string
  error?: string
}

interface IChainlinkNode {
  id: number
  name: string
}

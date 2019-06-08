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
  finishedAt?: string
  chainlinkNode: IChainlinkNode
  etherscanHost: string
  taskRuns: ITaskRun[]
}

interface ITaskRun {
  id: number
  type: string
  status: string
  transactionHash?: string
  transactionStatus?: string
  confirmations?: number
  minimumConfirmations?: number
  error?: string
}

interface IChainlinkNode {
  id: number
  name: string
}

declare module 'explorer' {}

declare module 'explorer/models' {
  interface ChainlinkNode {
    id: number
    name: string
    url?: string
    createdAt: string
  }

  interface JobRun {
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
    chainlinkNode: ChainlinkNode
    etherscanHost: string
    taskRuns: TaskRun[]
  }

  interface TaskRun {
    id: number
    type: string
    status: string
    transactionHash?: string
    transactionStatus?: string
    confirmations?: string
    minimumConfirmations?: string
    error?: string
  }

  interface Head {
    id: number
    createdAt: string
    parentHash: string
    uncleHash: string
    coinbase: string
    root: string
    txHash: string
    receiptHash: string
    bloom: string
    difficulty: string
    number: string
    gasLimit: string
    gasUsed: string
    time: string
    extra: string
    mixDigest: string
    nonce: string
  }

  interface Config {
    gaId: string
  }
}

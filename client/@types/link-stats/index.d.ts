interface IJobRun {
  id: string
  jobId: string
  status: string
  initiatorType: string
  error?: string
  createdAt: string
  completedAt?: string
}

import statusText from '../../utils/statusText'

const COMPLETED_ETHTX_WITHOUT_STATUS = {
  type: 'ethtx',
  status: 'completed',
  transactionStatus: undefined
} as ITaskRun
const COMPLETED_ETHTX_WITH_STATUS = {
  type: 'ethtx',
  status: 'completed',
  transactionStatus: '0x1'
} as ITaskRun

describe('utils/statusText', () => {
  it('is "Pending (Not Fullfilled)" for a status of "in_progress" without ethtx successful transaction', () => {
    const taskRuns: ITaskRun[] = [COMPLETED_ETHTX_WITHOUT_STATUS]
    const jobRun = { status: 'in_progress', taskRuns: taskRuns } as IJobRun
    expect(statusText(jobRun)).toEqual('Pending (but Not Fullfilled)')
  })

  it('is "Pending" for a status of "in_progress" with ethtx successful transaction', () => {
    const taskRuns: ITaskRun[] = [COMPLETED_ETHTX_WITH_STATUS]
    const jobRun = { status: 'in_progress', taskRuns: taskRuns } as IJobRun
    expect(statusText(jobRun)).toEqual('Pending')
  })

  it('is "Pending" for a status of "in_progress" when there is no ethtx', () => {
    const jobRun = { status: 'in_progress' } as IJobRun
    expect(statusText(jobRun)).toEqual('Pending')
  })

  it('is "Failed (but Not Fullfilled)" for a status of "error" without ethx successful transaction', () => {
    const taskRuns: ITaskRun[] = [COMPLETED_ETHTX_WITHOUT_STATUS]
    const jobRun = { status: 'error', taskRuns: taskRuns } as IJobRun
    expect(statusText(jobRun)).toEqual('Failed (but Not Fullfilled)')
  })

  it('is "Failed" for a status of "error" with ethx successful transaction', () => {
    const taskRuns: ITaskRun[] = [COMPLETED_ETHTX_WITH_STATUS]
    const jobRun = { status: 'error', taskRuns: taskRuns } as IJobRun
    expect(statusText(jobRun)).toEqual('Failed')
  })

  it('is "Failed" for a status of "error"', () => {
    const jobRun = { status: 'error' } as IJobRun
    expect(statusText(jobRun)).toEqual('Failed')
  })

  it('is "Complete (but Not Fullfilled)" for a status of "completed" without ethtx successful transaction', () => {
    const taskRuns: ITaskRun[] = [COMPLETED_ETHTX_WITHOUT_STATUS]
    const jobRun = { status: 'completed', taskRuns: taskRuns } as IJobRun
    expect(statusText(jobRun)).toEqual('Complete (but Not Fullfilled)')
  })

  it('is "Complete" for a status of "completed" without ethtx successful transaction', () => {
    const taskRuns: ITaskRun[] = [COMPLETED_ETHTX_WITH_STATUS]
    const jobRun = { status: 'completed', taskRuns: taskRuns } as IJobRun
    expect(statusText(jobRun)).toEqual('Complete')
  })

  it('is "Complete" for a status of "completed"', () => {
    const jobRun = { status: 'completed' } as IJobRun
    expect(statusText(jobRun)).toEqual('Complete')
  })

  it('returns the status as titlecase by default', () => {
    const jobRun = {
      status: 'pending_confirmations'
    } as IJobRun
    expect(statusText(jobRun)).toEqual('Pending Confirmations')
  })
})

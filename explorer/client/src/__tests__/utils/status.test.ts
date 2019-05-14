import status from '../../utils/status'

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

describe('utils/status', () => {
  it('is Pending not fulfilled for a status of "in_progress" without successful ethtx transaction', () => {
    const taskRuns: ITaskRun[] = [COMPLETED_ETHTX_WITHOUT_STATUS]
    const jobRun = { status: 'in_progress', taskRuns: taskRuns } as IJobRun

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Pending')
    expect(unfulfilledEthTx).toEqual(true)
  })

  it('is Pending for a status of "in_progress" with ethtx successful transaction', () => {
    const taskRuns: ITaskRun[] = [COMPLETED_ETHTX_WITH_STATUS]
    const jobRun = { status: 'in_progress', taskRuns: taskRuns } as IJobRun

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Pending')
    expect(unfulfilledEthTx).toEqual(false)
  })

  it('is Pending for a status of "in_progress" when there is no ethtx', () => {
    const jobRun = { status: 'in_progress' } as IJobRun

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Pending')
    expect(unfulfilledEthTx).toEqual(false)
  })

  it('is Failed not fulfilled for a status of "error" without successful ethx transaction', () => {
    const taskRuns: ITaskRun[] = [COMPLETED_ETHTX_WITHOUT_STATUS]
    const jobRun = { status: 'error', taskRuns: taskRuns } as IJobRun

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Failed')
    expect(unfulfilledEthTx).toEqual(true)
  })

  it('is Failed for a status of "error" with ethx successful transaction', () => {
    const taskRuns: ITaskRun[] = [COMPLETED_ETHTX_WITH_STATUS]
    const jobRun = { status: 'error', taskRuns: taskRuns } as IJobRun

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Failed')
    expect(unfulfilledEthTx).toEqual(false)
  })

  it('is Failed for a status of "error"', () => {
    const jobRun = { status: 'error' } as IJobRun

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Failed')
    expect(unfulfilledEthTx).toEqual(false)
  })

  it('is "Complete (but Not Fullfilled)" for a status of "completed" without ethtx successful transaction', () => {
    const taskRuns: ITaskRun[] = [COMPLETED_ETHTX_WITHOUT_STATUS]
    const jobRun = { status: 'completed', taskRuns: taskRuns } as IJobRun

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Complete')
    expect(unfulfilledEthTx).toEqual(true)
  })

  it('is "Complete" for a status of "completed" without ethtx successful transaction', () => {
    const taskRuns: ITaskRun[] = [COMPLETED_ETHTX_WITH_STATUS]
    const jobRun = { status: 'completed', taskRuns: taskRuns } as IJobRun

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Complete')
    expect(unfulfilledEthTx).toEqual(false)
  })

  it('is "Complete" for a status of "completed"', () => {
    const jobRun = { status: 'completed' } as IJobRun

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Complete')
    expect(unfulfilledEthTx).toEqual(false)
  })

  it('returns the status as titlecase by default', () => {
    const jobRun = { status: 'pending_confirmations' } as IJobRun

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Pending Confirmations')
    expect(unfulfilledEthTx).toEqual(false)
  })
})

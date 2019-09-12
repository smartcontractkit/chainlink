import status from '../../utils/status'

const COMPLETED_ETHTX_WITHOUT_STATUS = {
  type: 'ethtx',
  status: 'completed',
  transactionStatus: undefined,
} as TaskRun
const COMPLETED_ETHTX_WITH_STATUS = {
  type: 'ethtx',
  status: 'completed',
  transactionStatus: 'fulfilledRunLog',
} as TaskRun

describe('utils/status', () => {
  it('is Pending not fulfilled for a status of "in_progress" without successful ethtx transaction', () => {
    const taskRuns: TaskRun[] = [COMPLETED_ETHTX_WITHOUT_STATUS]
    const jobRun = { status: 'in_progress', taskRuns: taskRuns } as JobRun

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Pending')
    expect(unfulfilledEthTx).toEqual(true)
  })

  it('is Pending for a status of "in_progress" with ethtx successful transaction', () => {
    const taskRuns: TaskRun[] = [COMPLETED_ETHTX_WITH_STATUS]
    const jobRun = { status: 'in_progress', taskRuns: taskRuns } as JobRun

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Pending')
    expect(unfulfilledEthTx).toEqual(false)
  })

  it('is Pending for a status of "in_progress" when there is no ethtx', () => {
    const jobRun = { status: 'in_progress' } as JobRun

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Pending')
    expect(unfulfilledEthTx).toEqual(false)
  })

  it('is Errored not fulfilled for a status of "error" without successful ethx transaction', () => {
    const taskRuns: TaskRun[] = [COMPLETED_ETHTX_WITHOUT_STATUS]
    const jobRun = { status: 'error', taskRuns: taskRuns } as JobRun

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Errored')
    expect(unfulfilledEthTx).toEqual(true)
  })

  it('is Errored for a status of "error" with ethx successful transaction', () => {
    const taskRuns: TaskRun[] = [COMPLETED_ETHTX_WITH_STATUS]
    const jobRun = { status: 'error', taskRuns: taskRuns } as JobRun

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Errored')
    expect(unfulfilledEthTx).toEqual(false)
  })

  it('is Errored for a status of "error"', () => {
    const jobRun = { status: 'error' } as JobRun

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Errored')
    expect(unfulfilledEthTx).toEqual(false)
  })

  it('is Complete not fullfilled for a status of "completed" without ethtx successful transaction', () => {
    const taskRuns: TaskRun[] = [COMPLETED_ETHTX_WITHOUT_STATUS]
    const jobRun = { status: 'completed', taskRuns: taskRuns } as JobRun

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Complete')
    expect(unfulfilledEthTx).toEqual(true)
  })

  it('is Complete for a status of "completed" without ethtx successful transaction', () => {
    const taskRuns: TaskRun[] = [COMPLETED_ETHTX_WITH_STATUS]
    const jobRun = { status: 'completed', taskRuns: taskRuns } as JobRun

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Complete')
    expect(unfulfilledEthTx).toEqual(false)
  })

  it('is Complete for a status of "completed"', () => {
    const jobRun = { status: 'completed' } as JobRun

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Complete')
    expect(unfulfilledEthTx).toEqual(false)
  })

  it('returns the status as titlecase by default', () => {
    const jobRun = { status: 'pending_confirmations' } as JobRun

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Pending Confirmations')
    expect(unfulfilledEthTx).toEqual(false)
  })
})

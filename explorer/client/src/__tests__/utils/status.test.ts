import { partialAsFull } from '@chainlink/ts-helpers'
import status from '../../utils/status'
import { JobRun, TaskRun } from 'explorer/models'

const COMPLETED_ETHTX_WITHOUT_STATUS = partialAsFull<TaskRun>({
  type: 'ethtx',
  status: 'completed',
  transactionStatus: undefined,
})
const COMPLETED_ETHTX_WITH_STATUS = partialAsFull<TaskRun>({
  type: 'ethtx',
  status: 'completed',
  transactionStatus: 'fulfilledRunLog',
})

describe('utils/status', () => {
  it('is Pending not fulfilled for a status of "in_progress" without successful ethtx transaction', () => {
    const taskRuns: TaskRun[] = [COMPLETED_ETHTX_WITHOUT_STATUS]
    const jobRun = partialAsFull<JobRun>({
      status: 'in_progress',
      taskRuns,
    })

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Pending')
    expect(unfulfilledEthTx).toEqual(true)
  })

  it('is Pending for a status of "in_progress" with ethtx successful transaction', () => {
    const taskRuns: TaskRun[] = [COMPLETED_ETHTX_WITH_STATUS]
    const jobRun = partialAsFull<JobRun>({
      status: 'in_progress',
      taskRuns,
    })

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Pending')
    expect(unfulfilledEthTx).toEqual(false)
  })

  it('is Pending for a status of "in_progress" when there is no ethtx', () => {
    const jobRun = partialAsFull<JobRun>({ status: 'in_progress' })

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Pending')
    expect(unfulfilledEthTx).toEqual(false)
  })

  it('is Errored not fulfilled for a status of "error" without successful ethx transaction', () => {
    const taskRuns: TaskRun[] = [COMPLETED_ETHTX_WITHOUT_STATUS]
    const jobRun = partialAsFull<JobRun>({
      status: 'error',
      taskRuns,
    })

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Errored')
    expect(unfulfilledEthTx).toEqual(true)
  })

  it('is Errored for a status of "error" with ethx successful transaction', () => {
    const taskRuns: TaskRun[] = [COMPLETED_ETHTX_WITH_STATUS]
    const jobRun = partialAsFull<JobRun>({
      status: 'error',
      taskRuns,
    })

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Errored')
    expect(unfulfilledEthTx).toEqual(false)
  })

  it('is Errored for a status of "error"', () => {
    const jobRun = partialAsFull<JobRun>({ status: 'error' })

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Errored')
    expect(unfulfilledEthTx).toEqual(false)
  })

  it('is Complete not fulfilled for a status of "completed" without ethtx successful transaction', () => {
    const taskRuns: TaskRun[] = [COMPLETED_ETHTX_WITHOUT_STATUS]
    const jobRun = partialAsFull<JobRun>({
      status: 'completed',
      taskRuns,
    })

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Complete')
    expect(unfulfilledEthTx).toEqual(true)
  })

  it('is Complete for a status of "completed" without ethtx successful transaction', () => {
    const taskRuns: TaskRun[] = [COMPLETED_ETHTX_WITH_STATUS]
    const jobRun = partialAsFull<JobRun>({
      status: 'completed',
      taskRuns,
    })

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Complete')
    expect(unfulfilledEthTx).toEqual(false)
  })

  it('is Complete for a status of "completed"', () => {
    const jobRun = partialAsFull<JobRun>({ status: 'completed' })

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Complete')
    expect(unfulfilledEthTx).toEqual(false)
  })

  it('returns the status as titlecase by default', () => {
    const jobRun = partialAsFull<JobRun>({ status: 'pending_confirmations' })

    const [text, unfulfilledEthTx] = status(jobRun)
    expect(text).toEqual('Pending Confirmations')
    expect(unfulfilledEthTx).toEqual(false)
  })
})

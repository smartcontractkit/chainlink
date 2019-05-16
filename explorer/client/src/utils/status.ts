import { titleCase } from 'change-case'

const hasUnfulfilledEthTx = (jobRun: IJobRun) => {
  return (jobRun.taskRuns || []).some((tr: ITaskRun) => {
    return (
      tr.type === 'ethtx' &&
      tr.status === 'completed' &&
      tr.transactionStatus !== '0x1'
    )
  })
}

export default (jobRun: IJobRun): [string, boolean] => {
  let unfulfilledEthTx = false
  let text

  if (jobRun.status === 'in_progress') {
    text = 'Pending'
    unfulfilledEthTx = hasUnfulfilledEthTx(jobRun)
  } else if (jobRun.status === 'error') {
    text = 'Errored'
    unfulfilledEthTx = hasUnfulfilledEthTx(jobRun)
  } else if (jobRun.status === 'completed') {
    text = 'Complete'
    unfulfilledEthTx = hasUnfulfilledEthTx(jobRun)
  } else {
    text = titleCase(jobRun.status)
  }

  return [text, unfulfilledEthTx]
}

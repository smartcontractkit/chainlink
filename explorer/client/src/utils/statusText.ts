import { titleCase } from 'change-case'

const hasUnfullfilledEthTx = (jobRun: IJobRun) => {
  return (jobRun.taskRuns || []).some((tr: ITaskRun) => {
    return (
      tr.type === 'ethtx' &&
      tr.status === 'completed' &&
      tr.transactionStatus !== '0x1'
    )
  })
}

export default (jobRun: IJobRun) => {
  let unfullfilled = false
  let base

  if (jobRun.status === 'in_progress') {
    base = 'Pending'
    unfullfilled = hasUnfullfilledEthTx(jobRun)
  } else if (jobRun.status === 'error') {
    base = 'Failed'
    unfullfilled = hasUnfullfilledEthTx(jobRun)
  } else if (jobRun.status === 'completed') {
    base = 'Complete'
    unfullfilled = hasUnfullfilledEthTx(jobRun)
  } else {
    base = titleCase(jobRun.status)
  }

  if (unfullfilled) {
    return `${base} (but Not Fullfilled)`
  }
  return base
}

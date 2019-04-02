import React from 'react'
import Paper from '@material-ui/core/Paper'
import Table from '../Table'

interface IProps {
  jobRuns?: any[]
  className?: string
}

const HEADERS = ['Run ID', 'Job ID']

const List = ({ jobRuns, className }: IProps) => {
  let rows

  if (jobRuns) {
    rows = jobRuns.map((r: IJobRun) => [r.id, r.jobId])
  }

  return (
    <Paper className={className}>
      <Table headers={HEADERS} rows={rows} />
    </Paper>
  )
}

export default List

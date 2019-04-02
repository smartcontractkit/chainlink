import React from 'react'
import Paper from '@material-ui/core/Paper'
import Table from '../Table'
import { LinkColumn, TextColumn } from '../Table/TableCell'

interface IProps {
  jobRuns?: any[]
  className?: string
}

const HEADERS = ['Run ID', 'Job ID']

const List = ({ jobRuns, className }: IProps) => {
  let rows

  if (jobRuns) {
    rows = jobRuns.map((r: IJobRun) => {
      const idCol: LinkColumn = {
        type: 'link',
        text: r.id,
        to: `/job-runs/${r.id}`
      }
      const jobIdCol: TextColumn = { type: 'text', text: r.id }

      return [idCol, jobIdCol]
    })
  }

  return (
    <Paper className={className}>
      <Table headers={HEADERS} rows={rows} />
    </Paper>
  )
}

export default List

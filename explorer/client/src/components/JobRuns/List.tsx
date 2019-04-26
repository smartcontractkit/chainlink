import React from 'react'
import Paper from '@material-ui/core/Paper'
import Table from '../Table'
import { LinkColumn, TextColumn, TimeAgoColumn } from '../Table/TableCell'

interface IProps {
  jobRuns?: any[]
  className?: string
}

const HEADERS = ['Node', 'Run ID', 'Job ID', 'Created At']

const List = ({ jobRuns, className }: IProps) => {
  let rows

  if (jobRuns) {
    rows = jobRuns.map((r: IJobRun) => {
      const nodeCol: TextColumn = { type: 'text', text: r.chainlinkNode.name }
      const idCol: LinkColumn = {
        type: 'link',
        text: r.runId,
        to: `/job-runs/${r.id}`
      }
      const jobIdCol: TextColumn = { type: 'text', text: r.jobId }
      const createdAtCol: TimeAgoColumn = {
        type: 'time_ago',
        text: r.createdAt
      }

      return [nodeCol, idCol, jobIdCol, createdAtCol]
    })
  }

  return (
    <Paper className={className}>
      <Table headers={HEADERS} rows={rows} />
    </Paper>
  )
}

export default List

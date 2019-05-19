import React from 'react'
import Paper from '@material-ui/core/Paper'
import Table, { ChangePageEvent } from '../Table'
import { LinkColumn, TextColumn, TimeAgoColumn } from '../Table/TableCell'

interface IProps {
  currentPage: number
  onChangePage: (event: ChangePageEvent, page: number) => void
  jobRuns?: any[]
  count?: number
  emptyMsg?: string
  className?: string
}

const HEADERS = ['Node', 'Run ID', 'Job ID', 'Created At']

const List = (props: IProps) => {
  let rows

  if (props.jobRuns) {
    rows = props.jobRuns.map((r: IJobRun) => {
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
    <Paper className={props.className}>
      <Table
        headers={HEADERS}
        currentPage={props.currentPage}
        rows={rows}
        count={props.count}
        onChangePage={props.onChangePage}
        emptyMsg={props.emptyMsg}
      />
    </Paper>
  )
}

export default List

import React from 'react'
import Paper from '@material-ui/core/Paper'
import Hidden from '@material-ui/core/Hidden'
import Table, { ChangePageEvent } from '../Table'
import { LinkColumn, TextColumn, TimeAgoColumn } from '../Table/TableCell'

interface IProps {
  currentPage: number
  onChangePage: (event: ChangePageEvent, page: number) => void
  jobRuns?: IJobRun[]
  count?: number
  emptyMsg?: string
  className?: string
}

const DEFAULT_HEADERS = ['Node', 'Run ID', 'Job ID', 'Created At']
const MOBILE_HEADERS = ['Run ID', 'Job ID', 'Created At', 'Node']

const buildNodeCol = (jobRun: IJobRun): LinkColumn => {
  return {
    type: 'link',
    text: jobRun.chainlinkNode.name,
    to: jobRun.chainlinkNode.url || ''
  }
}

const buildIdCol = (jobRun: IJobRun): LinkColumn => {
  return {
    type: 'link',
    text: jobRun.runId,
    to: `/job-runs/${jobRun.id}`
  }
}

const buildJobIdCol = (jobRun: IJobRun): TextColumn => {
  return { type: 'text', text: jobRun.jobId }
}

const buildCreatedAtCol = (jobRun: IJobRun): TimeAgoColumn => {
  return {
    type: 'time_ago',
    text: jobRun.createdAt
  }
}

const mobileRows = (jobRuns: IJobRun[]) => {
  return jobRuns.map((r: IJobRun) => {
    return [
      buildIdCol(r),
      buildJobIdCol(r),
      buildCreatedAtCol(r),
      buildNodeCol(r)
    ]
  })
}

const defaultRows = (jobRuns: IJobRun[]) => {
  return jobRuns.map((r: IJobRun) => {
    return [
      buildNodeCol(r),
      buildIdCol(r),
      buildJobIdCol(r),
      buildCreatedAtCol(r)
    ]
  })
}

const List = (props: IProps) => {
  const jobRuns = props.jobRuns || []
  return (
    <Paper className={props.className}>
      <Hidden xsDown>
        <Table
          headers={DEFAULT_HEADERS}
          currentPage={props.currentPage}
          rows={defaultRows(jobRuns)}
          count={props.count}
          onChangePage={props.onChangePage}
          emptyMsg={props.emptyMsg}
        />
      </Hidden>
      <Hidden smUp>
        <Table
          headers={MOBILE_HEADERS}
          currentPage={props.currentPage}
          rows={mobileRows(jobRuns)}
          count={props.count}
          onChangePage={props.onChangePage}
          emptyMsg={props.emptyMsg}
        />
      </Hidden>
    </Paper>
  )
}

export default List

import React from 'react'
import Paper from '@material-ui/core/Paper'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Link from '../Link'

const Loading = () => {
  return (
    <TableRow>
      <TableCell component="th" scope="row" colSpan={2}>
        Loading...
      </TableCell>
    </TableRow>
  )
}

const Empty = () => {
  return (
    <TableRow>
      <TableCell component="th" scope="row" colSpan={2}>
        Hold the line! We&apos;re just getting started and haven&apos;t received
        any job runs yet.
      </TableCell>
    </TableRow>
  )
}

interface IRunsProps {
  jobRuns: IJobRun[]
}

const Runs = ({ jobRuns }: IRunsProps) => {
  return (
    <>
      {jobRuns.map((r: any, idx: number) => (
        <TableRow key={r.id}>
          <TableCell component="th" scope="row">
            <Link to={`/job-runs/${r.id}`}>{r.id}</Link>
          </TableCell>
          <TableCell>{r.jobId}</TableCell>
        </TableRow>
      ))}
    </>
  )
}

const renderBody = (jobRuns?: IJobRun[]) => {
  if (!jobRuns) {
    return <Loading />
  } else if (jobRuns.length === 0) {
    return <Empty />
  } else {
    return <Runs jobRuns={jobRuns} />
  }
}

interface IProps {
  jobRuns?: any[]
  className?: string
}

const List = ({ jobRuns, className }: IProps) => {
  return (
    <Paper className={className}>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>Run ID</TableCell>
            <TableCell>Job ID</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>{renderBody(jobRuns)}</TableBody>
      </Table>
    </Paper>
  )
}

export default List

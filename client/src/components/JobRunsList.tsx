import React from 'react'
import Paper from '@material-ui/core/Paper'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'

const Loading = () => {
  return <div>Loading...</div>
}

const Empty = () => {
  return (
    <div>
      Hold the line! We&apos;re just getting started and haven&apos;t received
      any job runs yet.
    </div>
  )
}

type RunsListProps = { jobRuns: IJobRun[] }

const RunsList = (props: RunsListProps) => {
  return (
    <Paper>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>Run ID</TableCell>
            <TableCell>Job ID</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {props.jobRuns.map((r: any, idx: number) => (
            <TableRow key={r.id}>
              <TableCell component="th" scope="row">
                {r.jobRunId}
              </TableCell>
              <TableCell>{r.jobId}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </Paper>
  )
}

type JobRunsListProps = { jobRuns?: any[] }

const JobRunsList = (props: JobRunsListProps) => {
  if (!props.jobRuns) {
    return <Loading />
  } else if (props.jobRuns.length === 0) {
    return <Empty />
  } else {
    return <RunsList jobRuns={props.jobRuns} />
  }
}

export default JobRunsList

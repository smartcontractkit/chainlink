import React from 'react'
import Paper from '@material-ui/core/Paper'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import { Link } from '@reach/router'

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

interface IRunsProps {
  jobRuns: IJobRun[]
}

const Runs = ({ jobRuns }: IRunsProps) => {
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
          {jobRuns.map((r: any, idx: number) => (
            <TableRow key={r.id}>
              <TableCell component="th" scope="row">
                <Link to={`/job-runs/${r.id}`}>{r.id}</Link>
              </TableCell>
              <TableCell>{r.jobId}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </Paper>
  )
}

interface IProps {
  jobRuns?: any[]
}

const List = (props: IProps) => {
  if (!props.jobRuns) {
    return <Loading />
  } else if (props.jobRuns.length === 0) {
    return <Empty />
  } else {
    return <Runs jobRuns={props.jobRuns} />
  }
}

export default List

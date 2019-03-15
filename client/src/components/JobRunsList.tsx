import React from 'react'
import Typography from '@material-ui/core/Typography'
import List from '@material-ui/core/List'
import ListItem from '@material-ui/core/ListItem'

const Loading = () => {
  return <div>Loading...</div>
}

const Empty = () => {
  return (
    <div>
      Hold the line! We're just getting started and haven't received any job
      runs yet.
    </div>
  )
}

type RunsListProps = { jobRuns: IJobRun[] }

const RunsList = (props: RunsListProps) => {
  return (
    <>
      <Typography variant="h3">Latest Runs</Typography>

      <List>
        {props.jobRuns.map((r, idx) => {
          return <ListItem key={idx} disableGutters>{r.requestId}</ListItem>
        })}
      </List>
    </>
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

import React from 'react'

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

type ListProps = { jobRuns: any[] }

const List = (props: ListProps) => {
  return (
    <>
      <h2>Latest Runs</h2>
      <ol>
        {props.jobRuns.map((r, idx) => {
          return <li key={idx}>{r.requestId}</li>
        })}
      </ol>
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
    return <List jobRuns={props.jobRuns} />
  }
}

export default JobRunsList

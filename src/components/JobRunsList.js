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

const List = ({ jobRuns }) => {
  return (
    <>
      <h2>Latest Runs</h2>
      <ol>
        {jobRuns.map((r, idx) => {
          return <li key={idx}>{r.requestId}</li>
        })}
      </ol>
    </>
  )
}

const JobRunsList = ({ jobRuns }) => {
  if (!jobRuns) {
    return <Loading />
  } else if (jobRuns.length === 0) {
    return <Empty />
  } else {
    return <List jobRuns={jobRuns} />
  }
}

export default JobRunsList

import React, { Component } from 'react'
import JobRunsList from '../components/JobRunsList'
import { getJobRuns } from '../api'

type Props = {}
type State = { jobRuns?: IJobRun[] }

class Home extends Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = { jobRuns: undefined }
  }

  componentDidMount() {
    getJobRuns().then((jrs: IJobRun[]) => {
      this.setState({ jobRuns: jrs })
    })
  }

  render() {
    return <JobRunsList jobRuns={this.state.jobRuns} />
  }
}

export default Home

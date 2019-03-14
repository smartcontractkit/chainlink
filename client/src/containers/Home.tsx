import React, { Component } from 'react'
import JobRunsList from '../components/JobRunsList'
import { getJobRuns } from '../api'

type Props = {}
type State = { jobRuns?: any[] }

class Home extends Component<Props, State> {
  constructor(props: any) {
    super(props)
    this.state = { jobRuns: undefined }
  }

  componentDidMount() {
    getJobRuns().then((jr: any) => {
      this.setState({ jobRuns: jr })
    })
  }

  render() {
    return <JobRunsList jobRuns={this.state.jobRuns} />
  }
}

export default Home

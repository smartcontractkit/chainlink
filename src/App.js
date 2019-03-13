import React, { Component } from 'react'
import './App.css'
import JobRunsList from './components/JobRunsList'

const getJobRuns = () => fetch('/api/v1/job_runs').then(r => r.json())

class App extends Component {
  constructor(props) {
    super(props)
    this.state = { jobRuns: null }
  }

  componentDidMount() {
    getJobRuns().then(jr => {
      this.setState({ jobRuns: jr })
    })
  }

  render() {
    return (
      <div className="App">
        <header className="App-header">
          <h1>LINK Stats</h1>
        </header>

        <JobRunsList jobRuns={this.state.jobRuns} />
      </div>
    )
  }
}

export default App

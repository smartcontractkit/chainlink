import React, { Component } from 'react'
import Card from '@material-ui/core/Card'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import { getJobs } from 'api'

const formatInitiators = (initiators) => (initiators.map(i => i.type).join(', '))

export class JobList extends Component {
  constructor (props) {
    super(props)
    this.state = {
      jobs: [],
      fetchError: false
    }
  }

  componentDidMount () {
    getJobs()
      .then(({data: jobs}) => this.setState({jobs: jobs}))
      .catch(_ => this.setState({fetchError: true}))
  }

  renderJobs () {
    if (this.state.fetchError) {
      return (
        <TableRow>
          <TableCell component='th' scope='row' colSpan={3}>
            There was an error fetching the jobs. Please reload the page.
          </TableCell>
        </TableRow>
      )
    } else {
      return this.state.jobs.map(j => {
        return (
          <TableRow key={j.id}>
            <TableCell component='th' scope='row'>
              {j.id}
            </TableCell>
            <TableCell>{j.attributes.createdAt}</TableCell>
            <TableCell>
              {formatInitiators(j.attributes.initiators)}
            </TableCell>
          </TableRow>
        )
      })
    }
  }

  render () {
    return (
      <Card>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Created</TableCell>
              <TableCell>Initiator</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {this.renderJobs()}
          </TableBody>
        </Table>
      </Card>
    )
  }
}

export default JobList

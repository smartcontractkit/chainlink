import React from 'react'
import PropType from 'prop-types'
import { withSiteData } from 'react-static'
import Typography from '@material-ui/core/Typography'
import Card from '@material-ui/core/Card'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import { withStyles } from '@material-ui/core/styles'
import url from 'url'
import 'isomorphic-unfetch'
import { parse as parseQueryString } from 'query-string'

const DEFAULT_CHAINLINK_PORT = 6688

const styles = theme => ({
  title: {
    marginTop: theme.spacing.unit * 5
  },
  card: {
    marginTop: theme.spacing.unit * 6
  }
})

const formatInitiators = (initiators) => (initiators.map(i => i.type).join(', '))

export class Jobs extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      jobs: [],
      fetchError: false
    }
  }

  componentDidMount () {
    const query = parseQueryString(this.props.location.search)
    const port = query.port || process.env.CHAINLINK_PORT || DEFAULT_CHAINLINK_PORT
    const jobsUrl = url.format({
      hostname: global.location.hostname,
      port: port,
      pathname: '/v2/specs'
    })

    global.fetch(jobsUrl, {credentials: 'include'})
      .then(response => response.json())
      .then(({data: jobs}) => this.setState({jobs: jobs}))
      .catch(_ => { this.setState({fetchError: true}) })
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
    const {classes} = this.props

    return (
      <div>
        <Typography variant='display2' color='inherit' className={classes.title}>
          Jobs
        </Typography>

        <Card className={classes.card}>
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
      </div>
    )
  }
}

Jobs.propTypes = {
  classes: PropType.object.isRequired,
  location: PropType.object.isRequired
}

export default withSiteData(withStyles(styles)(Jobs))

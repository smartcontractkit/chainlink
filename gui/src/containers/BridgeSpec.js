import React, { Component } from 'react'
import PropTypes from 'prop-types'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import PaddedCard from 'components/PaddedCard'
import Breadcrumb from 'components/Breadcrumb'
import BreadcrumbItem from 'components/BreadcrumbItem'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { withStyles } from '@material-ui/core/styles'
import { connect } from 'react-redux'
import { fetchBridgeSpec } from 'actions'

const styles = theme => ({
  title: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  },
  definitionTitle: {
    marginBottom: theme.spacing.unit * 3
  },
  breadcrumb: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  }
})

const renderBridgeSpec = props => (
  <Grid container spacing={40}>
    <Grid item xs={8}>
      <PaddedCard>
        <Typography variant='title' className={props.classes.definitionTitle}>
          Bridge Info:
        </Typography>
        <Grid>
          <Typography variant='subheading' color='textSecondary'>Name</Typography>
          <Typography variant='body1' color='inherit'>{props.name}</Typography>

          <Typography variant='subheading' color='textSecondary'>URL</Typography>
          <Typography variant='body1' color='inherit'>{props.url}</Typography>

          <Typography variant='subheading' color='textSecondary'>Confirmations</Typography>
          <Typography variant='body1' color='inherit'>{props.confirmations}</Typography>

          <Typography variant='subheading' color='textSecondary'>Minimum Contract Payment</Typography>
          <Typography variant='body1' color='inherit'>{props.minimumContractPayment}</Typography>

          <Typography variant='subheading' color='textSecondary'>Incoming Token</Typography>
          <Typography variant='body1' color='inherit'>{props.incomingToken}</Typography>

          <Typography variant='subheading' color='textSecondary'>Outgoing Token</Typography>
          <Typography variant='body1' color='inherit'>{props.outgoingToken}</Typography>
        </Grid>
      </PaddedCard>
    </Grid>
  </Grid>
)

const renderFetching = () => <div>Fetching...</div>

const renderDetails = props => {
  if (!props.fetching) {
    return (
      <React.Fragment>
        {renderBridgeSpec(props)}
      </React.Fragment>
    )
  } else {
    return renderFetching()
  }
}

export class BridgeSpec extends Component {
  componentDidMount () {
    this.props.fetchBridgeSpec(this.props.match.params.bridgeName)
  }

  render () {
    const { classes, name } = this.props
    return (
      <div>
        <Breadcrumb className={classes.breadcrumb}>
          <BreadcrumbItem href='/'>Dashboard</BreadcrumbItem>
          <BreadcrumbItem>></BreadcrumbItem>
          <BreadcrumbItem href='/bridges'>Bridges</BreadcrumbItem>
          <BreadcrumbItem>></BreadcrumbItem>
          <BreadcrumbItem>{name}</BreadcrumbItem>
        </Breadcrumb>
        <Typography variant='display2' color='inherit' className={classes.title}>
          Bridge Spec Details
        </Typography>
        {renderDetails(this.props)}
      </div>
    )
  }
}

BridgeSpec.propTypes = {
  classes: PropTypes.object.isRequired
}

const mapStateToProps = state => {
  return {
    name: state.bridgeSpec.name,
    url: state.bridgeSpec.url,
    confirmations: state.bridgeSpec.confirmations,
    minimumContractPayment: state.bridgeSpec.minimumContractPayment,
    incomingToken: state.bridgeSpec.incomingToken,
    outgoingToken: state.bridgeSpec.outgoingToken
  }
}

export const ConnectedBridgeSpec = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({fetchBridgeSpec})
)(BridgeSpec)

export default withStyles(styles)(ConnectedBridgeSpec)

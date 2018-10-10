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
import { bridgeSelector } from 'selectors'

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

const renderLoading = props => (
  <div>Loading...</div>
)

const renderLoaded = props => (
  <Grid container spacing={40}>
    <Grid item xs={8}>
      <PaddedCard>
        <Typography variant='title' className={props.classes.definitionTitle}>
          Bridge Info:
        </Typography>
        <Grid>
          <Typography variant='subheading' color='textSecondary'>Name</Typography>
          <Typography variant='body1' color='inherit'>{props.bridge.name}</Typography>

          <Typography variant='subheading' color='textSecondary'>URL</Typography>
          <Typography variant='body1' color='inherit'>{props.bridge.url}</Typography>

          <Typography variant='subheading' color='textSecondary'>Confirmations</Typography>
          <Typography variant='body1' color='inherit'>{props.bridge.confirmations}</Typography>

          <Typography variant='subheading' color='textSecondary'>Minimum Contract Payment</Typography>
          <Typography variant='body1' color='inherit'>{props.bridge.minimumContractPayment}</Typography>

          <Typography variant='subheading' color='textSecondary'>Incoming Token</Typography>
          <Typography variant='body1' color='inherit'>{props.bridge.incomingToken}</Typography>

          <Typography variant='subheading' color='textSecondary'>Outgoing Token</Typography>
          <Typography variant='body1' color='inherit'>{props.bridge.outgoingToken}</Typography>
        </Grid>
      </PaddedCard>
    </Grid>
  </Grid>
)

const renderDetails = props => props.bridge ? renderLoaded(props) : renderLoading(props)

export class BridgeSpec extends Component {
  componentDidMount () {
    this.props.fetchBridgeSpec(this.props.match.params.bridgeId)
  }

  render () {
    return (
      <div>
        <Breadcrumb className={this.props.classes.breadcrumb}>
          <BreadcrumbItem href='/'>Dashboard</BreadcrumbItem>
          <BreadcrumbItem>></BreadcrumbItem>
          <BreadcrumbItem href='/bridges'>Bridges</BreadcrumbItem>
          <BreadcrumbItem>></BreadcrumbItem>
          <BreadcrumbItem>{this.props.bridge && this.props.bridge.id}</BreadcrumbItem>
        </Breadcrumb>
        <Typography variant='display2' color='inherit' className={this.props.classes.title}>
          Bridge Spec Details
        </Typography>
        {renderDetails(this.props)}
      </div>
    )
  }
}

BridgeSpec.propTypes = {
  classes: PropTypes.object.isRequired,
  bridge: PropTypes.object
}

const mapStateToProps = (state, ownProps) => {
  return {
    bridge: bridgeSelector(state, ownProps.match.params.bridgeId)
  }
}

export const ConnectedBridgeSpec = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({fetchBridgeSpec})
)(BridgeSpec)

export default withStyles(styles)(ConnectedBridgeSpec)

import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Button } from '@material-ui/core'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import PaddedCard from 'components/PaddedCard'
import Breadcrumb from 'components/Breadcrumb'
import BreadcrumbItem from 'components/BreadcrumbItem'
import ReactStaticLinkComponent from 'components/ReactStaticLinkComponent'
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
  breadcrumb: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  },
  main: {
    marginTop: theme.spacing.unit * 5
  }
})

const renderLoading = props => (
  <div>Loading...</div>
)

const renderLoaded = props => (
  <PaddedCard>
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
  </PaddedCard>
)

const renderDetails = props => props.bridge ? renderLoaded(props) : renderLoading(props)

export class Show extends Component {
  componentDidMount () {
    this.props.fetchBridgeSpec(this.props.match.params.bridgeId)
  }

  render () {
    return (
      <Grid container>
        <Grid item xs={12}>
          <Breadcrumb className={this.props.classes.breadcrumb}>
            <BreadcrumbItem href='/'>Dashboard</BreadcrumbItem>
            <BreadcrumbItem>></BreadcrumbItem>
            <BreadcrumbItem href='/bridges'>Bridges</BreadcrumbItem>
            <BreadcrumbItem>></BreadcrumbItem>
            <BreadcrumbItem>{this.props.bridge && this.props.bridge.id}</BreadcrumbItem>
          </Breadcrumb>
        </Grid>
        <Grid item xs={12} md={12} xl={6}>
          <Grid container alignItems='center'>
            <Grid item xs={9}>
              <Typography variant='display2' color='inherit'>
                Bridge Info
              </Typography>
            </Grid>
            <Grid item xs={3}>
              <Grid container justify='flex-end'>
                <Grid item>
                  {this.props.bridge &&
                    <Button
                      variant='outlined'
                      color='primary'
                      component={ReactStaticLinkComponent}
                      to={`/bridges/${this.props.bridge.id}/edit`}
                    >
                      Edit
                    </Button>
                  }
                </Grid>
              </Grid>
            </Grid>
          </Grid>

          <div className={this.props.classes.main}>
            {renderDetails(this.props)}
          </div>
        </Grid>
      </Grid>
    )
  }
}

Show.propTypes = {
  classes: PropTypes.object.isRequired,
  bridge: PropTypes.object
}

const mapStateToProps = (state, ownProps) => {
  return {
    bridge: bridgeSelector(state, ownProps.match.params.bridgeId)
  }
}

export const ConnectedShow = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({fetchBridgeSpec})
)(Show)

export default withStyles(styles)(ConnectedShow)

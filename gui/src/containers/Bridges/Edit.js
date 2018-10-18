import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Link } from 'react-static'
import { connect } from 'react-redux'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import { Button } from '@material-ui/core'
import { withStyles } from '@material-ui/core/styles'
import PaddedCard from 'components/PaddedCard'
import BridgesForm from 'components/Bridges/Form'
import Breadcrumb from 'components/Breadcrumb'
import BreadcrumbItem from 'components/BreadcrumbItem'
import { bridgeSelector } from 'selectors'
import {
  fetchBridgeSpec,
  updateBridge
} from 'actions'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import ReactStaticLinkComponent from 'components/ReactStaticLinkComponent'

const styles = theme => ({
  title: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  },
  breadcrumb: {
    marginTop: theme.spacing.unit * 5
  }
})

const successNotification = ({name}) => (<React.Fragment>
  Successfully updated <Link to={`/bridges/${name}`}>{name}</Link>
</React.Fragment>)

const errorNotification = ({name}) => (
  <React.Fragment>Error updating {name}</React.Fragment>
)

export class Edit extends Component {
  componentDidMount () {
    const {fetchBridgeSpec, match} = this.props
    fetchBridgeSpec(match.params.bridgeId)
  }

  checkLoaded () {
    return this.props.bridge
  }

  onLoad (buildLoadedComponent) {
    if (this.checkLoaded()) {
      return buildLoadedComponent(this.props)
    }

    return <div>Loading...</div>
  }

  render () {
    const {bridge, classes, updateBridge} = this.props
    return (
      <Grid container>
        <Grid item xs={12}>
          <Breadcrumb className={classes.breadcrumb}>
            <BreadcrumbItem href='/'>Dashboard</BreadcrumbItem>
            <BreadcrumbItem>></BreadcrumbItem>
            <BreadcrumbItem href='/bridges'>Bridges</BreadcrumbItem>
            <BreadcrumbItem>></BreadcrumbItem>
            <BreadcrumbItem>{bridge && bridge.id}</BreadcrumbItem>
          </Breadcrumb>
        </Grid>
        <Grid item xs={12} md={12} xl={6}>
          <Grid container alignItems='center'>
            <Grid item xs={9}>
              <Typography variant='display2' color='inherit' className={classes.title}>
                Edit Bridge
              </Typography>
            </Grid>
            <Grid item xs={3}>
              <Grid container justify='flex-end'>
                <Grid item>
                  {bridge &&
                    <Button
                      variant='outlined'
                      color='primary'
                      component={ReactStaticLinkComponent}
                      to={`/bridges/${bridge.id}`}
                    >
                      Cancel
                    </Button>
                  }
                </Grid>
              </Grid>
            </Grid>
          </Grid>

          {this.onLoad(({bridge}) => (
            <PaddedCard>
              <BridgesForm
                actionText='Save Bridge'
                onSubmit={updateBridge}
                name={bridge.name}
                nameDisabled
                url={bridge.url}
                confirmations={bridge.confirmations}
                minimumContractPayment={bridge.minimumContractPayment}
                onSuccess={successNotification}
                onError={errorNotification}
              />
            </PaddedCard>
          ))}
        </Grid>
      </Grid>
    )
  }
}

Edit.propTypes = {
  classes: PropTypes.object.isRequired,
  bridge: PropTypes.object
}

const mapStateToProps = (state, ownProps) => {
  const bridge = bridgeSelector(state, ownProps.match.params.bridgeId)
  return {bridge}
}

export const ConnectedEdit = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({fetchBridgeSpec, updateBridge})
)(Edit)

export default withStyles(styles)(ConnectedEdit)

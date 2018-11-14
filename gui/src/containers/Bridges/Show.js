import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { withStyles } from '@material-ui/core/styles'
import { Button } from '@material-ui/core'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import Title from 'components/Title'
import PaddedCard from 'components/PaddedCard'
import Breadcrumb from 'components/Breadcrumb'
import BreadcrumbItem from 'components/BreadcrumbItem'
import ReactStaticLinkComponent from 'components/ReactStaticLinkComponent'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { fetchBridgeSpec } from 'actions'
import bridgeSelector from 'selectors/bridge'
import Content from 'components/Content'

const styles = theme => ({
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
      <Content>
        <Grid container>
          <Grid item xs={12}>
            <Breadcrumb>
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
                <Title>Bridge Info</Title>
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
      </Content>
    )
  }
}

Show.propTypes = {
  bridge: PropTypes.object
}

const mapStateToProps = (state, ownProps) => ({
  bridge: bridgeSelector(state, ownProps.match.params.bridgeId)
})

export const ConnectedShow = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({fetchBridgeSpec})
)(Show)

export default withStyles(styles)(ConnectedShow)

import React, { Component } from 'react'
import { withSiteData } from 'react-static'
import { connect } from 'react-redux'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'
import { Button } from '@material-ui/core'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import BridgeList from 'components/BridgeList'
import ReactStaticLinkComponent from 'components/ReactStaticLinkComponent'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { bridgesSelector } from 'selectors'
import { fetchBridges } from 'actions'

const styles = theme => ({
  title: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  }
})

export class Index extends Component {
  componentDidMount () {
    this.props.fetchBridges(1, this.props.pageSize)
  }

  render () {
    const { bridges, bridgeCount, classes, pageSize, bridgesError, fetchBridges } = this.props

    return (
      <div>
        <Grid container spacing={8} alignItems='center'>
          <Grid item xs={9}>
            <Typography variant='display2' color='inherit' className={classes.title}>
              Bridges
            </Typography>
          </Grid>
          <Grid item xs={3}>
            <Grid container justify='flex-end' >
              <Grid item>
                <Button
                  variant='outlined'
                  color='primary'
                  component={ReactStaticLinkComponent}
                  to={'/bridges/new'}
                >
                  New Bridge
                </Button>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
        <Grid container spacing={40}>
          <Grid item xs={12}>
            <BridgeList
              bridges={bridges}
              bridgeCount={bridgeCount}
              pageSize={pageSize}
              error={bridgesError}
              fetchBridges={fetchBridges}
            />
          </Grid>
        </Grid>
      </div>
    )
  }
}

Index.propTypes = {
  bridgeCount: PropTypes.number.isRequired,
  bridges: PropTypes.array.isRequired,
  bridgesError: PropTypes.string,
  classes: PropTypes.object.isRequired,
  pageSize: PropTypes.number
}

Index.defaultProps = {
  pageSize: 10
}

const mapStateToProps = state => {
  let bridgesError
  if (state.bridges.networkError) {
    bridgesError = 'There was an error fetching the bridges. Please reload the page.'
  }

  return {
    bridgeCount: state.bridges.count,
    bridges: bridgesSelector(state),
    bridgesError: bridgesError
  }
}

export const ConnectedIndex = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({fetchBridges})
)(Index)

export default withSiteData(withStyles(styles)(ConnectedIndex))

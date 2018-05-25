import React, { Component } from 'react'
import PropType from 'prop-types'
import Grid from '@material-ui/core/Grid'
import TokenBalance from 'components/TokenBalance'
import { withStyles } from '@material-ui/core/styles'
import { getAccountBalance } from 'api'

const styles = theme => ({
  tokenBalance: {
    paddingTop: theme.spacing.unit * 5,
    paddingBottom: theme.spacing.unit * 5,
    paddingLeft: theme.spacing.unit * 3,
    paddingRight: theme.spacing.unit * 3
  }
})

export class AccountBalance extends Component {
  constructor (props) {
    super(props)
    this.state = {
      linkBalance: null,
      ethBalance: null
    }
  }

  componentDidMount () {
    getAccountBalance()
      .then(({data}) => {
        this.setState({
          ethBalance: data.attributes.ethBalance,
          linkBalance: data.attributes.linkBalance
        })
      })
  }

  render () {
    const { classes } = this.props

    return (
      <Grid container spacing={24}>
        <Grid item xs={12}>
          <TokenBalance
            title='Ethereum'
            value={this.state.ethBalance}
            className={classes.tokenBalance}
          />
        </Grid>
        <Grid item xs={12}>
          <TokenBalance
            title='Link'
            value={this.state.linkBalance}
            className={classes.tokenBalance}
          />
        </Grid>
      </Grid>
    )
  }
}

AccountBalance.propTypes = {
  classes: PropType.object.isRequired
}

export default withStyles(styles)(AccountBalance)

import React, { Component } from 'react'
import { connect } from 'react-redux'
import { Redirect } from 'react-router'
import { withStyles } from '@material-ui/core/styles'
import Button from '@material-ui/core/Button'
import Typography from '@material-ui/core/Typography'
import TextField from '@material-ui/core/TextField'
import { Grid } from '@material-ui/core'
import { submitSignIn } from 'actions'
import Title from 'components/Title'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'

const styles = theme => ({
  button: {
    margin: theme.spacing.unit * 5
  }
})

export class SignIn extends Component {
  state = { email: '', password: '' }

  handleChange = name => event => {
    this.setState({
      [name]: event.target.value
    })
  }

  onSubmit = e => {
    e.preventDefault()
    const { email, password } = this.state
    this.props.submitSignIn({ email: email, password: password })
  }

  render () {
    const { classes, fetching, authenticated } = this.props
    const enabled = this.state.email.length > 0 && this.state.password.length > 0
    if (authenticated) {
      return <Redirect to='/' />
    }
    return (
      <form noValidate onSubmit={this.onSubmit}>
        <Grid container alignItems='center' direction='column'>
          <Title>Sign In to Chainlink</Title>
          <TextField
            id='email'
            label='Email'
            className={classes.textField}
            margin='normal'
            value={this.state.email}
            onChange={this.handleChange('email')}
          />
          <TextField
            id='password'
            label='Password'
            className={classes.textField}
            type='password'
            autoComplete='password'
            margin='normal'
            value={this.state.password}
            onChange={this.handleChange('password')}
          />
          <Button
            type='submit'
            disabled={!enabled}
            variant='contained'
            color='primary'
            className={classes.button}
          >
            Sign In
          </Button>
          {fetching && (
            <Typography variant='body1' color='textSecondary'>
              Signing in...
            </Typography>
          )}
        </Grid>
      </form>
    )
  }
}

const mapStateToProps = state => ({
  fetching: state.authentication.fetching,
  authenticated: state.authentication.allowed
})

export const ConnectedSignIn = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({submitSignIn})
)(SignIn)

export default withStyles(styles)(ConnectedSignIn)

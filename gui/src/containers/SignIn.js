import React, { Component } from 'react'
import { Redirect } from 'react-router'
import PropTypes from 'prop-types'
import Button from '@material-ui/core/Button'
import Typography from '@material-ui/core/Typography'
import { withStyles } from '@material-ui/core/styles'
import TextField from '@material-ui/core/TextField'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import { matchRoute, submitSignIn } from 'actions'

const styles = theme => ({
  title: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  },
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

  onSubmit = (e) => {
    e.preventDefault()
    const { email, password } = this.state
    this.props.submitSignIn({email: email, password: password})
  }

  render () {
    const {classes, fetching, authenticated} = this.props
    const enabled = this.state.email.length > 0 && this.state.password.length > 0
    if (authenticated) {
      return <Redirect to='/' />
    }
    return (
      <form className={classes.container}
        noValidate
        onSubmit={this.onSubmit}
        align='center'>
        <Typography variant='display2' color='inherit' align='center' className={classes.title}>
          Sign In to Chainlink
        </Typography>

        <div>
          <TextField id='email' label='Email' className={classes.textField} margin='normal'
            value={this.state.email} onChange={this.handleChange('email')} />
        </div>
        <div>
          <TextField id='password' label='Password' className={classes.textField} type='password'
            autoComplete='password' margin='normal'
            value={this.state.password} onChange={this.handleChange('password')} />
        </div>
        <div>
          <Button type='submit' disabled={!enabled} variant='contained' color='primary' className={classes.button}>
            Sign In
          </Button>
        </div>
        { fetching &&
        <Typography variant='body1' color='textSecondary'>
          Signing in...
        </Typography>
        }
      </form>
    )
  }
}

SignIn.propTypes = {
  classes: PropTypes.object.isRequired
}

const mapStateToProps = state => ({
  fetching: state.session.fetching,
  authenticated: state.session.authenticated
})

const mapDispatchToProps = (dispatch, ownProps) => {
  dispatch(matchRoute(ownProps.match))
  return bindActionCreators({ submitSignIn }, dispatch)
}

export const ConnectedSignIn = connect(mapStateToProps, mapDispatchToProps)(SignIn)

export default withStyles(styles)(ConnectedSignIn)

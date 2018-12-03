import React from 'react'
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
import { useHooks, useState } from 'use-react-hooks'

const styles = theme => ({
  button: {
    margin: theme.spacing.unit * 5
  },
  title: {
    marginTop: theme.spacing.unit * 5
  }
})

export const SignIn = useHooks((props) => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const handleChange = name => event => {
    if (name === 'email') setEmail(event.target.value)
    if (name === 'password') setPassword(event.target.value)
  }
  const onSubmit = e => {
    e.preventDefault()
    props.submitSignIn({ email, password })
  }
  const { classes, fetching, authenticated } = props
  const enabled = email.length > 0 && password.length > 0

  if (authenticated) return <Redirect to='/' />
  return (
    <form noValidate onSubmit={onSubmit}>
      <Grid container alignItems='center' direction='column'>
        <Title className={classes.title}>Sign In to Chainlink</Title>
        <TextField
          id='email'
          label='Email'
          className={classes.textField}
          margin='normal'
          value={email}
          onChange={handleChange('email')}
        />
        <TextField
          id='password'
          label='Password'
          className={classes.textField}
          type='password'
          autoComplete='password'
          margin='normal'
          value={password}
          onChange={handleChange('password')}
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
)

const mapStateToProps = state => ({
  fetching: state.authentication.fetching,
  authenticated: state.authentication.allowed
})

export const ConnectedSignIn = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({submitSignIn})
)(SignIn)

export default withStyles(styles)(ConnectedSignIn)

import React, { useState } from 'react'
import { connect } from 'react-redux'
import { Redirect } from 'react-router-dom'
import Button from 'components/Button'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Typography from '@material-ui/core/Typography'
import TextField from '@material-ui/core/TextField'
import { Grid } from '@material-ui/core'
import { hot } from 'react-hot-loader'
import { submitSignIn } from 'actions'
import HexagonLogo from 'components/Logos/Hexagon'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { get } from 'utils/storage'
import { makeStyles } from '@material-ui/styles'

const useStyles = makeStyles(theme => {
  return {
    container: {
      height: '100%'
    },
    cardContent: {
      paddingTop: theme.spacing(6),
      paddingLeft: theme.spacing(4),
      paddingRight: theme.spacing(4),
      '&:last-child': {
        paddingBottom: theme.spacing(6)
      }
    },
    headerRow: {
      textAlign: 'center'
    },
    error: {
      backgroundColor: theme.palette.error.light,
      marginTop: theme.spacing(2)
    },
    errorText: {
      color: theme.palette.error.main
    }
  }
})

export const SignIn = props => {
  document.title = 'Sign In'
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
  const { fetching, authenticated, errors } = props
  const classes = useStyles()

  const hasPrevState = Object.keys(get('persistURL')).length !== 0
  if (authenticated)
    return <Redirect to={(hasPrevState && get('persistURL')) || '/'} />
  return (
    <Grid
      container
      justify="center"
      alignItems="center"
      className={classes.container}
      spacing={0}
    >
      <Grid item xs={10} sm={6} md={4} lg={3} xl={2}>
        <Card>
          <CardContent className={classes.cardContent}>
            <form noValidate onSubmit={onSubmit}>
              <Grid container spacing={1}>
                <Grid item xs={12}>
                  <Grid container spacing={0}>
                    <Grid item xs={12} className={classes.headerRow}>
                      <HexagonLogo width={50} />
                    </Grid>
                    <Grid item xs={12} className={classes.headerRow}>
                      <Typography variant="h5">Operator</Typography>
                    </Grid>
                  </Grid>
                </Grid>

                {errors.length > 0 &&
                  errors.map(({ props }, idx) => {
                    return (
                      <Grid item xs={12} key={idx}>
                        <Card raised={false} className={classes.error}>
                          <CardContent>
                            <Typography
                              variant="body1"
                              className={classes.errorText}
                            >
                              {props.msg}
                            </Typography>
                          </CardContent>
                        </Card>
                      </Grid>
                    )
                  })}

                <Grid item xs={12}>
                  <TextField
                    id="email"
                    label="Email"
                    margin="normal"
                    value={email}
                    onChange={handleChange('email')}
                    error={errors.length > 0}
                    variant="outlined"
                    fullWidth
                  />
                </Grid>
                <Grid item xs={12}>
                  <TextField
                    id="password"
                    label="Password"
                    type="password"
                    autoComplete="password"
                    margin="normal"
                    value={password}
                    onChange={handleChange('password')}
                    error={errors.length > 0}
                    variant="outlined"
                    fullWidth
                  />
                </Grid>
                <Grid item xs={12}>
                  <Grid container spacing={0} justify="center">
                    <Grid item>
                      <Button type="submit" variant="primary">
                        Access Account
                      </Button>
                    </Grid>
                  </Grid>
                </Grid>
                {fetching && (
                  <Typography variant="body1" color="textSecondary">
                    Signing in...
                  </Typography>
                )}
              </Grid>
            </form>
          </CardContent>
        </Card>
      </Grid>
    </Grid>
  )
}

const mapStateToProps = state => ({
  fetching: state.authentication.fetching,
  authenticated: state.authentication.allowed,
  errors: state.notifications.errors
})

export const ConnectedSignIn = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ submitSignIn })
)(SignIn)

export default hot(module)(ConnectedSignIn)

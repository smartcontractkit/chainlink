import React, { Component } from 'react'
import { createGenerateClassName } from '@material-ui/core/styles'
import extractBuildInfo from './src/utils/extractBuildInfo'

const buildInfo = extractBuildInfo()
const generateClassName = createGenerateClassName()

export default {
  getSiteData: () => ({
    title: 'Chainlink'
  }),
  getRoutes: async () => {
    return [
      { path: '/' },
      { path: '/jobs' },
      { path: '/jobs/page/_jobPage_' },
      { path: '/jobs/new' },
      { path: '/jobs/_jobSpecId_' },
      { path: '/jobs/_jobSpecId_/definition' },
      { path: '/jobs/_jobSpecId_/runs' },
      { path: '/jobs/_jobSpecId_/runs/page/_jobRunsPage_' },
      { path: '/jobs/_jobSpecId_/runs/id/_jobRunId_' },
      { path: '/jobs/_jobSpecId_/runs/id/_jobRunId_/json' },
      { path: '/bridges' },
      { path: '/bridges/page/_bridgePage_' },
      { path: '/bridges/new' },
      { path: '/bridges/_bridgeId_' },
      { path: '/bridges/_bridgeId_/edit' },
      {
        path: '/config',
        getData: () => buildInfo
      },
      { path: '/signin' },
      { path: '/signout' },
      { path: '404', component: 'src/containers/404' }
    ]
  },
  beforeRenderToElement: (render, Comp) => render(Comp),
  plugins: [
    ['react-static-plugin-jss', { providerProps: { generateClassName } }],
    ['react-static-plugin-react-router']
  ],
  Document: class CustomHtml extends Component {
    render () {
      const { Html, Head, Body, children } = this.props
      return (
        <Html>
          <Head>
            <meta charSet='UTF-8' />
            <meta name='viewport' content='width=device-width, initial-scale=1' />
            <link
              href='https://fonts.googleapis.com/css?family=Roboto:300,400,500'
              rel='stylesheet'
            />
            <link href='https://fonts.googleapis.com/icon?family=Material+Icons' rel='stylesheet' />
            <link href='/favicon.ico' rel='shortcut icon' type='image/x-icon' />
          </Head>
          <Body>{children}</Body>
        </Html>
      )
    }
  }
}

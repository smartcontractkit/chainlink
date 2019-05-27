import React, { Component } from 'react'
import OS from 'os'

const MAX_EXPORT_HTML_THREADS =
  process.env.MAX_EXPORT_HTML_THREADS &&
  parseInt(process.env.MAX_EXPORT_HTML_THREADS, 10)
const CORES = Math.max(OS.cpus().length, 1)

export default {
  maxThreads: MAX_EXPORT_HTML_THREADS || CORES,
  getSiteData: () => ({
    title: 'Chainlink'
  }),
  getRoutes: async () => {
    return [{ path: '404', component: 'src/containers/404' }]
  },
  beforeRenderToElement: (render, Comp) => render(Comp),
  plugins: [
    'react-static-plugin-react-router',
    'react-static-plugin-typescript'
  ],
  Document: class CustomHtml extends Component {
    render() {
      const { Html, Head, Body, children } = this.props
      return (
        <Html>
          <Head>
            <meta charSet="UTF-8" />
            <meta
              name="viewport"
              content="width=device-width, initial-scale=1"
            />
            <link
              href="https://fonts.googleapis.com/css?family=Roboto:300,400,500"
              rel="stylesheet"
            />
            <link
              href="https://fonts.googleapis.com/icon?family=Material+Icons"
              rel="stylesheet"
            />
            <link href="/favicon.ico" rel="shortcut icon" type="image/x-icon" />
          </Head>
          <Body>{children}</Body>
        </Html>
      )
    }
  }
}

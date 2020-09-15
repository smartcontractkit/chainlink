import React from 'react'
import { Button } from 'antd'
import { withRouter, RouteComponentProps, Link } from 'react-router-dom'
import ReactGA from 'react-ga'
import ChainlinkLogo from '../shared/ChainlinkLogo'

interface Props extends RouteComponentProps {}

const Header: React.FC<Props> = ({ location }) => {
  return (
    <header className="header">
      <div className="header__main-nav">
        <Link to="/">
          <div className="header__logotype">
            <ChainlinkLogo />
            <h1>Chainlink</h1>
          </div>
        </Link>
        {location.pathname !== '/' && (
          <Link to={`/`}>
            <Button type="primary" ghost icon="left">
              Back to listing
            </Button>
          </Link>
        )}
      </div>
      <div className="header__secondary-nav">
        {location.pathname !== '/' && (
          <a
            onClick={() =>
              ReactGA.event({
                category: 'Form Conversion',
                action: 'Click on Button',
                label: 'Integrate with Chainlink',
              })
            }
            href="https://chainlinkcommunity.typeform.com/to/XcgLVP"
            target="_blank"
            rel="noopener noreferrer"
          >
            <Button type="primary" shape="round">
              Integrate with Chainlink
            </Button>
          </a>
        )}
      </div>
    </header>
  )
}

export default withRouter(Header)

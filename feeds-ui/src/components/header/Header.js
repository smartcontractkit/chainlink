import React from 'react'
import { Link } from 'react-router-dom'
import { Button } from 'antd'
import { withRouter } from 'react-router'
import ChainlinkLogo from 'components/shared/ChainlinkLogo'
import ReactGA from 'react-ga'

const Header = ({ location }) => {
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
      <div className="header__secondary-nav"></div>
      {location.pathname !== '/' && (
        <a
          onClick={() =>
            ReactGA.event({
              category: 'Conversion',
              action: 'Click on Email Button',
              label: 'Integrate with Chainlink',
            })
          }
          href="https://chainlink.typeform.com/to/gEwrPO"
          target="_blank"
          rel="noopener noreferrer"
        >
          <Button type="primary" shape="round">
            Integrate with Chainlink
          </Button>
        </a>
      )}
    </header>
  )
}

export default withRouter(Header)

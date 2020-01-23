import React from 'react'
import { Link } from 'react-router-dom'
import { Button } from 'antd'
import { withRouter } from 'react-router'
import ChainlinkLogo from 'components/shared/ChainlinkLogo'

const Header = ({ location }) => {
  return (
    <header className="header">
      <Link to="/">
        <div className="header-logotype">
          <ChainlinkLogo />
          <h1>Chainlink</h1>
        </div>
      </Link>
      <div className="header-menu">
        {location.pathname !== '/' && (
          <Link to={`/`}>
            <Button type="primary" ghost icon="left">
              Back to listing
            </Button>
          </Link>
        )}
      </div>
    </header>
  )
}

export default withRouter(Header)

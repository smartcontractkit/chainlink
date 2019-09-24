import React from 'react'
import { Link } from 'react-router-dom'
// import { Button } from 'antd'

import ChainlinkLogo from 'components/shared/ChainlinkLogo'

const Header = () => {
  return (
    <header className="header">
      <Link to="/">
        <div className="header-logotype">
          <ChainlinkLogo />
          <h1>
            Chainlink <span className="header-sub-name">Reference Data</span>
          </h1>
        </div>
      </Link>
      {/* <div className="header-menu">
        <div className="header-menu--link">
          <Link to={`/`}>ETH / USD</Link>
        </div>
        <Button ghost type="primary" style={{ marginLeft: 15 }}>
          <Link to={`/create`}>Create</Link>
        </Button>
      </div> */}
    </header>
  )
}

export default Header

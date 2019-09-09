import React from 'react'

import ChainlinkLogo from 'components/shared/ChainlinkLogo'

const Header = () => {
  return (
    <header className="header">
      <div className="header-logotype">
        <ChainlinkLogo />
        <h1>
          Chainlink <span className="header-sub-name">ETH/USD Aggregation</span>
        </h1>
      </div>
    </header>
  )
}

export default Header

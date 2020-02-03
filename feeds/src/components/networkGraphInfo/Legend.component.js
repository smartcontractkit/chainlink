import React from 'react'
import { Icon } from 'antd'

const Legend = () => (
  <div className="network-graph-legend">
    <h3 className="network-graph-legend__title">Chart Legend</h3>

    <div className="network-graph-legend__item">
      <div className="network-graph-legend__symbol">
        <div className="network-graph-legend__sc">
          <Icon type="safety-certificate" />
        </div>
      </div>
      <div className="network-graph-legend__description">
        Trusted answer in Aggregator Smart Contract
      </div>
    </div>

    <div className="network-graph-legend__item">
      <div className="network-graph-legend__symbol">
        <div className="network-graph-legend__oracle--fetching"></div>
      </div>
      <div className="network-graph-legend__description">
        Oracle - fetching external data
      </div>
    </div>

    <div className="network-graph-legend__item">
      <div className="network-graph-legend__symbol">
        <div className="network-graph-legend__oracle--fulfilled"></div>
      </div>
      <div className="network-graph-legend__description">
        Oracle - request fulfilled
      </div>
    </div>

    <div className="network-graph-legend__item">
      <div className="network-graph-legend__symbol">
        <svg width="20px" height="20px">
          <line
            className="network-graph-legend__line--fetching"
            x1="0"
            y1="10"
            x2="20"
            y2="10"
          />
        </svg>
      </div>
      <div className="network-graph-legend__description">
        Smart contract is waiting for response from oracle
      </div>
    </div>

    <div className="network-graph-legend__item">
      <div className="network-graph-legend__symbol">
        <svg width="20px" height="20px">
          <line
            className="network-graph-legend__line--fulfilled"
            x1="0"
            y1="10"
            x2="20"
            y2="10"
          />
        </svg>
      </div>
      <div className="network-graph-legend__description">
        Smart contract received answer
      </div>
    </div>
  </div>
)

export default Legend

import React from 'react'
import { Row, Col, Button } from 'antd'
import aave from 'assets/sponsors/aave.png'
import ampleforth from 'assets/sponsors/ampleforth.png'
import loopring from 'assets/sponsors/loopring.png'
import nexusmutual from 'assets/sponsors/nexusmutual.png'
import setprotocol from 'assets/sponsors/setprotocol.png'
import synthetix from 'assets/sponsors/synthetix.png'
import ReactGA from 'react-ga'

const grid = { xs: 24, sm: 12, md: 8 }
const list = [
  {
    name: 'Synthetix',
    url: 'https://www.synthetix.io/',
    src: synthetix,
  },
  {
    name: 'Loopring',
    url: 'https://loopring.org/',
    src: loopring,
  },
  {
    name: 'Aavve',
    url: 'https://aave.com/',
    src: aave,
  },
  {
    name: 'Ampleforth',
    url: 'https://www.ampleforth.org/',
    src: ampleforth,
  },
  {
    name: 'Set Protocol',
    url: 'https://www.tokensets.com/',
    src: setprotocol,
  },
  {
    name: 'Nexusmutual',
    url: 'https://nexusmutual.io',
    src: nexusmutual,
  },
]

const Logo = ({ node }) => (
  <a
    className="logo-item"
    href={node.url}
    target="_blank"
    rel="noopener noreferrer"
  >
    <img
      alt={node.name}
      title={node.name}
      src={node.src}
      className={node.name}
    />
  </a>
)

const SponsorsLogos = () => (
  <section className="logos-wrapper">
    <div className="cta-integrate">
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
        <Button type="primary" shape="round" size="large">
          Integrate with Chainlink
        </Button>
      </a>
    </div>
    <h3>Made possible and sponsored by</h3>
    <div className="logos">
      <Row gutter={18} type="flex" justify="space-around">
        {list.map((node, i) => (
          <Col key={i} {...grid}>
            <Logo node={node} />
          </Col>
        ))}
      </Row>
    </div>
  </section>
)

export default SponsorsLogos

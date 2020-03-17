import React from 'react'
import { Row, Col, Button } from 'antd'
import aave from 'assets/sponsors/aave.png'
import ampleforth from 'assets/sponsors/ampleforth.png'
import loopring from 'assets/sponsors/loopring.png'
import nexusmutual from 'assets/sponsors/nexusmutual.png'
import openlaw from 'assets/sponsors/openlaw.png'
import setprotocol from 'assets/sponsors/setprotocol.png'
import synthetix from 'assets/sponsors/synthetix.png'
import dmm from 'assets/sponsors/dmm.png'
import bzx from 'assets/sponsors/bzx.png'
import haven from 'assets/sponsors/haven.png'

import ReactGA from 'react-ga'

interface Sponsor {
  name: string
  url: string
  src: string
}

const grid = { xs: 24, sm: 12, md: 8 }
const list: Sponsor[] = [
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
  {
    name: 'OpenLaw',
    url: 'https://www.openlaw.io',
    src: openlaw,
  },
  {
    name: 'DMM',
    url: 'https://defimoneymarket.com',
    src: dmm,
  },
  {
    name: 'bZx',
    url: 'https://bzx.network',
    src: bzx,
  },
  {
    name: 'Haven Protocol',
    url: 'https://havenprotocol.org',
    src: haven,
  },
]

interface LogoProps {
  item: Sponsor
}

const Logo: React.FC<LogoProps> = ({ item }) => (
  <a
    className="logo-item"
    href={item.url}
    target="_blank"
    rel="noopener noreferrer"
  >
    <img
      alt={item.name}
      title={item.name}
      src={item.src}
      className={item.name}
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
    <h3>Used and Sponsored By</h3>
    <div className="logos">
      <Row gutter={18} type="flex" justify="space-around">
        {list.map((sponsor: Sponsor, i: number) => (
          <Col key={i} {...grid}>
            <Logo item={sponsor} />
          </Col>
        ))}
      </Row>
    </div>
  </section>
)

export default SponsorsLogos

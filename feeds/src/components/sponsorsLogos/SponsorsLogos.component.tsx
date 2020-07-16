import React from 'react'
import { Row, Col, Button } from 'antd'
import ReactGA from 'react-ga'
import { sponsorList, SponsorListItem } from '../../assets/sponsors'

const grid = { xs: 24, sm: 12, md: 8 }

interface LogoProps {
  item: SponsorListItem
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
      src={item.imageLg}
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
        {sponsorList.map((sponsor: SponsorListItem, i: number) => (
          <Col key={i} {...grid}>
            <Logo item={sponsor} />
          </Col>
        ))}
      </Row>
    </div>
  </section>
)

export default SponsorsLogos

import React from 'react'
import { Listing } from 'components/listing'
import { ReactComponent as Aave } from 'assets/aave.svg'
import { ReactComponent as Loopring } from 'assets/loopring.svg'
import { ReactComponent as Synthetix } from 'assets/synthetix.svg'
import ampleforth from 'assets/ampleforth.png'
import { Button } from 'antd'
import { Header } from 'components/header'
import ReactGA from 'react-ga'

const LangingPage = () => (
  <div className="page-wrapper landing-page">
    <div className="page-container">
      <Header />
    </div>
    <div className="head">
      <div className="head__title">
        <h1>
          PRICE <br />
          REFERENCE <br />
          DATA{' '}
        </h1>
        <div className="square"></div>
      </div>
    </div>
    <div className="page-container">
      <section>
        <h3>Decentralized Oracle Networks for Price Reference Data</h3>
        <p>
          The Chainlink Network provides the largest collection of secure and
          decentralized on-chain price reference data available. Composed of
          security reviewed, sybil resistant and fully independent nodes which
          are run by leading blockchain devops and security teams. Creating a
          shared global resource which is supported by a growing list of top
          DeFi Dapps.
        </p>
        <p>
          Please feel free to look into the details of each Decentralized Oracle
          Network listed below. You can easily use these oracle networks to
          quickly and securely launch, add more capabilities to and/or just
          greatly improve the security of your smart contracts.
        </p>
      </section>
      <section>
        <Listing />
      </section>
    </div>

    <section className="supporters-wrapper">
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
      <h3>Made possible and supported by</h3>
      <div className="supporters">
        <a
          href="https://www.synthetix.io/"
          target="_blank"
          rel="noopener noreferrer"
        >
          <Synthetix />
        </a>

        <a
          href="https://loopring.org/"
          target="_blank"
          rel="noopener noreferrer"
        >
          <Loopring className="loopring" />
        </a>

        <a href="https://aave.com/" target="_blank" rel="noopener noreferrer">
          <Aave />
        </a>
        <a
          href="https://www.ampleforth.org/"
          target="_blank"
          rel="noopener noreferrer"
        >
          <img
            alt="Ampleforth"
            title="Ampleforth"
            src={ampleforth}
            className="ampleforth"
          />
        </a>
      </div>
    </section>
  </div>
)

export default LangingPage

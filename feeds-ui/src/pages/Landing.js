import React from 'react'
import { Listing } from 'components/listing'
import { ReactComponent as Aave } from 'assets/aave.svg'
import { ReactComponent as Loopring } from 'assets/loopring.svg'
import { ReactComponent as Synthetix } from 'assets/synthetix.svg'
import ampleforth from 'assets/ampleforth.png'
import { Button } from 'antd'
import { Header } from 'components/header'

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
          decentralized on-chain price reference data available.
          Chainlink&apos;s decentralized oracle networks are composed of
          security reviewed, sybil resistant and fully independent nodes which
          are run by leading blockchain devops and security teams.
        </p>
        <p> </p>
        <p>
          Chainlink Decentralized Oracle Networks for Price Reference Data are a
          shared community resource supported by its users, who pay less for
          using these oracle networks than it would take for them to broadcast
          the same data individually, while benefiting from a twenty times
          increase in the security created by the decentralization of oracle
          networks.
        </p>
        <p>
          Please feel free to look into the operational details of each
          Chainlink Decentralized Oracle Network on this page and easily start
          using them here.
        </p>
        <br />
        <a href="mailto:support@smartcontract.com?subject=Price Reference Data">
          <Button type="primary" shape="round" size="large">
            Access Chainlink Oracles
          </Button>
        </a>
      </section>
      <section>
        <Listing />
      </section>
    </div>

    <section className="supporters-wrapper">
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
          <img alt="Ampleforth" src={ampleforth} className="ampleforth" />
        </a>
      </div>
    </section>
  </div>
)

export default LangingPage

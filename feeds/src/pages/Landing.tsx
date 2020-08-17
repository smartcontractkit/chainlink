import React from 'react'
import { useLocation } from 'react-router-dom'
import { Listing } from 'components/listing'
import { Header } from 'components/header'
import { NodesLogos } from 'components/nodesLogos'
import { SponsorsLogos } from 'components/sponsorsLogos'

const Page = () => {
  return (
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
            shared global resource which is sponsored by a growing list of top
            DeFi Dapps.
          </p>
          <p>
            Please feel free to look into the details of each Decentralized
            Oracle Network listed below. You can easily use these oracle
            networks to quickly and securely launch, add more capabilities to
            and/or just greatly improve the security of your smart contracts.
          </p>
        </section>
        <section>
          <Listing enableDetails={useDetails()} />
        </section>
      </div>

      <SponsorsLogos />
      <NodesLogos />
    </div>
  )
}

function useDetails(): boolean {
  const query = new URLSearchParams(useLocation().search)
  return query.get('details') === 'true'
}

export default Page

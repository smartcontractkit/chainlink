import React from 'react'
import { Row, Col } from 'antd'

const grid = { xs: 24, sm: 12, md: 6 }

const Footer = () => (
  <footer className="footer">
    <div className="footer__container">
      <Row gutter={36}>
        <Col {...grid}>
          <div className="footer__title">We can be mailed at:</div>
          <address>
            Strathvale House, 90 North Church Street, George Town, KY1-1102,
            Grand Cayman, Cayman Islands
          </address>
          <br />
          <p>
            Chainlink<sup>®</sup>
            <br />
            <br />© 2019 SmartContract Chainlink Ltd SEZC
            <br />
            <br />
            <a href="/terms"> Terms of Use</a>
            <br />
            <a href="/privacy-policy"> Privacy Policy</a>
          </p>
        </Col>
        <Col {...grid}>
          <div className="footer__title">Contact Us 24/7/365:</div>
          <ul className="contact-list">
            <li>
              24/7 Technical Support
              <br />
              <a href="mailto:support@chain.link"> support@chain.link</a>
            </li>
            <li>
              Custom Chainlinks
              <br />
              <a href="mailto:custom@chain.link"> custom@chain.link</a>
            </li>
            <li>
              Press Inquiries
              <br />
              <a href="mailto:press@chain.link"> press@chain.link</a>
            </li>
            <li>
              <a
                href="https://careers.chain.link"
                target="_blank"
                rel="noopener noreferrer"
              >
                Careers
              </a>
            </li>
          </ul>
        </Col>
        <Col {...grid}>
          <div className="footer__title">Learn More:</div>
          <ul className="social-list">
            <li>
              <a
                href="https://medium.com/chainlink"
                rel="noopener noreferrer"
                target="_blank"
              >
                <span className="icon-monogram"></span>Medium
              </a>
            </li>
            <li>
              <a
                href="https://www.reddit.com/r/Chainlink/"
                rel="noopener noreferrer"
                target="_blank"
              >
                <span className="icon-reddit"></span>Reddit
              </a>
            </li>
            <li>
              <a
                href="https://t.me/chainlinkofficial"
                rel="noopener noreferrer"
                target="_blank"
              >
                <span className="icon-telegram"></span>Telegram
              </a>
            </li>
            <li>
              <a
                href="https://twitter.com/chainlink"
                rel="noopener noreferrer"
                target="_blank"
              >
                <span className="icon-twitter"></span>Twitter
              </a>
            </li>
            <li>
              <a
                href="https://www.youtube.com/chainlinkofficial"
                rel="noopener noreferrer"
                target="_blank"
              >
                <span className="icon-you-tube"></span>YouTube
              </a>
            </li>
            <li>
              <div className="footer__title">Presentations:</div>
            </li>
            <li>
              <a
                href="https://chain.link/presentations/devcon5.pdf"
                rel="noopener noreferrer"
                target="_blank"
              >
                <span className="icon-new-window"></span>Devcon5
              </a>
            </li>
            <li>
              <a
                href="https://chain.link/presentations/english.pdf"
                rel="noopener noreferrer"
                target="_blank"
              >
                <span className="icon-new-window"></span>English
              </a>
            </li>
            <li>
              <a
                href="https://chain.link/presentations/chinese.pdf"
                rel="noopener noreferrer"
                target="_blank"
              >
                <span className="icon-new-window"></span>Chinese
              </a>
            </li>
            <li>
              <a
                href="https://chain.link/presentations/korean.pdf"
                rel="noopener noreferrer"
                target="_blank"
              >
                <span className="icon-new-window"></span>Korean
              </a>
            </li>
            <li>
              <a
                href="https://chain.link/presentations/japanese.pdf"
                rel="noopener noreferrer"
                target="_blank"
              >
                <span className="icon-new-window"></span>Japanese
              </a>
            </li>
          </ul>
        </Col>
        <Col {...grid}>
          <div className="footer__title">Build With Us:</div>
          <ul className="social-list">
            <li>
              <a
                href="https://github.com/smartcontractkit/chainlink"
                rel="noopener noreferrer"
                target="_blank"
              >
                <span className="icon-git"></span>GitHub
              </a>
            </li>
            <li>
              <a
                href="https://gitter.im/smartcontractkit-chainlink/Lobby"
                rel="noopener noreferrer"
                target="_blank"
              >
                <span className="icon-gitter"></span>Gitter
              </a>
            </li>
            <li>
              <a
                href="https://discord.gg/aSK4zew"
                rel="noopener noreferrer"
                target="_blank"
              >
                <span className="icon-discord"></span>Discord
              </a>
            </li>
            <li>
              <a
                href="https://eth-usd-aggregator.chain.link/"
                rel="noopener noreferrer"
                target="_blank"
              >
                <span className="icon-new-window"></span>Chainlink Reference
                Data
              </a>
            </li>
          </ul>
        </Col>
      </Row>
    </div>
  </footer>
)

export default Footer

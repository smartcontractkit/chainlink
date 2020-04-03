import React from 'react'

const YEAR = new Date().getFullYear()

const Footer = () => (
  <footer className="footer">
    <div className="footer__container">
      <span>Chainlink &copy; {YEAR}</span>
      <span>
        <a
          target="_blank"
          rel="noopener noreferrer"
          href="https://chain.link/terms/"
        >
          Terms of Use
        </a>
      </span>
      <span>
        <a
          target="_blank"
          rel="noopener noreferrer"
          href="https://chain.link/privacy-policy/"
        >
          Privacy Policy
        </a>
      </span>
    </div>
  </footer>
)

export default Footer

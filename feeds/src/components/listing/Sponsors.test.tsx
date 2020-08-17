import React from 'react'
import '@testing-library/jest-dom/extend-expect'
import { render } from '@testing-library/react'
import Sponsors from './Sponsors'

describe('components/listing/Sponsors', () => {
  it('renders thumbnail logo', () => {
    const { container } = render(<Sponsors sponsors={['Synthetix']} />)
    const logo = container.querySelector('img[alt="Synthetix"]')
    expect(logo).toBeTruthy()
  })

  it('renders max 5 thumbnails', () => {
    const sponsors = [
      'Synthetix',
      'Loopring',
      'OpenLaw',
      '1inch',
      'ParaSwap',
      'MCDEX',
      'Futureswap',
    ]
    const { container } = render(<Sponsors sponsors={sponsors} />)
    const logos = container.querySelectorAll('img')
    expect(logos.length).toEqual(5)
    expect(container).toHaveTextContent('Sponsored by (+2)')
  })

  it('does not render thumbnails when there are no sponsors', () => {
    const sponsors: string[] | undefined = []
    const { container } = render(<Sponsors sponsors={sponsors} />)
    const logos = container.querySelectorAll('img')
    expect(logos.length).toEqual(0)
  })
})

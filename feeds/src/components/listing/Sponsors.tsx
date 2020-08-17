import React, { useState, useEffect } from 'react'
import { sponsorList, SponsorListItem } from '../../assets/sponsors'
import { Popover } from 'antd'

interface LogoProps {
  sponsor: SponsorListItem
}

const Logo: React.FC<LogoProps> = ({ sponsor }) => (
  <div className="listing-grid__item--sponsors__logo">
    <img alt={sponsor.name} title={sponsor.name} src={sponsor.imageTn} />
  </div>
)

interface SponsoredProps {
  sponsors: string[] | undefined
}

const PopoverList: React.FC<SponsoredProps> = ({ sponsors }) => (
  <>
    {sponsors?.map((sponsor, i) => (
      <div className="listing-grid__item--sponsors-popover" key={i}>
        {sponsor}
      </div>
    ))}
  </>
)

const Sponsors: React.FC<SponsoredProps> = ({ sponsors = [] }) => {
  const [thumbnails, setThumbnails] = useState<SponsorListItem[]>([])
  const [sponsorsRemainder, setSponsorsRemainder] = useState(0)

  useEffect(() => {
    const sponsorsTn = sponsorList.filter(s => {
      return sponsors.includes(s.name)
    })
    setSponsorsRemainder(sponsors.length - 5)
    setThumbnails(sponsorsTn.slice(0, 5))
  }, [setThumbnails, sponsors])

  if (!sponsors || !sponsors.length) {
    return null
  }

  return (
    <Popover content={<PopoverList sponsors={sponsors} />} title="Sponsored by">
      <div>
        <div className="listing-grid__item--sponsors-title">
          Sponsored by {sponsorsRemainder > 0 && `(+${sponsorsRemainder})`}
        </div>

        <div className="listing-grid__item--sponsors">
          {thumbnails?.map((sponsor, i: number) => (
            <Logo key={i} sponsor={sponsor} />
          ))}
        </div>
      </div>
    </Popover>
  )
}

export default Sponsors

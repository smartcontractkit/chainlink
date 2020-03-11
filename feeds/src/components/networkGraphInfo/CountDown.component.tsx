import React, { useEffect, useState } from 'react'
import moment from 'moment'

let timer: number

const toHHMMSS = (secs: number) => {
  const hours = Math.floor(secs / 3600)
  const minutes = Math.floor(secs / 60) % 60
  const seconds = secs % 60

  return [hours, minutes, seconds]
    .map(v => (v < 10 ? '0' + v : v))
    .filter((v, i) => v !== '00' || i > 0)
    .join(':')
}

export interface Props {
  requestTime: number
  counter: number
}

const CountDown: React.FC<Props> = ({ requestTime, counter }) => {
  const [next, setNext] = useState<string>()

  useEffect(() => {
    if (!requestTime) {
      return setNext('00:00')
    }

    const finish = moment
      .unix(requestTime)
      .add(counter, 'seconds')
      .unix()

    clearInterval(timer)

    timer = window.setInterval(() => {
      const now = moment(new Date()).unix()
      const distance = finish - now

      if (distance <= 0) {
        setNext('00:00')
        return clearInterval(timer)
      }

      setNext(toHHMMSS(distance))
    }, 1000)

    return () => {
      clearInterval(timer)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [requestTime])

  return <span className="countdown">{next || '...'}</span>
}

export default CountDown

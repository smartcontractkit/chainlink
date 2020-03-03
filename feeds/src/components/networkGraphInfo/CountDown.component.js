import React, { useEffect, useState } from 'react'
import moment from 'moment'

let timer

const toHHMMSS = secs => {
  const secNum = parseInt(secs, 10)
  const hours = Math.floor(secNum / 3600)
  const minutes = Math.floor(secNum / 60) % 60
  const seconds = secNum % 60

  return [hours, minutes, seconds]
    .map(v => (v < 10 ? '0' + v : v))
    .filter((v, i) => v !== '00' || i > 0)
    .join(':')
}

function CountDown({ requestTime, counter }) {
  const [next, setNext] = useState()

  useEffect(() => {
    if (!requestTime) {
      return setNext('00:00')
    }

    const finish = moment
      .unix(requestTime)
      .add(counter, 'seconds')
      .unix()

    clearInterval(timer)

    timer = setInterval(() => {
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

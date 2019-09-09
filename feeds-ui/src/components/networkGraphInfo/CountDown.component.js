import React, { useEffect, useState } from 'react'
import moment from 'moment'

let timer

var toHHMMSS = secs => {
  var sec_num = parseInt(secs, 10)
  var hours = Math.floor(sec_num / 3600)
  var minutes = Math.floor(sec_num / 60) % 60
  var seconds = sec_num % 60

  return [hours, minutes, seconds]
    .map(v => (v < 10 ? '0' + v : v))
    .filter((v, i) => v !== '00' || i > 0)
    .join(':')
}

function CountDown({ requestTime }) {
  const [next, setNext] = useState()

  useEffect(() => {
    if (!requestTime) {
      return
    }

    const finish = moment
      .unix(requestTime)
      .add(5, 'm')
      .unix()

    clearInterval(timer)

    timer = setInterval(() => {
      let now = moment(new Date()).unix()
      const distance = finish - now

      if (distance <= 0) {
        setNext('00:00')
        return clearInterval(timer)
      }

      setNext(toHHMMSS(distance))
    }, 1000)
  }, [requestTime])

  return <span className="countdown">{next || '...'}</span>
}

export default CountDown

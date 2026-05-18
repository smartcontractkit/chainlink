export const validCrons = [
  '* * * * *', // every minute
  '*/2 * * * *', // every even minute
  '0 * * * *', // every hour
  '0 0 * * *', // every day at midnight
  '0 12 * * *', // every day at noon
  '0 12 * * 0-4', // week days at noon
  '0 0 1 * *', // every month on the first at midnight
  '0 0 1 7 *', // first of July at midnight
  '0 0 * * 1', // every monday at midnight
  '*/5 * * * *', // every 5 min
  '0 0 * * 2-4', // wed - friday at midnight
  '0 * 31 * 0', // 31st day of the month, mondays, at midnight
  '59 23 29 2 1', // last minute of tuesday leap days
  '0 12 1,3,5,7,11,13,17,19,23,27,29,31 * *', // prime days at noon
  '*/20 3,7,20 10-20 */2 5-6', // every 20 min b/t hours 3:4, 7:8, and 20:21 on the 10-20th days, even months, weekends
  '0 0 29 2 *', // every leap day
  '0 0 */2 2 *', // every even day in february
]

export const invalidCrons = [
  '60 * * * *', // invalid minute
  '0 24 * * *', // invalid hour
  '0 * 32 * *', // invalid day
  '0 * 0 * *', // invalid day
  '0 * * 13 *', // invalid month
  '0 * * 0 *', // invalid month
  '* * 30 2 *', // invalid day/month
  '* * 31 2,4,6,9,11 *', // invalid day/month
  '* * 31 */9 *', // invalid day/month
  '* * 20-31 2 *', // invalid day/month
  '* * 28,29,30 2 *', // invalid day/month
  '0 * * * 7', // invalid day of week
  '0 12-24 * * 7', // invalid hour range
  '0 * * * 5-10', // invalid day of week range
  '0 * * * 1-1', // invalid range
  '0 * * * 2-1', // invalid range
  '0 0,3,5,30 * * *', // invalid hour list
  '*/100 * * * *', // invalid interval
  '*/0 * * * *', // invalid interval
  '0****', // no spaces
  '0 * * **', // too few spaces
  '0 * * *  *', // too many spaces
  ' 0 * * * *', // leading whitespace
  '0 * * * * ', // trailing whitespace
  '0 * * * ', // field missing
  '0 1, * * *', // invalid list
  '0 * 1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27 * *', // list too big
]

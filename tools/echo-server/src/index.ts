import express from 'express'
import bodyParser from 'body-parser'
import chalk from 'chalk'

let echoes = 0
const app = express()
app.use(bodyParser.json())

app.get('/count', function (_, res) {
  res.json(echoes)
})

app.all('*', function (req, res) {
  echoes += 1
  const { headers, body } = req

  console.log({ headers, body })
  res.json({ headers, body })
})

const port = process.env.PORT ? parseInt(process.env.PORT, 10) : 6688
app.listen(port, function () {
  console.log(chalk.green(`echo_server listening on port ${port}`))
})

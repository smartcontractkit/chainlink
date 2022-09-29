import express from 'express'
import bodyParser from 'body-parser'

const app = express()

app.use(bodyParser.json())

app.use((req, _, next) => {
  console.log(`${req.method} request made to ${req.originalUrl}`)
  console.log('Request body...')
  console.log(JSON.stringify(req.body, null, 2))
  next()
})

let result = 100

app.post('/', (req, res) => {
  const jobRunID = req.body.id
  const responseBody = {
    jobRunID,
    data: {
      result,
    },
    error: null,
  }
  res.json(responseBody)
  console.log('Response body...')
  console.log(JSON.stringify(responseBody, null, 2))
})

app.patch('/result', (req, res) => {
  result = req.body.result
  res.status(200)
  res.send()
})

const port = process.env.EXTERNAL_ADAPTER_PORT || 6644
app.listen(port, () => console.log(`Listening on port ${port}!`))

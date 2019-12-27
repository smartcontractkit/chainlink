import express from 'express'
import bodyParser from 'body-parser'

const app = express()
app.use(bodyParser.json())

app.use((req, _, next) => {
  console.log(`Requested ${req.method} to ${req.path}`)
  console.log('Request Data:', req.body)
  next()
})

let result = 100.0

app.post('/', (req, res) => {
  const jobRunID = req.body.id
  res.status(200).json({
    jobRunID,
    data: {
      result,
    },
    error: null,
  })
})

app.patch('/result', (req, res) => {
  result = req.body.result
  res.status(200)
  res.send()
})

const port = process.env.EXTERNAL_ADAPTER_PORT || 6644
app.listen(port, () => console.log(`Listening on port ${port}!`))

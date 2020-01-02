import express from 'express'
import bodyParser from 'body-parser'
import morganBody from 'morgan-body'

const app = express()

app.use(bodyParser.json())
morganBody(app)

let result = 0

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

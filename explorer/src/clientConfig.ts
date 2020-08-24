import fs from 'fs'
import path from 'path'

const workingDir = process.cwd()
const clientDir = 'client/build/config.json'
const fileToUpdate = path.resolve(workingDir, clientDir)
const GA_ID = 'GA_ID'

export function updateClientGaId(id: string) {
  fs.readFile(fileToUpdate, 'utf8', (err, file) => {
    if (err) {
      return console.log(err)
    }

    const config = JSON.parse(file)
    config[GA_ID] = id

    fs.writeFile(fileToUpdate, JSON.stringify(config), 'utf8', err => {
      if (err) return console.log(err)
    })
  })
}

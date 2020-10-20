const execa = require('execa')

module.exports = {
  run,
  sleep,
  getConfig,
}

async function run(cmd, env) {
  try {
    if (cmd instanceof Array) {
      return (await execa(cmd[0], cmd.slice(1), { all: true, env: env })).all
    } else {
      return (await execa.command(cmd, { all: true, env: env })).all
    }
  } catch (err) {
    console.log('ERROR WITH RUN:', err)
  }
}

async function sleep(ms) {
  return new Promise((resolve) => {
    setTimeout(resolve, ms)
  })
}

function getConfig() {
  return JSON.parse(fs.readFileSync('../../config/config.json').toString())
}

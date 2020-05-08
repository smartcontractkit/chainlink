const ethwallet = require('ethereumjs-wallet')
const fs = require('fs')
const path = require('path')

const password = 'password'
let personas = []
let accounts = {}
try { personas = require(path.join(accountDirPath(), 'personas.json')) } catch (err) {}
try { accounts = require(path.join(accountDirPath(), 'accounts.json')) } catch (err) {}

module.exports = {
    accounts,
    password,
    addPersona,
    ensureAccounts,
}

function accountDirPath() {
    return path.join(__dirname, '..', 'accounts')
}

function keyfilePath(persona) {
    return path.join(accountDirPath(), 'keys', persona)
}

function addPersona(persona) {
    personas.push(persona)
    fs.writeFileSync(path.join(accountDirPath(), 'personas.json'), JSON.stringify(personas, null, 4))
    ensureAccounts([ persona ])
}

function rmPersona(persona) {
    try { fs.unlinkSync( keyfilePath(persona) ) } catch (err) {}

    personas = personas.filter(x => x !== persona)
    fs.writeFileSync(path.join(accountDirPath(), 'personas.json'), JSON.stringify(personas, null, 4))

    delete accounts[persona]
    fs.writeFileSync(path.join(accountDirPath(), 'accounts.json'), JSON.stringify(accounts, null, 4))
}

function ensureAccounts(personas) {
    // Generate missing accounts
    for (let persona of personas) {
        let wallet
        if (accounts[persona] === undefined) {
            console.log(persona, 'unknown, generating account...')
            wallet = ethwallet.generate()
            // keyfileName = wallet.getV3Filename(new Date().getTime())
        } else {
            wallet = ethwallet.fromV3(JSON.stringify(accounts[persona].keyfile), password)
            console.log(persona, `wallet valid (address: ${wallet.getAddressString()})`)
            // keyfileName = accounts[persona].keyfileName
        }
        accounts[persona] = {
            address: wallet.getAddressString(),
            privkey: wallet.getPrivateKeyString(),
            keyfile: JSON.parse(wallet.toV3String(password)),
            // keyfileName: persona,
        }
    }
    fs.writeFileSync(path.join(accountDirPath(), 'accounts.json'), JSON.stringify(accounts, null, 4))

    // Write keyfiles
    for (let persona of Object.keys(accounts)) {
        const account = accounts[persona]
        if (!account) {
            continue
        }
        const keyfileJSON = JSON.stringify(account.keyfile)
        const wallet = ethwallet.fromV3(keyfileJSON, password)
        const filename = account.keyfileName
        fs.writeFileSync(keyfilePath(persona), keyfileJSON)
    }
}

function clean() {
    // Remove from accounts.json
    for (let persona of Object.keys(accounts)) {
        if (!personas.includes(persona)) {
            delete accounts[persona]
        }
    }
    fs.writeFileSync(path.join(accountDirPath(), 'accounts.json'), JSON.stringify(accounts, null, 4))

    const knownFiles = Object.keys(accounts).map(x => accounts[x].keyfileName)
    const dirFiles = fs.readdirSync(__dirname)
                       .filter(x => x.startsWith('UTC--'))
                       .filter(x => !knownFiles.includes(x))
                       .map(x => path.join(accountDirPath(), x))
    for (let file of dirFiles) {
        fs.unlinkSync(file)
        console.log('deleting', file)
    }
}

const USAGE = `
    Usage:
        config <command>

    Commands:
        add-persona     - add a persona to personas.json and generate its account and keyfile
        rm-persona      - remove a persona from personas.json and delete its account data and keyfile
        ensure-accounts - ensure all personas (in personas.json) have accounts and keyfiles
        clean           - remove and account data and keyfiles that don't have a corresponding persona
`

function main() {
    if (__filename !== process.argv[1]) {
        // This is being imported by another script
        return
    } else if (process.argv.length < 3) {
        console.log(USAGE)
        process.exit(1)
        return
    }

    switch (process.argv[2]) {
    case 'ensure-accounts':
        ensureAccounts(personas)
        return

    case 'add-persona':
        if (process.argv.length < 4) {
            console.log(USAGE)
            process.exit(1)
            return
        }
        addPersona(process.argv[3])
        return

    case 'rm-persona':
        if (process.argv.length < 4) {
            console.log(USAGE)
            process.exit(1)
            return
        }
        rmPersona(process.argv[3])
        return

    case 'clean':
        clean()
        return

    default:
        console.log('Unknown command.')
        console.log(USAGE)
        process.exit(1)
        return
    }
}

main()

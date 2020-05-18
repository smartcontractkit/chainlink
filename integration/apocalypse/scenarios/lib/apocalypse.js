const ethers = require('ethers')
const { getConfig } = require('./utils')

module.exports = {
    mimicRegularTraffic,
}

async function mimicRegularTraffic({ funderPrivkey, numAccounts, ethereumRPCProviders }) {
    let accounts = await makeRandomAccounts(funderPrivkey, numAccounts, ethereumRPCProviders)
    for (let account of accounts) {
        sendRandomTransactions(account, accounts)
    }
}

async function makeRandomAccounts(funderPrivkey, num, ethereumRPCProviders) {
    let senders = []
    for (let providerURL of ethereumRPCProviders) {
        let wallet = new ethers.Wallet(funderPrivkey, new ethers.providers.JsonRpcProvider(providerURL))
        senders.push({
            providerURL: providerURL,
            nonce: await wallet.provider.getTransactionCount(wallet.address, 'pending'),
            wallet: wallet,
        })
    }
    let jobs = Array(num).fill(null).map((_, i) => {
        let sender = senders[i % senders.length]
        return {
            providerURL: sender.providerURL,
            wallet: ethers.Wallet.createRandom().connect(new ethers.providers.JsonRpcProvider(sender.providerURL)),
            sender: sender,
        }
    })
    // Fund the accounts
    await Promise.all(
        jobs.map(job => {
            let nonce = job.sender.nonce
            job.sender.nonce++
            return job.sender.wallet.sendTransaction({
                to: job.wallet.address,
                value: ethers.utils.parseUnits('5', 'ether'),
                gasPrice: ethers.utils.parseUnits('20', 'gwei'),
                nonce: nonce,
            }).catch(err => {
                console.log(err, 'nonce =', nonce, job.sender.wallet.address, job.sender.nonce, job.providerURL)
            })
        })
    )
    return jobs.map(job => job.wallet)
}

function sendRandomTransactions(fromAccount, toAccounts) {
    function randomAccount() {
        let i = Math.floor(Math.random() * Math.floor(toAccounts.length - 1))
        return toAccounts[i]
    }

    async function send() {
        let msBetweenTxs = 500
        try {
            // Re-read the config each time so that we can control the congestion dynamically
            msBetweenTxs = getConfig().randomTraffic.msBetweenTxs

            await fromAccount.sendTransaction({
                to: randomAccount().address,
                value: ethers.utils.parseUnits('1', 'wei'),
                gasPrice: ethers.utils.parseUnits('20', 'gwei'),
            })
        } catch (err) {}

        setTimeout(send, msBetweenTxs)
    }
    send()
}

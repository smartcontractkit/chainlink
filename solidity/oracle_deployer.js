const solc = require("solc");
const fs = require('fs');
const path = require('path');
const url = require('url');
const Web3 = require('web3');
const Eth = require('ethjs');
const Tx = require('ethereumjs-tx');
const Wallet = require('ethereumjs-wallet');

function toWei(eth) {
  return (parseInt(eth) * 10**18).toString();
}

function compile(filename) {
  const lookupPaths = ["./", "./node_modules/"];
  function lookupIncludeFile(filename) {
    for(let path of lookupPaths) {
      let fullPath = path + filename;
      if(fs.existsSync(fullPath)) {
        return {contents: fs.readFileSync(fullPath).toString()};
      }
    }
    console.log("Unable to load", filename);
    return null;
  }

  let inputBasename = path.basename(filename).toString();
  let input = {[inputBasename]: {"urls": [filename]}};

  let solInput = {
    language: "Solidity",
    sources: input,
    settings: {
      outputSelection: {
          [inputBasename]: {
            "*": [ "evm.bytecode" ]
          },
      },
    },
  };
  let output = solc.compileStandardWrapper(JSON.stringify(solInput), lookupIncludeFile);
  return JSON.parse(output);
};

async function main(filename) {
  let inputBasename = path.basename(filename).toString();
  let output = compile(filename);

  // Print any errors
  let failure = false;
  for (let error of output.errors) {
    if(error.type !== "Warning") {
      console.log(error.sourceLocation.file + ':' + error.sourceLocation.start + " - " + error.message)
      failure = true;
    }
  }
  if(failure) {
    console.log("Aborting because of errors.");
    process.exit(1);
  }

  let bytecode = null;
  for (const [key, module] of Object.entries(output.contracts)) {
    if(key === inputBasename) {
      for (const [contractName, contract] of Object.entries(module)) {
        bytecode = contract.evm.bytecode.object.toString();
      }
    }
  }

  let eth = new Eth(new Eth.HttpProvider('http://localhost:18545'));
  let from = '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f';
  let privateKey = new Buffer('4d6cf3ce1ac71e79aa33cf481dedf2e73acb548b1294a70447c960784302d2fb', 'hex');
  let wallet = Wallet.fromPrivateKey(privateKey);
  let deployer = wallet.getAddress().toString('hex');

  await eth.sendTransaction({
    to: deployer,
    from: from,
    value: toWei("1"),
    gas: 25000,
    data: "0x",
  });

  let tx = new Tx({
    gas: 4700000,
    data: "0x" + bytecode + "0000000000000000000000004b274dfcd56656742A55ad54549b3770c392aA87",
    nonce: await eth.getTransactionCount(deployer),
    chainId: 17,
  });
  tx.sign(privateKey);
  let txHash = await eth.sendRawTransaction(tx.serialize().toString("hex"));
  await setTimeout(async () => {
    let receipt = await eth.getTransactionReceipt(txHash);
    console.log("receipt:", receipt);
  }, 1000);
}

if(process.argv.length != 3) {
  console.error("Usage: node oracle_deployer.js <solidity contract>");
  process.exit(1);
}

let input = process.argv[2];
main(input);

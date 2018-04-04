#!/usr/bin/env node

// make sourcemaps work!
require('source-map-support/register')

// `yargs/yargs` required to work with webpack, see here.
// https://github.com/yargs/yargs/issues/781
var yargs = require('yargs/yargs');
var Ganache = require("ganache-core");
var pkg = require("./package.json");
var corepkg = require("./node_modules/ganache-core/package.json");
var URL = require("url");
var fs = require("fs");
var to = require("ganache-core/lib/utils/to");

var parser = yargs()
.option("unlock", {
  type: "string",
  alias: "u"
});

var argv = parser.parse(process.argv);

if (argv.help || argv['?']) {
  console.log("");
  console.log("ganache-cli: Fast Ethereum RPC client for testing and development");
  console.log("  Full docs: https://github.com/trufflesuite/ganache-cli");
  console.log("");
  console.log("Usage: ganache-cli [options]");
  console.log("  options:");
  console.log("  --port/-p <port to bind to, default 8545>");
  console.log("  --host/-h <host to bind to, default 127.0.0.1>");
  console.log("  --fork/-f <url>   (Fork from another currently running Ethereum client at a given block)");
  console.log("");
  console.log("  --db <db path>   (directory to save chain db)");
  console.log("  --seed/-s <seed value for PRNG, default random>");
  console.log("  --deterministic/-d     (uses fixed seed)");
  console.log("  --mnemonic/-m <mnemonic>");
  console.log("");
  console.log("  --accounts/-a <number of accounts to generate at startup>   (ignored when using --account flag)");
  console.log("  --defaultBalanceEther/-e <Amount of ether to assign each test account, default 100.0>   (ignored when using --account flag)");
  console.log("  --account <privatekey>,<balance>   (Can be specified multiple times. Note that private keys are 64 characters long,");
  console.log("                                       and must be input as a 0x-prefixed hex string. Balance can either be input as an");
  console.log("                                       integer or 0x-prefixed hex value specifying the amount of wei in that account.)");
  console.log("");
  console.log("  --acctKeys <path to file>   (saves generated accounts and private keys as JSON object in specified file)");
  console.log("  --secure/-n   (Lock accounts by default)");
  console.log("  --unlock <accounts>   (Comma-separated list of accounts or indices to unlock)");
  console.log("");
  console.log("  --noVMErrorsOnRPCResponse   (Do not transmit transaction failures as RPC errors. Enable this flag for error reporting behaviour which is compatible with other clients such as geth and Parity.)");
  console.log("");
  console.log("  --blockTime/-b <block time in seconds> (Will instamine if option omitted. Avoid using unless your test cases require interval mining.)");
  console.log("  --networkId/-i <network id> (default current time)");
  console.log("  --gasPrice/-g <gas price>   (default 20000000000)");
  console.log("  --gasLimit/-l <gas limit>   (default 90000)");
  console.log("");
  console.log("  --debug       (Output VM opcodes for debugging)");
  console.log("  --verbose/-v");
  console.log("  --mem         (Only show memory output, not tx history)");
  console.log("");
  console.log("");
  console.log("  --help / -?    (this output)");
  console.log("");
  process.exit(0);
}

function parseAccounts(accounts) {
  function splitAccount(account) {
    account = account.split(',')
    return {
      secretKey: account[0],
      balance: account[1]
    };
  }

  if (typeof accounts === 'string')
    return [ splitAccount(accounts) ];
  else if (!Array.isArray(accounts))
    return;

  var ret = []
  for (var i = 0; i < accounts.length; i++) {
    ret.push(splitAccount(accounts[i]));
  }
  return ret;
}

if (argv.d || argv.deterministic) {
  argv.s = "TestRPC is awesome!"; // Seed phrase; don't change to Ganache, maintain original determinism
}

if (typeof argv.unlock == "string") {
  argv.unlock = [argv.unlock];
}

var logger = console;

// If the mem argument is passed, only show memory output,
// not transaction history.
if (argv.mem === true) {
  logger = {
    log: function() {}
  };

  setInterval(function() {
    console.log(process.memoryUsage());
  }, 1000);
}

var options = {
  port: argv.p || argv.port || "8545",
  hostname: argv.h || argv.hostname,
  debug: argv.debug,
  seed: argv.s || argv.seed,
  mnemonic: argv.m || argv.mnemonic,
  total_accounts: argv.a || argv.accounts,
  default_balance_ether: argv.e || argv.defaultBalanceEther,
  blockTime: argv.b || argv.blockTime,
  gasPrice: argv.g || argv.gasPrice,
  gasLimit: argv.l || argv.gasLimit,
  accounts: parseAccounts(argv.account),
  unlocked_accounts: argv.unlock,
  fork: argv.f || argv.fork || false,
  network_id: argv.i || argv.networkId,
  verbose: argv.v || argv.verbose,
  secure: argv.n || argv.secure || false,
  db_path: argv.db || null,
  account_keys_path: argv.acctKeys || null,
  vmErrorsOnRPCResponse: !argv.noVMErrorsOnRPCResponse,
  logger: logger
}

var fork_address;

// If we're forking from another client, don't try to use the same port.
if (options.fork) {
  var split = options.fork.split("@");
  fork_address = split[0];
  var block;
  if (split.length > 1) {
    block = split[1];
  }

  if (URL.parse(fork_address).port == options.port) {
    options.port = (parseInt(options.port) + 1);
  }

  options.fork = fork_address + (block != null ? "@" + block : "");
}

var server = Ganache.server(options);

console.log("Ganache CLI v" + pkg.version + " (ganache-core: " + corepkg.version + ")");

server.listen(options.port, options.hostname, function(err, result) {
  if (err) {
    console.log(err);
    return;
  }

  var state = result ? result : server.provider.manager.state;

  console.log("");
  console.log("Available Accounts");
  console.log("==================");

  var accounts = state.accounts;
  var addresses = Object.keys(accounts);

  addresses.forEach(function(address, index) {
    var line = "(" + index + ") " + address;

    if (state.isUnlocked(address) == false) {
      line += " 🔒";
    }

    console.log(line);
  });

  console.log("");
  console.log("Private Keys");
  console.log("==================");

  addresses.forEach(function(address, index) {
    console.log("(" + index + ") " + accounts[address].secretKey.toString("hex"));
  });


  if (options.account_keys_path != null) {
    console.log("");
    console.log("Saving accounts and keys to " + options.account_keys_path);
    var obj = {}
    obj.addresses = accounts;
    obj.private_keys = {};
    addresses.forEach(function(address, index) {
       obj.private_keys[address] = accounts[address].secretKey.toString("hex");
    });
    var json = JSON.stringify(obj);
    fs.writeFile(options.account_keys_path, json, 'utf8',function(err){
      if(err) throw err;
    })
  }

  if (options.accounts == null) {
    console.log("");
    console.log("HD Wallet");
    console.log("==================");
    console.log("Mnemonic:      " + state.mnemonic);
    console.log("Base HD Path:  " + state.wallet_hdpath + "{account_index}")
  }

  if (options.gasPrice) {
    console.log("");
    console.log("Gas Price");
    console.log("==================");
    console.log(options.gasPrice);
  }

  if (options.gasLimit) {
    console.log("");
    console.log("Gas Limit");
    console.log("==================");
    console.log(options.gasLimit);
  }

  if (options.fork) {
    console.log("");
    console.log("Forked Chain");
    console.log("==================");
    console.log("Location:    " + fork_address);
    console.log("Block:       " + to.number(state.blockchain.fork_block_number));
    console.log("Network ID:  " + state.net_version);
    console.log("Time:        " + (state.blockchain.startTime || new Date()).toString());
  }

  console.log("");
  console.log("Listening on " + (options.hostname || "localhost") + ":" + options.port);
});

process.on('uncaughtException', function(e) {
  console.log(e.stack);
  process.exit(1);
})

// See http://stackoverflow.com/questions/10021373/what-is-the-windows-equivalent-of-process-onsigint-in-node-js
if (process.platform === "win32") {
  require("readline").createInterface({
    input: process.stdin,
    output: process.stdout
  })
  .on("SIGINT", function () {
    process.emit("SIGINT");
  });
}

process.on("SIGINT", function () {
  // graceful shutdown
  server.close(function(err) {
    if (err) {
      console.log(err.stack || err);
    }
    process.exit();
  });
});

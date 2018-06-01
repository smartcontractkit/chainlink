function checkAllBalances() {
  var totalBal = 0;
  for (var acctNum in eth.accounts) {
    var acct = eth.accounts[acctNum];
    var acctBal = web3.fromWei(eth.getBalance(acct), "ether");
    totalBal += parseFloat(acctBal);
    console.log("  eth.accounts[" + acctNum + "]: \t" + acct + " \tbalance: " + acctBal + " ether");
  }
  console.log("  Total balance: " + totalBal + " ether");
};

function fundAccount(amount) {
  // XXX: Strip off anything after the 18th LSD to prevent web3.toWei from
  // returning a value with a decimal place which sendTransaction will barf on
  amount = (amount || 1000).toFixed(18).replace(/0+$/, "");
  return eth.sendTransaction({
    from:eth.accounts[0],
    to:eth.accounts[1],
    value: web3.toWei(amount, "ether")
  });
};

function topOffAccount() {
  var acct = "0x9ca9d2d5e04012c9ed24c0e513c9bfaa4a2dd77f";
  var acctBal = web3.fromWei(eth.getBalance(acct), "ether");
  var fundAmount = 10000;

  var diff = 10000 - acctBal;
  if (diff > 0) {
    return fundAccount(diff);
  } else {
    console.log("Not performing top off, account has sufficient funds:", acctBal);
  }
};

function confirm(txHash) {
  if(!txHash) {
    return null;
  }
  var count = 0;
  do {
    receipt = eth.getTransactionReceipt(txHash);
    count += 1;
  } while (receipt == null && count < 50);
  return receipt;
}

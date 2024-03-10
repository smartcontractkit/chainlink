import { ethers } from 'hardhat'
// Suppress "Duplicate definition" error logs
ethers.utils.Logger.setLogLevel(ethers.utils.Logger.levels.ERROR)

import { Signer } from 'ethers'

export interface Contracts {
  contract1: Signer
  contract2: Signer
  contract3: Signer
  contract4: Signer
  contract5: Signer
  contract6: Signer
  contract7: Signer
  contract8: Signer
}

export interface Roles {
  defaultAccount: Signer
  oracleNode: Signer
  oracleNode1: Signer
  oracleNode2: Signer
  oracleNode3: Signer
  oracleNode4: Signer
  stranger: Signer
  consumer: Signer
  consumer2: Signer
}

export interface Personas {
  Default: Signer
  Carol: Signer
  Eddy: Signer
  Nancy: Signer
  Ned: Signer
  Neil: Signer
  Nelly: Signer
  Norbert: Signer
  Nick: Signer
}

export interface Users {
  contracts: Contracts
  roles: Roles
  personas: Personas
}

export async function getUsers() {
  const accounts = await ethers.getSigners()

  const personas: Personas = {
    Default: accounts[0],
    Neil: accounts[1],
    Ned: accounts[2],
    Nelly: accounts[3],
    Nancy: accounts[4],
    Norbert: accounts[5],
    Carol: accounts[6],
    Eddy: accounts[7],
    Nick: accounts[8],
  }

  const contracts: Contracts = {
    contract1: accounts[0],
    contract2: accounts[1],
    contract3: accounts[2],
    contract4: accounts[3],
    contract5: accounts[4],
    contract6: accounts[5],
    contract7: accounts[6],
    contract8: accounts[7],
  }

  const roles: Roles = {
    defaultAccount: accounts[0],
    oracleNode: accounts[1],
    oracleNode1: accounts[2],
    oracleNode2: accounts[3],
    oracleNode3: accounts[4],
    oracleNode4: accounts[5],
    stranger: accounts[6],
    consumer: accounts[7],
    consumer2: accounts[8],
  }

  const users: Users = {
    personas,
    roles,
    contracts,
  }
  return users
}

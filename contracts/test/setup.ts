import { ethers } from "hardhat";
import { Signer } from "ethers";

export interface Contracts {
  contract1: Signer;
  contract2: Signer;
  contract3: Signer;
  contract4: Signer;
  contract5: Signer;
  contract6: Signer;
  contract7: Signer;
  contract8: Signer;
}

export interface Roles {
  defaultAccount: Signer;
  oracleNode: Signer;
  oracleNode1: Signer;
  oracleNode2: Signer;
  oracleNode3: Signer;
  oracleNode4: Signer;
  stranger: Signer;
  consumer: Signer;
}

export interface Personas {
  Default: Signer;
  Carol: Signer;
  Eddy: Signer;
  Nancy: Signer;
  Ned: Signer;
  Neil: Signer;
  Nelly: Signer;
  Norbert: Signer;
}

export interface Users {
  contracts: Contracts;
  roles: Roles;
  personas: Personas;
}

export async function getUsers() {
  let accounts = await ethers.getSigners();

  const personas: Personas = {
    Default: accounts[0],
    Neil: accounts[1],
    Ned: accounts[2],
    Nelly: accounts[3],
    Nancy: accounts[4],
    Norbert: accounts[5],
    Carol: accounts[6],
    Eddy: accounts[7],
  };

  const contracts: Contracts = {
    contract1: accounts[8],
    contract2: accounts[9],
    contract3: accounts[10],
    contract4: accounts[11],
    contract5: accounts[12],
    contract6: accounts[13],
    contract7: accounts[14],
    contract8: accounts[15],
  };

  const roles: Roles = {
    defaultAccount: accounts[16],
    oracleNode: accounts[17],
    oracleNode1: accounts[18],
    oracleNode2: accounts[19],
    oracleNode3: accounts[20],
    oracleNode4: accounts[21],
    stranger: accounts[22],
    consumer: accounts[23],
  };

  const users: Users = {
    personas,
    roles,
    contracts,
  };
  return users;
}

import json
import subprocess
import os
from os import path
import shutil
import re

# A proof of concept / convenient script to quickly compile contracts and their go bindings
# Can be run from the Makefile with make compile_contracts

solc_versions = ["v0.4", "v0.6", "v0.7", "v0.8"]
rootdir = "./artifacts/contracts/ethereum/"
targetdir = "./contracts/ethereum"

# The names of the contracts that we're actually compiling to use.
used_contract_names = [
    "APIConsumer",
    "BlockhashStore",
    "DeviationFlaggingValidator",
    "Flags",
    "FluxAggregator",
    "KeeperConsumer",
    "KeeperConsumerPerformance",
    "KeeperRegistrar",
    "KeeperRegistrar2_0",
    # Note: when re generating wrappers you need to rollback changes made to old registries as they have manual changes to config and state struct names
    "KeeperRegistry1_1",
    "KeeperRegistry1_2",
    "KeeperRegistry1_3",
    "KeeperRegistry2_0",
    # Note: KeeperRegistryLogic 1.3/2.0 needs some manual changes in go wrapper after generation to avoid
    # conflict with KeeperRegistry. Hence it is commented out to not be regenerated every time
    # "KeeperRegistryLogic1_3",
    # "KeeperRegistryLogic2_0",
    "LinkToken",
    "MockETHLINKAggregator",
    "MockGASAggregator",
    "OffchainAggregator",
    "Oracle",
    "SimpleReadAccessController",
    "SimpleWriteAccessController",
    "UpkeepCounter",
    "UpkeepPerformCounterRestrictive",
    "UpkeepTranscoder"
    "VRF",
    "VRFConsumer",
    "VRFCoordinator",
    "VRFConsumerV2",
    "VRFCoordinatorV2",
    "KeeperConsumerBenchmark"
]

print("Locally installing hardhat...")
subprocess.run('npm install --save-dev hardhat', shell=True, check=True)

print("Locally installing openzeppelin contracts...")
subprocess.run(
    'npm install --save-dev @openzeppelin/contracts@^4.3.3', shell=True, check=True)

print("Modifying hardhat settings...")
with open("hardhat.config.js", "w") as hardhat_config:
    hardhat_config.write("""module.exports = {
solidity: {
    compilers: [
    {
        version: "0.8.6",
        settings: {
            optimizer: {
                enabled: true,
                runs: 50
            }
        }
    },
    {
        version: "0.8.13",
        settings: {
            optimizer: {
                enabled: true,
                runs: 50
            }
        }
    },
    {
        version: "0.7.6",
        settings: {
            optimizer: {
                enabled: true,
                runs: 50
            }
        }
    },
    {
        version: "0.6.6",
        settings: {
            optimizer: {
                enabled: true,
                runs: 50
            }
        }
    },
    {
        version: "0.6.0",
        settings: {
            optimizer: {
                enabled: true,
                runs: 50
            }
        }
    },
    {
        version: "0.4.24",
        settings: {
            optimizer: {
                enabled: true,
                runs: 50
            }
        }
    }
    ]
}
}""")

print("Compiling contracts...")
subprocess.run('npx hardhat compile', shell=True, check=True)

print("Creating contract go bindings...")
for version in solc_versions:
    for subdir, dirs, files in os.walk(rootdir + version):
        for f in files:
            if ".dbg." not in f:
                print(version, f)
                compile_contract = open(subdir + "/" + f, "r")
                data = json.load(compile_contract)
                contract_name = data["contractName"]

                abi_name = targetdir + "/" + version + "/abi/" + contract_name + ".abi"
                abi_file = open(abi_name, "w")
                abi_file.write(json.dumps(data["abi"], indent=2))

                bin_name = targetdir + "/" + version + "/bin/" + contract_name + ".bin"
                bin_file = open(bin_name, "w")
                bin_file.write(str(data["bytecode"]))
                abi_file.close()
                bin_file.close()

                if contract_name in used_contract_names:
                    go_file_name = targetdir + "/" + contract_name + ".go"
                    subprocess.run("abigen --bin=" + bin_name + " --abi=" + abi_name + " --pkg=" + contract_name + " --out=" +
                                   go_file_name, shell=True, check=True)
                    # Replace package name in file, abigen doesn't let you specify differently
                    with open(go_file_name, 'r+') as f:
                        text = f.read()
                        text = re.sub("package " + contract_name,
                                      "package ethereum", text)
                        f.seek(0)
                        f.write(text)
                        f.truncate()

print("Cleaning up Hardhat...")
subprocess.run('npm uninstall --save-dev hardhat', shell=True)
if path.exists("hardhat.config.js"):
    os.remove("hardhat.config.js")
if path.exists("package-lock.json"):
    os.remove("package-lock.json")
if path.exists("package.json"):
    os.remove("package.json")
if path.exists("node_modules/"):
    shutil.rmtree("node_modules/")
if path.exists("artifacts/"):
    shutil.rmtree("artifacts/")
if path.exists("cache/"):
    shutil.rmtree("cache/")

print("Done!")

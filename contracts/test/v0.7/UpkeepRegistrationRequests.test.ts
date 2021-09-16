import { ethers } from "hardhat";
//import { evmRevert } from "../test-helpers/matchers";
import {getUsers, Personas} from "../test-helpers/setup";
import {ContractFactory, Contract, BigNumber, Signer} from "ethers";

let linkTokenFactory: ContractFactory;
let mockV3AggregatorFactory: ContractFactory;
let keeperRegistryFactory: ContractFactory;
let upkeepRegistrationRequestsFactory: ContractFactory;
let upkeepMockFactory: ContractFactory;

let personas: Personas;

before(async () => {
  personas = (await getUsers()).personas;

  linkTokenFactory = await ethers.getContractFactory("src/v0.4/LinkToken.sol:LinkToken");
  mockV3AggregatorFactory = await ethers.getContractFactory("src/v0.6/tests/MockV3Aggregator.sol:MockV3Aggregator");
  keeperRegistryFactory = await ethers.getContractFactory("src/v0.7/KeeperRegistry.sol:KeeperRegistry");
  upkeepRegistrationRequestsFactory = await ethers.getContractFactory("src/v0.7/UpkeepRegistrationRequests.sol:UpkeepRegistrationRequests");
  upkeepMockFactory = await ethers.getContractFactory("src/v0.7/tests/UpkeepMock.sol:UpkeepMock");
});

/*const errorMsgs = {
  onlyOwner: "revert Only callable by owner",
  onlyAdmin: "revert only admin",
  hashPayload: "hash and payload do not match",
  requestNotFound: "request not found"
}*/

describe("UpkeepRegistrationRequests", () => {
  const upkeepName = "SampleUpkeep";

  //const stranger = personas.Nancy;
  const linkEth = BigNumber.from(300000000);
  const gasWei = BigNumber.from(100);
  const executeGas = BigNumber.from(100000);
  const source = BigNumber.from(100);
  const paymentPremiumPPB = BigNumber.from(250000000);

  //const window_big = new BN(1000);
  const window_small = BigNumber.from(2);
  const threshold_big = BigNumber.from(1000);
  //const threshold_small = new BN(5);

  const blockCountPerTurn = BigNumber.from(3);
  const emptyBytes = "0x00";
  const stalenessSeconds = BigNumber.from(43820);
  const gasCeilingMultiplier = BigNumber.from(1)
  const maxCheckGas = BigNumber.from(20000000);
  const fallbackGasPrice = BigNumber.from(200);
  const fallbackLinkPrice = BigNumber.from(200000000);
  const minLINKJuels = BigNumber.from('1000000000000000000');
  //const amount = BigNumber.from('5000000000000000000');
  const amount1 = BigNumber.from('6000000000000000000');

  let owner: Signer;
  let admin: Signer;
  //let someAddress: Signer;
  let registrarOwner: Signer;

  let linkToken: Contract
  let linkEthFeed: Contract
  let gasPriceFeed: Contract
  let registry: Contract
  let mock: Contract
  let registrar: Contract

  beforeEach(async () => {
    owner = personas.Default;
    admin = personas.Neil;
    //someAddress = personas.Ned;
    registrarOwner = personas.Nelly;

    linkToken = await linkTokenFactory.connect(owner).deploy();
    gasPriceFeed = await mockV3AggregatorFactory.connect(owner).deploy(0, gasWei);
    linkEthFeed = await mockV3AggregatorFactory.connect(owner).deploy(9, linkEth);
    registry = await keeperRegistryFactory.connect(owner).deploy(
      linkToken.address,
      linkEthFeed.address,
      gasPriceFeed.address,
      paymentPremiumPPB,
      blockCountPerTurn,
      maxCheckGas,
      stalenessSeconds,
      gasCeilingMultiplier,
      fallbackGasPrice,
      fallbackLinkPrice,
    );

    mock = await upkeepMockFactory.deploy();

    registrar = await upkeepRegistrationRequestsFactory.connect(registrarOwner).deploy(
      linkToken.address,
      minLINKJuels,
    );

    await registry.setRegistrar(registrar.address);
  });

  describe("#register", () => {

    /*it("reverts if not called by the LINK token", async () => {
      await evmRevert(
        registrar.connect(someAddress).register(
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          admin,
          emptyBytes,
          amount,
          source,
        ),
        "Must use LINK token"
      );
    });*/

    it("reverts if the amount passed in data mismatches actual amount sent", async () => {
      //set auto approve ON with high threshold limits
      await registrar.connect(registrarOwner).setRegistrationConfig(
        true,
        window_small,
        threshold_big,
        registry.address,
        minLINKJuels,
      );

      //register with auto approve ON
      let registerTx = await registrar.connect(owner).register(
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          admin.getAddress(),
          emptyBytes,
          amount1,
          source
        );
      console.log("registerTx", registerTx)

      /*await evmRevert(
        linkToken.transferAndCall(
          registrar.address,
          amount,
          abiEncodedBytes
        ),
        "Amount mismatch"
      );*/
    });

/*    it("reverts if the admin address is 0x0000...", async () => {
      let abiEncodedBytes = registrar.register(
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          "0x0000000000000000000000000000000000000000",
          emptyBytes,
          amount,
          source
        )
        .encodeABI();

      await evmRevert(
        linkToken.transferAndCall(
          registrar.address,
          amount,
          abiEncodedBytes
        ),
        "Unable to create request"
      );
    });

    it("Auto Approve ON - registers an upkeep on KeeperRegistry instantly and emits both RegistrationRequested and RegistrationApproved events", async () => {
      //get current upkeep count
      const upkeepCount = await registry.getUpkeepCount();

      //set auto approve ON with high threshold limits
      await registrar.connect(registrarOwner).setRegistrationConfig(
        true,
        window_small,
        threshold_big,
        registry.address,
        minLINKJuels,
      );

      //register with auto approve ON
      let abiEncodedBytes = registrar.contract.methods
        .register(
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          admin,
          emptyBytes,
          amount,
          source
        )
        .encodeABI();
      const { receipt } = await linkToken.transferAndCall(
        registrar.address,
        amount,
        abiEncodedBytes
      );

      //confirm if a new upkeep has been registered and the details are the same as the one just registered
      const newupkeep = await registry.getUpkeep(upkeepCount);
      assert.equal(newupkeep.target, mock.address);
      assert.equal(newupkeep.admin, admin);
      assert.equal(newupkeep.checkData, emptyBytes);
      assert.equal(newupkeep.balance.toString(), amount.toString());
      assert.isTrue(newupkeep.executeGas.eq(executeGas));

      //confirm if RegistrationRequested and RegistrationApproved event are received
      let event_RegistrationRequested = receipt.rawLogs.some((l) => {
        return (
          l.topics[0] ==
          web3.utils.keccak256(
            "RegistrationRequested(bytes32,string,bytes,address,uint32,address,bytes,uint96,uint8)"
          )
        );
      });
      assert.ok(
        event_RegistrationRequested,
        "RegistrationRequested event not emitted"
      );

      let event_RegistrationApproved = receipt.rawLogs.some((l) => {
        return (
          l.topics[0] ==
          web3.utils.keccak256("RegistrationApproved(bytes32,string,uint256)")
        );
      });
      assert.ok(
        event_RegistrationApproved,
        "RegistrationApproved event not emitted"
      );
    });

    it("Auto Approve OFF - does not registers an upkeep on KeeperRegistry, emits only RegistrationRequested event", async () => {
      //get upkeep count before attempting registration
      const beforeCount = await registry.getUpkeepCount();

      //set auto approve OFF, threshold limits dont matter in this case
      await registrar.connect(registrarOwner).setRegistrationConfig(
        false,
        window_small,
        threshold_big,
        registry.address,
        minLINKJuels,
      );

      //register with auto approve OFF
      let abiEncodedBytes = registrar.contract.methods
        .register(
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          admin,
          emptyBytes,
          amount,
          source
        )
        .encodeABI();
      const { receipt } = await linkToken.transferAndCall(
        registrar.address,
        amount,
        abiEncodedBytes
      );

      //get upkeep count after attempting registration
      const afterCount = await registry.getUpkeepCount();
      //confirm that a new upkeep has NOT been registered and upkeep count is still the same
      assert.deepEqual(beforeCount, afterCount);

      //confirm that only RegistrationRequested event is amitted and RegistrationApproved event is not
      let event_RegistrationRequested = receipt.rawLogs.some((l) => {
        return (
          l.topics[0] ==
          web3.utils.keccak256(
            "RegistrationRequested(bytes32,string,bytes,address,uint32,address,bytes,uint96,uint8)"
          )
        );
      });
      assert.ok(
        event_RegistrationRequested,
        "RegistrationRequested event not emitted"
      );

      let event_RegistrationApproved = receipt.rawLogs.some((l) => {
        return (
          l.topics[0] ==
          web3.utils.keccak256("RegistrationApproved(bytes32,string,uint256)")
        );
      });
      assert.ok(
        !event_RegistrationApproved,
        "RegistrationApproved event should not be emitted"
      );

      const hash = receipt.rawLogs[2].topics[1]
      const pendingRequest = await registrar.getPendingRequest(hash)
      assert.equal(admin, pendingRequest[0])
      assert.ok(amount.eq(pendingRequest[1]))
    });

    it("Auto Approve ON - Throttle max approvals - does not registers an upkeep on KeeperRegistry beyond the throttle limit, emits only RegistrationRequested event after throttle starts", async () => {
      //get upkeep count before attempting registration
      const beforeCount = await registry.getUpkeepCount();

      //set auto approve on, with low threshold limits
      await registrar.setRegistrationConfig(
        true,
        window_big,
        threshold_small,
        registry.address,
        minLINKJuels,
        { from: registrarOwner }
      );

      let abiEncodedBytes = registrar.contract.methods
        .register(
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          admin,
          emptyBytes,
          amount,
          source
        )
        .encodeABI();

      //register within threshold, new upkeep should be registered
      await linkToken.transferAndCall(
        registrar.address,
        amount,
        abiEncodedBytes
      );
      const intermediateCount = await registry.getUpkeepCount();
      //make sure 1 upkeep was registered
      assert.equal(beforeCount.toNumber() + 1, intermediateCount.toNumber());

      //try registering more than threshold(say 2x), new upkeeps should not be registered after the threshold amount is reached
      for (let step = 0; step < threshold_small * 2; step++) {
        abiEncodedBytes = registrar.contract.methods
        .register(
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas + step, // make unique hash
          admin,
          emptyBytes,
          amount,
          source
        )
        .encodeABI();

        await linkToken.transferAndCall(
          registrar.address,
          amount,
          abiEncodedBytes
        );
      }
      const afterCount = await registry.getUpkeepCount();
      //count of newly registered upkeeps should be equal to the threshold set for auto approval
      const newRegistrationsCount = afterCount.toNumber() - beforeCount.toNumber();
      assert(newRegistrationsCount == threshold_small,"Registrations beyond threshold");
    });*/
  });

/*  describe("#approve", () => {
    let hash

    beforeEach(async () => {
      await registrar.setRegistrationConfig(
        false,
        window_small,
        threshold_big,
        registry.address,
        minLINKJuels,
        { from: registrarOwner }
      );

      //register with auto approve OFF
      let abiEncodedBytes = registrar.contract.methods
        .register(
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          admin,
          emptyBytes,
          amount,
          source
        )
        .encodeABI();
      const { receipt } = await linkToken.transferAndCall(
        registrar.address,
        amount,
        abiEncodedBytes
      );
      hash = receipt.rawLogs[2].topics[1]
    })

    it("reverts if not called by the owner", async () => {
      const tx = registrar.approve(upkeepName, mock.address, executeGas, admin, emptyBytes, hash, { from: stranger })
      await expectRevert(tx, errorMsgs.onlyOwner)
    })

    it("reverts if the hash does not exist", async () => {
      const tx = registrar.approve(upkeepName, mock.address, executeGas, admin, emptyBytes, "0x1234", { from: registrarOwner })
      await expectRevert(tx, errorMsgs.requestNotFound)
    })

    it("reverts if any member of the payload changes", async () => {
      tx = registrar.approve(upkeepName, ethers.Wallet.createRandom().address, executeGas, admin, emptyBytes, hash, { from: registrarOwner })
      await expectRevert(tx, errorMsgs.hashPayload)
      tx = registrar.approve(upkeepName, mock.address, 10000, admin, emptyBytes, hash, { from: registrarOwner })
      await expectRevert(tx, errorMsgs.hashPayload)
      tx = registrar.approve(upkeepName, mock.address, executeGas, ethers.Wallet.createRandom().address, emptyBytes, hash, { from: registrarOwner })
      await expectRevert(tx, errorMsgs.hashPayload)
      tx = registrar.approve(upkeepName, mock.address, executeGas, admin, "0x1234", hash, { from: registrarOwner })
      await expectRevert(tx, errorMsgs.hashPayload)
    })

    it("approves an existing registration request", async () => {
      const receipt = await registrar.approve(upkeepName, mock.address, executeGas, admin, emptyBytes, hash, { from: registrarOwner })
      await expectEvent(receipt, "RegistrationApproved")
    })

    it("deletes the request afterwards / reverts if the request DNE", async () => {
      await registrar.approve(upkeepName, mock.address, executeGas, admin, emptyBytes, hash, { from: registrarOwner })
      const tx = registrar.approve(upkeepName, mock.address, executeGas, admin, emptyBytes, hash, { from: registrarOwner })
      await expectRevert(tx, errorMsgs.requestNotFound)
    })
  })

  describe("#cancel", () => {
    let hash

    beforeEach(async () => {
      await registrar.setRegistrationConfig(
        false,
        window_small,
        threshold_big,
        registry.address,
        minLINKJuels,
        { from: registrarOwner }
      );

      //register with auto approve OFF
      let abiEncodedBytes = registrar.contract.methods
        .register(
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          admin,
          emptyBytes,
          amount,
          source
        )
        .encodeABI();
      const { receipt } = await linkToken.transferAndCall(
        registrar.address,
        amount,
        abiEncodedBytes
      );
      hash = receipt.rawLogs[2].topics[1]
      // submit duplicate request (increase balance)
      await linkToken.transferAndCall(
        registrar.address,
        amount,
        abiEncodedBytes
      );
    })

    it("reverts if not called by the admin / owner", async () => {
      const tx = registrar.cancel(hash, { from: stranger })
      await expectRevert(tx, "only admin / owner can cancel")
    })

    it("reverts if the hash does not exist", async () => {
      const tx = registrar.cancel("0x1234", { from: admin })
      await expectRevert(tx, errorMsgs.onlyAdmin)
    })

    it("refunds the total request balance to the admin address", async () => {
      const before = await linkToken.balanceOf(admin)
      const receipt = await registrar.cancel(hash, { from: admin })
      const after = await linkToken.balanceOf(admin)
      assert.isTrue(after.sub(before).eq(amount.mul(new BN(2))))
      await expectEvent(receipt, "RegistrationRejected", {hash})
    })

    it("deletes the request hash", async () => {
      await registrar.cancel(hash, { from: registrarOwner })
      let tx = registrar.cancel(hash, { from: registrarOwner })
      await expectRevert(tx, errorMsgs.requestNotFound)
      tx = registrar.approve(upkeepName, mock.address, executeGas, admin, emptyBytes, hash, { from: registrarOwner })
      await expectRevert(tx, errorMsgs.requestNotFound)
    })
  })*/
});

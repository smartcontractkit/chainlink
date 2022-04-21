package handler

import (
	"context"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cron_factory "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/cron_upkeep_factory_wrapper"
	proxy "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/permissioned_forward_proxy_wrapper"
	registrar "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/upkeep_registration_requests_wrapper"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

func (k *Keeper) MigrateCron(ctx context.Context, inputFilePath string, proxyAbi abi.ABI) {
	outputFile := "migrate_cron_output_" + time.Now().String() + ".csv"
	o, err := os.Create(outputFile)
	if err != nil {
		log.Fatalln("Error writing to "+outputFile, err)
	}
	w := csv.NewWriter(o)
	defer o.Close()

	for _, inputLine := range readCsvFile(inputFilePath) {
		targetContractAddr := common.HexToAddress(inputLine[0])
		targetFunction := inputLine[1]
		if strings.Contains(targetFunction, "(") || strings.Contains(targetFunction, ")") {
			log.Fatalln("Error targetFunction should strictly be the function name and not contain parenthesis")
		}
		cronSchedule := inputLine[2]
		fundingAmountLink := new(big.Int)
		fundingAmountLink.SetString(inputLine[3], 10)
		if fundingAmountLink.Cmp(big.NewInt(1e18)) < 0 {
			log.Fatalln("Error fundingAmountLink should at least be 1 LINK (1e18)")
		}
		upkeepName := inputLine[4]
		if upkeepName == "" {
			log.Fatalln("Error upkeep name cannot be empty")
		}
		encryptedEmail, err := hex.DecodeString(inputLine[5])
		if err != nil {
			log.Fatalln("Error decoding encrypted email:", inputLine[5], err)
		}
		upkeepAdminAddr := common.HexToAddress(inputLine[6])
		gasLimit, err := strconv.ParseUint(inputLine[7], 10, 32)
		if err != nil {
			log.Fatalln("Error parsing gas limit:", inputLine[7], err)
		}

		fmt.Println("Processing:", targetContractAddr, targetFunction, cronSchedule)
		targetHandler := getTargetHandler(targetFunction)                                   // Encoding of function call on target
		cronByteHandler, err := proxyAbi.Pack("forward", targetContractAddr, targetHandler) // function call on target through proxy
		if err != nil {
			log.Fatalln("Error generating cron byte handler", err)
		}

		proxyAddr := common.HexToAddress(k.cfg.ProxyAddr)
		cronUpkeepAddr := k.deployNewCronUpkeep(ctx, proxyAddr, cronByteHandler, cronSchedule)
		k.setProxyPermission(ctx, cronUpkeepAddr, targetContractAddr)
		registrationHash, blockNum := k.registerUpkeep(ctx, upkeepName, encryptedEmail, cronUpkeepAddr, upkeepAdminAddr, uint32(gasLimit), fundingAmountLink)

		row := []string{
			targetContractAddr.String(), targetFunction, cronSchedule,
			upkeepName, cronUpkeepAddr.String(), strconv.FormatUint(gasLimit, 10),
			upkeepAdminAddr.String(), registrationHash, blockNum.String(),
		}
		if err := w.Write(row); err != nil {
			log.Fatalln("Error writing record to output file", err)
		}
		w.Flush()
	}
}

func (k *Keeper) deployNewCronUpkeep(ctx context.Context, targetAddr common.Address, targetHandler []byte, cronSchedule string) common.Address {
	log.Println("Deploying new cron upkeep")

	cronFactoryAddr := common.HexToAddress(k.cfg.CronFactoryAddr)
	cronFactoryInstance, err := cron_factory.NewCronUpkeepFactory(
		cronFactoryAddr,
		k.client,
	)
	if err != nil {
		log.Fatalln("Error while instantiating "+cronFactoryAddr.String()+" to cron factory", err)
	}

	callOpts := bind.CallOpts{
		Context: ctx,
	}
	encodedCronJob, err := cronFactoryInstance.EncodeCronJob(&callOpts, targetAddr, targetHandler, cronSchedule)
	if err != nil {
		log.Fatalln("Error getting encoded cron job", err)
	}

	cronJobTx, err := cronFactoryInstance.NewCronUpkeepWithJob(k.buildTxOpts(ctx), encodedCronJob)
	if err != nil {
		log.Fatalln("Error creating cron job", err)
	}
	txReceipt, err := bind.WaitMined(ctx, k.client, cronJobTx)
	if err != nil {
		log.Fatalln("Error getting receipt for cron job tx", err)
	}
	if txReceipt.Status != 1 {
		log.Fatalln("tx", cronJobTx.Hash(), "failed")
	}
	log.Println("Cron upkeep deployed:", helpers.ExplorerLink(k.cfg.ChainID, cronJobTx.Hash()))

	rawLog := *txReceipt.Logs[1]
	parsedLog, err := cronFactoryInstance.ParseLog(rawLog)
	if err != nil {
		log.Fatalln("Error parsing NewCronUpkeepCreated log", err)
	}
	cronUpkeepCreatedLog, ok := parsedLog.(*cron_factory.CronUpkeepFactoryNewCronUpkeepCreated)
	if !ok {
		log.Fatalln("Error type casting NewCronUpkeepCreated log", err)
	}
	log.Println("Cron upkeep address:", cronUpkeepCreatedLog.Upkeep)

	return cronUpkeepCreatedLog.Upkeep
}

func (k *Keeper) setProxyPermission(ctx context.Context, from, to common.Address) {
	log.Println("Setting permission on proxy")

	proxyAddr := common.HexToAddress(k.cfg.ProxyAddr)
	proxyInstance, err := proxy.NewPermissionedForwardProxy(
		proxyAddr,
		k.client,
	)
	if err != nil {
		log.Fatalln("Error while instantiating "+proxyAddr.String()+" to permissioned forward proxy", err)
	}

	proxyTx, err := proxyInstance.SetPermission(k.buildTxOpts(ctx), from, to)
	if err != nil {
		log.Fatalln("Error setting proxy permission", err)
	}
	txReceipt, err := bind.WaitMined(ctx, k.client, proxyTx)
	if err != nil {
		log.Fatalln("Error getting receipt for proxy tx", err)
	}
	if txReceipt.Status != 1 {
		log.Fatalln("tx", proxyTx.Hash(), "failed")
	}

	log.Println("Proxy permission from", from, "to", to, "set:", helpers.ExplorerLink(k.cfg.ChainID, proxyTx.Hash()))
}

func (k *Keeper) registerUpkeep(ctx context.Context, name string, encryptedEmail []byte, target, admin common.Address, gasLimit uint32, amount *big.Int) (string, *big.Int) {
	log.Println("Registering upkeep")

	registrarAddr := common.HexToAddress(k.cfg.RegistrarAddr)
	registrarInstance, err := registrar.NewUpkeepRegistrationRequests(
		registrarAddr,
		k.client,
	)
	if err != nil {
		log.Fatalln("Error while instantiating "+registrarAddr.String()+" to registrar", err)
	}
	registrarABI, err := abi.JSON(strings.NewReader(registrar.UpkeepRegistrationRequestsABI))
	if err != nil {
		log.Fatalln("Error generating Registrar ABI", err)
	}

	registrationData, err := registrarABI.Pack("register", name, encryptedEmail, target, gasLimit, admin, []byte{}, amount, uint8(0))
	if err != nil {
		log.Fatalln("Error generating registration data", err)
	}
	registrationTx, err := k.linkToken.TransferAndCall(k.buildTxOpts(ctx), registrarAddr, amount, registrationData)
	if err != nil {
		log.Fatalln("registering", err)
	}
	txReceipt, err := bind.WaitMined(ctx, k.client, registrationTx)
	if err != nil {
		log.Fatalln("Error getting receipt for upkeep register tx", err)
	}
	if txReceipt.Status != 1 {
		log.Fatalln("tx", registrationTx.Hash(), "failed")
	}
	log.Println("Upkeep registered:", helpers.ExplorerLink(k.cfg.ChainID, registrationTx.Hash()))

	rawLog := *txReceipt.Logs[2]
	parsedLog, err := registrarInstance.ParseLog(rawLog)
	if err != nil {
		log.Fatalln("Error parsing RegistrationRequested log", err)
	}
	registrationRequestedLog, ok := parsedLog.(*registrar.UpkeepRegistrationRequestsRegistrationRequested)
	if !ok {
		log.Fatalln("Error type casting RegistrationRequested log", err)
	}
	hash := hex.EncodeToString(registrationRequestedLog.Hash[:])
	log.Println("Registration request hash:", hash)

	return hash, txReceipt.BlockNumber
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func getTargetHandler(targetFunction string) []byte {
	targetAbi, err := abi.JSON(strings.NewReader(`[{"type":"function","inputs":[],"outputs":[],"name":"` + targetFunction + `"}]`))
	if err != nil {
		log.Fatalln("Error generating target ABI", err)
	}
	targetHandler, err := targetAbi.Pack(targetFunction)
	if err != nil {
		log.Fatalln("Error generating target handler", err)
	}
	return targetHandler
}

func (k *Keeper) FetchUpkeepIds(ctx context.Context, inputFilePath string) {
	outputFile := "migrate_cron_output_" + time.Now().String() + ".csv"
	o, err := os.Create(outputFile)
	if err != nil {
		log.Fatalln("Error writing to "+outputFile, err)
	}
	w := csv.NewWriter(o)
	defer o.Close()

	var inputHashes [][32]byte
	var startingBlock uint64

	startingBlock = math.MaxUint64
	inputLines := readCsvFile(inputFilePath)
	for _, inputLine := range inputLines {
		hexHash, err := hex.DecodeString(inputLine[7])
		if err != nil {
			log.Fatalln("Error parsing hash:", inputLine[7], err)
		}
		if len(hexHash) != 32 {
			log.Fatalln("Error, hash should be 32 bytes:", inputLine[7], err)
		}
		var hash [32]byte
		copy(hash[:], hexHash)
		inputHashes = append(inputHashes, hash)

		blockNum, err := strconv.ParseUint(inputLine[8], 10, 64)
		if err != nil {
			log.Fatalln("Error parsing block number:", inputLine[8], err)
		}
		if blockNum < startingBlock {
			startingBlock = blockNum
		}
	}

	if startingBlock == math.MaxUint64 {
		log.Fatalln("Error: No starting block numbers given in input")
	}

	registrarAddr := common.HexToAddress(k.cfg.RegistrarAddr)
	registrarInstance, err := registrar.NewUpkeepRegistrationRequests(
		registrarAddr,
		k.client,
	)
	if err != nil {
		log.Fatalln("Error while instantiating "+registrarAddr.String()+" to registrar", err)
	}

	filterOpts := bind.FilterOpts{
		Start:   startingBlock,
		End:     nil, // Go till latest
		Context: ctx,
	}

	registrarIterator, err := registrarInstance.FilterRegistrationApproved(&filterOpts, inputHashes, nil)
	if err != nil {
		log.Fatalln("Failed to get registrar iterator", err)
	}

	m := make(map[string]string)
	for registrarIterator.Next() {
		// Iterate over the events and store the upkeedID per hash
		hash := hex.EncodeToString(registrarIterator.Event.Hash[:])
		upkeepId := registrarIterator.Event.UpkeepId.String()
		m[hash] = upkeepId
	}

	for _, inputLine := range inputLines {
		hash := inputLine[7]
		row := []string{inputLine[0], inputLine[1], inputLine[2], inputLine[3], inputLine[4], inputLine[5], inputLine[6], inputLine[7], inputLine[8]}
		upkeepId, ok := m[hash]
		if ok {
			log.Println("Found upkeep ID for hash", hash, upkeepId)
			row = append(row, upkeepId)
		} else {
			log.Println("Failed to find upkeep ID for hash", hash)
		}
		if err := w.Write(row); err != nil {
			log.Fatalln("Error writing record to output file", err)
		}
		w.Flush()
	}
}

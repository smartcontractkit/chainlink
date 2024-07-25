import { route53RecordsExist } from "./lib/check-route53-records.js";

// Example usage
const hostedZoneId = "Z0701115F9JUQY0J2ISL"; // Your hosted zone ID here
const recordNames = [
  "cname-crib-chainlink-11635-geth-ws.",
  "crib-chainlink-11635-node4.",
]; // DNS record names you want to check
const maxRetries = 3; // Maximum number of retries

route53RecordsExist(hostedZoneId, recordNames, maxRetries)
  .then((result) => console.log("Records exist:", result))
  .catch((error) => console.error("Error:", error));

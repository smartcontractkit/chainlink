import {
  Route53Client,
  ListResourceRecordSetsCommand,
} from "@aws-sdk/client-route-53";

// us-east-1 is the global region used by Route 53.
const route53Client = new Route53Client({ region: "us-east-1" });

// Function to wait for a specified amount of time (ms).
const wait = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

/**
 * Check if Route 53 records exist for a given Route 53 zone.
 *
 * @param {string} hostedZoneId The ID of the hosted zone.
 * @param {string[]} recordNames An array of record names to check.
 * @param {number} maxRetries The maximum number of retries.
 * @param {number} initialBackoff The initial backoff time in milliseconds.
 * @returns {Promise<boolean>} True if records exist, false otherwise.
 */
async function route53RecordsExist(
  hostedZoneId,
  recordNames,
  maxRetries = 7,
  initialBackoff = 2000 // 2 seconds
) {
  let attempts = 0;

  while (attempts < maxRetries) {
    try {
      let isTruncated = true;
      let nextRecordName;
      let nextRecordType;
      let allRecordSets = [];

      while (isTruncated) {
        const params = {
          HostedZoneId: hostedZoneId,
          MaxItems: "300",
          ...(nextRecordName && { StartRecordName: nextRecordName }),
          ...(nextRecordType && { StartRecordType: nextRecordType }),
        };

        const data = await route53Client.send(
          new ListResourceRecordSetsCommand(params)
        );
        allRecordSets = allRecordSets.concat(data.ResourceRecordSets);
        isTruncated = data.IsTruncated;
        if (isTruncated) {
          nextRecordName = data.NextRecordName;
          nextRecordType = data.NextRecordType;
        }
      }

      const recordExists = recordNames.every((name) =>
        allRecordSets.some((r) => r.Name.includes(name))
      );

      if (recordExists) {
        return true; // All specified records are found
      } else {
        // If any record is not found, throw an error to trigger a retry
        throw new Error("One or more records not found, retrying...");
      }
    } catch (error) {
      console.error(`Attempt ${attempts + 1}:`, error.message);
      if (attempts === maxRetries - 1) {
        return false; // Return false after the last attempt
      }
      // Exponential backoff
      await wait(initialBackoff * 2 ** attempts);
      attempts++;
    }
  }
  // Should not reach here if retries are exhausted
  return false;
}

export { route53RecordsExist };

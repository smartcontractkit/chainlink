const fs = require("fs");
const request = require("request-promise");
const util = require('util');
const url = "http://chainlink:twochains@localhost:6688/v2/specs";

module.exports = {
  // Deploys chainlink jobs.
  job: async function(filename) {
    let readFile = util.promisify(fs.readFile);
    return await readFile(filename, 'utf8').then((file) => {
      let data = JSON.parse(file);
      console.log(`Posting to ${url}:\n`, data);
      return request.post(url, {json: data});
    });
  }
};

const request = require("request-promise");

module.exports = {
  // Deploys chainlink jobs.
  job: function(url, data) {
    console.log(`Posting to ${url}:\n`, data);
    return request.post(url, {json: data});
  }
};

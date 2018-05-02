let fs = require("fs");
let request = require("request");
let url = "http://"+process.env.USERNAME+":"+process.env.PASSWORD+"@localhost:6688/v2/specs";

module.exports = {
  // Deploys chainlink jobs.
  job: function(data, callback, callbackError) {
    let job = JSON.parse(data);
    console.log(`Posting to ${url}:\n`, job);
    request.post(url, {json: job}, function (error, response, body) {
        if (!error && response && response.statusCode == 200) {
          callback(error, response, body);
        } else {
          if (callbackError) {
            callbackError(error, response);
          }
        }
      }
    );
  }
};

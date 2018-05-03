let fs = require("fs");
let request = require("request");
let url = "http://chainlink:twochains@localhost:6688/v2/specs";

module.exports = {
  // Deploys chainlink jobs.
  job: function(filename, callback, callbackError) {
    fs.readFile(filename, 'utf8', (err, file) => {
      let data = JSON.parse(file);
      console.log(`Posting to ${url}:\n`, data);
      request.post(url, {json: data},
        function (error, response, body) {
          if (!error && response && response.statusCode == 200) {
            callback(error, response, body);
          } else {
            if (callbackError) {
              callbackError(error, response);
            }
          }
        }
      );
    });
  }
};

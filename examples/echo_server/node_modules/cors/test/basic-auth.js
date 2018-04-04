(function () {
  /*global describe, it*/

  'use strict';

  var should = require('should'),
    express = require('express'),
    supertest = require('supertest'),
    basicAuth = require('basic-auth-connect'),
    cors = require('../lib');

  var app;

  /* -------------------------------------------------------------------------- */

  app = express();
  app.use(basicAuth('username', 'password'));
  app.use(cors());
  app.post('/', function (req, res) {
    res.send('hello world');
  });

  /* -------------------------------------------------------------------------- */

  describe('basic auth', function () {
    it('POST works', function (done) {
      supertest(app)
        .post('/')
        .auth('username', 'password')
        .expect(200)
        .end(function (err, res) {
          should.not.exist(err);
          res.headers['access-control-allow-origin'].should.eql('*');
          res.text.should.eql('hello world');
          done();
        });
    });
  });

}());

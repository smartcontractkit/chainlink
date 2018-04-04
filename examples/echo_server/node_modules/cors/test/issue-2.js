(function () {
  /*global describe, it*/

  'use strict';

  var should = require('should'),
    express = require('express'),
    supertest = require('supertest'),
    cors = require('../lib');

  var app,
    corsOptions;

  /* -------------------------------------------------------------------------- */

  app = express();
  corsOptions = {
    origin: true,
    methods: ['POST'],
    credentials: true,
    maxAge: 3600
  };
  app.options('/api/login', cors(corsOptions));
  app.post('/api/login', cors(corsOptions), function (req, res) {
    res.send('LOGIN');
  });

  /* -------------------------------------------------------------------------- */

  describe('issue  #2', function () {
    it('OPTIONS works', function (done) {
      supertest(app)
        .options('/api/login')
        .expect(204)
        .set('Origin', 'http://example.com')
        .end(function (err, res) {
          should.not.exist(err);
          res.headers['access-control-allow-origin'].should.eql('http://example.com');
          done();
        });
    });
    it('POST works', function (done) {
      supertest(app)
        .post('/api/login')
        .expect(200)
        .set('Origin', 'http://example.com')
        .end(function (err, res) {
          should.not.exist(err);
          res.headers['access-control-allow-origin'].should.eql('http://example.com');
          res.text.should.eql('LOGIN');
          done();
        });
    });
  });

}());

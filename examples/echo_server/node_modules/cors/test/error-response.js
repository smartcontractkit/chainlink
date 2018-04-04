(function () {
  /*global describe, it*/

  'use strict';

  var should = require('should'),
    express = require('express'),
    supertest = require('supertest'),
    cors = require('../lib');

  var app;

  /* -------------------------------------------------------------------------- */

  app = express();
  app.use(cors());

  app.post('/five-hundred', function (req, res, next) {
    next(new Error('nope'));
  });

  app.post('/four-oh-one', function (req, res, next) {
    next(new Error('401'));
  });

  app.post('/four-oh-four', function (req, res, next) {
    next();
  });

  app.use(function (err, req, res, next) {
    if (err.message === '401') {
      res.status(401).send('unauthorized');
    } else {
      next(err);
    }
  });

  /* -------------------------------------------------------------------------- */

  describe('error response', function () {
    it('500', function (done) {
      supertest(app)
        .post('/five-hundred')
        .expect(500)
        .end(function (err, res) {
          should.not.exist(err);
          res.headers['access-control-allow-origin'].should.eql('*');
          res.text.should.containEql('Error: nope');
          done();
        });
    });

    it('401', function (done) {
      supertest(app)
        .post('/four-oh-one')
        .expect(401)
        .end(function (err, res) {
          should.not.exist(err);
          res.headers['access-control-allow-origin'].should.eql('*');
          res.text.should.eql('unauthorized');
          done();
        });
    });

    it('404', function (done) {
      supertest(app)
        .post('/four-oh-four')
        .expect(404)
        .end(function (err, res) {
          should.not.exist(err);
          res.headers['access-control-allow-origin'].should.eql('*');
          done();
        });
    });
  });

}());

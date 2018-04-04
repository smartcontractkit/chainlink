(function () {
  /*global describe, it*/

  'use strict';

  var should = require('should'),
    cors = require('../lib');

  var fakeRequest = function (headers) {
      return {
        headers: headers || {
          'origin': 'request.com',
          'access-control-request-headers': 'requestedHeader1,requestedHeader2'
        },
        pause: function () {
          // do nothing
          return;
        },
        resume: function () {
          // do nothing
          return;
        }
      };
    },
    fakeResponse = function () {
      var headers = {};
      return {
        allHeaders: function () {
          return headers;
        },
        getHeader: function (key) {
          return headers[key];
        },
        setHeader: function (key, value) {
          headers[key] = value;
          return;
        },
        get: function (key) {
          return headers[key];
        }
      };
    };

  describe('cors', function () {
    it('does not alter `options` configuration object', function () {
      var options = Object.freeze({
        origin: 'custom-origin'
      });
      (function () {
        cors(options);
      }).should.not.throw();
    });

    it('passes control to next middleware', function (done) {
      // arrange
      var req, res, next;
      req = fakeRequest();
      res = fakeResponse();
      next = function () {
        done();
      };

      // act
      cors()(req, res, next);
    });

    it('shortcircuits preflight requests', function (done) {
      // arrange
      var req, res, next;
      req = fakeRequest();
      req.method = 'OPTIONS';
      res = fakeResponse();
      res.end = function () {
        // assert
        res.statusCode.should.equal(204);
        done();
      };
      next = function () {
        // assert
        done('should not be called');
      };

      // act
      cors()(req, res, next);
    });

    it('can configure preflight success response status code', function (done) {
      // arrange
      var req, res, next;
      req = fakeRequest();
      req.method = 'OPTIONS';
      res = fakeResponse();
      res.end = function () {
        // assert
        res.statusCode.should.equal(200);
        done();
      };
      next = function () {
        // assert
        done('should not be called');
      };

      // act
      cors({optionsSuccessStatus: 200})(req, res, next);
    });

    it('doesn\'t shortcircuit preflight requests with preflightContinue option', function (done) {
      // arrange
      var req, res, next;
      req = fakeRequest();
      req.method = 'OPTIONS';
      res = fakeResponse();
      res.end = function () {
        // assert
        done('should not be called');
      };
      next = function () {
        // assert
        done();
      };

      // act
      cors({preflightContinue: true})(req, res, next);
    });

    it('normalizes method names', function (done) {
      // arrange
      var req, res, next;
      req = fakeRequest();
      req.method = 'options';
      res = fakeResponse();
      res.end = function () {
        // assert
        res.statusCode.should.equal(204);
        done();
      };
      next = function () {
        // assert
        done('should not be called');
      };

      // act
      cors()(req, res, next);
    });

    it('includes Content-Length response header', function (done) {
      // arrange
      var req, res, next;
      req = fakeRequest();
      req.method = 'options';
      res = fakeResponse();
      res.end = function () {
        // assert
        res.getHeader('Content-Length').should.equal('0');
        done();
      };
      next = function () {
        // assert
        done('should not be called');
      };

      // act
      cors()(req, res, next);
    });

    it('no options enables default CORS to all origins', function (done) {
      // arrange
      var req, res, next;
      req = fakeRequest();
      res = fakeResponse();
      next = function () {
        // assert
        res.getHeader('Access-Control-Allow-Origin').should.equal('*');
        should.not.exist(res.getHeader('Access-Control-Allow-Methods'));
        done();
      };

      // act
      cors()(req, res, next);
    });

    it('OPTION call with no options enables default CORS to all origins and methods', function (done) {
      // arrange
      var req, res, next;
      req = fakeRequest();
      req.method = 'OPTIONS';
      res = fakeResponse();
      res.end = function () {
        // assert
        res.statusCode.should.equal(204);
        done();
      };
      next = function () {
        // assert
        res.getHeader('Access-Control-Allow-Origin').should.equal('*');
        res.getHeader('Access-Control-Allow-Methods').should.equal('GET,PUT,PATCH,POST,DELETE');
        done();
      };

      // act
      cors()(req, res, next);
    });

    describe('passing static options', function () {
      it('overrides defaults', function (done) {
        // arrange
        var req, res, next, options;
        options = {
          origin: 'example.com',
          methods: ['FOO', 'bar'],
          headers: ['FIZZ', 'buzz'],
          credentials: true,
          maxAge: 123
        };
        req = fakeRequest();
        req.method = 'OPTIONS';
        res = fakeResponse();
        res.end = function () {
          // assert
          res.statusCode.should.equal(204);
          done();
        };
        next = function () {
          // assert
          res.getHeader('Access-Control-Allow-Origin').should.equal('example.com');
          res.getHeader('Access-Control-Allow-Methods').should.equal('FOO,bar');
          res.getHeader('Access-Control-Allow-Headers').should.equal('FIZZ,buzz');
          res.getHeader('Access-Control-Allow-Credentials').should.equal('true');
          res.getHeader('Access-Control-Max-Age').should.equal('123');
          done();
        };

        // act
        cors(options)(req, res, next);
      });

      it('matches request origin against regexp', function(done) {
        var req = fakeRequest();
        var res = fakeResponse();
        var options = { origin: /^(.+\.)?request.com$/ };
        cors(options)(req, res, function(err) {
          should.not.exist(err);
          res.getHeader('Access-Control-Allow-Origin').should.equal(req.headers.origin);
          should.exist(res.getHeader('Vary'));
          res.getHeader('Vary').should.equal('Origin');
          return done();
        });
      });

      it('matches request origin against array of origin checks', function(done) {
        var req = fakeRequest();
        var res = fakeResponse();
        var options = { origin: [ /foo\.com$/, 'request.com' ] };
        cors(options)(req, res, function(err) {
          should.not.exist(err);
          res.getHeader('Access-Control-Allow-Origin').should.equal(req.headers.origin);
          should.exist(res.getHeader('Vary'));
          res.getHeader('Vary').should.equal('Origin');
          return done();
        });
      });

      it('doesn\'t match request origin against array of invalid origin checks', function(done) {
        var req = fakeRequest();
        var res = fakeResponse();
        var options = { origin: [ /foo\.com$/, 'bar.com' ] };
        cors(options)(req, res, function(err) {
          should.not.exist(err);
          should.not.exist(res.getHeader('Access-Control-Allow-Origin'));
          should.exist(res.getHeader('Vary'));
          res.getHeader('Vary').should.equal('Origin');
          return done();
        });
      });

      it('origin of false disables cors', function (done) {
        // arrange
        var req, res, next, options;
        options = {
          origin: false,
          methods: ['FOO', 'bar'],
          headers: ['FIZZ', 'buzz'],
          credentials: true,
          maxAge: 123
        };
        req = fakeRequest();
        res = fakeResponse();
        next = function () {
          // assert
          should.not.exist(res.getHeader('Access-Control-Allow-Origin'));
          should.not.exist(res.getHeader('Access-Control-Allow-Methods'));
          should.not.exist(res.getHeader('Access-Control-Allow-Headers'));
          should.not.exist(res.getHeader('Access-Control-Allow-Credentials'));
          should.not.exist(res.getHeader('Access-Control-Max-Age'));
          done();
        };

        // act
        cors(options)(req, res, next);
      });

      it('can override origin', function (done) {
        // arrange
        var req, res, next, options;
        options = {
          origin: 'example.com'
        };
        req = fakeRequest();
        res = fakeResponse();
        next = function () {
          // assert
          res.getHeader('Access-Control-Allow-Origin').should.equal('example.com');
          done();
        };

        // act
        cors(options)(req, res, next);
      });

      it('includes Vary header for specific origins', function (done) {
        // arrange
        var req, res, next, options;
        options = {
          origin: 'example.com'
        };
        req = fakeRequest();
        res = fakeResponse();
        next = function () {
          // assert
          should.exist(res.getHeader('Vary'));
          res.getHeader('Vary').should.equal('Origin');
          done();
        };

        // act
        cors(options)(req, res, next);
      });

      it('appends to an existing Vary header', function (done) {
        // arrange
        var req, res, next, options;
        options = {
          origin: 'example.com'
        };
        req = fakeRequest();
        res = fakeResponse();
        res.setHeader('Vary', 'Foo');
        next = function () {
          // assert
          res.getHeader('Vary').should.equal('Foo, Origin');
          done();
        };

        // act
        cors(options)(req, res, next);
      });

      it('origin defaults to *', function (done) {
        // arrange
        var req, res, next, options;
        options = {
        };
        req = fakeRequest();
        res = fakeResponse();
        next = function () {
          // assert
          res.getHeader('Access-Control-Allow-Origin').should.equal('*');
          done();
        };

        // act
        cors(options)(req, res, next);
      });

      it('specifying true for origin reflects requesting origin', function (done) {
        // arrange
        var req, res, next, options;
        options = {
          origin: true
        };
        req = fakeRequest();
        res = fakeResponse();
        next = function () {
          // assert
          res.getHeader('Access-Control-Allow-Origin').should.equal('request.com');
          done();
        };

        // act
        cors(options)(req, res, next);
      });

      it('should allow origin when callback returns true', function (done) {
        var req, res, next, options;
        options = {
          origin: function (sentOrigin, cb) {
            sentOrigin.should.equal('request.com');
            cb(null, true);
          }
        };
        req = fakeRequest();
        res = fakeResponse();
        next = function () {
          res.getHeader('Access-Control-Allow-Origin').should.equal('request.com');
          done();
        };

        cors(options)(req, res, next);
      });

      it('should not allow origin when callback returns false', function (done) {
        var req, res, next, options;
        options = {
          origin: function (sentOrigin, cb) {
            sentOrigin.should.equal('request.com');
            cb(null, false);
          }
        };
        req = fakeRequest();
        res = fakeResponse();
        next = function () {
          should.not.exist(res.getHeader('Access-Control-Allow-Origin'));
          should.not.exist(res.getHeader('Access-Control-Allow-Methods'));
          should.not.exist(res.getHeader('Access-Control-Allow-Headers'));
          should.not.exist(res.getHeader('Access-Control-Allow-Credentials'));
          should.not.exist(res.getHeader('Access-Control-Max-Age'));
          done();
        };

        cors(options)(req, res, next);
      });

      it('should not override options.origin callback', function (done) {
        var req, res, next, options;
        options = {
          origin: function (sentOrigin, cb) {
            var isValid = sentOrigin === 'request.com';
            cb(null, isValid);
          }
        };

        req = fakeRequest();
        res = fakeResponse();
        next = function () {
          res.getHeader('Access-Control-Allow-Origin').should.equal('request.com');
        };

        cors(options)(req, res, next);

        req = fakeRequest({
          'origin': 'invalid-request.com'
        });
        res = fakeResponse();

        next = function () {
          should.not.exist(res.getHeader('Access-Control-Allow-Origin'));
          should.not.exist(res.getHeader('Access-Control-Allow-Methods'));
          should.not.exist(res.getHeader('Access-Control-Allow-Headers'));
          should.not.exist(res.getHeader('Access-Control-Allow-Credentials'));
          should.not.exist(res.getHeader('Access-Control-Max-Age'));
          done();
        };

        cors(options)(req, res, next);
      });


      it('can override methods', function (done) {
        // arrange
        var req, res, next, options;
        options = {
          methods: ['method1', 'method2']
        };
        req = fakeRequest();
        req.method = 'OPTIONS';
        res = fakeResponse();
        res.end = function () {
          // assert
          res.statusCode.should.equal(204);
          done();
        };
        next = function () {
          // assert
          res.getHeader('Access-Control-Allow-Methods').should.equal('method1,method2');
          done();
        };

        // act
        cors(options)(req, res, next);
      });

      it('methods defaults to GET, PUT, PATCH, POST, DELETE', function (done) {
        // arrange
        var req, res, next, options;
        options = {
        };
        req = fakeRequest();
        req.method = 'OPTIONS';
        res = fakeResponse();
        res.end = function () {
          // assert
          res.statusCode.should.equal(204);
          done();
        };
        next = function () {
          // assert
          res.getHeader('Access-Control-Allow-Methods').should.equal('GET,PUT,PATCH,POST,DELETE');
          done();
        };

        // act
        cors(options)(req, res, next);
      });

      it('can specify allowed headers', function (done) {
        // arrange
        var req, res, options;
        options = {
          allowedHeaders: ['header1', 'header2']
        };
        req = fakeRequest();
        req.method = 'OPTIONS';
        res = fakeResponse();
        res.end = function () {
          // assert
          res.getHeader('Access-Control-Allow-Headers').should.equal('header1,header2');
          should.not.exist(res.getHeader('Vary'));
          done();
        };

        // act
        cors(options)(req, res, null);
      });

      it('specifying an empty list or string of allowed headers will result in no response header for allowed headers', function (done) {
        // arrange
        var req, res, next, options;
        options = {
          allowedHeaders: []
        };
        req = fakeRequest();
        res = fakeResponse();
        next = function () {
          // assert
          should.not.exist(res.getHeader('Access-Control-Allow-Headers'));
          should.not.exist(res.getHeader('Vary'));
          done();
        };

        // act
        cors(options)(req, res, next);
      });

      it('if no allowed headers are specified, defaults to requested allowed headers', function (done) {
        // arrange
        var req, res, options;
        options = {
        };
        req = fakeRequest();
        req.method = 'OPTIONS';
        res = fakeResponse();
        res.end = function () {
          // assert
          res.getHeader('Access-Control-Allow-Headers').should.equal('requestedHeader1,requestedHeader2');
          should.exist(res.getHeader('Vary'));
          res.getHeader('Vary').should.equal('Access-Control-Request-Headers');
          done();
        };

        // act
        cors(options)(req, res, null);
      });

      it('can specify exposed headers', function (done) {
        // arrange
        var req, res, options, next;
        options = {
          exposedHeaders: ['custom-header1', 'custom-header2']
        };
        req = fakeRequest();
        res = fakeResponse();
        next = function () {
          // assert
          res.getHeader('Access-Control-Expose-Headers').should.equal('custom-header1,custom-header2');
          done();
        };

        // act
        cors(options)(req, res, next);
      });

      it('specifying an empty list or string of exposed headers will result in no response header for exposed headers', function (done) {
        // arrange
        var req, res, next, options;
        options = {
          exposedHeaders: []
        };
        req = fakeRequest();
        res = fakeResponse();
        next = function () {
          // assert
          should.not.exist(res.getHeader('Access-Control-Expose-Headers'));
          done();
        };

        // act
        cors(options)(req, res, next);
      });

      it('includes credentials if explicitly enabled', function (done) {
        // arrange
        var req, res, options;
        options = {
          credentials: true
        };
        req = fakeRequest();
        req.method = 'OPTIONS';
        res = fakeResponse();
        res.end = function () {
          // assert
          res.getHeader('Access-Control-Allow-Credentials').should.equal('true');
          done();
        };

        // act
        cors(options)(req, res, null);
      });

      it('does not includes credentials unless explicitly enabled', function (done) {
        // arrange
        var req, res, next, options;
        options = {
        };
        req = fakeRequest();
        res = fakeResponse();
        next = function () {
          // assert
          should.not.exist(res.getHeader('Access-Control-Allow-Credentials'));
          done();
        };

        // act
        cors(options)(req, res, next);
      });

      it('includes maxAge when specified', function (done) {
        // arrange
        var req, res, options;
        options = {
          maxAge: 456
        };
        req = fakeRequest();
        req.method = 'OPTIONS';
        res = fakeResponse();
        res.end = function () {
          // assert
          res.getHeader('Access-Control-Max-Age').should.equal('456');
          done();
        };

        // act
        cors(options)(req, res, null);
      });

      it('does not includes maxAge unless specified', function (done) {
        // arrange
        var req, res, next, options;
        options = {
        };
        req = fakeRequest();
        res = fakeResponse();
        next = function () {
          // assert
          should.not.exist(res.getHeader('Access-Control-Max-Age'));
          done();
        };

        // act
        cors(options)(req, res, next);
      });
    });

    describe('passing a function to build options', function () {
      it('handles options specified via callback', function (done) {
        // arrange
        var req, res, next, delegate;
        delegate = function (req2, cb) {
          cb(null, {
            origin: 'delegate.com'
          });
        };
        req = fakeRequest();
        res = fakeResponse();
        next = function () {
          // assert
          res.getHeader('Access-Control-Allow-Origin').should.equal('delegate.com');
          done();
        };

        // act
        cors(delegate)(req, res, next);
      });

      it('handles options specified via callback for preflight', function (done) {
        // arrange
        var req, res, delegate;
        delegate = function (req2, cb) {
          cb(null, {
            origin: 'delegate.com',
            maxAge: 1000
          });
        };
        req = fakeRequest();
        req.method = 'OPTIONS';
        res = fakeResponse();
        res.end = function () {
          // assert
          res.getHeader('Access-Control-Allow-Origin').should.equal('delegate.com');
          res.getHeader('Access-Control-Max-Age').should.equal('1000');
          done();
        };

        // act
        cors(delegate)(req, res, null);
      });

      it('handles error specified via callback', function (done) {
        // arrange
        var req, res, next, delegate;
        delegate = function (req2, cb) {
          cb('some error');
        };
        req = fakeRequest();
        res = fakeResponse();
        next = function (err) {
          // assert
          err.should.equal('some error');
          done();
        };

        // act
        cors(delegate)(req, res, next);
      });
    });
  });

}());

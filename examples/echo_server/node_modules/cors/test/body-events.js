(function () {
  /*global describe, it*/

  'use strict';

  var should = require('should'),
    express = require('express'),
    supertest = require('supertest'),
    bodyParser = require('body-parser'),
    cors = require('../lib');

  var dynamicOrigin,
    app1,
    app2,
    text = 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed justo turpis, tempor id sem fringilla, cursus tristique purus. Mauris a sollicitudin magna. Etiam dui lacus, vehicula non dictum at, cursus vitae libero. Curabitur lorem nulla, sollicitudin id enim ut, vehicula rhoncus felis. Ut nec iaculis velit. Vivamus at augue nulla. Fusce at molestie arcu. Duis at dui at tellus mattis tincidunt. Vestibulum sit amet dictum metus. Curabitur nec pretium ante. Proin vulputate elit ac lorem gravida, sit amet placerat lorem fringilla. Mauris fermentum, diam et volutpat auctor, ante enim imperdiet purus, sit amet tincidunt ipsum nulla nec est. Fusce id ipsum in sem malesuada laoreet vitae non magna. Praesent commodo turpis in nulla egestas, eu posuere magna venenatis. Integer in aliquam sem. Fusce quis lorem tincidunt eros rutrum lobortis.\n\nNam aliquam cursus ipsum, a hendrerit purus. Cras ultrices viverra nunc ac lacinia. Sed sed diam orci. Vestibulum ut orci a nibh scelerisque pretium. Sed suscipit vestibulum metus, ac ultricies leo sodales a. Aliquam erat volutpat. Vestibulum mauris massa, luctus et libero vel, cursus suscipit nulla. Cras sed erat quis massa fermentum congue. Mauris ultrices sem ligula, id malesuada lectus tincidunt eget. Donec sed nisl elit. Aenean ac lobortis massa. Phasellus felis nisl, dictum a dui volutpat, dictum sagittis diam. Vestibulum lacinia tellus vel commodo consequat.\n\nNulla at varius nibh, non posuere enim. Curabitur urna est, ultrices vel sem nec, consequat molestie nisi. Aliquam sed augue sit amet ante viverra pretium. Cras aliquam turpis vitae eros gravida egestas. Etiam quis dolor non quam suscipit iaculis. Sed euismod est libero, ac ullamcorper elit hendrerit vitae. Vivamus sollicitudin nulla dolor, vitae porta lacus suscipit ac.\n\nSed volutpat, magna in scelerisque dapibus, eros ante volutpat nisi, ac condimentum diam sem sed justo. Aenean justo risus, bibendum vitae blandit ac, mattis quis nunc. Quisque non felis nec justo auctor accumsan non id odio. Mauris vel dui feugiat dolor dapibus convallis in et neque. Phasellus fermentum sollicitudin tortor ac pretium. Proin tristique accumsan nulla eu venenatis. Cras porta lorem ac arcu accumsan pulvinar. Sed dignissim leo augue, a pretium ante viverra id. Phasellus blandit at purus a malesuada. Nam et cursus mauris. Vivamus accumsan augue laoreet lectus lacinia eleifend. Fusce sit amet felis nunc. Pellentesque eu turpis nisl.\n\nPellentesque vitae quam feugiat, volutpat lectus et, faucibus massa. Maecenas consectetur quis nisi eu aliquam. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Etiam laoreet condimentum laoreet. Praesent sit amet massa sit amet dui porta condimentum. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Sed volutpat massa nec risus malesuada hendrerit.';

  /* -------------------------------------------------------------------------- */

  dynamicOrigin = function (origin, cb) {
    setTimeout(function () {
      cb(null, true);
    }, 200);
  };

  /* -------------------------------------------------------------------------- */

  app1 = express();
  app1.use(cors({origin: dynamicOrigin}));
  app1.use(bodyParser.json());
  app1.post('/', function (req, res) {
    res.send(req.body);
  });

  /* -------------------------------------------------------------------------- */

  app2 = express();
  app2.use(bodyParser.json());
  app2.use(cors({origin: dynamicOrigin}));
  app2.post('/', function (req, res) {
    res.send(req.body);
  });

  /* -------------------------------------------------------------------------- */

  describe('body-parser-events', function () {
    describe('app1 (cors before bodyparser)', function () {
      it('POST works', function (done) {
        var body = {
          example: text
        };
        supertest(app1)
          .post('/')
          .send(body)
          .expect(200)
          .end(function (err, res) {
            should.not.exist(err);
            res.body.should.eql(body);
            done();
          });
      });
    });

    describe('app2 (bodyparser before cors)', function () {
      it('POST works', function (done) {
        var body = {
          example: text
        };
        supertest(app2)
          .post('/')
          .send(body)
          .expect(200)
          .end(function (err, res) {
            should.not.exist(err);
            res.body.should.eql(body);
            done();
          });
      });
    });
  });

}());

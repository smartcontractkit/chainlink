'use strict';
/*globals describe, it, beforeEach, afterEach */
var fsp = require('..'),
    path = require('path'),
    assert = require('assert'),
    Prom = require('any-promise'),
    testdir = path.join(__dirname, 'tmp'),
    testdir2 = path.join(__dirname, 'tmp2')


describe('basic', function(){
  beforeEach(function(){
    return fsp.mkdir(testdir).then(existstmp(true));
  });

  afterEach(function(){
    return fsp.remove(testdir).then(existstmp(false));
  });

  it('should create files and readdir', function(){
    return fsp.ensureFile(file('hello')).then(readtmp).then(function(files){
      assert.deepEqual(files.sort(), ['hello']);
      return fsp.ensureFile(file('world'));
    }).then(readtmp).then(function(files){
      assert.deepEqual(files.sort(), ['hello', 'world']);
      return fsp.exists(testdir2)
    }).then(function(exists){
      assert.equal(exists, false);
      return fsp.move(testdir, testdir2)
    }).then(function(){
      return Prom.all([fsp.exists(testdir), fsp.exists(testdir2)])
    }).then(function(exists){
      return assert.deepEqual(exists, [false, true])
    }).then(function(){
      return fsp.copy(testdir2, testdir)
    }).then(function(){
      return Prom.all([fsp.exists(testdir), fsp.exists(testdir2)])
    }).then(function(exists){
      return assert.deepEqual(exists, [true, true])
    }).then(readtmps).then(function(files){
      assert.deepEqual(files[0].sort(), files[1].sort());
      return fsp.emptyDir(testdir2);
    }).then(readtmp2).then(function(files){
      assert.deepEqual(files, []);
    }).then(function(){
      fsp.remove(testdir2);
    })
  });

  it('should pass through Sync as value', function(){
    return fsp.ensureFile(file('hello')).then(function(files){
      assert(fsp.existsSync(file('hello')));
      assert(!fsp.existsSync(file('world')));
      return fsp.ensureFile(file('world'));
    }).then(readtmp).then(function(files){
      assert(fsp.existsSync(file('hello')));
      assert(fsp.existsSync(file('world')));
    });
  });

  it('should copy with pipe read/write stream', function(){
    return fsp.writeFile(file('hello1'), 'hello world').then(function(){
      return fsp.readFile(file('hello1'), {encoding:'utf8'});
    }).then(function(contents){
      assert.equal(contents, 'hello world');
      var read = fsp.createReadStream(file('hello1')),
          write = fsp.createWriteStream(file('hello2')),
          promise = new Prom(function(resolve, reject){
            read.on('end', resolve);
            write.on('error', reject);
            read.on('error', reject);
          });
      read.pipe(write);
      return promise;
    }).then(function(){
      return fsp.readFile(file('hello2'), {encoding:'utf8'});
    }).then(function(contents){
      assert.equal(contents, 'hello world');
    });
  });

  it('should pass third argument from write #7', function testWriteFsp() {
    return fsp.open(file('some.txt'), 'w+').then(function (fd){
      return fsp.write(fd, "hello fs-promise").then(function(result) {
        var written = result[0];
        var text = result[1];
        assert.equal(text.substring(0, written), "hello fs-promise".substring(0, written))
        return fsp.close(fd);
      })
    })
  });

});

function file(){
  var args = [].slice.call(arguments);
  args.unshift(testdir);
  return path.join.apply(path, args);
}

function existstmp(shouldExist){
  return function(){
    return fsp.exists(testdir).then(function(exists){
        assert.equal(exists, shouldExist);
      });
  };
}

function readtmp(){
  return fsp.readdir(testdir);
}

function readtmp2(){
  return fsp.readdir(testdir2);
}

function readtmps(){
  return Prom.all([readtmp(), readtmp2()]);
}

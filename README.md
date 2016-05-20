# docker-path

[![Build Status](https://travis-ci.org/sh19910711/docker-path.svg?branch=master)](https://travis-ci.org/sh19910711/docker-path)

## Sketch of usage

```shell
$ cat example.yml
bundle:
  image: ruby:2.3.1
  command: bundle
npm:
  image: nodejs
  command: npm
$ docker-path example.yml
/tmp/path/to/mnt
$ ls /tmp/path/to/mnt
bundle npm
$ export PATH=$(docker-path example.yml):$PATH
$ bundle --version
Bundler version *.*.*
$ npm --version
*.*.*
```

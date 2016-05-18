# docker-memoize

## Sketch of usage

```shell
$ cat docker-memoize.yml
bundle:
  image: ruby-bundler
  command: bundle
npm:
  image: nodejs
  command: npm

$ docker-memoize config.yml
export PATH=$PATH:/tmp/path/to/mnt

$ eval $(docker-memoize config.yml)
$ bundle --version
Bundler version *.*.*
$ npm --version
*.*.*
```

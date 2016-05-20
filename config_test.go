package main

import (
	"testing"
)

var EXAMPLE_CONFIG = `
bundler:
  image: ruby-bundler
  command: bundle
  git: true
npm:
  image: nodejs
  command: npm
`

func TestParseImage(t *testing.T) {
	conf := Parse(EXAMPLE_CONFIG)
	if len(conf) != 2 {
		t.Fatalf("size should be 2, not %v", len(conf))
	}
	if conf["bundler"].Image != "ruby-bundler" {
		t.Fatalf("bundler.Image should be nodejs, not %v", conf["bundler"].Image)
	}
	if conf["npm"].Image != "nodejs" {
		t.Fatalf("npm.Image should be nodejs, not %v", conf["npm"].Image)
	}
}

func TestParseCommand(t *testing.T) {
	conf := Parse(EXAMPLE_CONFIG)
	if len(conf) != 2 {
		t.Fatalf("size should be 2, not %v", len(conf))
	}
	if conf["bundler"].Command != "bundle" {
		t.Fatalf("bundler.Command should be nodejs, not %v", conf["bundler"].Command)
	}
	if conf["npm"].Command != "npm" {
		t.Fatalf("npm.Command should be npm, not %v", conf["npm"].Command)
	}
}

func TestParseGit(t *testing.T) {
	conf := Parse(EXAMPLE_CONFIG)
	if len(conf) != 2 {
		t.Fatalf("size should be 2, not %v", len(conf))
	}
	if !conf["bundler"].Git {
		t.Fatalf("bundler.Git should be truethy, not %v", conf["bundler"].Git)
	}
	if conf["npm"].Git {
		t.Fatalf("npm.Command should be falesy, not %v", conf["npm"].Git)
	}
}

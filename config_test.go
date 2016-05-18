package main_test

import (
	"."
	"testing"
)

var EXAMPLE_CONFIG = `
bundler:
  image: ruby-bundler
  command: bundle
npm:
  image: nodejs
  command: npm
`

func TestParse(t *testing.T) {
	conf := main.Parse(EXAMPLE_CONFIG)
	if len(conf) != 2 {
		t.Fatalf("size should be 2, but %v", len(conf))
	}
	if conf["bundler"].Image != "ruby-bundler" {
		t.Fatalf("bundler.Image should be nodejs, but %v", conf["bundler"].Image)
	}
	if conf["bundler"].Command != "bundle" {
		t.Fatalf("bundler.Command should be nodejs, but %v", conf["bundler"].Command)
	}
	if conf["npm"].Image != "nodejs" {
		t.Fatalf("npm.Image should be nodejs, but %v", conf["npm"].Image)
	}
	if conf["npm"].Command != "npm" {
		t.Fatalf("npm.Command should be npm, but %v", conf["npm"].Command)
	}
}

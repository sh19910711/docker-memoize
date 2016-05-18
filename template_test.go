package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRenderImage(t *testing.T) {
	locals := Command{Command: "my-command", Image: "my-image"}
	script := Render(&locals)
	assert.Regexp(t, "my-image", script, "should contain image name")
}

func TestRenderCommand(t *testing.T) {
	locals := Command{Command: "my-command", Image: "my-image"}
	script := Render(&locals)
	assert.Regexp(t, "my-command", script, "should contain command")
}

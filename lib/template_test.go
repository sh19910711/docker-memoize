package lib_test

import (
	"."
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRender(t *testing.T) {
	locals := lib.Locals{Image: "my-image"}
	script := lib.Render(&locals)
	assert.Regexp(t, "my-image", script, "should contain image name")
}

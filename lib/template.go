package lib

import (
	"bytes"
	"text/template"
)

const TEMPLATE_SCRIPT string = `
#!/bin/bash
# image = {{.Image}}

DOCKER_BIN=$(which docker)

msg() {
  echo $@ >&2
}

__sudo() {
  if sudo -v >/dev/null 2>&1; then
    sudo "$@"
  else
    "$@"
  fi
}

docker() {
  __sudo ${DOCKER_BIN} "$@"
}

container_names() {
  docker ps --format {{"{{.Names}}"}}
}

container_exist() {
  container_names | grep $1 >/dev/null 2>&1
}

command() {
  local name=docker_memoize_$(echo $PWD | sed -e 's;[^a-zA-Z0-9_.-];_;g' | cut -c 2-)
  if container_exist $name; then
    : nothing to do
  else
    msg create $name
    docker run -d \
      -u "$UID:$GROUPS" \
      -e "HOME=/home/user" \
      -v "$HOME:/home/user" \
      -w "/wrk" \
      -v "$PWD:/wrk" \
      --name $name \
      '{{.Image}}' /bin/sh -c 'while true; do echo hello; sleep 1000; done'
  fi
  docker exec -i $name '{{.Command}}' "$@"
}

command $@
`

func Render(locals *Command) string {
	tmpl, err := template.New("script").Parse(TEMPLATE_SCRIPT)
	if err != nil {
		panic(err)
	}
	var doc bytes.Buffer
	if err := tmpl.Execute(&doc, locals); err != nil {
		panic(err)
	}
	return doc.String()
}

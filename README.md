An espncricinfo wrapper written in go to track scores in real time and in the cmd/cli you can find the Command Line Interface wrapped over this wrapper.

### CLI Demo
[![asciicast](https://asciinema.org/a/PoSKzpzArg6XdSi62X5OmvyHM.svg)](https://asciinema.org/a/PoSKzpzArg6XdSi62X5OmvyHM?speed=2)

Choose the match that you want to follow by the match number

Command line flags available are ```refresh```. The default ```refresh``` time is 1 second.

Use makefile to build CLI for your desktop.
```sh
make
```

This will return the available commands for the makefile. The help is self explanatory.
```
Usage:
  make [target...]

Useful commands:
  build                          build the binary(binary name- score) in the current working directory
  move                           move to /usr/bin so that you can use this binary anywhere.
  run                            run the CLI
```

Once you move the binary in your path so that it can be used from anywhere, you can use it like ```score -refresh 1 ```. The CLI is in very initial stage and has witnessed just a few hours of code after the idea.


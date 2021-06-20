An espncricinfo wrapper written in go to track scores in real time and in the cmd/cli you can find the Command Line Interface wrapped over this wrapper.

### CLI Demo
[![asciicast](https://asciinema.org/a/AHvD2sVXYy1FblQrS1WrpeqhR.svg)](https://asciinema.org/a/AHvD2sVXYy1FblQrS1WrpeqhR?speed=2)

Right now the CLI is not listing all the matches to choose from, but you can choose to track a particular match by entering a specific matchid by using command line flags.

Flags available are ```matchid``` and ```refresh```. The ```matchid``` flag is the espncricinfo description for a specific match and you can find it in the match url. The default ```matchid``` is the wtc final which is between NZ and India matchid. The default ```refresh``` time is 1 second.

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

Once you move the binary in your path so that it can be used from anywhere, you can use it like ```score -refresh 1 -matchid 1249875```. The CLI is in very initial stage and has witnessed just a couple of hours of code after the idea. 


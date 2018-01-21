# DND-SHUB

A nifty little bot for rolling dice, managing players, etc. written by jamesbyrnes in Go.

## Installation

Download the repo to your $GOPATH by using:

```sh
go get github.com/jamesbyrnes/dnd-shub
```

As this bot uses [godotenv](https://github.com/joho/godotenv), you will need to set up a Discord bot token (see [here](https://discordapp.com/developers/docs/topics/oauth2#bots)) in a file called ```.env```. Inside the file, you only need to put the following:

```
DISC_TOKEN=YOUR-BOT-TOKEN-HERE
```

If you want to run the provided tests, for the time being, you will have to symlink this file into the cmd/ directory as well. I will be moving this to a config-based solution shortly, so this won't be a thing for very long.

Build the executable using:
```sh
go build github.com/jamesbyrnes/dnd-shub
```

## Starting the bot
The bare command will start the server. Terminate the command with CTRL-C.

### Flags
* ```-p PREFIX```: Restrict DND-SHUB to only responding in commands beginning with the $PREFIX. For example, if you run the bot with ```-p bot``` it will respond to commands in #botchannel but not #general.

## Using the bot
The following commands are available through DND-SHUB:
* ```/dice``` or ```/d``` is a dice roller. It accepts a standard DND dice rolling string (e.g. 2d20 for two twenty-sided dice), plus or minus a modifier (optional) and returns a formatted string with your result. You can use multiple sets of dice separated by a space. For example:
 * ```/d 2d20+5``` will roll two twenty-sided dice and add 5 to the total of that roll.
 * ```/d 1d6 2d10``` will roll one six-sided die and two ten-sided dice.
 * ```/d 4d2 1d20-6``` will roll four two-sided dice (i.e. coins) and one twenty-sided die, subtracting 6 from the roll of the twenty-sided die only.

Player management coming soon!
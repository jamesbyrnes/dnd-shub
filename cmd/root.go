package cmd

import (
    "os"
    "os/signal"
    "syscall"
    "log"
    "fmt"
    "strings"

    "github.com/joho/godotenv"
    "github.com/bwmarrin/discordgo"
    "github.com/spf13/cobra"

    "github.com/jamesbyrnes/dnd-shub/utils"
)

var RootCmd = &cobra.Command {
    Use:    "dnd-shub",
    Short:  "DND-SHUB is a Dungeons and Dragons bot for Discord",
    Long:   `A nifty little bot for rolling dice, managing players,
            etc written by jamesbyrnes in Go.`,
    Run: func(cmd *cobra.Command, args []string) {
        runBot()
    },
}

var ChanPrefix string

func init() {
    RootCmd.Flags().StringVarP(&ChanPrefix, "prefix", "p", "", "restrict DND-SHUB to only work in channels with a given prefix")
}

// chanPrefixMatches checks to see if the variable ChanPrefix matches the 'left' side of the 
// provided chanName (provided as discordgo.Session.State.Channel.Name below in runBot()).
func chanPrefixMatches(chanName string) bool {
    if len(ChanPrefix) > 0 {
        if len(ChanPrefix) > len(chanName) { 
            return false
        } else if ChanPrefix != chanName[:len(ChanPrefix)] {
            return false
        }
    }
    return true
}

// EnvVar uses a .env file in the application root to get the API token for the 
// Discort bot.
type EnvVar struct {
    Token string
}

// EnvVar.Get will load the .env file and its value for the API token ("DISC_TOKEN").
func (e *EnvVar) Get() error {
    // Heads up, running 'go test' on this only works if the .env file is in the
    // same directory as the PWD of the code being run. However, if you move the
    // .env file to the cmd/ subdir, 'go run' won't work because then .env is 
    // expected to be in the root... quick and dirty solution is to make a 
    // symlink from root to cmd/... github issue seems to indicate this is not 
    // resolvable, may be best to use a config file instead...
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Failed to load environment variables file.")
        return err
    }

    e.Token = os.Getenv("DISC_TOKEN")
    return nil
}

// getNewDGSession will set up DiscordGo with the token we provided as an EnvVar struct.
func getNewDGSession(ev *EnvVar) (discordgo.Session, error) {
    dg, err := discordgo.New("Bot " + ev.Token)
    if err != nil {
        log.Fatal("Failed to initialize DiscordGo session.")
        return *dg, err
    }

    return *dg, nil
}

// messageHandler does the bulk of the work in routing the commands to the bot and ensuring
// that, e.g. the ChanPrefix matches the channel from which the message is sent.
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID {
        return
    }

    channel, err := s.State.Channel(m.ChannelID)
    if err != nil {
        log.Fatal("Failed to get channel information")
    }

    if !chanPrefixMatches(channel.Name) {return}

    // The command /dice or /d will route to dice.go and get a random roll of the dice
    // based on a parsed 'traditional' Dungeons and Dragons-style die string (e.g. 2d20).
    if m.Content[:3] == "/d " || m.Content[:6] == "/dice " {
        diceStrings := strings.SplitAfterN(m.Content, " ", 2)[1]
        var res string
        res = "Rolling the dice!"

        for i, v := range strings.Split(diceStrings, " ") {
            d := new(utils.Dice)
            err := d.Build(v)
            if len(res) > 0 {
                res += "\n"
            }

            res += fmt.Sprintf("Set #%d (%s): ", i+1, v)

            if err != nil {
                res += "Format incorrect - should be XdY(+/-Z)"
                break
            }

            str, err := d.GetFullString()
            if err != nil {
                res += "Dice string turned out invalid"
                break
            }

            res += str
        }
        s.ChannelMessageSend(m.ChannelID, res)
    }
}

// runBot is where Cobra gets routed in order to start the initialization of DiscordGo.
func runBot() {
    envVar := new(EnvVar)
    err := envVar.Get()

    if err != nil {
        os.Exit(1)
    }

    dg, err := getNewDGSession(envVar)
    if err != nil {
        os.Exit(1)
    }

    dg.AddHandler(messageHandler)

    err = dg.Open()
    if err != nil {
        log.Fatal("Failed to open a Discord session.")
        os.Exit(1)
    }
    defer dg.Close()

    // This error check is required since dg.Open() will happily not give an error as long as 
    // it can connect to Discord, but it won't tell you if, e.g., the token doesn't work. We need to run
    // a "throwaway" command to make sure everything has been set up properly.
    _, err = dg.User("@me")
    if err != nil {
        log.Fatal("Failed to connect to Discord.")
        os.Exit(1)
    }

    fmt.Println("DND-SHUB is now running. Press CTRL-C to exit.")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc
}

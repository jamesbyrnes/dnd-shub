package cmd

import (
    "testing"
)

func TestEnvVarNotError (t *testing.T) {
    e := new(EnvVar)
    err := e.Get()

    if err != nil {
        t.Error("Tried to load .env file but failed.")
    }
}

func TestEnvVarNotBlank (t *testing.T) {
    e := new(EnvVar)
    _ = e.Get()

    if e.Token == "" {
        t.Error("Loaded .env file but desired keys are blank")
    }
}

func TestDGSessionInitializes (t *testing.T) {
    e := new(EnvVar)
    _ = e.Get()

    _, err := getNewDGSession(e)
    if err != nil {
        t.Error("Failed to initialize a DiscordGo session")
    }
}

func TestSGSessionConnects (t *testing.T) {
    e := new(EnvVar)
    _ = e.Get()

    dg, _ := getNewDGSession(e)

    dg.AddHandler(messageHandler)

    err := dg.Open()
    if err != nil {
        t.Error("Failed to open Discord session")
    }
    _, err = dg.User("@me")
    if err != nil {
        t.Error("Failed to connect to Discord")
    }
    dg.Close()
}

func TestChanPrefixMatches (t *testing.T) {
    ChanPrefix = "test"
    chanMatch := chanPrefixMatches("testchan")
    chanMismatch := chanPrefixMatches("xchan")
    chanBlank := chanPrefixMatches("")

    if !chanMatch {
        t.Error("Expected channel prefix match; it failed!")
    } else if chanMismatch || chanBlank {
        t.Error("Expected channel prefix mismatch; it failed!")
    }
}

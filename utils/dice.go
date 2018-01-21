package utils

import (
    "math/rand"
    "time"
    "errors"
    "strings"
    "strconv"
    "fmt"
)

// Dice represents a single set of dice with the same number of sides, plus (or minus) a modifier
// which should be combined with the total value of the dice.
type Dice struct {
    NumSides int
    Modifier int
    Values []int
}

// Build parses a 'traditional' D&D dice string (i.e. 2d20+1 means 'roll two twenty-sided dice and add one to the total)
// and turns it into a Dice struct. Returns an error if there is one during the parsing process, else returns
// nil.
func (d *Dice) Build(sz string) error {
    var err error
    var numDice int
    var modType string

    d.Modifier = 0

    if !strings.Contains(sz, "d") ||
       strings.Count(sz, "d") > 1 ||
       strings.Count(sz, "+") > 1 ||
       strings.Count(sz, "-") > 1 {
        return errors.New(fmt.Sprintf("Error: Incorrect argument provided (should be XdY(+/-Z)"))
    }

    if strings.Index(sz, "+") != -1 && strings.Index(sz, "-") != -1 {
        return errors.New("Error: Both positive and negative modifier provided")
    } else if strings.Index(sz, "-") != -1 {
        modType = "-"
    } else {
        modType = "+"
    }

    modSplit := strings.Split(sz, modType)

    if len(modSplit) > 1 {
        d.Modifier, err = strconv.Atoi(modSplit[1])
        if err != nil {
            return err
        }
        if modType == "-" {
            d.Modifier *= -1
        }
        sz = modSplit[0]
    }

    dSplit := strings.Split(sz, "d")

    numDice, err = strconv.Atoi(dSplit[0])
    if err != nil {
        return err
    }

    d.NumSides, err = strconv.Atoi(dSplit[1])
    if err != nil {
        return err
    }

    rando := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < numDice; i ++ {
        d.Values = append(d.Values, rando.Int() % d.NumSides + 1)
    }

    return nil
}

// Roll is a programmatic way to roll a new Dice object, without having to provide
// a Dice string.
func (d *Dice) Roll(numDice int, numSides int, mod int) {
    d.NumSides = numSides
    d.Modifier = mod

    rando := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < numDice; i ++ {
        d.Values = append(d.Values, rando.Int() % d.NumSides + 1)
    }
}

// Sum returns the sum of all of the dice in the set, as long as the dice are initialized.
func (d *Dice) Sum(includeMod bool) (int, error) {
    if len(d.Values) == 0 {
        return 0, errors.New("Error: Dice set not initialized")
    } 

    var sum int
    
    if len(d.Values) == 1 {
        sum = d.Values[0]
    } else {
        for _, v := range d.Values {
            sum += v
        }
    }

    if includeMod {
        return (sum + d.Modifier), nil
    }

    return sum, nil
}

// NumDice returns the number of dice in the set.
func (d *Dice) NumDice() int {
    return len(d.Values)
}

// GetDiceString returns a comma-separated list of dice values as a string.
func (d *Dice) GetDiceString() (string, error) {
    if len(d.Values) == 0 {
        return "", errors.New(fmt.Sprintf("Error: Dice set not initialized"))
    } else if len(d.Values) == 1 {
        return strconv.Itoa(d.Values[0]), nil
    }

    var res string

    for _, v := range d.Values {
        if len(res) > 0 {
            res += ", "
        }
        res += strconv.Itoa(v)
    }

    return res, nil
}

// GetFullString returns the string from GetDiceString, along with the pre-modded sum, modifier and sum+mod
func (d *Dice) GetFullString() (string, error) {
    var modType string
    var modString string
    var err error

    modString = ""

    if d.Modifier < 0 {
        modType = "-"
    } else {
        modType = "+"
    }

    var sum int
    sum, err = d.Sum(false)
    if err != nil {
        return "", err
    }

    //TODO - this seems silly - better to just do sum + d.Mod?
    var sumMod int
    sumMod, err = d.Sum(true)
    if err != nil {
        return "", err
    }

    if d.Modifier != 0 {
        modString = fmt.Sprintf(" %s %d => %d", modType, d.Modifier, sumMod)
    }

    var dString string
    dString, err = d.GetDiceString()
    if err != nil {
        return "", err
    }

    resString := fmt.Sprintf("%s = %d%s", dString, sum, modString)

    return resString, nil
}
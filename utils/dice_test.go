package utils

import (
    "testing"
    "strings"
)

func TestDiceCorrectSides(t *testing.T) {
    d := new(Dice)
    d.Build("1d20")
    if d.Values[0] > 20 || d.Values[0] < 1 {
        t.Errorf("Expected a value for a %d sided die but got %d instead.", d.NumSides, d.Values[0])
    }
}

func TestInvalidDiceFormatD(t *testing.T) {
    d := new(Dice)
    e := d.Build("1dd20")
    if e == nil {
        t.Error("Expected error for invalid dice format: too many 'd' chars")
    }
}

func TestInvalidDiceFormatPlus(t *testing.T) {
    d := new(Dice)
    e := d.Build("1d20++1")
    if e == nil {
        t.Error("Expected error for invalid dice format: too many '+' chars")
    }
}

func TestInvalidDiceFormatMinus(t *testing.T) {
    d := new(Dice)
    e := d.Build("1d20--1")
    if e == nil {
        t.Error("Expected error for invalid dice format: too many '-' chars")
    }
}

func TestInvalidDiceFormatPlusMinus(t *testing.T) {
    d := new(Dice)
    e := d.Build("1d20+-1")
    if e == nil {
        t.Error("Expected error for invalid dice format: both pos and neg")
    }
}

func TestUninitializedDiceError(t *testing.T) {
    d := new(Dice)
    _, e := d.GetDiceString()
    if e == nil {
        t.Error("Expected error: dice set not initialized")
    }
}

func TestNumDiceOneToOne(t *testing.T) {
    d := new(Dice)
    d.Build("1d20")
    s, _ := d.GetDiceString()
    l := len(strings.Split(s, " "))
    if l > 1 {
       t.Errorf("Expected 1 die; got %d instead", l)
    }
}

func TestNumDiceManyToMany(t *testing.T) {
    d := new(Dice)
    d.Build("5d20")
    s, _ := d.GetDiceString()
    l := len(strings.Split(s, " "))
    if l != 5 {
        t.Errorf("Expected 5 dice; got %d instead", l)
    }
}

func TestSumIsCorrect(t *testing.T) {
    d := new(Dice)
    d.Build("5d1")
    s, _ := d.Sum()
    if s != 5 {
        t.Errorf("Expected sum of 5; got %d instead", s)
    }
}

func TestNegModIsNegative(t *testing.T) {
    d := new(Dice)
    d.Build("1d20-5")
    if d.Modifier > 0 {
        t.Error("Negative modifier is not negative")
    }
}

func TestModifier(t *testing.T) {
    d := new(Dice)
    d.Build("1d1+5")
    sm, err := d.SumMod()
    if err != nil {
        t.Errorf("Sum + mod gave an error %s", err)
    }
    if sm != 6 {
        t.Errorf("Expected sum + mod to give 6 but got %d", sm)
    }
}

package commands

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func roll(c *Command, s *discordgo.Session, m *discordgo.MessageCreate) error {
	var (
		args = c.Args()
	)
	roller := m.Author.DisplayName()

	if len(args) < 1 {
		MessageSend(s, m, rollDuality(roller))
		return nil
	}

	var results []string

	for _, roll := range args {

		if strings.ToLower(roll) == "duality" || strings.ToLower(roll) == "duelity" {
			results = append(results, rollDuality(roller))
			continue
		}

		isMultiDiceRoll, _ := regexp.MatchString("^[0-9]*d[0-9]+.*$", roll)

		diceNum, err := strconv.ParseFloat(roll, 64)
		if err != nil && !isMultiDiceRoll {
			results = append(results, fmt.Sprintf("%s is not a valid roll. A roll is a number, duality, or dice abbreviation", roll))
			continue
		}
		if isMultiDiceRoll {
			results = append(results, rollMultiDice(roll, roller))
		} else {
			_, diceRollResultString := rollDice(diceNum)

			if diceRollResultString == "1" {
				results = append(results, fmt.Sprintf("%s's d%s result is %s :cry:\n", roller, roll, diceRollResultString))
				continue
			}

			results = append(results, fmt.Sprintf("%s's d%s result is %s\n", roller, roll, diceRollResultString))
		}
	}

	MessageSend(s, m, strings.Join(results, "\n"))
	return nil

}

func init() {
	RegisterCommand(NewCommand("Roll", "Replies with Roll!", roll))
}

func rollDuality(roller string) string {

	hope, hopeString := rollDice(12)
	fear, fearString := rollDice(12)

	result := hope + fear
	resultString := strconv.FormatFloat(result, 'f', -1, 64)

	var dualityResult string

	if hope > fear {
		dualityResult = fmt.Sprintf("%s rolled with Hope :heart:", roller)
	} else if fear > hope {
		dualityResult = fmt.Sprintf("%s rolled with Fear :dagger:", roller)
	} else {
		dualityResult = fmt.Sprintf("# %s CRIT!!!! :dagger: :heart:", strings.ToUpper(roller))
	}

	return fmt.Sprintf("> %s \n> The Hope roll was %s \n> The Fear roll was %s \n> The total was %s\n", dualityResult, hopeString, fearString, resultString)

}

func rollDice(diceSides float64) (float64, string) {
	diceRollResult := math.Ceil(rand.Float64() * diceSides)
	diceRollResultString := strconv.FormatFloat(diceRollResult, 'f', -1, 64)
	return diceRollResult, diceRollResultString
}

func rollMultiDice(diceDesignation string, roller string) string {
	var positiveModifier float64 = 0
	var negativeModifier float64 = 0
	var diceTotal float64 = 0
	var diceDesignationString = diceDesignation

	if strings.Contains(diceDesignation, "+") {
		modifierSplitStrings := strings.Split(diceDesignation, "+")
		positiveModifier, _ = strconv.ParseFloat(modifierSplitStrings[1], 64)
		diceDesignation = modifierSplitStrings[0]
	} else if strings.Contains(diceDesignation, "-") {
		modifierSplitStrings := strings.Split(diceDesignation, "-")
		negativeModifier, _ = strconv.ParseFloat(modifierSplitStrings[1], 64)
		diceDesignation = modifierSplitStrings[0]
	}

	splitString := strings.Split(diceDesignation, "d")
	numRolls, err := strconv.Atoi(splitString[0])
	if err != nil {
		numRolls = 1
	}
	rollValue, err := strconv.ParseFloat(splitString[1], 64)
	if err != nil {
		return fmt.Sprintf("%s is not a valid roll. A roll is a number, duality, or dice abbreviation", diceDesignationString)
	}
	for i := 0; i < numRolls; i++ {
		dicevalue, _ := rollDice(rollValue)
		diceTotal += dicevalue
	}

	diceTotal += (positiveModifier - negativeModifier)
	diceTotalString := strconv.FormatFloat(diceTotal, 'f', -1, 64)

	return fmt.Sprintf("%s's %s result is %s\n", roller, diceDesignationString, diceTotalString)
}

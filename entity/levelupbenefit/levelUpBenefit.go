package levelupbenefit

import (
	"errors"
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
)

// Type defines the expected sources the LevelUpBenefit could be conjured from.
type Type string

const (
	// Small is for small improvements (stats mostly)
	Small Type = "small"
	// Big is for substantial changes to character (new powers, movement changes)
	Big = "big"
)

// Identification is used to uniquely identify this LevelUpBenefit.
type Identification struct {
	ClassID            string             `json:"class_id" yaml:"class_id"`
	LevelUpBenefitType Type               `json:"level_up_benefit_type" yaml:"level_up_benefit_type"`
	ID                 string             `json:"id" yaml:"id"`
}

// Defense describes all of the defensive benefits of leveling up.
type Defense struct {
	MaxHitPoints       int                `json:"max_hit_points" yaml:"max_hit_points"`
	Dodge              int                `json:"dodge" yaml:"dodge"`
	Deflect            int                `json:"deflect" yaml:"deflect"`
	MaxBarrier         int                `json:"max_barrier" yaml:"max_barrier"`
	Armor              int                `json:"armor" yaml:"armor"`
}

// Offense describes all of the offensive benefits of leveling up.
type Offense struct {
	Aim                int                `json:"aim" yaml:"aim"`
	Strength           int                `json:"strength" yaml:"strength"`
	Mind               int                `json:"mind" yaml:"mind"`
}

// LevelUpBenefit describes how a Squaddie improves upon levelling up.
type LevelUpBenefit struct {
	Identification				*Identification              `json:"identification" yaml:"identification"`
	Defense						*Defense						`json:"defense" yaml:"defense"`
	Offense						*Offense						`json:"offense" yaml:"offense"`
	//PowerCollection				PowerCollection				`json:"powers" yaml:"powers"`
	PowerIDGained      []*power.Reference `json:"powers" yaml:"powers"`
	PowerIDLost        []*power.Reference `json:"powers_lost" yaml:"powers_lost"`
	Movement           *squaddie.Movement `json:"Movement" yaml:"Movement"`
}

// CheckForErrors ensures the LevelUpBenefit has valid fields
func (benefit *LevelUpBenefit) CheckForErrors() error {
	if benefit.Identification.LevelUpBenefitType != Small && benefit.Identification.LevelUpBenefitType != Big {
		return fmt.Errorf(`unknown level up benefit type: "%s"`, benefit.Identification.LevelUpBenefitType)
	}

	if benefit.Identification.ClassID == "" {
		return errors.New(`no classID found for LevelUpBenefit`)
	}
	return nil
}

// FilterLevelUpBenefits filters a slice of LevelUpBenefits.
func FilterLevelUpBenefits(sliceToFilter []*LevelUpBenefit, condition func(benefit *LevelUpBenefit) bool) []*LevelUpBenefit {
	keptLevelUpBenefits := []*LevelUpBenefit{}
	for _, benefitToTest := range sliceToFilter {
		if condition(benefitToTest) {
			keptLevelUpBenefits = append(keptLevelUpBenefits, benefitToTest)
		}
	}
	return keptLevelUpBenefits
}

// AnyLevelUpBenefits returns true if at least LevelUpBenefit in the sliceToAnalyze satisfies the condition.
func AnyLevelUpBenefits(sliceToAnalyze []*LevelUpBenefit, condition func(benefit *LevelUpBenefit) bool) bool {
	for _, benefitToTest := range sliceToAnalyze {
		if condition(benefitToTest) {
			return true
		}
	}
	return false
}

// CountLevelUpBenefits returns the number of LevelUpBenefit that satisfy the given condition.
func CountLevelUpBenefits(sliceToAnalyze []*LevelUpBenefit, condition func(benefit *LevelUpBenefit) bool) int {
	count := 0
	for _, benefitToTest := range sliceToAnalyze {
		if condition(benefitToTest) {
			count = count + 1
		}
	}
	return count
}

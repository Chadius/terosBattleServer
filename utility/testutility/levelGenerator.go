package testutility

import (
	"github.com/cserrant/terosBattleServer/entity/levelupbenefit"
	"github.com/cserrant/terosBattleServer/utility"
	"strconv"
)

// LevelGeneratorInstruction tells the LevelGenerator how to build the levels.
type LevelGeneratorInstruction struct {
	NumberOfLevels int
	ClassID        string
	PrefixLevelID  string
	Type           levelupbenefit.Type
}

// LevelGenerator is a Builder pattern that can be programmed to generate a slice of LevelUpBenefits.
type LevelGenerator struct {
	Instructions *LevelGeneratorInstruction
}

// NewLevelGenerator returns a pointer to LevelGenerator following the given instructions.
func NewLevelGenerator(instructions *LevelGeneratorInstruction) *LevelGenerator {
	if instructions == nil {
		instructions = &LevelGeneratorInstruction{
			NumberOfLevels: 0,
			PrefixLevelID:  "defaultGeneratorLevel" + utility.StringWithCharset(8, "abcdefg0123456789"),
			Type:           levelupbenefit.Small,
			ClassID:        "defaultGeneratorClass" + utility.StringWithCharset(8, "abcdefg0123456789"),
		}
	}

	return &LevelGenerator{
		Instructions: instructions,
	}
}

// Build follows the instructions to create a slice of levelUpBenefits.
func (generator *LevelGenerator) Build() []*levelupbenefit.LevelUpBenefit {
	levels := []*levelupbenefit.LevelUpBenefit{}
	for i:= 0; i < generator.Instructions.NumberOfLevels; i++ {
		newLevelUpBenefit := &levelupbenefit.LevelUpBenefit{
			Identification: &levelupbenefit.Identification{
				LevelUpBenefitType: generator.Instructions.Type,
				ClassID:            generator.Instructions.ClassID,
				ID:                 generator.Instructions.PrefixLevelID + strconv.Itoa(i),
			},
		}
		levels = append(levels, newLevelUpBenefit)
	}
	return levels
}
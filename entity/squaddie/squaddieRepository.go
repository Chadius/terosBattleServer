package squaddie

import (
	"encoding/json"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/utility"
	"gopkg.in/yaml.v2"
)

// Repository will interact with external devices to manage Squaddies.
type Repository struct {
	squaddiesByName map[string]*Squaddie
}

// NewSquaddieRepository generates a pointer to a new Repository.
func NewSquaddieRepository() *Repository {
	repository := Repository{
		map[string]*Squaddie{},
	}
	return &repository
}

// AddJSONSource consumes a given bytestream and tries to analyze it.
func (repository *Repository) AddJSONSource(data []byte) (bool, error) {
	return repository.addSource(data, json.Unmarshal)
}

// AddYAMLSource consumes a given bytestream and tries to analyze it.
func (repository *Repository) AddYAMLSource(data []byte) (bool, error) {
	return repository.addSource(data, yaml.Unmarshal)
}

// AddSquaddies adds a slice of Squaddie to the repository.
func (repository *Repository) AddSquaddies(squaddies []*Squaddie) (bool, error) {
	for _, squaddieToAdd := range squaddies {
		_, err := repository.tryToAddSquaddie(squaddieToAdd)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// AddSource consumes a given bytestream of the given sourceType and tries to analyze it.
func (repository *Repository) addSource(data []byte, unmarshal utility.UnmarshalFunc) (bool, error) {
	var unmarshalError error
	var listOfSquaddies []Squaddie
	unmarshalError = unmarshal(data, &listOfSquaddies)

	if unmarshalError != nil {
		return false, unmarshalError
	}

	for index := range listOfSquaddies {
		success, err := repository.tryToAddSquaddie(&listOfSquaddies[index])
		if success == false {
			return false, err
		}
	}
	return true, nil
}

func (repository *Repository) tryToAddSquaddie(squaddieToAdd *Squaddie) (bool, error) {
	squaddieErr := CheckSquaddieForErrors(squaddieToAdd)
	if squaddieErr != nil {
		return false, squaddieErr
	}
	squaddieToAdd.SetHPToMax()
	repository.squaddiesByName[squaddieToAdd.Name] = squaddieToAdd
	return true, nil
}

// GetNumberOfSquaddies returns the number of Squaddies ready to retrieve.
func (repository *Repository) GetNumberOfSquaddies() int {
	return len(repository.squaddiesByName)
}

// GetByName retrieves a Squaddie by name
func (repository *Repository) GetByName(squaddieName string) *Squaddie {
	squaddie, squaddieExists := repository.squaddiesByName[squaddieName]
	if !squaddieExists {
		return nil
	}

	clonedSquaddie, cloneErr := repository.CloneSquaddie(squaddie, "")
	if cloneErr != nil {
		return nil
	}
	return clonedSquaddie
}

// MarshalSquaddieIntoJSON converts the given Squaddie into JSON.
func (repository *Repository) MarshalSquaddieIntoJSON(squaddie *Squaddie) ([]byte, error) {
	type Alias Squaddie

	return json.Marshal(&struct {
		*Alias
		PowerIDNames []*power.Reference `json:"powers" yaml:"powers"`
	}{
		Alias:        (*Alias)(squaddie),
		PowerIDNames: squaddie.GetInnatePowerIDNames(),
	})
}

//CloneSquaddie uses the base Squaddie to create a new one.
//  All fields will be the same except the ID.
//  If newID isn't empty, the clone ID is set to that.
//  Otherwise it is randomly generated.
func (repository *Repository) CloneSquaddie(base *Squaddie, newID string) (*Squaddie, error) {
	clone := NewSquaddie(base.Name)
	clone.Affiliation = base.Affiliation
	if newID != "" {
		clone.ID = newID
	}

	clone.CurrentHitPoints = base.CurrentHitPoints
	clone.MaxHitPoints = base.MaxHitPoints
	clone.Aim = base.Aim
	clone.Strength = base.Strength
	clone.Mind = base.Mind
	clone.Dodge = base.Dodge
	clone.Deflect = base.Deflect
	clone.CurrentBarrier = base.CurrentBarrier
	clone.MaxBarrier = base.MaxBarrier
	clone.Armor = base.Armor

	clone.Movement.Distance = base.Movement.Distance
	clone.Movement.Type = base.Movement.Type
	clone.Movement.HitAndRun = base.Movement.HitAndRun

	clone.PowerReferences = append([]*power.Reference{}, base.PowerReferences...)

	clone.BaseClassID = base.BaseClassID
	clone.CurrentClass = base.CurrentClass
	for classID, progress := range base.ClassLevelsConsumed {
		newProgress := ClassProgress{
			ClassID:        classID,
			ClassName:      progress.ClassName,
			LevelsConsumed: append([]string{}, progress.LevelsConsumed...),
		}

		clone.ClassLevelsConsumed[classID] = &newProgress
	}

	return clone, nil
}

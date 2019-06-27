package dlshow

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// TVEpisodeInfo contains information about an individual
// TV show episode
type TVEpisodeInfo struct {
	ShowName       string
	SeasonNumber   int
	EpisodeNumber  int
	EpisodeTitle   string
	EpisodeSummary string

	AiredYear  int
	AiredMonth int
	AiredDay   int

	ParseType int
}

const (
	// ParseTypeUnknown represents and unknown parse type
	ParseTypeUnknown = iota

	// ParseTypeSE represents a season/episode parse type (original)
	ParseTypeSE = iota

	// ParseTypeDate represents a date parse type
	ParseTypeDate = iota

	// ParseTypeSE2 represents a season/episode parse type (alternate)
	ParseTypeSE2 = iota
)

// GetEpisodeInfo returns TV show information for a given
// downloaded filename.
func GetEpisodeInfo(filename string) (TVEpisodeInfo, error) {
	retval := TVEpisodeInfo{}

	//	Make sure we have just a filename -- not an entire path, not a directory:
	_, filename = filepath.Split(filename)

	if strings.TrimSpace(filename) == "" {
		return retval, fmt.Errorf("Filename does not appear to be a valid filename")
	}

	/******* PARSE THE FILENAME ********/

	//	Season / episode parsers
	rxSE := regexp.MustCompile(`(?i)^((?P<series_name>.+?)[. _-]+)?s(?P<season_num>\d+)[. _-]*e(?P<ep_num>\d+)(([. _-]*e|-)(?P<extra_ep_num>(!(1080|720)[pi])\d+))*[. _-]*((?P<extra_info>.+?)((![. _-])-(?P<release_group>[^-]+))?)?$`)
	rxSE2 := regexp.MustCompile(`(?P<series_name>.*?)\.S?(?P<season_num>\d{1,2})[Ex-]?(?P<ep_num>\d{2})\.(.*)`)

	//	Airdate parsers:
	rxD := regexp.MustCompile(`^((?P<series_name>.+?)[. _-]+)(?P<year>\d{4}).(?P<month>\d{1,2}).(?P<day>\d{1,2})`)

	//	Show name formatter
	rxShow := regexp.MustCompile(`[\W]|_`)

	//	See how we made out:
	if rxSE.MatchString(filename) {
		seMatches := getMatches(rxSE, filename)
		retval.ParseType = ParseTypeSE
		showName := rxShow.ReplaceAllString(strings.TrimSpace(seMatches["series_name"]), " ")
		retval.ShowName = showName
		retval.SeasonNumber, _ = strconv.Atoi(seMatches["season_num"])
		retval.EpisodeNumber, _ = strconv.Atoi(seMatches["ep_num"])
	} else if rxD.MatchString(filename) {
		seMatches := getMatches(rxD, filename)
		retval.ParseType = ParseTypeDate
		showName := rxShow.ReplaceAllString(strings.TrimSpace(seMatches["series_name"]), " ")
		retval.ShowName = showName
		retval.AiredYear, _ = strconv.Atoi(seMatches["year"])
		retval.AiredMonth, _ = strconv.Atoi(seMatches["month"])
		retval.AiredDay, _ = strconv.Atoi(seMatches["day"])
	} else if rxSE2.MatchString(filename) {
		seMatches := getMatches(rxSE2, filename)
		retval.ParseType = ParseTypeSE2
		showName := rxShow.ReplaceAllString(strings.TrimSpace(seMatches["series_name"]), " ")
		retval.ShowName = showName
		retval.SeasonNumber, _ = strconv.Atoi(seMatches["season_num"])
		retval.EpisodeNumber, _ = strconv.Atoi(seMatches["ep_num"])
	}

	return retval, nil
}

// getMatches return the named match groups and
// their associated values for a given compiled
//	regex, and test string
func getMatches(rx *regexp.Regexp, findString string) map[string]string {
	rxMatches := make(map[string]string)

	rxMatchArray := rx.FindStringSubmatch(findString)
	for i, name := range rx.SubexpNames() {
		if i > 0 && i <= len(rxMatchArray) {
			rxMatches[name] = rxMatchArray[i]
		}
	}

	return rxMatches
}

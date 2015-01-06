package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var Tag_Helper map[string]map[string]struct{}

func init() {
	Tag_Helper = map[string]map[string]struct{}{
		"NFL": map[string]struct{}{
			//Teams
			"cowboys":    struct{}{},
			"patriots":   struct{}{},
			"eagles":     struct{}{},
			"49ers":      struct{}{},
			"seahawks":   struct{}{},
			"packers":    struct{}{},
			"broncos":    struct{}{},
			"steelers":   struct{}{},
			"bears":      struct{}{},
			"raiders":    struct{}{},
			"giants":     struct{}{},
			"vikings":    struct{}{},
			"ravens":     struct{}{},
			"dolphins":   struct{}{},
			"chargers":   struct{}{},
			"browns":     struct{}{},
			"lions":      struct{}{},
			"bills":      struct{}{},
			"redkings":   struct{}{},
			"jets":       struct{}{},
			"cardinals":  struct{}{},
			"panthers":   struct{}{},
			"colts":      struct{}{},
			"falcons":    struct{}{},
			"bengals":    struct{}{},
			"buccaneers": struct{}{},
			"chiefs":     struct{}{},
			"jaguars":    struct{}{},
			"titans":     struct{}{},
			"rams":       struct{}{},
			"texans":     struct{}{},

			"touchdowns": struct{}{},
			"football":   struct{}{},
		},
		"Finance": map[string]struct{}{
			"dow":    struct{}{},
			"s&p":    struct{}{},
			"nasdaq": struct{}{},
		},
	}
}

func GetSuggestedTagHandler(w http.ResponseWriter, r *http.Request) {
	predictionText := "ram and p"
	predictionText = strings.ToLower(predictionText)
	searchStrings := strings.Split(predictionText, " ")
	tags := GetSuggestedTags(searchStrings)
	j, _ := json.Marshal(tags)
	fmt.Fprintln(w, string(j))
}

func GetSuggestedTags(searchStrings []string) []PtTag {
	tags := []PtTag{}
	suffix_list := []string{"", "s"}
	for _, word := range searchStrings {
		for k, v := range Tag_Helper {
			for _, suffix := range suffix_list {
				_, exists := v[word+suffix]
				if exists {
					tags = append(tags,
						PtTag{
							Tag: k,
						})
				}
			}
		}
	}
	return tags
}

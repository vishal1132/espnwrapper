package espn

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

/*
	Match URL to hit for json.
	https://www.espncricinfo.com/<anything>/engine/match/<matchid>.json
	For example https://www.espncricinfo.com/anything/engine/match/<matchid>.json
	for india vs newzealand wtc 1249875
*/

// GetMatchSummary returns the match summary in form of ESPNMatch struct. Refer to ESPNMatch struct for more.
func (e *ESPN) GetMatchSummary(matchid string) (*ESPNMatch, error) {
	resp, err := e.c.Get(fmt.Sprintf("https://www.espncricinfo.com/something/engine/match/%v.json", matchid))
	if err != nil {
		return nil, err
	}
	// response body close.
	defer resp.Body.Close()
	// read response body
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// kk := fmt.Sprintf("%v}", string(b)[:strings.Index(string(b), "comms")-2])
	// l := fmt.Sprintf("%v}", kk)
	// log.Println(l)
	espnMatch := &ESPNMatch{}
	if err := json.Unmarshal([]byte(fmt.Sprintf("%v}", string(b)[:strings.Index(string(b), "comms")-2])), espnMatch); err != nil {
		return nil, err
	}
	return espnMatch, nil
}

// CommonBatting is the complete batting scorecard for the batting team including the batsmen who are already in the dug out.
type CommonBatting struct {
	BallsFaced      string `json:"balls_faced"`
	Hand            string `json:"hand"`
	ImagePath       string `json:"image_path"`
	KnownAs         string `json:"known_as"`
	Notout          int    `json:"notout"`
	PlayerID        int    `json:"player_id"`
	PopularName     string `json:"popular_name"`
	Position        int    `json:"position"`
	PositionGroup   string `json:"position_group"`
	Runs            int    `json:"runs"`
	LiveCurrentName string `json:"live_current_name,omitempty"`
}

// CommonBowling is the complete bowling scorecard for the bowling team including bowlers who are not the one having the spell currently.
type CommmonBowling struct {
	Conceded        int    `json:"conceded"`
	Hand            string `json:"hand"`
	ImagePath       string `json:"image_path"`
	KnownAs         string `json:"known_as"`
	LiveCurrentName string `json:"live_current_name,omitempty"`
	Maidens         int    `json:"maidens"`
	Overs           string `json:"overs"`
	Pacespin        string `json:"pacespin"`
	PlayerID        string `json:"player_id"`
	PopularName     string `json:"popular_name"`
	Position        int    `json:"position"`
	Wickets         int    `json:"wickets"`
}

// CommonInnings is the innings summary description struct.
type CommonInnings struct {
	ControlPercentage int         `json:"control_percentage"`
	DotBallPercentage int         `json:"dot_ball_percentage"`
	Event             int         `json:"event"`
	EventName         string      `json:"event_name"`
	OverLimit         string      `json:"over_limit"`
	Overs             string      `json:"overs"`
	RunRate           interface{} `json:"run_rate"` // NOTE: RunRate is sometimes integer and sometimes string, will have to look
	Runs              int         `json:"runs"`
	RunsSummary       []string    `json:"runs_summary"`
	Target            int         `json:"target"`
	Wickets           int         `json:"wickets"`
}

// CommonInningsList contains the information of the specific innings in a match.
type CommonInningsList struct {
	Current          int    `json:"current"`
	Description      string `json:"description"`
	DescriptoinShort string `json:"descriptoin_short"`
	InningsNumber    int    `json:"innings_number"`
	Selected         int    `json:"selected"`
	TeamID           int    `json:"team_id"`
}

// Match is the description of the current match.
type Match struct {
	ControlPercentage int      `json:"control_percentage"`
	DotBallPercentage int      `json:"dot_ball_percentage"`
	ResultString      string   `json:"result_string"`
	RunsSummary       []string `json:"runs_summary"`
}

// FOW is the fall of wickets.
type FOW struct {
	Notout int         `json:"notout"`
	Overs  string      `json:"overs"`
	Player []FOWPlayer `json:"player"`
	Runs   int         `json:"runs"`
}

// FOWPlayer is each player's fall of wickets records.
type FOWPlayer struct {
	KnownAs     string `json:"known_as"`
	PlayerID    string `json:"player_id"`
	PopularName string `json:"popular_name"`
	Runs        int    `json:"runs"`
}

// Centre is the driving struct of the match summary.
// Everything is unmarshalled in {"centre":...}
type Centre struct {
	Ba            []Batting `json:"batting"`
	Bo            []Bowling `json:"bowling"`
	Co            Common    `json:"common"`
	InningsNumber string    `json:"innings_number"`
	Match         Match     `json:"match"`
	Fow           []FOW     `json:"fow"`
}

// ESPNMatch is returned by the GetSummary method defined on ESPN struct.
type ESPNMatch struct {
	Centre Centre `json:"centre"`
}

type Common struct {
	Batting     []CommonBatting     `json:"batting"`
	Bowling     []CommmonBowling    `json:"bowling"`
	Innings     CommonInnings       `json:"innings"`
	InningsList []CommonInningsList `json:"innings_list"`
}

type Batting struct {
	BallsFaced        string               `json:"balls_faced"`
	BattingStyle      string               `json:"batting_style"`
	ControlPercentage int                  `json:"control_percentage"`
	DismissalName     string               `json:"dismissal_name"`
	DotBallPercentage int                  `json:"dot_ball_percentage"`
	KnownAs           string               `json:"known_as"`
	LiveCurrentName   string               `json:"live_current_name"`
	MatchAward        int                  `json:"match_award"`
	Notout            int                  `json:"notout"`
	PlayerID          string               `json:"player_id"`
	PopularName       string               `json:"popular_name"`
	PreferredShot     PreferredShotBatting `json:"preferred_shot"`
	Runs              int                  `json:"runs"`
	ScoringShots      string               `json:"scoring_shots"`
	StrikeRate        string               `json:"strike_rate"`
	// WagonZone         []WagonZoneBatting   `json:"wagon_zone"`
	// RunsSummary  []string `json:"runs_summary"`
}

// WagonZoneBatting is the wagon wheel/run summary description of the player currently batting.
type WagonZoneBatting struct {
	Runs         int   `json:"runs"`
	RunsSummary  []int `json:"runs_summary"`
	ScoringShots int   `json:"scoring_shots"`
}

// PreferredShotBatting is the Batsman's preferred shot description embedded in Batting struct
type PreferredShotBatting struct {
	BallsFaced  string   `json:"balls_faced"`
	Runs        string   `json:"runs"`
	RunsSummary []string `json:"runs_summary"`
	ShotName    string   `json:"shot_name"`
}

// BattingStyleBowling is the description of the bowler bowling to Right Hand Batsman(rhb) or lhb.
type BattingStyleBowling struct {
	Balls       int    `json:"balls"`
	Conceded    string `json:"conceded"`
	EconomyRate string `json:"economy_rate"`
	Wickets     string `json:"wickets"`
}

// Bowling is the description of the bowler.
type Bowling struct {
	BowlingStyle    string              `json:"bowling_style"`
	Conceded        int                 `json:"conceded"`
	EconomyRate     string              `json:"economy_rate"`
	KnownAs         string              `json:"known_as"`
	LiveCurrentName string              `json:"live_current_name"`
	Maidens         int                 `json:"maidens"`
	MatchAward      int                 `json:"match_award"`
	OverallLhb      BattingStyleBowling `json:"overall_lhb"`
	OverallRhb      BattingStyleBowling `json:"overall_rhb"`
	Overs           string              `json:"overs"`
	PlayerID        string              `json:"player_id"`
	PopularName     string              `json:"popular_name"`
	Wickets         int                 `json:"wickets"`
	// PitchMapLhb [][][]int `json:"pitch_map_lhb"`
	// PitchMapRhb [][][]int     `json:"pitch_map_rhb"`
}

// all matches summary.
package espn

import (
	"strings"

	"golang.org/x/net/html"
)

// Node Types
/*
	ErrorNode int = iota    //0
	TextNode
	DocumentNode
	ElementNode
	CommentNode
	DoctypeNode
	// RawNode nodes are not returned by the parser, but can be part of the
	// Node tree passed to func Render to insert raw HTML (without escaping).
	// If so, this package makes no guarantee that the rendered HTML is secure
	// (from e.g. Cross Site Scripting attacks) or well-formed.
	RawNode
*/

// ESPNMatchDescription is the match representation of this package.
type ESPNMatchDescription struct {
	MatchID     string
	TeamA       string
	TeamB       string
	Description string
	Link        string
	Status      string
}

// GetAllMatches is going to return a map of string to some struct or slice of struct or maybe something else, not exactly sure for now.
// the struct is going to contain the match id.
// This will have to parse raw html from a page and scrap the desired html tags.
// Some div tag having a specific class, that i am forgetting right now at the time of righting this comment. But that div will have all the information.
// url to hit for is https://www.espncricinfo.com/live-cricket-score
func (e *ESPN) GetAllMatches() (*[]ESPNMatchDescription, error) {
	res, err := e.c.Get("https://www.espncricinfo.com/live-cricket-score")
	if err != nil {
		// error handle, for now just return
		return nil, err
	}
	defer res.Body.Close()
	n, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	matches := traverseMatches(n, make([]ESPNMatchDescription, 0))
	validMatches := make([]ESPNMatchDescription, 0, len(matches)/2)
	for _, v := range matches {
		if v.Description != "" {
			validMatches = append(validMatches, v)
		}
	}
	return &validMatches, nil
}

// func printNode(n *html.Node) {
// 	log.Println("------------------------------------")
// 	log.Printf("%p", n)
// 	for c := n.NextSibling; c != nil; c = c.NextSibling {
// 		log.Printf("node pointer %p", c)
// 		log.Println(c.Data, c.Type)
// 	}
// 	log.Println("------------------------------------")
// }

// traverseMatches looks for anchor tags with class set to - match-info-link-FIXTURES
func traverseMatches(n *html.Node, matches []ESPNMatchDescription) []ESPNMatchDescription {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			// check if class matches - match-info-link-FIXTURES
			if attr.Key == "class" && attr.Val == "match-info-link-FIXTURES" {
				// log.Println("class equals", attr.Key, attr.Val)
				// printNode(n)
				matches = append(matches, extractMatch(n))
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		matches = traverseMatches(c, matches)
	}
	return matches
}

func extractMatch(n *html.Node) ESPNMatchDescription {
	// log.Printf("extractMatch %p", n)
	var match = ESPNMatchDescription{}
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			match.Link = attr.Val
			match.extractMatchIDFromLink()
			break
		}
	}
	match.Description = extractDescriptionInfo(n, "")
	// match.Teams = extractTeamInfo(n, make([]string, 0))
	return match
}

func extractDescriptionInfo(n *html.Node, description string) string {
	if n.Type == html.ElementNode && n.Data == "div" {
		// div now check for attrs for class = description for the description of the message.
		for _, v := range n.Attr {
			if v.Key == "class" && v.Val == "description" {
				description = extractText(n, "")
				// log.Println("description", description)
				return description
			}
		}
	}
	// for c := n.NextSibling; c != nil; c = c.FirstChild is giving a lot more.
	for c := n.NextSibling; c != nil; c = c.NextSibling {
		// log.Println(c.Data)
		description = extractDescriptionInfo(c, description)
	}
	return strings.TrimSpace(description)
}

func extractText(n *html.Node, text string) string {
	if n.Type == html.TextNode {
		if len(text) > 0 && !strings.HasSuffix(text, " ") {
			text = text + " "
		}
		text = text + strings.TrimSpace(strings.Trim(n.Data, "\n"))
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text = extractText(c, text)
	}
	return strings.TrimSpace(text)
}

func (m *ESPNMatchDescription) extractMatchIDFromLink() {
	paths := strings.Split(m.Link, "/")
	dashes := strings.Split(paths[len(paths)-2], "-")
	m.MatchID = dashes[len(dashes)-1]
	for i, v := range dashes {
		if v == "vs" {
			m.TeamA = strings.Join(dashes[:i], " ")
			m.TeamB = strings.Join(dashes[i+1:len(dashes)-2], " ")
			break
		}
	}
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/vishal1132/espnwrapper/espn"
)

type CLI struct {
	s       tcell.Screen
	e       *espn.ESPN
	wg      *sync.WaitGroup
	refresh int
	matchid string
}

func main() {
	c := CLI{
		e: espn.New(),
		// wg: &sync.WaitGroup{},
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	c.initScreen()
	// flags for the cli
	interval := flag.Int("refresh", 2, "refresh interval to refresh the score. Default interval is 2 seconds")
	flag.Parse()
	c.refresh = *interval
	// flags for cli over
	c.s.Clear()
	matches, err := c.e.GetAllMatches()
	if err != nil {
		log.Println(err)
	}
	c.showMatches(matches)

	// c.wg.Wait()
}

func (c *CLI) showMatches(matches *[]espn.ESPNMatchDescription) {
	c.s.Clear()
	nh := 0
	w, _ := c.s.Size()
	for i, v := range *matches {
		emitStr(c.s, 0, nh, tcell.StyleDefault, fmt.Sprintf("Match Number %d:", i+1))
		nh++
		emitStr(c.s, 0, nh, tcell.StyleDefault, fmt.Sprintf("%s vs %s", v.TeamA, v.TeamB))
		nh++
		emitStr(c.s, 0, nh, tcell.StyleDefault, v.Description)
		nh += 2
	}
	emitStr(c.s, w/2-15, nh, tcell.StyleDefault, "Enter match number to follow or escape to exit.")
	c.s.Show()
	for {
		switch ev := c.s.PollEvent().(type) {
		case *tcell.EventResize:
			c.s.Sync()
			// displayHelloWorld(s)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyESC {
				c.s.Fini()
				// c.wg.Done()
				os.Exit(0)
			}

			if ev.Rune() == 0 {
				//no op
				continue
			}
			if int(ev.Rune())-49 >= len(*matches) || ev.Rune()-49 < 0 {
				emitStr(c.s, 0, nh+1, tcell.StyleDefault, "Please enter a valid match number")
				continue
			}
			c.matchid = (*matches)[ev.Rune()-49].MatchID
			c.showScore()
		}
	}
}

func (c *CLI) showScore() {
	t := time.Tick(time.Duration(c.refresh) * time.Second)
	go func() {
		for {
			select {
			case <-t:
				c.printMatchSummary(c.matchid)

			}
		}
	}()
	for {
		switch ev := c.s.PollEvent().(type) {
		case *tcell.EventResize:
			c.s.Sync()
			c.printMatchSummary(c.matchid)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyESC {
				c.s.Fini()
				// c.wg.Done()
				os.Exit(0)
			}
		}
	}

}

func (c *CLI) printMatchSummary(matchID string) {
	summary, err := c.e.GetMatchSummary(matchID)
	if err != nil {
		// no op, continue for the next ticker
		log.Println(err)
		return
	}
	c.s.Clear()
	// next height is the next height for the line to be printed
	nh := 0
	w, h := c.s.Size()
	{
		// Score
		x, _ := emitStr(c.s, 0, nh, tcell.StyleDefault, "Match Summary")
		nh += 1
		x, _ = emitStr(c.s, 0, nh, tcell.StyleDefault, fmt.Sprintf("%v Overs", summary.Centre.Co.Innings.Overs))
		x, _ = emitStr(c.s, x+5, nh, tcell.StyleDefault, fmt.Sprintf("%v Runs", summary.Centre.Co.Innings.Runs))
		x, _ = emitStr(c.s, x+1, nh, tcell.StyleDefault, fmt.Sprintf("/%v Wickets", summary.Centre.Co.Innings.Wickets))
		x, _ = emitStr(c.s, x+5, nh, tcell.StyleDefault, fmt.Sprintf("%v RunRate", summary.Centre.Co.Innings.RunRate))
		nh += 1
	}
	// log.Println(summary.Centre)
	{
		// for batsmen
		nh += 1
		emitStr(c.s, 0, nh, tcell.StyleDefault, "Batsmen")
	}
	for _, v := range summary.Centre.Ba {
		// same height but different widths
		nh += 1
		x, _ := emitStr(c.s, 0, nh, tcell.StyleDefault, v.KnownAs)
		x, _ = emitStr(c.s, x+5, nh, tcell.StyleDefault, fmt.Sprintf("%v Runs", v.Runs))
		x, _ = emitStr(c.s, x+5, nh, tcell.StyleDefault, fmt.Sprintf("%v Balls", v.BallsFaced))
	}
	{
		nh += 2
		emitStr(c.s, 0, nh, tcell.StyleDefault, "Bowlers")
	}
	for _, v := range summary.Centre.Bo {
		nh += 1
		x, _ := emitStr(c.s, 0, nh, tcell.StyleDefault, v.KnownAs)
		x, _ = emitStr(c.s, x+5, nh, tcell.StyleDefault, fmt.Sprintf("%v Runs", v.Conceded))
		x, _ = emitStr(c.s, x+5, nh, tcell.StyleDefault, fmt.Sprintf("%v Overs", v.Overs))
	}
	{
		if summary.Centre.Match.ResultString != "" {
			nh += 2
			emitStr(c.s, 0, nh, tcell.StyleDefault, summary.Centre.Match.ResultString)
		}
	}

	// for _, v := range summary.Centre. {
	// 	nh += 1
	// 	x, _ := emitStr(c.s, 0, nh, tcell.StyleDefault, v.KnownAs)
	// 	x, _ = emitStr(c.s, x+5, nh, tcell.StyleDefault, fmt.Sprintf("%v Runs", v.Conceded))
	// 	x, _ = emitStr(c.s, x+5, nh, tcell.StyleDefault, fmt.Sprintf("%v Overs", v.Overs))
	// }
	emitStr(c.s, w/2-9, h/2+1, tcell.StyleDefault, "Press ESC to exit.")
	c.s.Show()
	/*
		// working
		if len(summary.Centre.Ba) >= 2 {
			fmt.Printf("\r%v  -  %v Runs %v Balls\t%v   -   %v Runs %v Balls", summary.Centre.Ba[0].KnownAs, summary.Centre.Ba[0].Runs, summary.Centre.Ba[0].BallsFaced,
				summary.Centre.Ba[1].KnownAs, summary.Centre.Ba[1].Runs, summary.Centre.Ba[1].BallsFaced)
		}
	*/
}

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) (int, int) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
	return x, y
}

func (c *CLI) initScreen() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackUTF8)
	var err error
	c.s, err = tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	c.s.SetStyle(defStyle)
	if err := c.s.Init(); err != nil {
		log.Fatal(err)
	}
	c.s.Clear()
}

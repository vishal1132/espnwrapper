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
	s  tcell.Screen
	e  *espn.ESPN
	wg *sync.WaitGroup
}

func main() {
	c := CLI{
		e: espn.New(),
		// wg: &sync.WaitGroup{},
	}
	c.initScreen()
	// flags for the cli
	interval := flag.Int("refresh", 2, "refresh interval to refresh the score. Default interval is 2 seconds")
	matchID := flag.String("matchid", "1249875", "matchid extract from espncricinfo")
	flag.Parse()
	// flags for cli over
	t := time.Tick(time.Duration(*interval) * time.Second)
	c.s.Clear()
	// c.wg.Add(1)
	// var matchID string = "1263150"
	go func() {
		for {
			select {
			case <-t:
				c.printMatchSummary(*matchID)

			}
		}
	}()
	for {
		switch ev := c.s.PollEvent().(type) {
		case *tcell.EventResize:
			c.s.Sync()
			c.printMatchSummary(*matchID)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyESC {
				c.s.Fini()
				// c.wg.Done()
				os.Exit(0)
			}
		}
	}

	// c.wg.Wait()
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

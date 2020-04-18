package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Do not forget to normalize files:
// ffmpeg -i in.flac -filter:a loudnorm -sample_fmt s16 -ar 44100 out.flac
func main() {
	n := flag.Int("n", 10, "Number of listenings")
	flag.Parse()

	if flag.NArg() < 2 {
		_, _ = fmt.Fprintln(os.Stderr, "Usage:", "compare-tracks", "track1", "track2")
		os.Exit(1)
	}
	tracks := [...]string{flag.Arg(0), flag.Arg(1)}
	results := [2]trackResults{}

	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < *n; i++ {
		trackI := rand.Intn(2)

		fmt.Println("Playing", tracks[trackI])
		cmd := exec.Command("ffplay", "-nodisp", tracks[trackI])
		err := cmd.Start()
		if err != nil {
			panic(err)
		}

		var answer int
		for {
			fmt.Print("Which track is this: 1 or 2? ")
			line, _ := reader.ReadString('\n')
			answer, err = strconv.Atoi(strings.TrimSpace(line))
			if err == nil && answer >= 1 && answer <= 2 {
				break
			}
			fmt.Println("Bad input")
		}

		err = cmd.Process.Signal(os.Interrupt)
		if err != nil {
			panic(err)
		}

		if answer-1 == trackI {
			results[trackI].right++
		} else {
			results[trackI].wrong++
		}
	}

	fmt.Println("Played track  Right answers  Wrong answers")
	for i, trackResults := range results {
		fmt.Printf("%-13d %-14d %d\n", i+1, trackResults.right, trackResults.wrong)
	}
}

type trackResults struct {
	right int
	wrong int
}

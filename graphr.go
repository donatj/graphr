package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var glyphs = []string{" ", "▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

type tstack struct {
	size int

	Amax int
	Amin int

	Stack []int
}

func (t *tstack) Push(i int) {
	if i > t.Amax {
		t.Amax = i
	}

	if i < t.Amin {
		t.Amin = i
	}

	t.Stack = append(t.Stack, i)

	if len(t.Stack) > t.size {
		t.Stack = t.Stack[len(t.Stack)-t.size:]
	}
}

func (t *tstack) Max() (i int) {
	for _, v := range t.Stack {
		if v > i {
			i = v
		}
	}

	return
}

func (t *tstack) Min() (i int) {
	i = t.Max()
	for _, v := range t.Stack {
		if v < i {
			i = v
		}
	}

	return
}

func init() {
	flag.Parse()
}

func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	fmt.Print("\033[?7l")
	fmt.Println("\033[?25l")

	go func() {
		<-c

		fmt.Println("\033[?7h")
		fmt.Println("\033[?25h\033[?0c")

		os.Exit(1)
	}()
}

func getVal(scmd string) int {
	cmd := exec.Command("bash", "-c", scmd)
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	val := string(out)
	val = strings.TrimSpace(val)

	i, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}

	return i
}

func main() {
	t := tstack{
		size: width(),
		Amin: 1<<31 - 1,
		Amax: -1 << 31,
	}

	args := flag.Args()
	cmd := strings.Join(args, " ")
	cmd = strings.TrimSpace(cmd)

	if cmd == "" {
		log.Fatal("missing value providing command argument")
	}

	/// move down 3 lines and then back up
	// making space if we're at the bottom of the terminal
	fmt.Print("\n\n\n\033[3A")
	fmt.Print("\0337")
	//	ii := 1
	for {
		t.size = width()

		val := getVal(cmd)
		t.Push(val)

		min := t.Min()
		max := t.Max()
		dif := max - min

		fmt.Print("\033[K")
		for _, i := range t.Stack {
			gi := 0
			if dif != 0 {
				gi = ((i - min) * 8) / dif
			}

			fmt.Print(glyphs[gi])
		}

		fmt.Print("\n\033[J", val, " of ", max, "(", t.Amin, "/", t.Amax, ")")
		fmt.Print("\0338")

		time.Sleep(5000 * time.Millisecond)
	}
}

func width() (width int) {
	width = 30 //default

	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return
	}

	parts := strings.Split(string(out), " ")
	switch len(parts) {
	case 2:
		w, err := strconv.Atoi(strings.Trim(parts[1], " \n"))
		if err != nil {
			return
		}

		width = w
	case 1:
	default:
		log.Fatal("Super weird response")
	}

	return
}

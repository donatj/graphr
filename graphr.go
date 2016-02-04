package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
		t.Stack = t.Stack[1:]
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
		size: width() - 1,
		Amin: 1<<31 - 1,
		Amax: -1 << 31,
	}

	args := flag.Args()
	cmd := strings.Join(args, " ")

	fmt.Print("\n\n\n\033[3A\0337")
	//	ii := 1
	for {
		val := rand.Intn(1000)
		val = getVal(cmd)
		t.Push(val)

		min := t.Min()
		max := t.Max()
		dif := max - min

		for _, i := range t.Stack {

			gi := 0
			if dif != 0 {
				gi = ((i - min) * 8) / dif
			}
			_ = gi

			//			fmt.Print(glyphs[gi])
			fmt.Print(glyphs[gi])

			//			fmt.Println(glyphs[], val, len(t.Stack))
		}

		fmt.Print("\n", val, " of ", max, "(", t.Amin, "/", t.Amax, ")")
		//		fmt.Printf("%#v", t.Stack)
		time.Sleep(5000 * time.Millisecond)
		fmt.Print("\0338")
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

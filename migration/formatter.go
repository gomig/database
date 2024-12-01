package migration

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/windows"
)

const (
	RESET     = "\033[0m"
	BOLD      = "\033[1m"
	UNDERLINE = "\033[4m"
	STRIKE    = "\033[9m"
	ITALIC    = "\033[3m"
)

const (
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
	BLUE   = "\033[34m"
	PURPLE = "\033[35m"
	CYAN   = "\033[36m"
	GRAY   = "\033[37m"
	WHITE  = "\033[37m"
)

func init() {
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32

	windows.GetConsoleMode(stdout, &originalMode)
	windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}

// styling patterns
//
// {R}: RESET, {B}: BOLD ,{U}: UNDERLINE ,{S}: STRIKE
// {I}: ITALIC ,{r}: RED ,{g}: GREEN ,{y}: YELLOW
//
// {b}: BLUE ,{p}: PURPLE ,{c}: CYAN ,{m}: GRAY
// {w}: WHITE
func Formatter(pattern string, args ...any) {
	replacer := strings.NewReplacer(
		"{R}", RESET,
		"{B}", BOLD,
		"{U}", UNDERLINE,
		"{S}", STRIKE,
		"{I}", ITALIC,
		"{r}", RED,
		"{g}", GREEN,
		"{y}", YELLOW,
		"{b}", BLUE,
		"{p}", PURPLE,
		"{c}", CYAN,
		"{m}", GRAY,
		"{w}", WHITE,
	)
	fmt.Printf(replacer.Replace(pattern), args...)
}

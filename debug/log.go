package debug

import (
	"fmt"
	"time"
)

func Log(filename string, message string) {
	colourMap := map[string]string{
		"main":     "\033[33m", // yellow
		"chatRoom": "\033[32m", // green
		"server":   "\033[34m", // blue
		"p2p":      "\033[35m", // magenta
		"keys":     "\033[36m", // cyan
		"raft":     "\033[92m", // bright green
		"db":       "\033[94m", // bright blue
		"err":      "\033[91m", // bright red
		"reset":    "\033[0m",
	}
	if filename == "error" {
		fmt.Println(colourMap[filename] + "[" + filename + "] [" + time.Now().Format("15:04:05") + "] " + colourMap["reset"] + message)
	} else {
		fmt.Println(colourMap[filename] + "[" + filename + ".go] [" + time.Now().Format("15:04:05") + "] " + colourMap["reset"] + message)
	}
}

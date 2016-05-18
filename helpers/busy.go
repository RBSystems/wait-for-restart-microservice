package helpers

import (
	"regexp"
	"strings"
	"time"

	"github.com/ziutek/telnet"
)

func IsSystemBusy(curReq Request) bool {
	var conn *telnet.Conn

	conn, err := telnet.Dial("tcp", curReq.MachineAddress+":23")
	if err != nil {
		return true
	}

	conn.SetUnixWriteMode(true) // Convert any '\n' (LF) to '\r\n' (CR LF)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	_, err = conn.Write([]byte("\n\n"))
	if err != nil {
		return true
	}

	// Dynamically get the prompt
	conn.SkipUntil(">")
	promptBytes, err := conn.ReadUntil(">")
	if err != nil {
		return true
	}

	regex := "\\S.*?>"

	re := regexp.MustCompile(regex)

	prompt := string(re.Find(promptBytes))

	_, err = conn.Write([]byte("hostname\n\n")) // Send a second newline so we get the prompt
	if err != nil {
		return true
	}

	err = conn.SkipUntil(prompt)
	if err != nil {
		return true
	}

	response, err := conn.ReadUntil(prompt) // Read until the second prompt delimiter
	if err != nil {
		return true
	}

	conn.Close()

	if strings.Contains(string(response), "system is busy") {
		return true
	}

	return false
}

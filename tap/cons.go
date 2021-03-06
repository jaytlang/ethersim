package tap

import (
	"bufio"
	"fmt"
	"os"
)

func readMsgChars() string {
	fmt.Print("Type your message here: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

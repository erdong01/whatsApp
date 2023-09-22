package whatsapp

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
)

func NewConn() {
	wac, err := whatsapp.NewConn(20 * time.Second)
	if err != nil {
		slog.Error("new connn", err)
	}
	qrChan := make(chan string)
	go func() {
		terminal := qrcodeTerminal.New()
		terminal.Get(<-qrChan).Print()
		fmt.Printf("qr code: %v\n", <-qrChan)
		//show qr code or save it somewhere to scan
	}()
	sess, err := wac.Login(qrChan)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error during login: %v\n", err)
		return
	}

	fmt.Printf("login successful, session: %v\n", sess)

}

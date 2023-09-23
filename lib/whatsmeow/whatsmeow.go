package whatsmeow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
	"whatsApp/service"
	"whatsApp/service/chatLogic"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/skip2/go-qrcode"
	"google.golang.org/protobuf/proto"

	"os"
	"os/signal"
	"syscall"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/appstate"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	waProto "go.mau.fi/whatsmeow/binary/proto"
)

var log waLog.Logger
var cli *whatsmeow.Client

type Handler struct {
	UserId uint
}

func (h *Handler) eventHandler(evt any) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())
		fmt.Println("events.Message ---------------------- start")
		b, _ := json.Marshal(v)
		fmt.Println("events.Message", string(b))

		if v.Message.GetConversation() != "" {
			chatLogic.ReceiverMessageStore(h.UserId, v.Info.Sender.User, v.Message.GetConversation(), v.Info.ID)
		} else {
			user, err := service.ServiceApp.WhatsAppUserService.FindByPhone(v.Info.Sender.User)
			if err != nil {
				user.Phone = v.Info.Sender.User
				user, err = service.ServiceApp.WhatsAppUserService.Carete(user)
				if err != nil {
					return
				}
			}
			h.UserId = user.ID
		}

		fmt.Println("events.Message ---------------------- end")
	case *events.HistorySync:
		fmt.Println("HistorySync ------------------ start")
		//入库
		chatLogic.HistorySync(h.UserId, v.Data.Conversations)

		id := atomic.AddInt32(&historySyncID, 1)
		fileName := fmt.Sprintf("history-%d-%d.json", startupTime, id)
		file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Errorf("Failed to open file to write history sync: %v", err)
			return
		}
		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		err = enc.Encode(v.Data)

		if err != nil {
			log.Errorf("Failed to write history sync: %v", err)
			return
		}
		log.Infof("Wrote history sync to %s", fileName)
		_ = file.Close()
		fmt.Println("HistorySync ------------------ end")
	case *events.AppState:
		fmt.Println("AppState ----------------------------- start")
		log.Debugf("App state event: %+v / %+v", v.Index, v.SyncActionValue)
		fmt.Println("AppState ----------------------------- end")
	case *events.Connected, *events.PushNameSetting:
		fmt.Println("events.Connected =============================== start")
		// Send presence available when connecting and when the pushname is changed.
		// This makes sure that outgoing messages always have the right pushname.
		fmt.Println("events.Connected =============================== end")
	}
}
func GetQRChannel(client *whatsmeow.Client, ginCtx *gin.Context) (png []byte, err error) {
	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				fmt.Println("QR code:", evt.Code)
				ginCtx.Header("Content-Type", "image/png")
				// qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, ginCtx.Writer)
				// fmt.Println("ginCtx.Writer", ginCtx.Writer)
				//func Encode(content string, level RecoveryLevel, size int) ([]byte, error):
				//content:内容, level 图片质量, size" 大小
				png, err := qrcode.Encode(evt.Code, qrcode.Medium, 256)
				if err != nil {
					ginCtx.JSON(http.StatusOK, "二维生成失败")
					return nil, err
				}
				return png, err

			} else {
				fmt.Println("Login event:", evt.Event)
			}

		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}
	return
}
func NewConn(userId uint) (client *whatsmeow.Client, err error) {
	log = waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", log)
	if err != nil {
		return nil, err
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, err
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)
	var handled = Handler{
		UserId: userId,
	}
	client.AddEventHandler(handled.eventHandler)
	return
}

// 退出
func Logout(client *whatsmeow.Client) (err error) {
	return cli.Logout()
}

func NewConn2() (cli *whatsmeow.Client) {
	log = waLog.Stdout("Database", "DEBUG", true)
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("mysql", "root:Ykt!142857!_ykt@tcp(185.132.125.227:3306)/hello", dbLog)
	if err != nil {
		panic(err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	cli = client
	var handled Handler
	client.AddEventHandler(handled.eventHandler)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				fmt.Println("QR code:", evt.Code)
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}

		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("登陆成功")

	cli.AddEventHandler(handler)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
	return
}

func Send(client *whatsmeow.Client, phone string, cotent string) (msgId string, err error) {
	parseJid, ok := parseJID(phone)
	if !ok {
		return "", errors.New("Send() parseJID ok false")
	}
	msg := &waProto.Message{Conversation: proto.String(strings.Join([]string{cotent}, " "))}
	result, err := client.SendMessage(context.Background(), parseJid, msg)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

func parseJID(arg string) (types.JID, bool) {
	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), true
	} else {
		recipient, err := types.ParseJID(arg)
		if err != nil {
			log.Errorf("Invalid JID %s: %v", arg, err)
			return recipient, false
		} else if recipient.User == "" {
			log.Errorf("Invalid JID %s: no server specified", arg)
			return recipient, false
		}
		return recipient, true
	}
}

var historySyncID int32
var startupTime = time.Now().Unix()

func handler(rawEvt interface{}) {
	switch evt := rawEvt.(type) {
	case *events.AppStateSyncComplete:
		if len(cli.Store.PushName) > 0 && evt.Name == appstate.WAPatchCriticalBlock {
			err := cli.SendPresence(types.PresenceAvailable)
			if err != nil {
				log.Warnf("Failed to send available presence: %v", err)
			} else {
				log.Infof("Marked self as available")
			}
		}
	case *events.Connected, *events.PushNameSetting:
		if len(cli.Store.PushName) == 0 {
			return
		}
		// Send presence available when connecting and when the pushname is changed.
		// This makes sure that outgoing messages always have the right pushname.
		err := cli.SendPresence(types.PresenceAvailable)
		if err != nil {
			log.Warnf("Failed to send available presence: %v", err)
		} else {
			log.Infof("Marked self as available")
		}
	case *events.StreamReplaced:
		os.Exit(0)
	case *events.Message:
		fmt.Println("start Message-------------------")
		metaParts := []string{fmt.Sprintf("pushname: %s", evt.Info.PushName), fmt.Sprintf("timestamp: %s", evt.Info.Timestamp)}
		if evt.Info.Type != "" {
			metaParts = append(metaParts, fmt.Sprintf("type: %s", evt.Info.Type))
		}
		if evt.Info.Category != "" {
			metaParts = append(metaParts, fmt.Sprintf("category: %s", evt.Info.Category))
		}
		if evt.IsViewOnce {
			metaParts = append(metaParts, "view once")
		}
		if evt.IsViewOnce {
			metaParts = append(metaParts, "ephemeral")
		}
		if evt.IsViewOnceV2 {
			metaParts = append(metaParts, "ephemeral (v2)")
		}
		if evt.IsDocumentWithCaption {
			metaParts = append(metaParts, "document with caption")
		}
		if evt.IsEdit {
			metaParts = append(metaParts, "edit")
		}

		log.Infof("Received message %s from %s (%s): %+v", evt.Info.ID, evt.Info.SourceString(), strings.Join(metaParts, ", "), evt.Message)

		if evt.Message.GetPollUpdateMessage() != nil {
			decrypted, err := cli.DecryptPollVote(evt)
			if err != nil {
				log.Errorf("Failed to decrypt vote: %v", err)
			} else {
				log.Infof("Selected options in decrypted vote:")
				for _, option := range decrypted.SelectedOptions {
					log.Infof("- %X", option)
				}
			}
		} else if evt.Message.GetEncReactionMessage() != nil {
			decrypted, err := cli.DecryptReaction(evt)
			if err != nil {
				log.Errorf("Failed to decrypt encrypted reaction: %v", err)
			} else {
				log.Infof("Decrypted reaction: %+v", decrypted)
			}
		}

		img := evt.Message.GetImageMessage()
		if img != nil {
			data, err := cli.Download(img)
			if err != nil {
				log.Errorf("Failed to download image: %v", err)
				return
			}
			exts, _ := mime.ExtensionsByType(img.GetMimetype())
			path := fmt.Sprintf("%s%s", evt.Info.ID, exts[0])
			err = os.WriteFile(path, data, 0600)
			if err != nil {
				log.Errorf("Failed to save image: %v", err)
				return
			}
			log.Infof("Saved image in message to %s", path)
		}
		fmt.Println("end Message-------------------")
	case *events.Receipt:
		if evt.Type == events.ReceiptTypeRead || evt.Type == events.ReceiptTypeReadSelf {
			log.Infof("%v was read by %s at %s", evt.MessageIDs, evt.SourceString(), evt.Timestamp)
		} else if evt.Type == events.ReceiptTypeDelivered {
			log.Infof("%s was delivered to %s at %s", evt.MessageIDs[0], evt.SourceString(), evt.Timestamp)
		}
	case *events.Presence:
		if evt.Unavailable {
			if evt.LastSeen.IsZero() {
				log.Infof("%s is now offline", evt.From)
			} else {
				log.Infof("%s is now offline (last seen: %s)", evt.From, evt.LastSeen)
			}
		} else {
			log.Infof("%s is now online", evt.From)
		}
	case *events.HistorySync:
		fmt.Println("HistorySync ------------------ start")
		id := atomic.AddInt32(&historySyncID, 1)
		fileName := fmt.Sprintf("history-%d-%d.json", startupTime, id)
		file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Errorf("Failed to open file to write history sync: %v", err)
			return
		}
		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		err = enc.Encode(evt.Data)
		if err != nil {
			log.Errorf("Failed to write history sync: %v", err)
			return
		}
		log.Infof("Wrote history sync to %s", fileName)
		_ = file.Close()
		fmt.Println("HistorySync ------------------ end")
	case *events.AppState:
		log.Debugf("App state event: %+v / %+v", evt.Index, evt.SyncActionValue)
	case *events.KeepAliveTimeout:
		log.Debugf("Keepalive timeout event: %+v", evt)
	case *events.KeepAliveRestored:
		log.Debugf("Keepalive restored")
	case *events.Blocklist:
		log.Infof("Blocklist event: %+v", evt)
	}
}

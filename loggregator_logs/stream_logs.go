package loggregator_logs

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cloudfoundry/loggregatorlib/logmessage"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/pivotal-cf-experimental/veritas/say"
)

func StreamLogs(loggregatorAddr string, appGuid string, out io.Writer) error {
	var ws *websocket.Conn
	i := 0
	for {
		var err error
		ws, _, err = websocket.DefaultDialer.Dial(
			fmt.Sprintf("ws://%s/tail/?app=%s", loggregatorAddr, appGuid),
			http.Header{},
		)
		if err != nil {
			i++
			if i > 10 {
				return fmt.Errorf("Unable to connect to Server in 100ms, giving up.\n")
			}

			time.Sleep(10 * time.Millisecond)
			continue
		} else {
			break
		}
	}

	errorChan := make(chan error)
	messageChan := make(chan *logmessage.LogMessage)

	go func() {
		//keep the connection alive
		for {
			err := ws.WriteMessage(websocket.BinaryMessage, []byte{42})
			if err != nil {
				errorChan <- err
				break
			}

			<-time.After(1 * time.Second)
		}
	}()

	go func() {
		for {
			_, data, err := ws.ReadMessage()

			if err != nil {
				errorChan <- err
				break
			}

			receivedMessage := &logmessage.LogMessage{}
			err = proto.Unmarshal(data, receivedMessage)
			if err != nil {
				errorChan <- err
				break
			}

			messageChan <- receivedMessage
		}
	}()

	for {
		select {
		case message := <-messageChan:
			printLog(message, out)
		case err := <-errorChan:
			return err
		}
	}
}

func printLog(message *logmessage.LogMessage, out io.Writer) {
	color := say.DefaultStyle
	if message.GetMessageType() == logmessage.LogMessage_ERR {
		color = say.RedColor
	}
	say.Println(0,
		"[%s] %s - %s",
		say.Green(message.GetSourceName()),
		say.Yellow("%s", time.Unix(0, message.GetTimestamp()).Format(time.StampMilli)),
		say.Colorize(color, string(message.GetMessage())),
	)
}

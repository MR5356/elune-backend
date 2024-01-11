package machine

import (
	"github.com/MR5356/elune-backend/pkg/response"
	"github.com/MR5356/elune-backend/pkg/utils/sshutil"
	"github.com/MR5356/elune-backend/pkg/utils/terminal"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

func (c *Controller) handleTerminal(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}

	host, err := c.service.machinePersistence.Detail(&Machine{ID: uint(id)})
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	} else {
		sshClient, err := sshutil.NewSSHClient(host.HostInfo)
		if err != nil {
			response.Error(ctx, response.CodeParamError, err.Error())
			return
		}
		t := terminal.NewTerminal()
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			Subprotocols: []string{"webssh"},
		}

		webConn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			response.Error(ctx, response.CodeParamError, err.Error())
			return
		}

		webConn.SetCloseHandler(t.CloseHandler)

		t.Websocket = webConn
		t.Session, err = sshClient.GetSession()
		if err != nil {
			response.Error(ctx, response.CodeParamError, err.Error())
			return
		}

		// 处理输入流
		t.Stdin, err = t.Session.StdinPipe()
		if err != nil {
			response.Error(ctx, response.CodeParamError, err.Error())
			return
		}

		// 处理输出流
		sshOut := new(terminal.WsBufferWriter)
		t.Session.Stdout = sshOut
		t.Session.Stderr = sshOut
		t.Stdout = sshOut

		if err := t.Session.RequestPty("xterm-256color", 30, 120, terminal.Modes); err != nil {
			response.Error(ctx, response.CodeParamError, err.Error())
			return
		}

		err = t.Session.Shell()
		if err != nil {
			response.Error(ctx, response.CodeParamError, err.Error())
			return
		}

		go t.Send2SSH()
		go t.Send2Web()
	}
}

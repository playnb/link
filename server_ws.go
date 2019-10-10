package link

import (
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/playnb/util/log"
	"github.com/unrolled/secure"
	"net"
	"net/http"
	"sync"
	"time"
)

var gins map[string]*gin.Engine = make(map[string]*gin.Engine)

func getGinEngine(addr string) (*gin.Engine, bool) {
	if e, ok := gins[addr]; ok {
		return e, false
	} else {
		e = gin.New()
		gins[addr] = e
		return e, true
	}
}

type ServerOption struct {
	Addr            string
	MaxMsgLen       int
	MaxConnNum      int
	PendingWriteNum int

	RelativePath string
	HTTPTimeout  time.Duration
	CertFile     string
	KeyFile      string
	GinLogger    gin.HandlerFunc
}

func (opt *ServerOption) Check() {
	if opt.MaxConnNum <= 0 {
		opt.MaxConnNum = 100
	}
	if opt.PendingWriteNum <= 0 {
		opt.PendingWriteNum = 200
	}
	if opt.MaxMsgLen <= 0 {
		opt.MaxMsgLen = 4096
	}
	if opt.HTTPTimeout <= 0 {
		opt.HTTPTimeout = 10 * time.Second
	}
	if len(opt.RelativePath) == 0 {
		opt.RelativePath = "/"
	}
	if opt.GinLogger == nil {
		opt.GinLogger = gin.Logger()
	}
}

type WSServer struct {
	sync.Mutex

	upgrader    websocket.Upgrader
	option      *ServerOption
	clientConns map[uint64]Conn

	ln     net.Listener
	engine *gin.Engine

	wg sync.WaitGroup

	OnAccept func(agent *Agent)
}

func (serv *WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	conn, err := serv.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Debug("upgrade error: %v", err)
		return
	}
	conn.SetReadLimit(int64(serv.option.MaxMsgLen))

	serv.wg.Add(1)
	defer serv.wg.Done()

	wsc := func() Conn {
		serv.Lock()
		defer serv.Unlock()
		if serv.clientConns == nil {
			conn.Close()
			return nil
		}
		if len(serv.clientConns) >= serv.option.MaxConnNum {
			conn.Close()
			log.Error("too many connections")
			return nil
		}

		wsc := newWSConn(conn, serv.option.PendingWriteNum, serv.option.MaxMsgLen)
		serv.clientConns[wsc.GetUniqueID()] = wsc
		return wsc
	}()
	if wsc == nil {
		return
	}

	agent := &Agent{}
	agent.init(wsc, serv.option.MaxMsgLen)
	if serv.OnAccept != nil {
		serv.OnAccept(agent)
	}
	for {
		data, err := wsc.ReadMsg()
		if err != nil {
			break
		}
		if len(agent.msgChan) == cap(agent.msgChan) {
			log.Error("Agent msgChan full")
			continue
		}
		agent.msgChan <- data
	}
	close(agent.msgChan)
	if agent.OnClose != nil {
		agent.OnClose()
	}

	serv.Lock()
	delete(serv.clientConns, wsc.GetUniqueID())
	serv.Unlock()
}

func (serv *WSServer) Start(option *ServerOption) error {
	serv.option = option
	serv.option.Check()
	serv.clientConns = make(map[uint64]Conn)
	serv.upgrader = websocket.Upgrader{
		HandshakeTimeout: serv.option.HTTPTimeout,
		CheckOrigin:      func(_ *http.Request) bool { return true },
	}
	if false {
		ln, err := net.Listen("tcp", serv.option.Addr)
		if err != nil {
			return err
		}

		if serv.option.CertFile != "" || serv.option.KeyFile != "" {
			config := &tls.Config{}
			config.NextProtos = []string{"http/1.1"}

			var err error
			config.Certificates = make([]tls.Certificate, 1)
			config.Certificates[0], err = tls.LoadX509KeyPair(serv.option.CertFile, serv.option.KeyFile)
			if err != nil {
				return err
			}

			ln = tls.NewListener(ln, config)
		}

		serv.ln = ln

		httpServer := &http.Server{
			Addr:           serv.option.Addr,
			Handler:        serv,
			ReadTimeout:    serv.option.HTTPTimeout,
			WriteTimeout:   serv.option.HTTPTimeout,
			MaxHeaderBytes: 1024,
		}
		go func() {
			err := httpServer.Serve(ln)
			if err != nil {
				log.Fatal(err.Error())
			}
		}()
	} else {
		gin.SetMode(gin.DebugMode)
		e, newGin := getGinEngine(serv.option.Addr) //gin.New()
		serv.engine = e
		if newGin {
			if serv.option.GinLogger != nil {
				serv.engine.Use(serv.option.GinLogger)
			}
		}

		serv.engine.Handle(http.MethodGet, serv.option.RelativePath, func(c *gin.Context) {
			serv.ServeHTTP(c.Writer, c.Request)
		})
		if newGin {
			if serv.option.CertFile != "" || serv.option.KeyFile != "" {
				secureMiddleware := secure.New(secure.Options{
					FrameDeny: true,
				})
				secureFunc := func() gin.HandlerFunc {
					return func(c *gin.Context) {
						err := secureMiddleware.Process(c.Writer, c.Request)
						// If there was an error, do not continue.
						if err != nil {
							c.Abort()
							return
						}
						// Avoid header rewrite if response is a redirection.
						if status := c.Writer.Status(); status > 300 && status < 399 {
							c.Abort()
						}
					}
				}()
				serv.engine.Use(secureFunc)
				go func() {
					err := serv.engine.RunTLS(serv.option.Addr, serv.option.CertFile, serv.option.KeyFile)
					if err != nil {
						log.Error(err.Error())
					}
				}()
			} else {
				go func() {
					err := serv.engine.Run(serv.option.Addr)
					if err != nil {
						log.Error(err.Error())
					}
				}()
			}
		}
	}
	return nil
}

func (serv *WSServer) Close() {
	serv.ln.Close()
	func() {
		serv.Lock()
		defer serv.Unlock()
		for _, conn := range serv.clientConns {
			conn.Close()
		}
		serv.clientConns = nil
	}()
	serv.wg.Wait()
}

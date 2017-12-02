package insight

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/dynamicgo/config"
	"github.com/dynamicgo/slf4go"
	"github.com/inwecrypto/neo-insight/claim"
	"github.com/inwecrypto/neogo"
	"github.com/julienschmidt/httprouter"
	"github.com/ybbus/jsonrpc"
)

var logger slf4go.Logger

// OpenLogger .
func OpenLogger() {
	logger = slf4go.Get("neo-insight")
}

type handler func(params []interface{}) (interface{}, *JSONRPCError)

// Server insight api jsonrpc 2.0 server
type Server struct {
	cnf      *config.Config
	router   *httprouter.Router
	remote   *url.URL
	dispatch map[string]handler
	db       *sql.DB
	utxo     *utxoModel
	blockfee *blockFeeModel
}

type loggerHandler struct {
	handler http.Handler
}

func (l *loggerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.DebugF("http route: %s %s", r.Method, r.URL)
	l.handler.ServeHTTP(w, r)
}

// NewServer create new server
func NewServer(cnf *config.Config) (*Server, error) {

	remote, err := url.Parse(cnf.GetString("insight.neonode", "http://xxxxxx:10332"))

	if err != nil {
		return nil, err
	}

	db, err := openDB(cnf)

	if err != nil {
		return nil, err
	}

	return &Server{
		cnf:      cnf,
		router:   httprouter.New(),
		remote:   remote,
		dispatch: make(map[string]handler),
		db:       db,
		utxo:     newUTXOModel(cnf, db),
		blockfee: newBlockFeeModel(cnf, db),
	}, nil
}

func openDB(cnf *config.Config) (*sql.DB, error) {
	driver := cnf.GetString("insight.database.driver", "xxxx")
	username := cnf.GetString("insight.database.username", "xxx")
	password := cnf.GetString("insight.database.password", "xxx")
	port := cnf.GetString("insight.database.port", "6543")
	host := cnf.GetString("insight.database.host", "localhost")
	schema := cnf.GetString("insight.database.schema", "postgres")
	maxconn := cnf.GetInt64("insight.database.maxconn", 10)

	db, err := sql.Open(driver, fmt.Sprintf("user=%v password=%v host=%v dbname=%v port=%v sslmode=disable", username, password, host, schema, port))

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(int(maxconn))

	return db, nil
}

func makeJSONRPCError(w http.ResponseWriter, id uint, code int, message string, data interface{}) {
	response := &jsonrpc.RPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &jsonrpc.RPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}

	jsonresponse, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "server internal error", http.StatusInternalServerError)
		logger.ErrorF("marshal response error :%s", err)
		return
	}

	w.WriteHeader(200)

	if _, err := w.Write(jsonresponse); err != nil {
		logger.ErrorF("write response error :%s", err)
	}
}

func makeRPCRequest(r *http.Request) (*jsonrpc.RPCRequest, error) {
	request := jsonrpc.RPCRequest{}
	decoder := json.NewDecoder(r.Body)
	decoder.UseNumber()

	err := decoder.Decode(&request)

	if err != nil {
		return nil, err
	}

	return &request, nil
}

func makeJSONRPCResponse(w http.ResponseWriter, id uint, data interface{}) {
	response := &jsonrpc.RPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  data,
	}

	jsonresponse, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "server internal error", http.StatusInternalServerError)
		logger.ErrorF("marshal response error :%s", err)
		return
	}

	w.WriteHeader(200)

	if _, err := w.Write(jsonresponse); err != nil {
		logger.ErrorF("write response error :%s", err)
	}
}

// Run insight server
func (server *Server) Run() {

	server.router.POST(server.cnf.GetString("insight.extend", "/extend"), server.DipsatchJSONRPC)

	server.router.POST(server.cnf.GetString("insight.proxy", "/"), server.ReverseProxy)

	server.dispatch["balance"] = server.getBalance
	server.dispatch["claim"] = server.getClaim

	logger.Fatal(http.ListenAndServe(
		server.cnf.GetString("insight.listen", ":10332"),
		&loggerHandler{
			handler: server.router,
		},
	))
}

// ReverseProxy reverse proxy handler
func (server *Server) ReverseProxy(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	reverseProxy := httputil.NewSingleHostReverseProxy(server.remote)

	reverseProxy.ServeHTTP(w, r)
}

// DipsatchJSONRPC .
func (server *Server) DipsatchJSONRPC(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	logger.DebugF("call extend api :%s", r.RemoteAddr)

	request, err := makeRPCRequest(r)

	if err != nil {
		makeJSONRPCError(w, 0, JSONRPCParserError, "parse error", nil)
		return
	}

	if method, ok := server.dispatch[request.Method]; ok {
		result, err := method(request.Params.([]interface{}))

		if err != nil {
			makeJSONRPCError(w, request.ID, err.ID, err.Message, result)
		} else {
			makeJSONRPCResponse(w, request.ID, result)
		}

	} else {
		makeJSONRPCError(w, request.ID, JSONRPCMethodNotFound, fmt.Sprintf("method %s not found", request.Method), nil)
	}
}

func (server *Server) getBalance(params []interface{}) (interface{}, *JSONRPCError) {
	if len(params) < 2 {
		return nil, errorf(JSONRPCInvalidParams, "expect address and asset parameters")
	}

	address, ok := params[0].(string)

	if !ok {
		return nil, errorf(JSONRPCInvalidParams, "address parameter must be string")
	}

	asset, ok := params[1].(string)

	if !ok {
		return nil, errorf(JSONRPCInvalidParams, "asset parameter must be string")
	}

	utxos, err := server.utxo.Unspent(address, asset)

	if err != nil {
		return nil, errorf(JSONRPCInnerError, "get %s balance %s err:\n\t%s", address, asset, err)
	}

	return utxos, nil
}

func (server *Server) getClaim(params []interface{}) (interface{}, *JSONRPCError) {

	if len(params) < 1 {
		return nil, errorf(JSONRPCInvalidParams, "expect address parameters")
	}

	address, ok := params[0].(string)

	if !ok {
		return nil, errorf(JSONRPCInvalidParams, "address parameter must be string")
	}

	logger.DebugF("start get claim :%s", address)

	utxos, err := server.utxo.Unclaimed(address)

	if err != nil {
		return nil, errorf(JSONRPCInnerError, "get %s get unclaimed utxo err:\n\t%s", address, err)
	}

	unavailable, available, err := claim.GetUnClaimedGas(utxos, server.blockfee.GetBlocksFee)

	if err != nil {
		return nil, errorf(JSONRPCInnerError, "get %s get unclaimed gas fee err:\n\t%s", address, err)
	}

	claims := make([]*neogo.UTXO, 0)

	for _, utxo := range utxos {
		if utxo.SpentBlock != -1 {
			claims = append(claims, utxo)
		}
	}

	unclaimed := &neogo.Unclaimed{
		Available:   fmt.Sprintf("%.8f", round(available, 8)),
		Unavailable: fmt.Sprintf("%.8f", round(unavailable, 8)),
		Claims:      claims,
	}

	logger.DebugF("finish get claim: %s available: %.8f unavailable: %.8f", address, round(available, 8), round(unavailable, 8))

	return unclaimed, nil
}

func round(f float64, n int) float64 {
	pow10n := math.Pow10(n)
	return math.Trunc(f*pow10n) / pow10n
}

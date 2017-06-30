package api

import (
	"github.com/jiusanzhou/pdf2html/pkg/util"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"os"
	"time"
	"github.com/jiusanzhou/pdf2html/pkg/server/backend"
)

const (
	minSupportedVersion = 1
)

var (
	supportedConvertType = []string{"2simple"}
)

type T struct {
	method string
	url    string
	do     func(ctx context.Context)
}

type response struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
	RequestId string      `json:"request_id"`
	Time      int64       `json:"time"`
}

func getResp(ctx context.Context, data interface{}) *response {
	return &response{0, "success", data, ctx.Value("_request_id").(string), time.Now().UnixNano()}
}

func getError(ctx context.Context, err string, code int) *response {
	return &response{code, err, nil, ctx.Value("_request_id").(string), time.Now().UnixNano()}
}

func registerRouter(app *iris.Application) {
	_api := app.Party("/api/{v:string}", handlerMiddlewareApi)
	{
		_apiConvert := _api.Party("/convert/{type:string}", handlerMiddlewareApiConvert)
		{
			_apiConvert.Get("/{job_id:string}", handlerApiConvertGet)
			_apiConvert.Post("/", handlerApiConvertPost)
		}
	}

	_view := app.Party("/")
	{
		_view.Get("/test/{name:string}", func(ctx context.Context) {
			ctx.HTML("测试页面 -> " + ctx.Params().Get("name"))
		})
	}
}

type ApiHttp struct {
	backend backend.Backend
	config *HttpConfig
}

var apihttp *ApiHttp

func NewApiHttp(c *HttpConfig, bkend backend.Backend)(*ApiHttp, error){
	api := &ApiHttp{
		config: c,
		backend: bkend,
	}

	apihttp = api
	return api, nil
}

func (apihttp *ApiHttp)Serve() error {

	app := iris.New()

	app.AttachLogger(os.Stdout)

	app.StaticWeb(apihttp.config.StaticPath, apihttp.config.StaticSystemPath)

	// add uuid to request
	app.Use(func(ctx context.Context) {
		ctx.Values().Set("_request_id", util.RandId(8))
		ctx.Application().Log("Begin request for path: %s", ctx.Path())
		ctx.Next()
	})

	// Register custom handler for specific http errors.
	app.OnErrorCode(iris.StatusInternalServerError, func(ctx context.Context) {
		// .Values are used to communicate between handlers, middleware.
		errMessage := ctx.Values().GetString("error")
		if errMessage != "" {
			ctx.JSON(getError(ctx, errMessage, codeServerError))
			return
		}

		ctx.JSON(getError(ctx, "(Unexpected) internal server error", codeServerErrorUnknow))
	})

	registerRouter(app)

	// Listen for incoming HTTP/1.x & HTTP/2 clients
	app.Run(iris.Addr(apihttp.config.Addr), iris.WithCharset("UTF-8"))
	return nil
}

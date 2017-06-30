package api

import (
	"github.com/kataras/iris/context"
)

func handlerMiddlewareApi(ctx context.Context) {

	// log this request for api

	// check now supported min version
	if ctx.Params().Get("v") != "v1" {
		ctx.JSON(getError(ctx, "api version not be supported", codeServerError))
	}

	ctx.Next()
}

func handlerMiddlewareApiConvert(ctx context.Context) {

	// log again ? NO

	// check whether we supported this type
	t := ctx.Params().Get("type")
	switch t {
	case "2simple":
	default:
		// no supported this type
		ctx.JSON(getError(ctx, "sorry, we now not supported convert type: "+t, codeServerError))
	}

	ctx.Next()
}

func handlerApiConvertGet(ctx context.Context) {

	// get job or download directly
	id := ctx.Params().Get("job_id")
	if id == "" {
		ctx.JSON(getError(ctx, "you must offer a job id.", codeServerError))
		return
	}

	// get job from backend by job id

	hang := ctx.Params().Get("hang")
	var _hang bool
	if hang != "" {
		_hang = true
	}

	// hang != "", we should wait job is ok to return
	job, _ := apihttp.backend.GetJob(id, _hang)

	if job == nil {
		ctx.JSON(getError(ctx, "no such job with id: "+id, codeServerError))
		return
	}

	download := ctx.Params().Get("download")

	// if download != "", we should download file directly
	// use location in headers
	if download != "" {
		ctx.Header("Location", apihttp.config.StaticPath+"/"+job.OutputFileInfo.Name())
	} else {
		ctx.JSON(getResp(ctx, job))
	}
}

func handlerApiConvertPost(ctx context.Context) {
	ctx.JSON(struct{ A string }{"POST"})
}

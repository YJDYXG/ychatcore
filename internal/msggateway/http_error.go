package msggateway

import "../../pkg/apiresp"

func httpError(ctx *UserConnContext, err error) {
	apiresp.HttpError(ctx.RespWriter, err)
}

package handler

import (
	"crypto-data-pipeline/generated"
	"crypto-data-pipeline/src/handler/example"
	"crypto-data-pipeline/src/handler/pipeline"

	"github.com/labstack/echo"
	"google.golang.org/grpc"
)

func RegisterHttpHandler(e *echo.Echo) {
	exampleH := example.NewExampleHandler(e)
	exampleH.RegisterExampleHandler()
}

func RegisterGrpcHandler(gs *grpc.Server, pipelineCtrl *pipeline.PipelineController) {
	generated.RegisterPipelineServiceServer(gs, pipelineCtrl)
}

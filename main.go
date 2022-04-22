package main

import (
	"log"
	"nft-portal/model"
	"os"

	_ "nft-portal/docs"

	"github.com/DevilsTear/opensea-go-api"
	"github.com/iris-contrib/swagger/swaggerFiles"
	"github.com/iris-contrib/swagger/v12"
	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
)

// @title           NFT PORTAL API
// @version         1.0
// @description     This api serves a hub to complate 3th party nft portal integrations
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  haluk.a.turan@gmail.com

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ac := makeAccessLog()
	defer ac.Close() // Close the underline file.

	app := iris.New()

	apiV1 := app.Party("/api/v1")
	{
		// apiV1.Use(iris.Compression)

		// GET: http://localhost:8080/opensea
		apiV1.Post("/listAssets/{walletAddress}", listAssets)
		apiV1.Post("/listAssetsTest/{walletAddress}", listAssetsTest)
		// // POST: http://localhost:8080/books
		// apiV1.Post("/", create)

	}
	config := swagger.Config{
		URL:         "http://localhost:8080/swagger/doc.json",
		DeepLinking: true,
	}
	swaggerUI := swagger.CustomWrapHandler(&config, swaggerFiles.Handler)
	app.Get("/swagger", swaggerUI)
	app.Get("/swagger/{any:path}", swaggerUI)
	app.Listen(":8080")
}

// listAssets
// @Description fetches the assets by user wallet address
// @Description one of the parameters walletAddress and x-wallet is mandatory and must be provided
// @ID list_assets

// @Accept  json
// @Produce  json
// @Param   walletAddress  			path    string     	false       "Wallet Address"
// @Param   x-wallet  				header  string     	false       "Wallet Address"
// @Param   offset     				query  	int     	false        "Offset"
// @Param   limit      				query  	int     	false        "Limit"
// @Param   asset-params	body  	opensea.GetAssetsParams     	true        "asset params"
// @Success 200 {object} model.Result	"status = true, Code = 200"
// @Failure 400 {object} model.Result "status = true, Code = 400, Message = Provide a valid walletAddress"
// @Router /listAssets/{walletAddress} [post]
func listAssets(ctx iris.Context) {
	walletAddress := ctx.Params().Get("walletAddress")
	cursor, err := ctx.URLParamInt("offset")
	if err != nil || cursor < 0 {
		cursor = 1
	}
	limit, err := ctx.URLParamInt("limit")
	if err != nil || limit < 0 {
		limit = 20
	}

	if walletAddress == "" {
		walletAddress = ctx.GetHeader("x-wallet")
	}

	response := model.Result{
		Paging: model.Paging{
			Offset: &cursor,
			Limit:  &limit,
		},
		APIStatus: model.APIStatus{
			Status:  true,
			Code:    200,
			Message: "Success",
		},
	}

	_, err = opensea.ParseAddress(walletAddress)
	if err != nil {
		// ctx.StopWithError(iris.StatusBadRequest, err)
		// response.APIStatus = model.APIStatus{
		// 	Status:  true,
		// 	Code:    400,
		// 	Message: err.Error(),
		// }
		// ctx.JSON(response)
		// return
		walletAddress = ""
	}

	params := opensea.GetAssetsParams{}
	err = ctx.ReadJSON(&params)
	if err != nil {
		response.APIStatus = model.APIStatus{
			Status:  true,
			Code:    400,
			Message: err.Error(),
		}
		ctx.JSON(response)
		return
	}

	// if params.Owner == "" || params.Owner == opensea.NullAddress {
	// 	params.Owner = opensea.Address(walletAddress)
	// }

	if *params.Limit <= 0 {
		*params.Limit = 20
	}

	if params.Cursor != nil && *params.Cursor <= 0 {
		*params.Cursor = 1
	}

	openseaAPI, err := opensea.NewOpensea(os.Getenv("OPENSEA_API_KEY"))
	if err != nil {
		// ctx.StopWithError(iris.StatusBadRequest, err)
		response.APIStatus = model.APIStatus{
			Status:  true,
			Code:    400,
			Message: err.Error(),
		}
		ctx.JSON(response)
		return
	}

	assets, err := openseaAPI.GetAssets(params)
	if err != nil {
		// ctx.StopWithError(iris.StatusBadRequest, err)
		response.APIStatus = model.APIStatus{
			Status:  true,
			Code:    400,
			Message: err.Error(),
		}
		ctx.JSON(response)
		return
	}

	response.Data = assets

	ctx.JSON(response)
}

// listAssetsTest
// @Description fetches the assets by user wallet address
// @Description one of the parameters walletAddress and x-wallet is mandatory and must be provided
// @ID list_assets

// @Accept  json
// @Produce  json
// @Param   walletAddress  			path    string     	false       "Wallet Address"
// @Param   x-wallet  				header  string     	false       "Wallet Address"
// @Param   offset     				query  	int     	false        "Offset"
// @Param   limit      				query  	int     	false        "Limit"
// @Param   asset-params	body  	opensea.GetAssetsParams     	true        "asset params"
// @Success 200 {object} model.Result	"status = true, Code = 200"
// @Failure 400 {object} model.Result "status = true, Code = 400, Message = Provide a valid walletAddress"
// @Router /listAssetsTest/{walletAddress} [post]
func listAssetsTest(ctx iris.Context) {
	walletAddress := ctx.Params().Get("walletAddress")
	cursor, err := ctx.URLParamInt("offset")
	if err != nil || cursor < 0 {
		cursor = 1
	}
	limit, err := ctx.URLParamInt("limit")
	if err != nil || limit < 0 {
		limit = 20
	}

	if walletAddress == "" {
		walletAddress = ctx.GetHeader("x-wallet")
	}

	response := model.Result{
		Paging: model.Paging{
			Offset: &cursor,
			Limit:  &limit,
		},
		APIStatus: model.APIStatus{
			Status:  true,
			Code:    200,
			Message: "Success",
		},
	}

	_, err = opensea.ParseAddress(walletAddress)
	if err != nil {
		// ctx.StopWithError(iris.StatusBadRequest, err)
		// response.APIStatus = model.APIStatus{
		// 	Status:  true,
		// 	Code:    400,
		// 	Message: err.Error(),
		// }
		// ctx.JSON(response)
		// return
		walletAddress = ""
	}

	params := opensea.GetAssetsParams{}
	err = ctx.ReadJSON(&params)
	if err != nil {
		response.APIStatus = model.APIStatus{
			Status:  true,
			Code:    400,
			Message: err.Error(),
		}
		ctx.JSON(response)
		return
	}

	// if params.Owner == "" || params.Owner == opensea.NullAddress {
	// 	params.Owner = opensea.Address(walletAddress)
	// }

	if params.Limit != nil && *params.Limit <= 0 {
		*params.Limit = 20
	}

	if params.Cursor != nil && *params.Cursor <= 0 {
		*params.Cursor = 1
	}

	openseaAPI, err := opensea.NewOpensea(os.Getenv("OPENSEA_API_KEY"))
	if err != nil {
		// ctx.StopWithError(iris.StatusBadRequest, err)
		response.APIStatus = model.APIStatus{
			Status:  true,
			Code:    400,
			Message: err.Error(),
		}
		ctx.JSON(response)
		return
	}

	assets, err := openseaAPI.GetAssets(params)
	if err != nil {
		// ctx.StopWithError(iris.StatusBadRequest, err)
		response.APIStatus = model.APIStatus{
			Status:  true,
			Code:    400,
			Message: err.Error(),
		}
		ctx.JSON(response)
		return
	}

	response.Data = assets

	ctx.JSON(response)
}

func makeAccessLog() *accesslog.AccessLog {
	// Initialize a new access log middleware.
	ac := accesslog.File("./access.log")
	// Remove this line to disable logging to console:
	ac.AddOutput(os.Stdout)

	// The default configuration:
	ac.Delim = '|'
	ac.TimeFormat = "2006-01-02 15:04:05"
	ac.Async = false
	ac.IP = true
	ac.BytesReceivedBody = true
	ac.BytesSentBody = true
	ac.BytesReceived = false
	ac.BytesSent = false
	ac.BodyMinify = true
	ac.RequestBody = true
	ac.ResponseBody = false
	ac.KeepMultiLineError = true
	ac.PanicLog = accesslog.LogHandler

	// Default line format if formatter is missing:
	// Time|Latency|Code|Method|Path|IP|Path Params Query Fields|Bytes Received|Bytes Sent|Request|Response|
	//
	// Set Custom Formatter:
	ac.SetFormatter(&accesslog.JSON{
		Indent:    "  ",
		HumanTime: true,
	})
	// ac.SetFormatter(&accesslog.CSV{})
	// ac.SetFormatter(&accesslog.Template{Text: "{{.Code}}"})

	return ac
}

func Ptr[T any](v T) *T {
	return &v
}

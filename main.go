package main

import (
	"nft-portal/model"

	opensea "github.com/DevilsTear/opensea-go-api"
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()

	booksAPI := app.Party("/openseaapi")
	{
		booksAPI.Use(iris.Compression)

		// GET: http://localhost:8080/opensea
		booksAPI.Get("/listAssets/{walletAddress}", listAssets)
		// // POST: http://localhost:8080/books
		// booksAPI.Post("/", create)
	}

	app.Listen(":8080")
}

// listAssets
// @Description fetches the assets by user wallet address
// @Description one of the parameters walletAddress and x-wallet is mandatory and must be provided
// @ID get-struct-array-by-string
// @Accept  json
// @Produce  json
// @Param   walletAddress  	path    string     	false        "Wallet Address"
// @Param   x-wallet  		header  string     	false        "Wallet Address"
// @Param   offset     		query  	int     	true        "Offset"
// @Param   limit      		query  	int     	true        "Limit"
// @Success 200 {object} model.APIError	"status = true, Code = 200"
// @Failure 400 {object} model.APIError "status = true, Code = 400, Message = Provide a valid walletAddress"
// @Router /listAssets/{walletAddress} [get]
func listAssets(ctx iris.Context) {
	walletAddress := ctx.Params().Get("walletAddress")
	// offset := ctx.Params().Get("offset")
	// limit := strconv.Atoi(ctx.qet().Get("limit"))

	if walletAddress == "" {
		walletAddress = ctx.GetHeader("x-wallet")
	}

	_, err := opensea.ParseAddress(walletAddress)
	if err != nil {
		// ctx.StopWithError(iris.StatusBadRequest, err)
		ctx.JSON(model.Result{
			APIStatus: model.APIStatus{
				Status:  true,
				Code:    400,
				Message: err.Error(),
			},
		})
		return
	}

	params := opensea.GetAssetsParams{}
	ctx.ReadJSON(params)
	openseaAPI, err := opensea.NewOpenseaRinkeby(walletAddress)
	if err != nil {
		// ctx.StopWithError(iris.StatusBadRequest, err)
		ctx.JSON(model.Result{
			APIStatus: model.APIStatus{
				Status:  true,
				Code:    400,
				Message: err.Error(),
			},
		})
		return
	}

	assets, err := openseaAPI.GetAssets(params)
	if err != nil {
		// ctx.StopWithError(iris.StatusBadRequest, err)
		ctx.JSON(model.Result{
			APIStatus: model.APIStatus{
				Status:  true,
				Code:    400,
				Message: err.Error(),
			},
		})
		return
	}

	ctx.JSON(assets)
}

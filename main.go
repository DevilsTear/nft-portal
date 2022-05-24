package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"nft-portal/model"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "nft-portal/docs"

	moralis "github.com/DevilsTear/moralis-go-client/apis"
	"github.com/DevilsTear/moralis-go-client/models"
	"github.com/DevilsTear/opensea-go-api"
	"github.com/iris-contrib/swagger/swaggerFiles"
	"github.com/iris-contrib/swagger/v12"
	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
)

const (
	// Path to the AWS CA file
	caFilePath = "certificates/rds-combined-ca-bundle.pem"

	// Timeout operations after N seconds
	connectTimeout  = 5
	queryTimeout    = 30
	username        = "fundletestdbusr"
	password        = "fundletestdbpwd"
	clusterEndpoint = "docdb-2022-04-26-23-15-40.citmarx0zdgq.us-east-1.docdb.amazonaws.com:27017"

	// Which instances to read from
	readPreference = "secondaryPreferred"

	connectionStringTemplate = "mongodb://%s:%s@%s/sample-database?tls=true&replicaSet=rs0&readpreference=%s"
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

	// AWSDocDB()
	MoralisTest()

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

func MoralisTest() models.NftMetadataCollection {
	auth := context.WithValue(context.Background(), moralis.ContextAPIKey, moralis.APIKey{
		Key: os.Getenv("MORALIS_API_KEY"),
		// Prefix: "Bearer", // Omit if not necessary.
	})
	// r, err := client.Service.Operation(auth, args)
	config := moralis.NewConfiguration()

	client := moralis.NewAPIClient(config)
	//search?chain=eth&format=decimal&q=ape&filter=name
	searchOptions := moralis.TokenApiSearchNFTsOpts{}
	res, httpRes, err := client.TokenApi.SearchNFTs(auth, "APE", &searchOptions)

	if err != nil {
		fmt.Println(err)
	}

	if httpRes != nil {
		fmt.Println(httpRes)
	}

	return res
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
		PageObject: model.PageObject{
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
		response.Status = true
		response.Code = 400
		response.Message = err.Error()
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
		response.Status = true
		response.Code = 400
		response.Message = err.Error()
		ctx.JSON(response)
		return
	}

	assets, err := openseaAPI.GetAssets(params)
	if err != nil {
		// ctx.StopWithError(iris.StatusBadRequest, err)
		response.Status = true
		response.Code = 400
		response.Message = err.Error()
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
		PageObject: model.PageObject{
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
		response.Status = true
		response.Code = 400
		response.Message = err.Error()
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
		response.Status = true
		response.Code = 400
		response.Message = err.Error()
		ctx.JSON(response)
		return
	}

	assets, err := openseaAPI.GetAssets(params)
	if err != nil {
		// ctx.StopWithError(iris.StatusBadRequest, err)
		response.Status = true
		response.Code = 400
		response.Message = err.Error()
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

func AWSDocDB() {
	connectionURI := fmt.Sprintf(connectionStringTemplate, username, password, clusterEndpoint, readPreference)

	tlsConfig, err := getCustomTLSConfig(caFilePath)
	if err != nil {
		log.Fatalf("Failed getting TLS configuration: %v", err)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI).SetTLSConfig(tlsConfig))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to cluster: %v", err)
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping cluster: %v", err)
	}

	fmt.Println("Connected to DocumentDB!")

	collection := client.Database("sample-database").Collection("sample-collection")

	ctx, cancel = context.WithTimeout(context.Background(), queryTimeout*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159})
	if err != nil {
		log.Fatalf("Failed to insert document: %v", err)
	}

	id := res.InsertedID
	log.Printf("Inserted document ID: %s", id)

	ctx, cancel = context.WithTimeout(context.Background(), queryTimeout*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})

	if err != nil {
		log.Fatalf("Failed to run find query: %v", err)
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		log.Printf("Returned: %v", result)

		if err != nil {
			log.Fatal(err)
		}
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
}

func getCustomTLSConfig(caFile string) (*tls.Config, error) {
	tlsConfig := new(tls.Config)
	certs, err := ioutil.ReadFile(caFile)

	if err != nil {
		return tlsConfig, err
	}

	tlsConfig.RootCAs = x509.NewCertPool()
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs)

	if !ok {
		return tlsConfig, errors.New("Failed parsing pem file")
	}

	return tlsConfig, nil
}

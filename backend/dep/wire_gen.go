// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package dep

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/short-d/app/fw"
	"github.com/short-d/app/modern/mdanalytics"
	"github.com/short-d/app/modern/mdcli"
	"github.com/short-d/app/modern/mddb"
	"github.com/short-d/app/modern/mdenv"
	"github.com/short-d/app/modern/mdhttp"
	"github.com/short-d/app/modern/mdio"
	"github.com/short-d/app/modern/mdlogger"
	"github.com/short-d/app/modern/mdmetrics"
	"github.com/short-d/app/modern/mdnetwork"
	"github.com/short-d/app/modern/mdrequest"
	"github.com/short-d/app/modern/mdrouting"
	"github.com/short-d/app/modern/mdruntime"
	"github.com/short-d/app/modern/mdservice"
	"github.com/short-d/app/modern/mdtimer"
	"github.com/short-d/app/modern/mdtracer"
	"github.com/short-d/short/app/adapter/db"
	"github.com/short-d/short/app/adapter/facebook"
	"github.com/short-d/short/app/adapter/github"
	"github.com/short-d/short/app/adapter/google"
	"github.com/short-d/short/app/adapter/graphql"
	"github.com/short-d/short/app/adapter/kgs"
	"github.com/short-d/short/app/adapter/request"
	"github.com/short-d/short/app/usecase/account"
	"github.com/short-d/short/app/usecase/changelog"
	"github.com/short-d/short/app/usecase/repository"
	"github.com/short-d/short/app/usecase/requester"
	"github.com/short-d/short/app/usecase/service"
	"github.com/short-d/short/app/usecase/url"
	"github.com/short-d/short/app/usecase/validator"
	"github.com/short-d/short/dep/provider"
)

// Injectors from wire.go:

func InjectCommandFactory() fw.CommandFactory {
	cobraFactory := mdcli.NewCobraFactory()
	return cobraFactory
}

func InjectDBConnector() fw.DBConnector {
	postgresConnector := mddb.NewPostgresConnector()
	return postgresConnector
}

func InjectDBMigrationTool() fw.DBMigrationTool {
	postgresMigrationTool := mddb.NewPostgresMigrationTool()
	return postgresMigrationTool
}

func InjectEnvironment() fw.Environment {
	goDotEnv := mdenv.NewGoDotEnv()
	return goDotEnv
}

func InjectGraphQLService(name string, serverEnv fw.ServerEnv, prefix provider.LogPrefix, logLevel fw.LogLevel, sqlDB *sql.DB, graphqlPath provider.GraphQlPath, secret provider.ReCaptchaSecret, jwtSecret provider.JwtSecret, bufferSize provider.KeyGenBufferSize, kgsRPCConfig provider.KgsRPCConfig, tokenValidDuration provider.TokenValidDuration, dataDogAPIKey provider.DataDogAPIKey, segmentAPIKey provider.SegmentAPIKey, ipStackAPIKey provider.IPStackAPIKey) (mdservice.Service, error) {
	timer := mdtimer.NewTimer()
	buildIn := mdruntime.NewBuildIn()
	stdOut := mdio.NewBuildInStdOut()
	client := mdhttp.NewClient()
	http := mdrequest.NewHTTP(client)
	entryRepository := provider.NewEntryRepositorySwitch(serverEnv, stdOut, dataDogAPIKey, http)
	logger := provider.NewLogger(prefix, logLevel, timer, buildIn, entryRepository)
	tracer := mdtracer.NewLocal()
	urlSql := db.NewURLSql(sqlDB)
	userURLRelationSQL := db.NewUserURLRelationSQL(sqlDB)
	retrieverPersist := url.NewRetrieverPersist(urlSql, userURLRelationSQL)
	rpc, err := provider.NewKgsRPC(kgsRPCConfig)
	if err != nil {
		return mdservice.Service{}, err
	}
	keyGenerator, err := provider.NewKeyGenerator(bufferSize, rpc)
	if err != nil {
		return mdservice.Service{}, err
	}
	longLink := validator.NewLongLink()
	customAlias := validator.NewCustomAlias()
	creatorPersist := url.NewCreatorPersist(urlSql, userURLRelationSQL, keyGenerator, longLink, customAlias, timer)
	changeLogSQL := db.NewChangeLogSQL(sqlDB)
	persist := changelog.NewPersist(keyGenerator, timer, changeLogSQL)
	reCaptcha := provider.NewReCaptchaService(http, secret)
	verifier := requester.NewVerifier(reCaptcha)
	cryptoTokenizer := provider.NewJwtGo(jwtSecret)
	authenticator := provider.NewAuthenticator(cryptoTokenizer, timer, tokenValidDuration)
	short := graphql.NewShort(logger, tracer, retrieverPersist, creatorPersist, persist, verifier, authenticator)
	server := provider.NewGraphGophers(graphqlPath, logger, tracer, short)
	service := mdservice.New(name, server, logger)
	return service, nil
}

func InjectRoutingService(name string, serverEnv fw.ServerEnv, prefix provider.LogPrefix, logLevel fw.LogLevel, sqlDB *sql.DB, githubClientID provider.GithubClientID, githubClientSecret provider.GithubClientSecret, facebookClientID provider.FacebookClientID, facebookClientSecret provider.FacebookClientSecret, facebookRedirectURI provider.FacebookRedirectURI, googleClientID provider.GoogleClientID, googleClientSecret provider.GoogleClientSecret, googleRedirectURI provider.GoogleRedirectURI, jwtSecret provider.JwtSecret, bufferSize provider.KeyGenBufferSize, kgsRPCConfig provider.KgsRPCConfig, webFrontendURL provider.WebFrontendURL, tokenValidDuration provider.TokenValidDuration, dataDogAPIKey provider.DataDogAPIKey, segmentAPIKey provider.SegmentAPIKey, ipStackAPIKey provider.IPStackAPIKey) (mdservice.Service, error) {
	timer := mdtimer.NewTimer()
	buildIn := mdruntime.NewBuildIn()
	stdOut := mdio.NewBuildInStdOut()
	client := mdhttp.NewClient()
	http := mdrequest.NewHTTP(client)
	entryRepository := provider.NewEntryRepositorySwitch(serverEnv, stdOut, dataDogAPIKey, http)
	logger := provider.NewLogger(prefix, logLevel, timer, buildIn, entryRepository)
	tracer := mdtracer.NewLocal()
	dataDog := provider.NewDataDogMetrics(dataDogAPIKey, http, timer, serverEnv)
	segment := provider.NewSegment(segmentAPIKey, timer, logger)
	rpc, err := provider.NewKgsRPC(kgsRPCConfig)
	if err != nil {
		return mdservice.Service{}, err
	}
	keyGenerator, err := provider.NewKeyGenerator(bufferSize, rpc)
	if err != nil {
		return mdservice.Service{}, err
	}
	proxy := mdnetwork.NewProxy()
	ipStack := provider.NewIPStack(ipStackAPIKey, http, logger)
	requestClient := request.NewClient(proxy, ipStack)
	instrumentationFactory := request.NewInstrumentationFactory(serverEnv, logger, tracer, timer, dataDog, segment, keyGenerator, requestClient)
	urlSql := db.NewURLSql(sqlDB)
	userURLRelationSQL := db.NewUserURLRelationSQL(sqlDB)
	retrieverPersist := url.NewRetrieverPersist(urlSql, userURLRelationSQL)
	identityProvider := provider.NewGithubIdentityProvider(http, githubClientID, githubClientSecret)
	graphQL := mdrequest.NewGraphQL(http)
	githubAccount := github.NewAccount(graphQL)
	api := github.NewAPI(identityProvider, githubAccount)
	facebookIdentityProvider := provider.NewFacebookIdentityProvider(http, facebookClientID, facebookClientSecret, facebookRedirectURI)
	facebookAccount := facebook.NewAccount(http)
	facebookAPI := facebook.NewAPI(facebookIdentityProvider, facebookAccount)
	googleIdentityProvider := provider.NewGoogleIdentityProvider(http, googleClientID, googleClientSecret, googleRedirectURI)
	googleAccount := google.NewAccount(http)
	googleAPI := google.NewAPI(googleIdentityProvider, googleAccount)
	featureToggleSQL := db.NewFeatureToggleSQL(sqlDB)
	decisionMakerFactory := provider.NewFeatureDecisionMakerFactorySwitch(serverEnv, featureToggleSQL)
	cryptoTokenizer := provider.NewJwtGo(jwtSecret)
	authenticator := provider.NewAuthenticator(cryptoTokenizer, timer, tokenValidDuration)
	userSQL := db.NewUserSQL(sqlDB)
	accountProvider := account.NewProvider(userSQL, timer)
	v := provider.NewShortRoutes(instrumentationFactory, webFrontendURL, timer, retrieverPersist, api, facebookAPI, googleAPI, decisionMakerFactory, authenticator, accountProvider)
	server := mdrouting.NewBuiltIn(logger, tracer, v)
	service := mdservice.New(name, server, logger)
	return service, nil
}

// wire.go:

var authSet = wire.NewSet(provider.NewJwtGo, provider.NewAuthenticator)

var observabilitySet = wire.NewSet(wire.Bind(new(fw.StdOut), new(mdio.StdOut)), wire.Bind(new(fw.Logger), new(mdlogger.Logger)), wire.Bind(new(fw.Metrics), new(mdmetrics.DataDog)), wire.Bind(new(fw.Analytics), new(mdanalytics.Segment)), wire.Bind(new(fw.Network), new(mdnetwork.Proxy)), mdio.NewBuildInStdOut, provider.NewEntryRepositorySwitch, provider.NewLogger, mdtracer.NewLocal, provider.NewDataDogMetrics, provider.NewSegment, mdnetwork.NewProxy, request.NewClient, request.NewInstrumentationFactory)

var githubAPISet = wire.NewSet(provider.NewGithubIdentityProvider, github.NewAccount, github.NewAPI)

var facebookAPISet = wire.NewSet(provider.NewFacebookIdentityProvider, facebook.NewAccount, facebook.NewAPI)

var googleAPISet = wire.NewSet(provider.NewGoogleIdentityProvider, google.NewAccount, google.NewAPI)

var keyGenSet = wire.NewSet(wire.Bind(new(service.KeyFetcher), new(kgs.RPC)), provider.NewKgsRPC, provider.NewKeyGenerator)

var featureDecisionSet = wire.NewSet(wire.Bind(new(repository.FeatureToggle), new(db.FeatureToggleSQL)), db.NewFeatureToggleSQL, provider.NewFeatureDecisionMakerFactorySwitch)

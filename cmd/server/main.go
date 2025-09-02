package main

import (
	"cmp"
	"database/sql"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/pflag"
	"github.com/vilasle/gokeep/internal/encryption"
	"github.com/vilasle/gokeep/internal/logger"
	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/repository/server/postgres"
	"github.com/vilasle/gokeep/internal/server"
	"github.com/vilasle/gokeep/internal/service/auth"
	"github.com/vilasle/gokeep/internal/service/private"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	logger.Make(os.Stdout, logger.DebugLevel)
	logger.Info("metadata",
		"version", cmp.Or(buildVersion, "N/A"),
		"date", cmp.Or(buildDate, "N/A"),
		"commit", cmp.Or(buildCommit, "N/A"))

	masterKeyPath := "master.key"

	var serverCertificatePath, serverKeyPath, addr, dbUrl, jwtKeyPath string
	var showHelp bool
	pflag.StringVarP(&addr, "addr", "a", ":9092", "Listen port")
	pflag.StringVarP(&serverCertificatePath, "certificate", "c", "", "Server certificate file path")
	pflag.StringVarP(&serverKeyPath, "private-key", "k", "", "Server key file path")

	pflag.StringVarP(&dbUrl, "db-url", "d", "",
		"Database URL. Format postgres://user:password@127.0.0.1:5432/gokeep?sslmode=disable")

	pflag.StringVarP(&jwtKeyPath, "salt", "s", "", "JWT key path")
	pflag.BoolVarP(&showHelp, "help", "h", false, "Show help")
	pflag.Parse()

	if showHelp {
		pflag.Usage()
		os.Exit(0)
	}

	if dbUrl == "" {
		logger.Error("database url is not set")
		os.Exit(1)
	}

	logger.Debug("running arguments",
		"cert", serverCertificatePath,
		"key", serverKeyPath,
		"addr", addr,
		"db-url", dbUrl,
		"jwt-key", jwtKeyPath)

	_, err := os.Stat(masterKeyPath)
	if os.IsNotExist(err) {
		masterKey, err := encryption.GenerateNewAESKey()
		if err != nil {
			logger.Error("can not generate master key", "err", err)
			os.Exit(1)
		}
		content := []byte(masterKey.JSON())
		err = os.WriteFile(masterKeyPath, content, 0644)
		if err != nil {
			logger.Error("can not write master key", "err", err, "path", masterKeyPath)
			os.Exit(1)
		}
	}

	content, err := os.ReadFile(masterKeyPath)
	if err != nil {
		if err != os.ErrNotExist {
			logger.Error("can not read master key file", "err", err, "path", masterKeyPath)
			os.Exit(1)
		}
	}

	masterKey, err := encryption.NewAESKeyFromJSON(content)
	if err != nil {
		logger.Error("can not create master key from file", "err", err)
		os.Exit(1)
	}

	conn, err := sql.Open("pgx/v5", dbUrl)
	if err != nil {
		logger.Error("can not create database connection", "err", err)
		os.Exit(1)
	}
	defer conn.Close()

	userRepository, err := postgres.NewUserRepository(conn)
	if err != nil {
		logger.Error("can not create user repository", "err", err)
		os.Exit(1)
	}

	sessionRepository, err := postgres.NewSessionRepository(conn)
	if err != nil {
		logger.Error("can not create session repository", "err", err)
		os.Exit(1)
	}

	privateRepository, err := postgres.NewPrivateDataRepository(conn)
	if err != nil {
		logger.Error("can not create private data repository", "err", err)
		os.Exit(1)
	}

	collector := repositoryCollector{
		users:   userRepository,
		private: privateRepository,
	}

	manager := model.NewModelManager(collector)

	jwt, err := os.ReadFile(jwtKeyPath)
	if err != nil {
		logger.Error("can not read jwt key file", "err", err)
		os.Exit(1)
	}

	authSvc := auth.NewAuthService(manager, sessionRepository, jwt)

	loginPasswordSvc := private.NewUsepassService(manager, masterKey)
	bankCardSvc := private.NewBankCardService(manager, masterKey)
	textSvc := private.NewTextService(manager, masterKey)
	binarySvc := private.NewBinaryDataService(manager, masterKey)

	config := server.Config{
		Addr:                 addr,
		AuthService:          authSvc,
		LoginPasswordService: loginPasswordSvc,
		BankCardService:      bankCardSvc,
		TextDataService:      textSvc,
		BinaryDataService:    binarySvc,
		CertificatePath:      serverCertificatePath,
		PrivateKeyPath:       serverKeyPath,
	}

	srv, err := server.NewServer(config, server.WithLogger)
	if err != nil {
		logger.Error("can not create server", "err", err)
		os.Exit(1)
	}
	if err := srv.Listen(); err != nil {
		logger.Error("can not start server", "err", err)
		os.Exit(1)
	}
}

type repositoryCollector struct {
	users   model.UserRepository
	private model.PrivateDataRepository
}

func (c repositoryCollector) User() model.UserRepository {
	return c.users
}

func (c repositoryCollector) Private() model.PrivateDataRepository {
	return c.private
}

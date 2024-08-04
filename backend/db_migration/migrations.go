package ignore

//
//
//import (
//	"github.com/spf13/viper"
//
//	"github.com/golang-migrate/migrate/v4"
//	_ "github.com/golang-migrate/migrate/v4/database/postgres"
//	_ "github.com/golang-migrate/migrate/v4/source/file"
//	"go.uber.org/zap"
//)
//
//const dbMigrationsPath = "./db_migration/migrations"
//
//func RunDatabaseMigrations(logger *zap.Logger) error {
//	var err error
//	appMigrate, err := migrate.New("file://"+dbMigrationsPath, viper.GetString("postgres.url"))
//	if err != nil {
//		logger.Error("migrations error", zap.Error(err))
//		panic(err)
//	}
//	err = appMigrate.Up()
//	if err != nil && err != migrate.ErrNoChange {
//		logger.Error("migrations error", zap.Error(err))
//		return err
//	}
//	logger.Info("migrations successful")
//	return nil
//}

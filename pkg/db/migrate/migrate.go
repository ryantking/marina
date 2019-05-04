package migrate

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/ryantking/marina/pkg/config"

	"github.com/golang-migrate/migrate"
	log "github.com/sirupsen/logrus"

	// Driver for file migrations
	_ "github.com/golang-migrate/migrate/source/file"
)

// GetLatest reads the migration directory and finds the latest migration
func GetLatest() (int, error) {
	files, err := ioutil.ReadDir(config.Get().DB.MigrationsDir)
	if os.IsNotExist(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	latest := 0
	for _, f := range files {
		versionStr := strings.Split(f.Name(), "_")[0]
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			return 0, err
		}

		if version > latest {
			latest = version
		}
	}

	return latest, nil
}

// Up will migrate the database to the latest version
func Up(adapterName string, dsn string) error {
	m, err := getMigration(adapterName, dsn)
	if err != nil {
		return err
	}

	defer m.Close()
	return m.Up()
}

// Down will migrate the database all the way down
func Down(adapterName string, dsn string) error {
	m, err := getMigration(adapterName, dsn)
	if err != nil {
		return err
	}

	defer m.Close()
	return m.Down()
}

// To migrates the database to the specified version
func To(adapterName string, dsn string, n uint) error {
	m, err := getMigration(adapterName, dsn)
	if err != nil {
		return err
	}

	defer m.Close()
	return m.Migrate(n)
}

// Start the migration process if necessary
func Start() {
	cfg := config.Get()
	m, err := getMigration(cfg.DB.Type, cfg.DB.DSN)
	if err != nil {
		log.WithError(err).Fatal("error getting migration adapter")
	}
	defer m.Close()
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.WithError(err).Fatal("error checking database version")
	}
	if dirty {
		log.Fatal("dirty database version")
	}

	latest, err := GetLatest()
	if err != nil {
		log.WithError(err).Fatal("error getting latest database version")
	}

	if version < uint(latest) {
		log.Infof("migrating database from %d to %d", version, latest)
		if err := m.Up(); err != nil {
			log.WithError(err).Error("error migrating database")
		}

		log.Info("migration succeeded")
		return
	}

	log.Info("database version at latest, skipping migration")
}

func getMigration(adapterName, dsn string) (*migrate.Migrate, error) {
	cfg := config.Get()
	srcURL := fmt.Sprintf("file://%s", cfg.DB.MigrationsDir)
	dbURL := fmt.Sprintf("%s://%s", adapterName, dsn)
	m, err := migrate.New(srcURL, dbURL)
	if err != nil {
		return nil, err
	}

	return m, nil
}

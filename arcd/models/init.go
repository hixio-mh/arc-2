/*
 * Arc - Copyleft of Simone 'evilsocket' Margaritelli.
 * evilsocket at protonmail dot com
 * https://www.evilsocket.net/
 *
 * See LICENSE.
 */
package models

import (
	"github.com/evilsocket/arc/arcd/config"
	"github.com/evilsocket/arc/arcd/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var db *gorm.DB

func expand(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		path = strings.Replace(path, "~", usr.HomeDir, -1)
	}
	return filepath.Abs(path)
}

func Setup() (created bool, err error) {
	if config.Conf.Database, err = expand(config.Conf.Database); err != nil {
		return false, err
	}

	if _, err = os.Stat(config.Conf.Database); os.IsNotExist(err) {
		created = true
		log.Infof("Creating database %s ...", log.Bold(config.Conf.Database))
	} else {
		created = false
		log.Infof("Loading database %s ...", log.Bold(config.Conf.Database))
	}

	if db, err = gorm.Open("sqlite3", config.Conf.Database); err != nil {
		return false, err
	}

	log.Debugf("Migrating models ...")

	db.AutoMigrate(&Store{})
	db.AutoMigrate(&Record{})

	return created, nil
}

func Save(obj interface{}) error {
	return db.Save(obj).Error
}

func Create(obj interface{}) error {
	return db.Create(obj).Error
}

func Delete(obj interface{}) error {
	return db.Delete(obj).Error
}
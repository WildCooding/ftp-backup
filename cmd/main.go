package main

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/robfig/cron/v3"
	"github.com/wildcooding/ftp-backup/config"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		panic(err)
	}

	printConfig(c)

	ftp, err := connectToFtpSerer(c.Host, c.Username, c.Password, c.Port)

	if err != nil {
		panic(err.Error())
	}

	jobs := cron.New()
	jobs.AddFunc(c.CrontabInterval, func() { backup(ftp, c) })
	jobs.Start()

	for {
		time.Sleep(10 * time.Minute)
	}
}

func backup(ftp *ftp.ServerConn, config *config.Config) {
	for _, dir := range config.Directorys {
		go processDirectory(ftp, dir, config.DestinationDirectory)
	}
}

func processDirectory(ftp *ftp.ServerConn, dir string, destinationDirectory string) {
	log.Print("Processing dir: ", dir)
	filename, err := zipFolder(dir)
	if err != nil {
		log.Print("Error creating zip for directory ", dir)
		return
	}

	err = uploadFile(ftp, destinationDirectory, filename)
	if err != nil {
		log.Print("Error uploading file ", filename)
		return
	}

	deleteFile(filename)
}

func connectToFtpSerer(host string, username string, password string, port string) (*ftp.ServerConn, error) {
	c, err := ftp.Dial(host+":"+port, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	err = c.Login(username, password)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func zipFolder(folderToZip string) (string, error) {
	lastPartOfDir := strings.Split(folderToZip, "/")[len(strings.Split(folderToZip, "/"))-1]
	filename := "backup-" + lastPartOfDir + "-" + time.Now().Format("20060102150405") + ".zip"

	zipFile, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return filename, filepath.Walk(folderToZip, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Der relative Pfad zum Dateinamen im ZIP-Archiv
		relativePath, err := filepath.Rel(folderToZip, filePath)
		if err != nil {
			return err
		}

		// Wenn es sich um einen Ordner handelt, füge einen neuen Ordner im ZIP-Archiv hinzu
		if info.IsDir() {
			_, err := zipWriter.Create(relativePath + "/")
			return err
		}

		// Andernfalls, füge die Datei zum ZIP-Archiv hinzu
		zipFileWriter, err := zipWriter.Create(relativePath)
		if err != nil {
			return err
		}

		// Kopiere den Inhalt der Datei in die ZIP-Datei
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(zipFileWriter, file)
		return err
	})
}

func uploadFile(ftp *ftp.ServerConn, destinationDirecory string, filename string) error {
	data, err := os.Open(filename)
	if err != nil {
		return err
	}

	err = ftp.Stor(destinationDirecory+"/"+filename, data)
	return err
}

func deleteFile(filename string) {
	os.Remove(filename)
}

func printConfig(config *config.Config) {
	log.Print("Host:", config.Host)
	log.Print("Username:", config.Username)
	log.Print("Port:", config.Port)
	log.Print("Interval:", config.CrontabInterval)
	log.Print("Dirs:", config.Directorys)
	log.Print("Dest:", config.DestinationDirectory)
}

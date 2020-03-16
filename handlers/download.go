package handlers

import (
	"bufio"
	"fmt"
	"github.com/google/logger"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var log * logger.Logger

func Download(cmd *cobra.Command, args []string)  {
	log = logger.Init("DownloadLogger", true, true, ioutil.Discard)

	if err := cmd.Flags().Parse(args); err != nil {
		log.Fatal(err)
	}

	patchUri, _ := cmd.Flags().GetString("patch-hive")

	if strings.HasPrefix(patchUri, "http") {
		log.Infof("downloading external patch list from: %v", patchUri)
		if err :=  downloadPatchList("PatchHive.txt", patchUri); err != nil {
			log.Fatal(err)
		} else {
			patchUri = "PatchHive.txt"
		}
	}

	var resources []string
	file, err := os.Open(patchUri)

	if err != nil {
		log.Fatal(err)
	}
	gdpRoot := "nil"
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lineTxt := scanner.Text()
		if strings.HasPrefix(lineTxt,"#ROOT") {
			gdpRoot = strings.Split(scanner.Text(),"\t")[1]
			log.Infof("gdp root %v:", gdpRoot)
		}
		if strings.HasPrefix(lineTxt,"#PATCH") {
			resources = append(resources, strings.Split(scanner.Text(),"\t")[2])
		}
	}

	if gdpRoot == "nil" {
		log.Fatalf("no gdp #ROOT found")
	}

	destinationPath, _ := cmd.Flags().GetString("destination-path")
	destinationPath, err = filepath.Abs(destinationPath)
	overwrite, _ := cmd.Flags().GetBool("overwrite")

	if _, err := os.Stat(destinationPath); os.IsNotExist(err) {
		if err := os.Mkdir(destinationPath,0700); err != nil {
			log.Fatal(err)
		}
	}

	wg := new(sync.WaitGroup)

	for _, r := range resources {
		wg.Add(1)
		dst, _ := filepath.Abs(fmt.Sprintf("%v/%v", destinationPath, r))
		gdpUri := fmt.Sprintf("%v/%v", gdpRoot, r)
		fmt.Println(gdpUri)
		go downloadFile(wg, dst, gdpUri, overwrite)
	}
	wg.Wait()
}

func  downloadPatchList(filePath string, url string) error {
	// Get the data
	log.Infof("downloading file %v  at destination %v",url, filePath)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func downloadFile(wg * sync.WaitGroup, filePath string, url string, overwrite bool) {
	defer wg.Done()
	if _, err := os.Stat(filePath); err == nil {
		// file exists
		if !overwrite {
			log.Infof("file %v already exists, skipping.", filePath)

			return
		}

	} else if os.IsNotExist(err) {
		//if os.IsNotExist(err) {
			// Get the data
			resp, err := http.Get(url)
			log.Infof("downloading file %v  at destination %v",url, filePath)

			if err != nil {
				log.Error(err)
				return
			}

			defer resp.Body.Close()

			if out, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0700); err != nil {
				log.Error(err)
			} else {

				defer out.Close()

				// Write the body to file
				if _, err = io.Copy(out, resp.Body); err != nil {
					log.Error(err)
				}
			}
		//}
	} else {
		log.Error(err)
	}
}
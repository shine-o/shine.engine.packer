package handlers

import (
	"bufio"
	"context"
	"fmt"
	"github.com/google/logger"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var log * logger.Logger

func init()  {
	log = logger.Init("download logger", true, true, ioutil.Discard)
}

func Download(cmd *cobra.Command, args []string)  {
	runtime.GOMAXPROCS(runtime.NumCPU())

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

	destinationPath, _ := cmd.Flags().GetString("destination")
	destinationPath, err = filepath.Abs(destinationPath)
	overwrite, _ := cmd.Flags().GetBool("overwrite")

	if _, err := os.Stat(destinationPath); os.IsNotExist(err) {
		if err := os.Mkdir(destinationPath,0700); err != nil {
			log.Fatal(err)
		}
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	dc := &DownloadCmd{
		Overwrite:   overwrite,
		DownloadPath: destinationPath,
		ExtractPath: destinationPath,
	}
	wg := new(sync.WaitGroup)

	for _, r := range resources {
		wg.Add(1)
		uri := fmt.Sprintf("%v/%v", gdpRoot, r)
		go dc.download(wg, uri, r)
		time.Sleep(time.Millisecond * 200)
	}
	wg.Wait()
}

type DownloadCmd struct {
	Overwrite bool
	Extract bool
	DownloadPath string
	ExtractPath string
}

func (dc *DownloadCmd) download(wg * sync.WaitGroup, uri, fileName string) {
	log.Infof("downloading file from %v", uri)
	defer wg.Done()

	path, _ := filepath.Abs(fmt.Sprintf("%v/%v", dc.DownloadPath, fileName))
	_, err := os.Stat(path)
	if err == nil {
		if !dc.Overwrite {
			log.Infof("file %v already exists, skipping.", path)
			return
		}
	}
	if os.IsNotExist(err) {
		r, err := http.Get(uri)
		defer r.Body.Close()

		if err != nil {
			log.Error(err)
			return
		}

		if r.StatusCode >= 400  {
			log.Errorf("bad request  with status code %v for file %v", r.StatusCode, uri)
			return
		}

		dc.persist(path, r.Body)
	} else {
		log.Error(err)
	}
}

func (dc *DownloadCmd) persist(path string, body io.Reader) {
	defer logOperation(fmt.Sprintf("persisted file  %v", path))
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0700)
	defer out.Close()
	if err != nil {
		log.Error(err)
		return
	}
	if _, err = io.Copy(out, body); err != nil {
		log.Error(err)
		return
	}
}

func logOperation(msg string){
	log.Info(msg)
}

func  downloadPatchList(filePath string, url string) error {
	// Get the data
	defer logOperation(fmt.Sprintf("downloading file %v  at destination %v",url, filePath))
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
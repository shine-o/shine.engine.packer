package handlers

import (
	"encoding/binary"
	"fmt"
	"github.com/go-restruct/restruct"
	"github.com/google/logger"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type ExtractCmd struct {
	Path string
}

type GDP struct {
	Type string `struct:"[3]byte"`
	UnkData1 [5]byte
	Version string `struct:"[260]byte"`
	Num1 uint32
	NumOfFiles uint32
	UnkData3 [40]byte
	Files []File `struct-size:"NumOfFiles - 1"`
}


type File struct {
	Num0 uint64
	Name string `struct:"[264]byte"`
	Offset uint64
	Size uint64
	Num3 uint64
	UnkData4 [20]byte
	Data []byte
}

func Extract(cmd *cobra.Command, args []string) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	restruct.EnableExprBeta()

	log = logger.Init("extract logger", true, true, ioutil.Discard)

	if err := cmd.Flags().Parse(args); err != nil {
		log.Fatal(err)
	}

	sourcePath, _ := cmd.Flags().GetString("source")
	destinationPath, _ := cmd.Flags().GetString("destination")

	absSourcePath, err := filepath.Abs(sourcePath)
	if err != nil {
		log.Fatal(err)
	}
	absDestinationPath, err := filepath.Abs(destinationPath)

	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(absDestinationPath)
	if err == os.ErrNotExist {
		err := os.Mkdir(absDestinationPath, 0700)
		if err != nil {
			log.Fatal(err)
		}
	} else if err == os.ErrExist {
		err := os.RemoveAll(absDestinationPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	ec := ExtractCmd{
		Path:      absDestinationPath,
	}

	files, err := ioutil.ReadDir(absSourcePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fPath := fmt.Sprintf("%v/%v", absSourcePath, f.Name())
		absFPath, err := filepath.Abs(fPath)
		if err != nil {
			log.Fatal(err)
		}
		ec.extract(absFPath, f.Name())
	}
}

func (ec *ExtractCmd) extract(path string, folderName string) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var gf GDP

	err = restruct.Unpack(data, binary.LittleEndian, &gf)

	if err != nil {
		log.Errorf("%v %v", folderName, err)
	}
	eof := false
	for _, f := range gf.Files {
		if eof {
			return
		}
		directories :=  ec.Path + "/" + folderName + "/"
		log.Info(folderName)
		segments := strings.Split(f.Name, "\\")
		for i := 1; i < len(segments)-1; i++ {
			if i != 1 {
				directories += "/"
			}
			directories += segments[i]
		}

		absPath, err := filepath.Abs(directories)

		if err != nil {
			log.Error(err)
			return
		}

		//_, err = os.Stat(absPath)
		//if err == nil {
		//	log.Infof("directory path %v already exists, skipping.", absPath)
		//	return
		//}

		err = os.MkdirAll(absPath, 0700)
		if err != nil {
			log.Error(err)
			return
		}

		directories += "/"
		directories += segments[len(segments)-1]

		fileAbsPath, err := filepath.Abs(directories)
		if err != nil {
			log.Error(err)
			return
		}

		_, err = os.Stat(fileAbsPath)

		if err == nil {
			log.Infof("file %v already exists, skipping.", fileAbsPath)
			return
		}

		file, err := os.OpenFile(fileAbsPath, os.O_RDONLY|os.O_CREATE, 0700)

		if err != nil {
			fmt.Println(err)
		}

		var b []byte
		if f.Offset+f.Size > uint64(len(data)) {
			b = append(b, data[f.Offset:]...)
			// files are listed but there is no data to read from
			// this means the full patch is the next gdp
			// no clue why this was made like this x.x
			eof = true
		} else {
			b = append(b, data[f.Offset:f.Offset+f.Size]...)
		}

		_, err = file.Write(b)
		if err != nil {
			fmt.Println(err)
		}
	}
}
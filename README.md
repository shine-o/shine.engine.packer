![](shine.png)
# Shine Engine Packer

> utility for handling .gdp files.

CDNs: 

* US: http://patch.cdn.gamigo.com/fous/gdp/PatchHive.txt

* ES: http://patch.cdn.gamigo.com/fo/es/PatchHive.txt

* DE: http://patch.cdn.gamigo.com/fo/de/PatchHive.txt

* FR: http://patch.cdn.gamigo.com/fo/fr/PatchHive.txt
    
  ## build
        #or packer.exe if on windows 
        $ go build -o packer
    
    
  
  ## packer download
  
  Download .gdp files which contain patch data
  
  ### Synopsis
  
  Download .gdp files which contain patch data
  
  ```
  packer download [flags]
  ```
  
  ### Options
  
  ```
        --destination string   Path on system where downloaded .gdp files should be persisted (default "./downloaded")
    -h, --help                 help for download
        --overwrite            Download and replace if the file already exists
        --patch-hive string    list of .gdp files, e.g: PatchHive.txt, http://patch.cdn.gamigo.com/fous/gdp/PatchHive.txt (default "http://patch.cdn.gamigo.com/fo/fr/PatchHive.txt")
  ```
  
  ### Options inherited from parent commands
  
  ```
        --config string   config file (default is $HOME/.packer.yaml)
  ```

  
  ## packer extract
  
  Extract .gdp files into separate folders
  
  ### Synopsis
  
  Extract .gdp files into separate folders
  
  ```
  packer extract [flags]
  ```
  
  ### Options
  
  ```
        --destination string   path where extracted data should be stored (default is ./extracted) (default "./extracted")
    -h, --help                 help for extract
        --source string        path where raw .gdp files are stored (default is ./downloaded) (default "./downloaded")
  ```
  
  ### Options inherited from parent commands
  
  ```
        --config string   config file (default is $HOME/.packer.yaml)
  ```

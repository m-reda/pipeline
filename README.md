# Pipeline
Pipeline is node-based automation server.

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/6d74fab6f1f148c68dce4d285339f5d2)](https://www.codacy.com/app/m-reda/pipeline?utm_source=github.com&utm_medium=referral&utm_content=m-reda/pipeline&utm_campaign=badger)
[![Build Status](https://travis-ci.org/m-reda/pipeline.svg?branch=master)](https://travis-ci.org/m-reda/pipeline)
[![Go Report](https://goreportcard.com/badge/github.com/m-reda/pipeline)](https://goreportcard.com/report/github.com/m-reda/pipeline)
[![Docker Build Statu](https://img.shields.io/badge/docker-pipeline-green.svg?style=flat)](https://hub.docker.com/r/mreda/pipeline/)


Video
-----
[![Pipeline Video](http://i.imgur.com/GRMRsQx.png)](https://youtu.be/A_upilYkrok)

Features
--------
- Node-based tasks
- Ready to use units
- Easy to defined new units
- Remote build trigger
- Scheduled builds
- Real-time build logs
- Elegant user interface
- Responsive UI

Operating system
----------------
`Linux` and `macOS` are supported, for windows you can use docker.

Installation (Pre-release)
----------------
`Note: this release is non-production ready.`
#### Docker

	$ docker run -d --name pipeline -p 8080:80 mreda/pipeline

#### wget 

	$ wget https://github.com/m-reda/pipeline/releases/download/0.1/pipeline-linux.zip
	$ unzip pipeline.zip && cd pipeline
	$ PORT=8080 ./pipeline

#### Download
- v0.1-alpha
	- [pipeline-linux.zip](https://github.com/m-reda/pipeline/releases/download/v0.1/pipeline-alpha-linux.zip)
	- [pipeline-macOS.zip](https://github.com/m-reda/pipeline/releases/download/v0.1/pipeline-alpha-macOS.zip)



Built-in Units
-------------
- Filesystem
  - Copy
  - Move
  - Remove
  - Make file
  - Make directory
- Git
  - Init
  - Add
  - Commit
  - Push
  - Clone
  - Checkout
  - Merge
  - Add remote
  - Pull
- FTP
  - List contents
  - Make directory
  - Remove directory
  - Upload file
  - Remove file
  - Rename
- General
  - Run command
  - Send email
  - Sleep x seconds
  - Request URL
  - SSH Command
  - Zip / unzip

Add New Unit
--------
1. Create new directory under data/units `./data/units/[unit-id]`

2. Create unit definition file `data/units/[unit-id]/unit.js`
```json
{
	"ID": "fs_copy",
	"Name": "FS Copy",
	"Group": "filesystem",
	"Version ": "0.0.1",
	"Creator": "Mahmoud Reda",
	"Command": "bin:/builtin fs copy {source} {destination}",
	"Inputs": {
		"source": "Source",
		"destination": "Destination"
	},
	"Outputs": {
		"destination": "Destination"
	},
	"Setting": {
		"flag": {"Name": "Flag Name", "Type": "text", "Value": ""}
	}
}
```
- The inputs keys must match the names in the command:
```json
{
	"Command": "bin:/filesystem delete {file_path}",
	"Inputs": {
		"file_path": "File Path"
	}
}
```

- Command can be global or prefixed with `unit:` or `bin:`
	- `bin:` equals ./data/units/bin
	- `unit:` equals ./data/units/[unit-id]
	

- The setting values will passed to the unit command as flags.
- Unit directory can contain custom scripts.
- Each output should be printed in a separate line staring with output's key:
	
```
output1:sometext
output2:/path/to/file
output3:{"key":"value"}
```

TODO
----
- [ ] Authentication
- [ ] Upload unit

Custom Build
------------
You can customize the build setting from Makefile under `release` command, and rebuild using:
```
$ make release
```
the new build will be under `bin` directory, or you can build new docker image using:
```
$ make docker
```


Community
---------
Contributions, questions, and comments are welcomed and encouraged.

The Node Editor
----------
I'm using my library [Linker](https://github.com/m-reda/linker).

[![Linker](https://github.com/m-reda/linker/raw/master/dist/logo.png)](https://github.com/m-reda/linker)


Dependencies
------------
[mux](https://github.com/gorilla/mux) / 
[websocket](https://github.com/gorilla/websocket) / 
[cron](https://github.com/robfig/cron) / 
[cli](https://github.com/urfave/cli) /
[ftp](https://github.com/jlaffaye/ftp) /
[go.uuid](https://github.com/satori/go.uuid) /
[testify](https://github.com/stretchr/testify)

License
-------
This code is distributed under the MIT license found in the [LICENSE](https://github.com/m-reda/pipeline/LICENSE) file.

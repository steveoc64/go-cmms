all: sassgen templegen app-assets appjs sv run

build: sassgen templegen app-assets appjs sv 

help: 
	# sassgen    - make SASS files
	# templegen  - make Templates
	# app-assets - make Asset copy to dist	
	# appjs      - make Frontend app
	# sv         - make Server
	# run        - run  Server

clean:	
	# Delete existing build
	@mplayer -quiet audio/trash-empty.oga 2> /dev/null > /dev/null &
	rm -rf dist

sassgen: dist/public/css/app.css

dist/public/css/app.css: scss/*
	@mplayer -quiet audio/attention.oga 2> /dev/null > /dev/null
	@mkdir -p dist/public/css
	cd scss && node-sass --output-style compressed app.sass ../dist/public/css/app.css

templegen: app/template.go 

app/template.go: templates/*.tmpl 	
	@mplayer -quiet audio/attention.oga 2> /dev/null > /dev/null
	temple build templates app/template.go --package main

app-assets: dist/assets.log

dist/assets.log: assets/index.html assets/img/*  assets/fonts/* assets/css/*
	@mplayer -quiet audio/attention.oga 2> /dev/null > /dev/null
	@mkdir -p dist/public/css dist/public/font dist/public/js
	cp assets/index.html dist/public
	cp -R assets/img dist/public
	cp -R assets/css dist/public
	cp -R assets/fonts dist/public
	cp -R assets/js dist/public
	#cp bower_components/normalize.css/normalize.css dist/public/css
	cp server/config.json dist	
	@date > dist/assets.log

appjs: dist/public/app.js

dist/public/app.js: app/*.go shared/*.go
	@mplayer -quiet audio/frontend-compile.ogg 2> /dev/null > /dev/null &
	@mkdir -p dist/public/js
	@gosimple app
	cd app && gopherjs build *.go -o ../dist/public/app.js -m
	@ls -l dist/public/app.js

remake: 
	@mplayer -quiet audio/server-compile.oga 2> /dev/null > /dev/null &
	rm -f dist/cmms-server
	@gosimple server
	cd server && go build -o ../dist/cmms-server
	@ls -l dist/cmms-server

sv: dist/cmms-server 

dist/cmms-server: server/*.go shared/*.go
	@mplayer -quiet audio/server-compile.oga 2> /dev/null > /dev/null &
	@gosimple server
	cd server && go build -o ../dist/cmms-server
	@ls -l dist/cmms-server

run: 
	./terminate
	@mplayer -quiet audio/running.oga 2> /dev/null > /dev/null &
	@cd dist && ./cmms-server
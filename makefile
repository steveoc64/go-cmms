all: sassgen templegen app-assets appjs sv run

help: 
	# sassgen    - make SASS files
	# templegen  - make Templates
	# app-assets - make Asset copy to dist	
	# appjs      - make Frontend app
	# sv         - make Server
	# run        - run  Server

clean:	
	# Delete existing build
	rm -rf dist

sassgen: dist/public/css/app.css

dist/public/css/app.css: assets/scss/*
	@mkdir -p dist/public/css
	cd assets/scss && node-sass app.sass ../../dist/public/css/app.css

templegen: app/template.go 

app/template.go: templates/*.tmpl 	
	temple build templates app/template.go --package main

app-assets: dist/assets.log

dist/assets.log: assets/index.html assets/img/*  assets/fonts/* assets/css/*
	@mkdir -p dist/public/css dist/public/font dist/public/js
	cp assets/index.html dist/public
	cp -R assets/img dist/public
	cp -R assets/css dist/public
	cp -R assets/fonts dist/public
	cp -R assets/js dist/public
	cp bower_components/normalize.css/normalize.css dist/public/css
	cp server/config.json dist	
	@date > dist/assets.log

appjs: dist/public/app.js

dist/public/app.js: app/*.go
	@mkdir -p dist/public/js
	cd app && gopherjs build *.go -o ../dist/public/app.js -m
	@ls -l dist/public/app.js
	@mplayer -quiet audio/alldone.ogg 2> /dev/null > /dev/null &

sv: dist/cmms-server 

dist/cmms-server: server/*.go
	cd server && go build -o ../dist/cmms-server
	@mplayer -quiet audio/camera.oga 2> /dev/null > /dev/null &
	@ls -l dist/cmms-server

run: 
	./terminate
	@cd dist && ./cmms-server
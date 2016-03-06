all: clean run

clean:
	./terminate
	rm -rf dist server/cmms

content:
	cp assets/index.html dist/public
	cp -R assets/img dist/public
	cd assets/scss && sass --style compressed app.sass ../css/app.css
	cp -R assets/css dist/public
	cp -R assets/fonts dist/public
	temple build templates app/template.go
	cd app && gopherjs build *.go -o ../dist/public/app.js

dist: 
	##### Clean Out Dist Directory
	rm -rf dist
	mkdir -p dist/public
	mkdir -p dist/public/css dist/public/font dist/public/js
	##### Copy Our Assets
	cp assets/index.html dist/public
	cp -R assets/img dist/public
	# cd assets/scss && sass --style compressed app.sass ../css/app.css
	cd assets/scss && sass app.sass ../css/app.css
	cp -R assets/css dist/public
	cp -R assets/fonts dist/public
	cp -R assets/js dist/public
	# cp bower_components/milligram/dist/milligram.css dist/public/css
	cp bower_components/normalize.css/normalize.css dist/public/css
	cp server/config.json dist
	##### Building Client App
	temple build templates app/template.go --package main
	cd app && gopherjs build *.go -o ../dist/public/app.js -m
	# cd app && gopherjs build *.go -o ../dist/public/app.js -m
	##### Building Server App
	#cd server && go build -o ../dist/cmms-server.exe
	cd server && go build -o ../dist/cmms-server
	##### Dist directory looks like this	
	cd dist && ls -l && ls -l public/app.js && du -k .

run: dist
	###################################################################################################
	#  !!! All code passed compile and build stage !!!
	###################################################################################################
	cd dist && ./cmms-server

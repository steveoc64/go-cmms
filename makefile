all: clean dist

clean:
	rm -rf dist server/cmms


dist: 
	##### Clean Out Dist Directory
	rm -rf dist
	mkdir -p dist/public
	mkdir -p dist/public/css dist/public/font dist/public/js
	##### Copy Our Assets
	cp assets/index.html dist/public
	cp -R assets/img dist/public
	# cp -R assets/fonts dist/public
	cp -R assets/css dist/public
	##### Copy 3rd Party Assets
	cp bower_components/Materialize/dist/css/materialize.css dist/public/css
	cp bower_components/Materialize/dist/js/materialize.js dist/public/js
	cp bower_components/jquery/dist/jquery.js dist/public/js
	cp -R bower_components/Materialize/dist/font dist/public
	cp server/config.json dist
	##### Building Client App
	cd app && gopherjs build *.go -o ../dist/public/app.js
	# cd app && gopherjs build *.go -o ../dist/public/app.js -m
	##### Building Server App
	# cd server && go build -o ../dist/cmms-server.exe
	cd server && go build -o ../dist/cmms-server
	##### Dist directory looks like this	
	cd dist && ls -l && ls -l public/app.js && du -k .

run: dist
	###################################################################################################
	#  !!! All code passed compile and build stage !!!
	###################################################################################################
	cd dist && ./cmms-server

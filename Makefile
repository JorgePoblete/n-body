build-android:
	gogio -target android .

build:
	go build .

run: build
	./n-body

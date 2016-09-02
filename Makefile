SERVER_OUT := taut
DB_INIT_OUT := dbInit

run: build.server build.client
	./build/${SERVER_OUT}

build.db:
	cd dbInit && \
	go build -v -o ../build/${DB_INIT_OUT}

build.client:
	cd client && \
	npm i && \
	npm run build

build.server:
	cd server && \
	go build -i -v -o ../build/${SERVER_OUT}

install.go_deps:
	go get github.com/gorilla/websocket && \
	go get github.com/joho/godotenv && \
	go get github.com/lib/pq && \
	go get github.com/twinj/uuid

init: clean install.go_deps build.db
	./build/$(DB_INIT_OUT)

clean:
	-@rm ./build/${SERVER_OUT} ./build/${DB_INIT_OUT}

.PHONY: build.db build.client build.server init run clean

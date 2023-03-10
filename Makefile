install_dev:
	GOOS=js GOARCH=wasm go build -o frontend/public/chip-8.wasm backend/wasm/main.go

	if test -d frontend/node_modules; \
    	then echo "Node modules already installed\n"; \
    	else cd frontend; npm install --silent; \
    fi

build:
	docker buildx build --platform=linux/amd64 . --tag 273011490881.dkr.ecr.eu-west-1.amazonaws.com/chip-m8:latest
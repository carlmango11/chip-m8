build:
	GOOS=js GOARCH=wasm go build -o frontend/public/chip-8.wasm backend/main.go

	if test -d frontend/node_modules; \
    	then echo "Node modules already installed\n"; \
    	else cd frontend; npm install --silent; \
    fi

	cd frontend; npm test -- --watchAll=false;
	cd frontend; npm run build;
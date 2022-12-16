build:
	GOOS=js GOARCH=wasm go build -o frontend/public/chip-8.wasm backend/main.go

	#cd frontend; npm run test;
	cd frontend; npm run build;
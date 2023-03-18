$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -o main main.go
~\Go\Bin\build-lambda-zip.exe -o main.zip main
aws lambda update-function-code --function-name wenda --zip-file fileb://C:/Users/cakos/Code/projects/wenda/backend/main.zip

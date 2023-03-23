# Notes on building a module

## check the code
go mod tidy
go test go test ./...

## update 

git add <whatever changed>
git commit -m "simpleProperties: module for vx.y.z"
git tag vx.y.z
git push origin vx.y.z

## ...and next

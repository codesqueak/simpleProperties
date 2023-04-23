# Notes on building a module

## check the code
go mod tidy
go test ./...

## update 

git add <whatever changed> // use git add -u to add all updated files
git commit -m "simpleProperties: module for vx.y.z"
git tag vx.y.z

By default, the git push command does not transfer tags to remote servers. You need to specify it / them

git push origin vx.y.z

## ...and next

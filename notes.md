# Notes on building a module

## check the code
go mod tidy
go test ./...

## update 

create branch, add changes & merge as usual

git tag -a vx.y.z -m "This is my new version ..."

By default, the git push command does not transfer tags to remote servers. You need to specify it / them

git push origin vx.y.z

## ...and next

# Generate api go file
under app/api
## star
goctl api go -api bffs.api -dir .
goctl rpc protoc productadmin.proto --go_out=. --go-grpc_out=. --zrpc_out=.
goctl model mysql datasource --url "root:test123@tcp(localhost:3306)/ymir" --table "product_color_detail" -c .
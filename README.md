# About this project
A high-concurrency service based on gozero, using Cache Aside as the caching strategy.
The architecture follows a BFF (Backend for Frontend) + microservice (RPC) pattern.
Distributed transactions are handled using DTM(https://dtm.pub/guide/start.html)

# Development Guidelines:
- A microservice should only control table resources related to its business domain.
- Make full use of goroutines to reduce service call chains. Refer to app/bffd/internal/logic/product/productdetaillogic.go for implementation.

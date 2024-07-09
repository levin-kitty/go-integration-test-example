An example of a Go language integration test, inspired by the book [Zero To Production in Rust](https://www.zero2prod.com/index.html?country=Korea&discount_code=SEA60)

It includes the following situations
* Endpoint calls rely on external API servers and SSH servers.

For this, we use Mock HTTP Server and Mock SSH Server to perform internal integration testing.
* The Mock Servers start on random ports during testing to enable parallel testing
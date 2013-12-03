# The Router

Router is a modern replacent for [Pow from 37signals](http://pow.cx) written in Go.

It starts an http proxy on localhost:80 and depends on the host it's proxying requests to your Rails app.

It's using `/etc/hosts/` to get `your-app-name.dev` host under localhost scope.

# TODOs

* Flexible restart with txt file
* Try rackup for Rails
* Run under home user (config param)


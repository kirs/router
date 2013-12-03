# The Router

Router is a modern replacement for [Pow from 37signals](http://pow.cx), written in Go.

It starts an http proxy on localhost:80 and depends on then starts proxying requests to your Rails app.

It's using `/etc/hosts/` to get `your-app-name.dev` host under local scope.

In future it will support any applications, not only Rails: Django, Node.js or Scala.

# TODOs

* Flexible restart with txt file
* Try rackup for Rails
* Run under home user (config param)


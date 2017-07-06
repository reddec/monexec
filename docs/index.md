# MONEXEC

It's tool for controlling processes like a **supervisord** but with some important features:
* Easy to use - no dependencies. Just a single binary file pre-compilled for most major platforms
* Easy to hack - monexec can be used as a Golang library with clean and simple architecture
* Integrated with Consul - optionally, monexec can register all running processes as services and deregister on fail
* Supports gracefull and fast shutdown by signals
* Developed for used inside Docker containers
* Different strategies for processes


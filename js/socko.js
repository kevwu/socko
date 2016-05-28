// not sure how I'm going to organize this yet. none of this code has even been run, let alone tested

// example of how the interface should behave:
socko.start(function(socket, connectEvent) {
	socket.onOpen(function() {
		console.log("opened")
	})

	socket.onEvent("test-event", function(data) {
		console.log(data)
	})

	socket.onClose(function() {

	})
})
// end example

// actual object:
socko = {
	start: function(endPoint, callback) {
		// creates new socket
		var sock = new socket(endPoint)
		callback(sock, evt) // let user register event handlers first
		sock.startUp // start socket
	}
}

function socket(endPoint) {

}

socket.prototype.onOpen = function(fn) {
	var onOpen = fn
}

socket.prototype.onClose = function(fn) {
	var onClose = fn
}

socket.prototype.onEvent = function(eventName, fn) {
	this.events[eventName] = fn()
}
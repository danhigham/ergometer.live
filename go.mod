module github.com/danhigham/ergometer.live

go 1.24.0

replace github.com/danhigham/pm5 => ../usb-interface

require (
	github.com/danhigham/pm5 v0.0.0-00010101000000-000000000000
	github.com/gorilla/websocket v1.5.3
)

require (
	github.com/sstallion/go-hid v0.15.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
)

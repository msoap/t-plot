run-with-du:
	du -s ../* | go run t-plot.go -c ■

run-with-ls:
	ls -l | go run t-plot.go -c ■ -k 5

run-with-du:
	du -s ../* | go run . -c ■

run-with-ls:
	ls -l | go run . -c ■ -k 5

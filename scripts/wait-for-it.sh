RETRIES=30

until eval "${1}" &> /dev/null; do
	if [ $RETRIES -eq 0 ]; then
		1>2 echo "timeout"
		exit 1
	fi

	sleep 1
done


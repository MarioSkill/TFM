#!/bin/bash

j=$1
for (( i=0; i<j; i++ ))
do  
	docker run -d eidasclient ./test-eIDAS.sh $i &
done

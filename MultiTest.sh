#!/bin/bash
start=$(go run show_info.go | grep "r.tail" | awk '/r.tail =  / {print $NF}')
#echo $start

for i in $(seq 1 4)
do
   go run test.go &
done

wait

end=$(go run show_info.go | grep "r.tail" | awk '/r.tail =  / {print $NF}')
result=`expr $end - $start`
echo $result
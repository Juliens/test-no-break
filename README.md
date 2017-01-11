# test-no-break

```
./test-no-break -X POST -d "postdata" -interval 2s -H "Content-Type: application/json" -H "OtherHeader: Value"  http://requestb.in/xxxxx
```

Data from postdata file
```
./test-no-break -X POST -d "@postdata" -interval 2s -H "Content-Type: application/json" -H "OtherHeader: Value"  http://requestb.in/xxxxx
```

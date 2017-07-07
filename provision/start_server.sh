#!/bin/bash

if [ "x$CURRENT_DIRECTORY" == "x" ] ; then
  echo "CURRENT_DIRECTORY is not defined"
  exit 3
fi

echo [`date '+%F %T'`]: "Starting server host1-server$1:4040$1..."
mkdir -p $CURRENT_DIRECTORY/GemFire/cluster/host1/server$1
gfsh start server \
    --mcast-port=0 --server-port=4040$1 --locators=$LOCATORS_CONNECTION_STRING \
    --name=host1-server$1 --dir=$CURRENT_DIRECTORY/GemFire/cluster/host1/server$1 \
    --start-rest-api=true --http-service-port=8080 --http-service-bind-address=localhost \
    --J=-Dlog4j.configurationFile=$CURRENT_DIRECTORY/GemFire/cluster/config/log4j2-custom.xml \
    --J=-Dgemfire.statistic-sample-rate=1000 --J=-Dgemfire.statistic-sampling-enabled=true --J=-Dgemfire.statistic-archive-file=$CURRENT_DIRECTORY/GemFire/cluster/host1/statistics/host1-server$1.gfs \
    $JMX_OPTS $HEAP_OPTIONS $JAVA_GC_OPTS $GC_PRINT_OPT
echo [`date '+%F %T'`]: "Starting server host1-server$1:4040$1... Done!"


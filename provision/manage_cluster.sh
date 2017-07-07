#!/bin/bash
# Environment Variables
source /vagrant/provision/setenv.txt

# Java Options
HEAP_OPTIONS="--J=-Xmx512m --J=-Xms512m --J=-XX:+AlwaysPreTouch"
JAVA_GC_OPTS="--J=-XX:+UseParNewGC --J=-XX:+UseConcMarkSweepGC --J=-XX:CMSInitiatingOccupancyFraction=70 --J=-XX:+UseCMSInitiatingOccupancyOnly --J=-XX:+DisableExplicitGC --J=-XX:+CMSClassUnloadingEnabled"
GC_PRINT_OPT="--J=-verbose:gc --J=-Xloggc:server-gc.log --J=-XX:+PrintGCDateStamps --J=-XX:+PrintGCDetails --J=-XX:+PrintTenuringDistribution --J=-XX:+PrintGCApplicationConcurrentTime --J=-XX:+PrintGCApplicationStoppedTime"
LOCATORS_CONNECTION_STRING=""

# Save current directory
CURRENT_DIRECTORY=$GEMFIRE

# Build locators connection string
for i in $LOCATORS; do
	LOCATORS_CONNECTION_STRING="$LOCATORS_CONNECTION_STRING localhost[1010$i],"
done
LOCATORS_CONNECTION_STRING="${LOCATORS_CONNECTION_STRING//[[:blank:]]/}"

function start_cluster() {
	mkdir -p $CURRENT_DIRECTORY/GemFire/cluster/host1/statistics

	# Start Locators
	for i in $LOCATORS; do
		echo [`date '+%F %T'`]: "Starting locator host1-locator$i:[1010$i]..."
		mkdir -p $CURRENT_DIRECTORY/GemFire/cluster/host1/locator$i
		gfsh start locator \
		--name=host1-locator$i --dir=$CURRENT_DIRECTORY/GemFire/cluster/host1/locator$i --initial-heap=256m --max-heap=256m \
		--port=1010$i --locators=$LOCATORS_CONNECTION_STRING --mcast-port=0 --connect=false \
		--J=-Dgemfire.http-service-port=7070 \
		--J=-Dgemfire.jmx-manager=true --J=-Dgemfire.jmx-manager-start=true --J=-Dgemfire.jmx-manager-port=1011$i
		echo [`date '+%F %T'`]: "Starting locator host1-locator$i:[1010$i]... Done!"
	done

	# Start Servers
	export HEAP_OPTIONS=$HEAP_OPTIONS
	export JAVA_GC_OPTS=$JAVA_GC_OPTS
	export GC_PRINT_OPT=$GC_PRINT_OPT
	export CURRENT_DIRECTORY=$CURRENT_DIRECTORY
	export LOCATORS_CONNECTION_STRING=$LOCATORS_CONNECTION_STRING

	for i in $SERVERS; do
		/vagrant/provision/start_server.sh $i &
	done

	wait
}

function configure_cluster() {
	gfsh -e "connect --locator=localhost[10101]" -e "configure pdx --read-serialized=true" -e "create region --name=replicatedRegion --type=REPLICATE"
}

function stop_cluster() {
	# Stop Servers
	gfsh -e "connect --locator=localhost[10101]" -e "shutdown"

	# Stop Locators
	for i in $LOCATORS; do
		gfsh stop locator --dir=$CURRENT_DIRECTORY/GemFire/cluster/host1/locator$i
	done
}

function clean_cluster() {
	# Kill Processes
	for KILLPID in `ps ax | grep 'gemfire' | awk ' { print $1; }'`; do
		kill -9 $KILLPID;
	done

	# GFSH Files
	rm $CURRENT_DIRECTORY/GemFire/gfsh/*

	# Remove Locators Directories
	for i in $LOCATORS; do
		rm -R $CURRENT_DIRECTORY/GemFire/cluster/host1/locator$i
	done

	# Remove statistics
	rm -R $CURRENT_DIRECTORY/GemFire/cluster/host1/statistics

	# Remove Server Directories (Host1)
	for i in $SERVERS; do
		rm -R $CURRENT_DIRECTORY/GemFire/cluster/host1/server$i
	done
}

selection=
until [ "$selection" = "5" ]; do
	echo "##############################################################################################################################################################################"
	echo "Select an option:"
	echo "	1 - Start Cluster."
	echo "	2 - Configure Cluster."
	echo "	3 - Stop Cluster."
	echo "	4 - Clean Cluster."
	echo "	5 - Exit."
	echo "##############################################################################################################################################################################"
	echo -n "Enter a choice: "
    read selection
    echo ""
	case $selection in
        1)
			start_cluster
            ;;
        2)
			configure_cluster
            ;;
		3)
			stop_cluster
            ;;
		4)
			clean_cluster
            ;;
        5)
            exit
            ;;
        *) echo 'Invalid option, please select an option between 1 and 5';;
    esac
done


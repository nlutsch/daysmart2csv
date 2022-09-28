echo "Stopping Service"
service daysmart2csv stop
echo "Service Stopped"
echo "Pulling latest version"
git pull
echo "Building Application"
go build .
echo "Starting Service"
service daysmart2csv start
echo "Service Started"
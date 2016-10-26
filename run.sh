#!/bin/bash
#sed 's,127.0.0.1:9001,'"${DB_PORT#*//}"',g' -i conf/app.conf
sed 's,runmode.*,'runmode="${RUN_MODE:-dev}"',g' -i conf/app.conf
sed 's,httpport.*,'httpport="${HTTP_PORT:-8010}"',g' -i conf/app.conf
sed 's,mongodb.*,'mongodb="${MONGO_DB:-uic}"',g' -i conf/app.conf
sed 's,mongohost.*,'mongohost="${MONGO_HOST:-db:27017}"',g' -i conf/app.conf
sed 's,mongouser.*,'mongouser="${MONGO_USER:-uic}"',g' -i conf/app.conf
sed 's,mongopassword.*,'mongopassword="${MONGO_PASSWORD:-Goo7ieho}"',g' -i conf/app.conf
sed 's,uicdomain.*,'uicdomain="${UIC_DOMAIN:-accounts.igenetech.com}"',g' -i conf/app.conf
ssoauth

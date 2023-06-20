#!/bin/bash
set -e

# get commit_id
COMMIT_ID="main"
if [ $RISINGWAVE_DASHBOARD_COMMIT_ID ];then
	COMMIT_ID=$RISINGWAVE_DASHBOARD_COMMIT_ID
else
	echo "Environment variable \"RISINGWAVE_DASHBOARD_COMMIT_ID\" not found, use \"main\" as commit_id"
fi

common_url="https://raw.githubusercontent.com/risingwavelabs/risingwave/${COMMIT_ID}/grafana/common.py"
shell_url="https://raw.githubusercontent.com/risingwavelabs/risingwave/${COMMIT_ID}/grafana/generate.sh"
python_url="https://raw.githubusercontent.com/risingwavelabs/risingwave/${COMMIT_ID}/grafana/risingwave-dev-dashboard.dashboard.py"
python_url_user="https://raw.githubusercontent.com/risingwavelabs/risingwave/${COMMIT_ID}/grafana/risingwave-user-dashboard.dashboard.py"

# download ./generate.sh and risingwave-dev-dashboard.dashboard.py
wget $common_url -O "./common.py"
wget $shell_url -O "./generate_ori.sh"
wget $python_url -O "risingwave-dev-dashboard.dashboard.py"
wget $python_url_user -O "risingwave-user-dashboard.dashboard.py"

sed '/cp/d' ./generate_ori.sh > ./generate.sh

chmod +x ./generate.sh

# generate risingwave-dev-dashboard.json
DASHBOARD_NAMESPACE_FILTER_ENABLED=true DASHBOARD_RISINGWAVE_NAME_FILTER_ENABLED=true DASHBOARD_SOURCE_UID="prometheus" ./generate.sh

# replace for cloud deployment, will read risingwave-dev-dashboard.json and write the result into risingwave-dev-dashboard_new.json
python3 ./convert.py

# remove intermediate files
rm risingwave-dev-dashboard.dashboard.py
rm risingwave-user-dashboard.dashboard.py
rm risingwave-dev-dashboard.gen.json
rm risingwave-user-dashboard.gen.json
rm risingwave-dev-dashboard.json
rm risingwave-user-dashboard.json
rm ./generate_ori.sh
rm ./generate.sh
rm ./common.py

# rename
mv risingwave-dev-dashboard_new.json risingwave-dev-dashboard.json

echo "risingwave-dev-dashboard.json updated"

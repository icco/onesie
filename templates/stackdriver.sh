## Install Monitoring Agent
REPO_HOST=${STACKDRIVER_REPO_HOST:-repo.stackdriver.com}
APP_HOST=app.stackdriver.com
CODENAME="$(lsb_release -sc)"

if [[ -f /etc/os-release ]]; then
  . /etc/os-release
fi

curl -s -S -f -o /etc/apt/sources.list.d/stackdriver.list "https://${REPO_HOST}/${CODENAME}.list"
curl -s -f https://${APP_HOST}/RPM-GPG-KEY-stackdriver | apt-key add -

## Now Install Logging
REPO_HOST='packages.cloud.google.com'
CLOUD_LOGGING_DOCS_URL="https://cloud.google.com/logging/docs/agent"
REPO_NAME="google-cloud-logging-wheezy${REPO_SUFFIX+-${REPO_SUFFIX}}"
cat > /etc/apt/sources.list.d/google-cloud-logging.list <<EOM
deb http://${REPO_HOST}/apt ${REPO_NAME} main
EOM
curl -s -f https://${REPO_HOST}/apt/doc/apt-key.gpg | apt-key add -

aptitude update
DEBIAN_FRONTEND=noninteractive aptitude -y -q install stackdriver-agent google-fluentd google-fluentd-catch-all-config

FROM debian
WORKDIR .
EXPOSE 9090/tcp
# Copy our static executable.
RUN apt-get update && apt-get install -y curl jq zip &&  apt-get clean
RUN export GHRUNNER=$(curl -s https://api.github.com/repos/mkrakowitzer/githubrunner_exporter/releases/latest | jq -r '.assets[].browser_download_url') && curl -L $GHRUNNER -O
RUN unzip githubrunner_exporter.zip && rm -f githubrunner_exporter.zip && chmod 700 ./githubrunner_exporter
ENTRYPOINT ["/githubrunner_exporter"]

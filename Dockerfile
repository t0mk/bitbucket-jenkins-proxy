FROM busybox

ADD bitbucket-jenkins-proxy /bitbucket-jenkins-proxy
RUN chmod +x /bitbucket-jenkins-proxy
CMD ["/bitbucket-jenkins-proxy"]


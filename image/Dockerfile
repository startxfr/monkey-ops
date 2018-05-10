FROM scratch

MAINTAINER Startx - Christophe LARUE <dev@startx.fr>

ENV CHAOS_VERSION="0.1.0" \
    SUMMARY="Chaos monkey container for kubernetes or openshift environements" \
    DESCRIPTION="Simple tool used for testing application or microservices resiliency by applying the chaos monkey technique."

LABEL name="startx/monkey-ops" \
      summary="$SUMMARY" \
      description="$SUMMARY" \
      version="$CHAOS_VERSION" \
      maintainer="Startx - Christophe LARUE <dev@startx.fr>" \
      usage="docker run startx/monkey-ops" \
      io.k8s.description="$SUMMARY" \
      io.k8s.display-name="monkey-ops" \
      io.openshift.tags="chaos,test,monkey" \
      io.openshift.wants="kubernetes,k8s" \
      io.openshift.non-scalable="true" \
      io.openshift.min-memory="20Mi" \
      io.openshift.min-cpu="100m" \
      fr.startx.component="monkey-ops"

ADD monkey-ops /

CMD ["/monkey-ops"]
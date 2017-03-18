FROM jbonachera/arch
COPY release/homie-controller /usr/bin/
RUN useradd -r homie -d /var/lib/homie && \
    mkdir /var/cache/homie-controller && \
    chown homie: /var/cache/homie-controller
USER homie
CMD ["/usr/bin/homie-controller"]

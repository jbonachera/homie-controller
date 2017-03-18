FROM jbonachera/arch
COPY release/homie-controller /usr/bin/
RUN useradd -r homie -d /var/lib/homie
USER homie
CMD ["/usr/bin/homie-controller"]

FROM joeshaw/busybox-nonroot
EXPOSE 9999

COPY build/tcpdam /bin/tcpdam

USER nobody
CMD ["tcpdam"]

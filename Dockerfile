FROM alpine

WORKDIR /root

ENV WDIR_DOCKER=true

# Should be change
ENV SHOWHIDDENFILES=false

# Should not be changed
ENV FILEPATH=/mnt \
    PORT=8080 \ 
    LOGPATH=/log \
    TEMPLATE=compress

COPY compress ./compress
COPY wdir .

EXPOSE 8080

CMD [ "./wdir" ]

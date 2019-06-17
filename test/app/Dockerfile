FROM iron/go

WORKDIR /bin

ADD bin/app /bin/
RUN chmod +x ./app
ENTRYPOINT ["./app"]
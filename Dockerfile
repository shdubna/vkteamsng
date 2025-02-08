FROM alpine
COPY vkteamsng /bin/
ENTRYPOINT ["vkteamsng"]
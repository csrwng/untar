FROM scratch
EXPOSE 9080
COPY Untar /Untar
CMD ["/Untar"]

FROM golang:onbuild
# Copy the local package files to the container's workspace.
#ADD . /go/src/github.com/orenmazor/triscuits*/

# the HMAC your clients will need
ENV TRISCUITS_HMAC="foo"
ENV TABLEAU_URL=https://zombo.com

#RUN go install github.com/orenmazor/triscuits

# Run the outyet command by default when the container starts.
#ENTRYPOINT /go/bin/triscuits

# Document that the service listens on port 8080.
EXPOSE 31337

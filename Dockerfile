FROM golang:onbuild
# the HMAC your clients will need
ENV TRISCUITS_HMAC="foo"
ENV TABLEAU_URL=https://zombo.com
EXPOSE 31337
